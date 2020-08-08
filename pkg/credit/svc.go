package credit

import (
	"context"
	"droplist/pkg/ptr"
	"errors"
	"sort"
	"time"

	"github.com/google/uuid"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
	"github.com/tusharsoni/copper/cerror"
	"go.uber.org/fx"
)

var ErrNoValidPack = errors.New("user does not have any valid packs")

type Svc interface {
	PurchaseIntent(ctx context.Context, userUUID, productID string) (*Pack, string, error)
	CompletePurchase(ctx context.Context, packUUID string) error

	GetValidPacks(ctx context.Context, userUUID string) ([]Pack, error)
	UseBestAvailableCredit(ctx context.Context, userUUID, campaignUUID string) error
}

type SvcParams struct {
	fx.In

	Repo   Repo
	Config Config
}

func NewSvc(p SvcParams) Svc {
	stripe.Key = p.Config.StripeKey

	return &svc{
		repo:   p.Repo,
		config: p.Config,
	}
}

type svc struct {
	repo   Repo
	config Config
}

func (s *svc) PurchaseIntent(ctx context.Context, userUUID, productID string) (*Pack, string, error) {
	product, ok := s.config.Products[productID]
	if !ok {
		return nil, "", cerror.New(nil, "invalid product id", map[string]interface{}{
			"productID": productID,
		})
	}

	intent, err := paymentintent.New(&stripe.PaymentIntentParams{
		Description: stripe.String(product.Description),
		Amount:      stripe.Int64(product.PriceUSD),
		Currency:    stripe.String(string(stripe.CurrencyUSD)),
	})
	if err != nil {
		return nil, "", cerror.New(err, "failed to create stripe payment intent", map[string]interface{}{
			"productID": productID,
		})
	}

	pack := &Pack{
		UUID:            uuid.New().String(),
		UserUUID:        userUUID,
		PaymentID:       intent.ID,
		PaymentComplete: false,
		ProductID:       productID,
		UseLimit:        product.UseLimit,
	}

	if product.Duration != nil {
		pack.ExpiresAt = ptr.Time(time.Now().Add(*product.Duration))
	}

	err = s.repo.AddPack(ctx, pack)
	if err != nil {
		return nil, "", cerror.New(err, "failed to create pack", map[string]interface{}{
			"userUUID":  userUUID,
			"paymentID": intent.ID,
			"productID": productID,
		})
	}

	return pack, intent.ClientSecret, nil
}

func (s *svc) CompletePurchase(ctx context.Context, packUUID string) error {
	pack, err := s.repo.GetPackByUUID(ctx, packUUID)
	if err != nil {
		return cerror.New(err, "failed to get pack", map[string]interface{}{
			"uuid": packUUID,
		})
	}

	intent, err := paymentintent.Get(pack.PaymentID, nil)
	if err != nil {
		return cerror.New(err, "failed to get payment intent", map[string]interface{}{
			"paymentID": pack.PaymentID,
		})
	}

	if intent.Status != stripe.PaymentIntentStatusSucceeded {
		return cerror.New(nil, "payment status must be 'succeeded'", map[string]interface{}{
			"paymentID":     pack.PaymentID,
			"paymentStatus": intent.Status,
		})
	}

	pack.PaymentComplete = true

	err = s.repo.AddPack(ctx, pack)
	if err != nil {
		return cerror.New(err, "failed to update pack payment status", map[string]interface{}{
			"packUUID": pack.UUID,
		})
	}

	return nil
}

func (s *svc) UseBestAvailableCredit(ctx context.Context, userUUID, campaignUUID string) error {
	if !s.config.Enabled {
		return nil
	}

	pack, err := s.GetBestValidPack(ctx, userUUID)
	if err != nil {
		return cerror.New(err, "failed to get best valid pack", map[string]interface{}{
			"userUUID": userUUID,
		})
	}

	log := &UseLog{
		UUID:         uuid.New().String(),
		PackUUID:     pack.UUID,
		CampaignUUID: campaignUUID,
	}

	err = s.repo.AddUseLog(ctx, log)
	if err != nil {
		return cerror.New(err, "failed to insert use log", map[string]interface{}{
			"packUUID":     pack.UUID,
			"campaignUUID": campaignUUID,
		})
	}

	return nil
}

func (s *svc) GetValidPacks(ctx context.Context, userUUID string) ([]Pack, error) {
	packs, err := s.repo.FindPacksByUserUUID(ctx, userUUID)
	if err != nil {
		return nil, cerror.New(err, "failed to get packs", map[string]interface{}{
			"userUUID": userUUID,
		})
	}

	packUUIDs := make([]string, len(packs))
	for i := range packs {
		packUUIDs[i] = packs[i].UUID
	}

	useCountByPackUUID, err := s.repo.UseLogCountByPackUUIDs(ctx, packUUIDs)
	if err != nil {
		return nil, cerror.New(err, "failed to get use log count", map[string]interface{}{
			"packUUIDs": nil,
		})
	}

	// a pack is considered valid if:
	// - has completed payment
	// - has not reached its use limit
	// - is not expired
	validPacks := make([]Pack, 0)
	for i := range packs {
		useCount, _ := useCountByPackUUID[packs[i].UUID]
		didReachUseLimit := packs[i].UseLimit != nil && useCount >= *packs[i].UseLimit
		isExpired := packs[i].ExpiresAt != nil && packs[i].ExpiresAt.Before(time.Now())

		if packs[i].PaymentComplete && !isExpired && !didReachUseLimit {
			validPacks = append(validPacks, packs[i])
		}
	}

	return validPacks, nil
}

func (s *svc) GetBestValidPack(ctx context.Context, userUUID string) (*Pack, error) {
	validPacks, err := s.GetValidPacks(ctx, userUUID)
	if err != nil {
		return nil, cerror.New(err, "failed to get valid packs", nil)
	}

	if len(validPacks) == 0 {
		return nil, ErrNoValidPack
	}

	sort.Slice(validPacks, func(i, j int) bool {
		// if a pack is expiring, they get priority
		if validPacks[i].ExpiresAt != nil && validPacks[j].ExpiresAt != nil {
			return validPacks[i].ExpiresAt.Before(*validPacks[j].ExpiresAt)
		}
		if validPacks[i].ExpiresAt != nil {
			return true
		}
		if validPacks[j].ExpiresAt != nil {
			return false
		}

		// the pack with a smaller use limit gets priority
		if validPacks[i].UseLimit != nil && validPacks[j].UseLimit != nil {
			return *validPacks[i].UseLimit < *validPacks[j].UseLimit
		}
		if validPacks[i].UseLimit != nil {
			return true
		}
		if validPacks[j].UseLimit != nil {
			return false
		}

		// both packs don't have any limitations, order doesn't matter
		return true
	})

	return &validPacks[0], nil
}
