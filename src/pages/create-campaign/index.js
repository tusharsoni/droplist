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
import type { Template } from "../../lib/types/content";
import { KIND as NotificationKind, Notification } from "baseui/notification";
import useFetch from "use-http";
import { useHistory } from "react-router-dom";

const CreateCampaignPage = () => {
  const history = useHistory();
  const createCampaignAPI = useFetch<Template>(`/campaigns`);
  const [css] = useStyletron();
  const [name, setName] = React.useState("");
  const [fromName, setFromName] = React.useState("");
  const [fromEmail, setFromEmail] = React.useState("");
  const [template, setTemplate] = React.useState<?Template>(null);

  const isValid = React.useCallback(
    () => name && fromName && fromEmail && template,
    [fromEmail, fromName, name, template]
  );

  const onSubmit = React.useCallback(async () => {
    if (!isValid() || !template) {
      return;
    }
    const campaign = await createCampaignAPI.post({
      name,
      template_uuid: template.UUID,
      from_name: fromName,
      from_email: fromEmail,
    });

    if (createCampaignAPI.response.ok) {
      history.push(`/campaigns/${campaign.UUID}/review`);
    }
  }, [
    createCampaignAPI,
    fromEmail,
    fromName,
    history,
    isValid,
    name,
    template,
  ]);

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
        <Spacer20 />
        <FormControl label="Pick a Template">
          <TemplatePicker onSelect={setTemplate} />
        </FormControl>
        {createCampaignAPI.error && (
          <Notification
            kind={NotificationKind.negative}
            overrides={{
              Body: { style: { width: "auto" } },
            }}
          >
            Failed to create your campaign. Please try again.
          </Notification>
        )}
        <Button
          disabled={!isValid() || createCampaignAPI.loading}
          isLoading={createCampaignAPI.loading}
          onClick={onSubmit}
        >
          Review
        </Button>
        <Spacer40 />
      </div>
    </PageLayout>
  );
};

export default CreateCampaignPage;
