// @flow

import React from "react";
import PageLayout from "../../style-guide/page-layout";
import { Display3 } from "baseui/typography";
import { useStyletron } from "baseui";
import { Spacer20, Spacer40 } from "../../style-guide/spacer";
import { Input } from "baseui/input";
import { FormControl } from "baseui/form-control";
import TemplatePicker from "./../../components/template-picker";
import { Button } from "baseui/button";
import { KIND as NotificationKind, Notification } from "baseui/notification";
import useFetch from "use-http";
import { useHistory, useParams } from "react-router-dom";
import type { Campaign } from "../../lib/types/campaign";
import { StyledSpinnerNext as Spinner } from "baseui/spinner";
import DeleteCampaignButton from "./delete-button";

const EditCampaignPage = () => {
  const { uuid: campaignUUID } = useParams();
  const history = useHistory();
  const [css] = useStyletron();
  const campaignAPI = useFetch<Campaign>(`/campaigns/${campaignUUID}`);

  const [name, setName] = React.useState("");
  const [fromName, setFromName] = React.useState("");
  const [fromEmail, setFromEmail] = React.useState("");
  const [templateUUID, setTemplateUUID] = React.useState("");

  React.useEffect(() => {
    async function loadData() {
      const campaign: Campaign = await campaignAPI.get();

      if (campaignAPI.response.ok) {
        setName(campaign.Name);
        setFromName(campaign.FromName);
        setFromEmail(campaign.FromEmail);
        setTemplateUUID(campaign.TemplateUUID);
      }
    }

    loadData();
  }, []); /* eslint-disable-line react-hooks/exhaustive-deps */

  const isValid = React.useCallback(
    () => name && fromName && fromEmail && templateUUID,
    [fromEmail, fromName, name, templateUUID]
  );

  const onSubmit = React.useCallback(async () => {
    if (!isValid()) {
      return;
    }
    const campaign = await campaignAPI.post({
      name,
      template_uuid: templateUUID,
      from_name: fromName,
      from_email: fromEmail,
    });

    if (campaignAPI.response.ok) {
      history.push(`/campaigns/${campaign.UUID}/review`);
    }
  }, [campaignAPI, fromEmail, fromName, history, isValid, name, templateUUID]);

  if (campaignAPI.loading) {
    return (
      <PageLayout>
        <Spinner />
      </PageLayout>
    );
  }

  if (campaignAPI.error) {
    return <PageLayout>Failed to load this page. Please try again.</PageLayout>;
  }

  if (!campaignAPI.data) {
    return null;
  }

  return (
    <PageLayout>
      <Display3>Campaign Info</Display3>
      <Spacer40 />
      <div className={css({ maxWidth: "400px" })}>
        <FormControl label="Campaign Name">
          <Input
            placeholder="Welcome Email #1"
            value={name}
            onChange={(e) => {
              setName(e.target.value);
            }}
          />
        </FormControl>
        <Spacer20 />
        <FormControl label="From Name">
          <Input
            placeholder="Jane Doe"
            value={fromName}
            onChange={(e) => {
              setFromName(e.target.value);
            }}
          />
        </FormControl>
        <Spacer20 />
        <FormControl label="From Email Address">
          <Input
            placeholder="jane@company.com"
            value={fromEmail}
            onChange={(e) => {
              setFromEmail(e.target.value);
            }}
          />
        </FormControl>
      </div>
      <Spacer20 />

      <FormControl label="Pick a Template">
        <TemplatePicker
          initialTemplateUUID={campaignAPI.data.TemplateUUID}
          onSelect={(template) => setTemplateUUID(template.UUID)}
        />
      </FormControl>
      {campaignAPI.error && (
        <Notification
          kind={NotificationKind.negative}
          overrides={{
            Body: { style: { width: "auto" } },
          }}
        >
          Failed to update your campaign. Please try again.
        </Notification>
      )}
      <div
        className={css({ display: "flex", justifyContent: "space-between" })}
      >
        <DeleteCampaignButton campaign={campaignAPI.data} />
        <Button
          disabled={!isValid() || campaignAPI.loading}
          isLoading={campaignAPI.loading}
          onClick={onSubmit}
        >
          Save Changes & Review
        </Button>
      </div>
      <Spacer40 />
    </PageLayout>
  );
};

export default EditCampaignPage;
