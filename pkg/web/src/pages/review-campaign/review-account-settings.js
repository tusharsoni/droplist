// @flow
import * as React from "react";
import useFetch from "use-http";
import type { Profile } from "../../lib/types/profile";
import type { CreditPack } from "../../lib/types/credit";
import { Label2, ParagraphSmall } from "baseui/typography";
import { KIND as NotificationKind, Notification } from "baseui/notification";
import { StyledLink } from "baseui/link";
import { Link } from "react-router-dom";

const ReviewAccountSettings = () => {
  const getProfileAPI = useFetch<Profile>(`/profile`, {}, []);
  const getCreditPacksAPI = useFetch<CreditPack[]>("/credit/packs", {}, []);
  const getCreditProductsAPI = useFetch<{ Enabled: boolean }>(
    "/credit/products",
    {},
    []
  );

  if (
    getProfileAPI.loading ||
    getCreditPacksAPI.loading ||
    getCreditProductsAPI.loading
  ) {
    return (
      <>
        <Label2>Account Settings</Label2>
        <ParagraphSmall>Checking</ParagraphSmall>
      </>
    );
  }

  if (
    (!getProfileAPI.response.ok && getProfileAPI.response.status !== 404) ||
    !getCreditPacksAPI.response.ok ||
    !getCreditProductsAPI.response.ok
  ) {
    return (
      <>
        <Label2>Account Settings</Label2>
        <ParagraphSmall>Failed to verify account settings</ParagraphSmall>
      </>
    );
  }

  if (
    getCreditProductsAPI.response.data &&
    getCreditProductsAPI.response.data.Enabled &&
    !(getCreditPacksAPI.response.data || []).length
  ) {
    return (
      <>
        <Label2>Account Settings</Label2>
        <ParagraphSmall>
          <Notification
            kind={NotificationKind.negative}
            overrides={{
              Body: { style: { width: "auto" } },
            }}
          >
            <>
              You don't have any credits to publish this campaign.{" "}
              <StyledLink $as={Link} to={`/credits`}>
                Buy here
              </StyledLink>
            </>
          </Notification>
        </ParagraphSmall>
      </>
    );
  }

  if (getProfileAPI.response.status === 404) {
    return (
      <>
        <Label2>Account Settings</Label2>
        <ParagraphSmall>
          <Notification
            kind={NotificationKind.negative}
            overrides={{
              Body: { style: { width: "auto" } },
            }}
          >
            <>
              AWS settings have not been configured.{" "}
              <StyledLink $as={Link} to={`/profile`}>
                Fix here.
              </StyledLink>
            </>
          </Notification>
        </ParagraphSmall>
      </>
    );
  }

  return (
    <>
      <Label2>Account Settings</Label2>
      <ParagraphSmall>All good!</ParagraphSmall>
    </>
  );
};

export default ReviewAccountSettings;
