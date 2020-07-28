// @flow

import React from "react";
import PageLayout from "../../style-guide/page-layout";
import { Display3, Label1, Label2, ParagraphSmall } from "baseui/typography";
import { useStyletron } from "baseui";
import { Button } from "baseui/button";
import { useHistory, useParams } from "react-router-dom";
import { Spacer, Spacer20, Spacer40, Spacer8 } from "../../style-guide/spacer";
import useFetch from "use-http";
import { StyledSpinnerNext as Spinner } from "baseui/spinner";
import type { Campaign } from "../../lib/types/campaign";
import type { Template } from "../../lib/types/content";
import HR from "../../style-guide/hr";
import type { AudienceSummary } from "../../lib/types/audience";
import { KIND as NotificationKind, Notification } from "baseui/notification";
import TemplatePreview from "../../components/template-preview";
import DeleteCampaignButton from "./delete-button";

const ReviewCampaignPage = () => {
  const history = useHistory();
  const [css] = useStyletron();
  const { uuid: campaignUUID } = useParams();

  const publishCampaignAPI = useFetch(`/campaigns/${campaignUUID}/publish`);
  const {
    data: audienceSummary,
    ...getAudienceSummaryAPI
  } = useFetch<AudienceSummary>("/audience/summary", {}, []);
  const { data: campaign, ...getCampaignAPI } = useFetch<Campaign>(
    `/campaigns/${campaignUUID}`
  );
  const { data: template, ...getTemplateAPI } = useFetch<Template>(
    "/content/templates"
  );

  React.useEffect(() => {
    async function loadData() {
      const campaign = await getCampaignAPI.get();

      if (!getCampaignAPI.response.ok) {
        return;
      }

      await getTemplateAPI.get(campaign.TemplateUUID);
    }

    loadData();
  }, []); /* eslint-disable-line react-hooks/exhaustive-deps */

  const onSend = React.useCallback(async () => {
    await publishCampaignAPI.post();

    if (publishCampaignAPI.response.ok) {
      history.push("/campaigns");
    }
  }, [history, publishCampaignAPI]);

  if (
    getAudienceSummaryAPI.loading ||
    getCampaignAPI.loading ||
    getTemplateAPI.loading
  ) {
    return (
      <PageLayout>
        <Spinner />
      </PageLayout>
    );
  }

  if (
    getAudienceSummaryAPI.error ||
    getCampaignAPI.error ||
    getTemplateAPI.error
  ) {
    return <PageLayout>Failed to load this page. Please try again.</PageLayout>;
  }

  if (!audienceSummary || !campaign || !template) {
    return null;
  }

  return (
    <PageLayout>
      <Display3>Review & Send</Display3>
      <Spacer20 />
      <Label1>Review the details and hit the Send button</Label1>
      <Spacer40 />
      <div className={css({ maxWidth: "500px" })}>
        <Label2>Audience</Label2>
        <ParagraphSmall>
          {audienceSummary.SubscribedContacts === 1
            ? `This email will deliver to 1 contact`
            : `This email will deliver to ${audienceSummary.SubscribedContacts.toLocaleString()} contacts`}
        </ParagraphSmall>

        <Spacer8 />
        <HR />
        <Spacer8 />

        <Label2>From</Label2>
        <ParagraphSmall>
          {campaign.FromName} &lt;{campaign.FromEmail}&gt;
        </ParagraphSmall>

        <Spacer8 />
        <HR />
        <Spacer8 />

        <Label2>Subject</Label2>
        <ParagraphSmall>{template.Subject}</ParagraphSmall>

        <Spacer8 />
        <HR />
        <Spacer8 />

        <Label2>Template</Label2>
        <Spacer size={14} />
        <TemplatePreview template={template} />

        <Spacer40 />
        {publishCampaignAPI.error && (
          <Notification
            kind={NotificationKind.negative}
            overrides={{
              Body: { style: { width: "auto" } },
            }}
          >
            Failed to publish your campaign. Please try again.
          </Notification>
        )}
        <div
          className={css({ display: "flex", justifyContent: "space-between" })}
        >
          <DeleteCampaignButton campaign={campaign} />

          <Button
            disabled={publishCampaignAPI.loading}
            isLoading={publishCampaignAPI.loading}
            onClick={onSend}
          >
            Send
          </Button>
        </div>
        <Spacer40 />
      </div>
    </PageLayout>
  );
};

export default ReviewCampaignPage;
