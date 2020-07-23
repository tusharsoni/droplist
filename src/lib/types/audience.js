// @flow

export type AudienceSummary = {
  TotalContacts: number,
  SubscribedContacts: number,
};

type ContactStatus = "SUBSCRIBE" | "UNSUBSCRIBED";

export type Contact = {
  UUID: string,
  CreatedAt: string,
  UpdatedAt: string,
  CreatedBy: string,
  Email: string,
  Status: ContactStatus,
  Params: string,
};

export type CreateContactResult = {
  email: string,
  success: boolean,
};
