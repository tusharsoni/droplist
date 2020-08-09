// @flow
import { Duration } from "luxon";
import type { CreditPackProduct } from "../../lib/types/credit";

export function getProductDescription(product: CreditPackProduct): string {
  const price = `$${product.PriceUSD / 100}`;
  const useLimit = product.UseLimit
    ? product.UseLimit === 1
      ? "1 campaign"
      : `${product.UseLimit} campaigns`
    : "unlimited campaigns";
  const expiresAt = product.DurationMS
    ? `${Duration.fromMillis(product.DurationMS).toFormat("d")} day validity`
    : "never expires";

  return `${price}, ${useLimit}, ${expiresAt}`;
}
