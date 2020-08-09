// @flow

export type CreditPack = {
  UUID: string,
};

export type CreditPackProduct = {
  ID: string,
  Description: string,
  UseLimit: ?number,
  DurationMS: ?number,
  PriceUSD: number,
};
