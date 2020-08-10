// @flow
import * as React from "react";
import PageLayout from "../../style-guide/page-layout";
import { Card, StyledAction, StyledBody } from "baseui/card";
import useFetch from "use-http";
import type { CreditPack, CreditPackProduct } from "../../lib/types/credit";
import { StyledSpinnerNext as Spinner } from "baseui/spinner";
import { KIND as NotificationKind, Notification } from "baseui/notification";
import { Radio, RadioGroup } from "baseui/radio";
import { Paragraph3 } from "baseui/typography";
import { Spacer20 } from "../../style-guide/spacer";
import { getProductDescription } from "./product-description";
import { Button } from "baseui/button";

import { loadStripe } from "@stripe/stripe-js";
import { Elements } from "@stripe/react-stripe-js";
import PaymentForm from "./payment-form";
import { Link, useLocation } from "react-router-dom";

const Step = {
  SELECT_CREDIT_PACK: 1,
  PAY: 2,
  COMPLETE: 3,
};

let stripe;

const CreditsPage = () => {
  const location = useLocation();
  const queryParams = new URLSearchParams(location.search);

  const [step, setStep] = React.useState(Step.SELECT_CREDIT_PACK);
  const [product, setProduct] = React.useState<?CreditPackProduct>(null);

  const fetchProducts = useFetch<{ Products: CreditPackProduct[] }>(
    `/credit/products`
  );
  const purchaseAPI = useFetch<{
    StripePublicKey: string,
    StripeClientSecret: string,
    Pack: CreditPack,
  }>("/credit/packs/purchase");

  React.useEffect(() => {
    async function loadData() {
      const { Products: products } = await fetchProducts.get();

      if (fetchProducts.response.ok && queryParams.has("product")) {
        setProduct(products.find((p) => p.ID === queryParams.get("product")));
      }
    }

    loadData();
  }, []); /* eslint-disable-line react-hooks/exhaustive-deps */

  if (fetchProducts.loading) {
    return (
      <PageLayout>
        <Spinner />
      </PageLayout>
    );
  }

  if (!fetchProducts.response.ok) {
    return (
      <PageLayout>
        <Notification kind={NotificationKind.negative}>
          Failed to this page
        </Notification>
      </PageLayout>
    );
  }

  return (
    <PageLayout>
      <Card
        title={"Credit Packs"}
        overrides={{ Root: { style: { maxWidth: "500px", margin: "0 auto" } } }}
      >
        {step === Step.SELECT_CREDIT_PACK && (
          <>
            <StyledBody>
              <Paragraph3>Pick the pack that best suits your needs</Paragraph3>
              <Spacer20 />
              <RadioGroup
                name="credit-pack-products"
                value={product ? product.ID : ""}
                onChange={(e) =>
                  setProduct(
                    fetchProducts.response.data.Products.find(
                      (p) => p.ID === e.target.value
                    )
                  )
                }
              >
                {fetchProducts.response.data.Products.map(
                  (product: CreditPackProduct) => (
                    <Radio
                      key={product.ID}
                      value={product.ID}
                      description={getProductDescription(product)}
                      overrides={{
                        Description: {
                          style: { whiteSpace: "nowrap", marginBottom: "20px" },
                        },
                      }}
                    >
                      {product.Description}
                    </Radio>
                  )
                )}
              </RadioGroup>
            </StyledBody>
            <StyledAction>
              <Button
                isLoading={purchaseAPI.loading}
                onClick={async () => {
                  if (!product) {
                    throw new Error("product is not set");
                  }

                  await purchaseAPI.post({
                    product_id: product.ID,
                  });

                  if (!purchaseAPI.response.ok) {
                    // todo: show error
                    return;
                  }

                  stripe = await loadStripe(
                    purchaseAPI.response.data.StripePublicKey
                  );

                  setStep(Step.PAY);
                }}
                overrides={{ BaseButton: { style: { width: "100%" } } }}
              >
                Next
              </Button>
            </StyledAction>
          </>
        )}

        {step === Step.PAY && product ? (
          <Elements stripe={stripe}>
            <PaymentForm
              product={product}
              pack={purchaseAPI.response.data.Pack}
              stripeClientSecret={purchaseAPI.response.data.StripeClientSecret}
              onSuccess={() => setStep(Step.COMPLETE)}
            />
          </Elements>
        ) : null}

        {step === Step.COMPLETE && (
          <>
            <StyledBody>
              <Spacer20 />
              <Notification
                kind={NotificationKind.positive}
                overrides={{
                  Body: { style: { width: "auto" } },
                }}
              >
                The payment has been processed successfully.
              </Notification>
            </StyledBody>
            <StyledAction>
              <Button
                $as={Link}
                to={"/campaigns"}
                overrides={{ BaseButton: { style: { width: "100%" } } }}
              >
                Campaigns
              </Button>
            </StyledAction>
          </>
        )}
      </Card>
    </PageLayout>
  );
};

export default CreditsPage;
