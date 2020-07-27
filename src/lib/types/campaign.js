// @flow

export type Campaign = {
  UUID: string,
  UpdatedAt: string,
  TemplateUUID: string,
  Name: string,
  FromName: string,
  FromEmail: string,
  State: "DRAFT" | "PUBLISHED",
};

export type CampaignStats = {
  Queued: number,
  Sent: number,
  Failed: number,
  Opens: number,
  Clicks: number,
  OpenRate: number,
  ClickRate: number,
};
