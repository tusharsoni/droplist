// @flow
import * as React from "react";
import { StyledAction, StyledBody } from "baseui/card";
import { LabelSmall, Paragraph3 } from "baseui/typography";
import { Spacer20, Spacer8 } from "../../style-guide/spacer";
import { Radio, RadioGroup } from "baseui/radio";
import { getProductDescription } from "./product-description";
import { FormControl } from "baseui/form-control";
import { CardElement, useElements, useStripe } from "@stripe/react-stripe-js";
import { Input } from "baseui/input";
import { Select } from "baseui/select";
import countries from "../../lib/countries";
import { Button } from "baseui/button";
import type { CreditPack, CreditPackProduct } from "../../lib/types/credit";
import { useStyletron } from "baseui";
import useFetch from "use-http";
import { KIND as NotificationKind, Notification } from "baseui/notification";

type Props = {
  pack: CreditPack,
  product: CreditPackProduct,
  stripeClientSecret: string,

  onSuccess: () => void,
};

const PaymentForm = (props: Props) => {
  const [css] = useStyletron();
  const stripe = useStripe();
  const elements = useElements();

  const [fullName, setFullName] = React.useState("");
  const [address1, setAddress1] = React.useState("");
  const [address2, setAddress2] = React.useState("");
  const [city, setCity] = React.useState("");
  const [state, setState] = React.useState("");
  const [zipCode, setZipCode] = React.useState("");
  const [country, setCountry] = React.useState([
    { id: "US", label: countries.find((c) => c.iso2 === "US").name },
  ]);

  const completePurchaseAPI = useFetch(
    `/credit/packs/${props.pack.UUID}/purchase`
  );

  const [loading, setLoading] = React.useState(false);
  const [error, setError] = React.useState<?string>(null);

  return (
    <>
      <StyledBody>
        <Paragraph3>Review and add your credit card details</Paragraph3>
        <Spacer20 />
        <RadioGroup value={props.product.ID}>
          {/* $FlowFixMe: Flow error for only having a single Radio element even though it works */}
          <Radio
            value={props.product.ID}
            description={getProductDescription(props.product)}
            overrides={{
              Description: {
                style: { whiteSpace: "nowrap", marginBottom: "20px" },
              },
            }}
          >
            {props.product.Description}
          </Radio>
        </RadioGroup>

        <Spacer20 />

        <FormControl label="Credit Card">
          <div
            className={css({
              height: "44px",
              padding: "14px 20px",
              backgroundColor: "#EEEEEE",
            })}
          >
            <CardElement id="card-element" />
          </div>
        </FormControl>

        <Spacer20 />

        <LabelSmall>Billing Address</LabelSmall>
        <Spacer8 />
        <Input
          placeholder="Full Name"
          value={fullName}
          onChange={(e) => setFullName(e.target.value)}
        />
        <Spacer8 />
        <Input
          placeholder="Street Address"
          value={address1}
          onChange={(e) => setAddress1(e.target.value)}
        />
        <Spacer8 />
        <Input
          placeholder="Apt, Building (Optional)"
          value={address2}
          onChange={(e) => setAddress2(e.target.value)}
        />
        <Spacer8 />
        <Input
          placeholder="City"
          value={city}
          onChange={(e) => setCity(e.target.value)}
          overrides={{ Root: { style: { width: "50%" } } }}
        />
        <Spacer8 />
        <Input
          placeholder="State"
          value={state}
          onChange={(e) => setState(e.target.value)}
          overrides={{ Root: { style: { width: "50%" } } }}
        />
        <Spacer8 />
        <div className={css({ display: "flex" })}>
          <div className={css({ flex: 2 })}>
            <Input
              placeholder="Zip Code"
              value={zipCode}
              onChange={(e) => setZipCode(e.target.value)}
            />
          </div>
          <Spacer8 />
          <div className={css({ flex: 3 })}>
            <Select
              options={countries.map((c) => ({
                id: c.iso2,
                label: c.name,
              }))}
              value={country}
              placeholder="Country"
              onChange={(params) => setCountry(params.value)}
              clearable={false}
            />
          </div>
        </div>

        <Spacer20 />
      </StyledBody>

      <StyledAction>
        {error && (
          <Notification
            kind={NotificationKind.negative}
            overrides={{
              Body: { style: { width: "auto" } },
            }}
          >
            {error}
          </Notification>
        )}

        <Button
          isLoading={loading}
          disabled={
            !(
              fullName &&
              address1 &&
              city &&
              state &&
              zipCode &&
              country.length
            )
          }
          onClick={async () => {
            setLoading(true);
            setError(null);

            const stripeResp = await stripe.confirmCardPayment(
              props.stripeClientSecret,
              {
                payment_method: {
                  card: elements.getElement(CardElement),
                  billing_details: {
                    name: fullName,
                    address: {
                      line1: address1,
                      line2: address2,
                      city,
                      state,
                      postal_code: zipCode,
                      country: country[0].id,
                    },
                  },
                },
              }
            );

            if (stripeResp.error) {
              setError(stripeResp.error.message);
              setLoading(false);
              return;
            }

            await completePurchaseAPI.post();

            if (!completePurchaseAPI.response.ok) {
              setError(
                "Failed to complete your purchase. Please contact us to continue."
              );
              setLoading(false);
              return;
            }

            props.onSuccess();
          }}
          overrides={{ BaseButton: { style: { width: "100%" } } }}
        >
          Pay ${props.product.PriceUSD / 100}
        </Button>
      </StyledAction>
    </>
  );
};

export default PaymentForm;
