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
