// @flow

import React from "react";
import PageLayout from "../../style-guide/page-layout";
import { Display3 } from "baseui/typography";
import { useStyletron } from "baseui";
import { Button, KIND, SIZE } from "baseui/button";
import { Link } from "react-router-dom";
import useFetch from "use-http";
import type { Campaign, CampaignStats } from "../../lib/types/campaign";
import { StyledSpinnerNext as Spinner } from "baseui/spinner";
import { Spacer40 } from "../../style-guide/spacer";
import { Table } from "baseui/table-semantic";
import { Tag } from "baseui/tag";
import { StyledLink } from "baseui/link";

const CampaignsPage = () => {
  const [css] = useStyletron();
  const { data: campaigns, ...getCampaignsAPI } = useFetch<Campaign[]>(
    "/campaigns"
  );
  const { data: stats, ...getStatsAPI } = useFetch<{ [string]: CampaignStats }>(
    "/campaigns/stats"
  );

  React.useEffect(() => {
    async function loadData() {
      const campaigns = await getCampaignsAPI.get();
      if (!getCampaignsAPI.response.ok) {
        return;
      }

      getStatsAPI.post({
        campaign_uuids: campaigns.map((c) => c.UUID),
      });
    }

    loadData();
  }, []); /* eslint-disable-line react-hooks/exhaustive-deps */

  if (getCampaignsAPI.loading || (!campaigns && !getCampaignsAPI.error)) {
    return (
      <PageLayout>
        <Spinner />
      </PageLayout>
    );
  }

  if (getCampaignsAPI.error) {
    return <PageLayout>Failed to load this page. Please try again.</PageLayout>;
  }

  return (
    <PageLayout>
      <div
        className={css({
          display: "flex",
          justifyContent: "space-between",
          alignItems: "flex-end",
        })}
      >
        <div>
          <Display3>Campaigns</Display3>
        </div>
        <Link
          className={css({ textDecoration: "none" })}
          to={"/campaigns/create"}
        >
          <Button kind={KIND.secondary} size={SIZE.compact}>
            Create Campaign
          </Button>
        </Link>
      </div>

      <Spacer40 />

      <Table
        columns={[
          "",
          "",
          "Sent",
          "Failed",
          "Opens",
          "Clicks",
          "Open Rate",
          "Click Rate",
        ]}
        data={campaigns.map((campaign: Campaign) => [
          campaign.State === "DRAFT" ? (
            <StyledLink $as={Link} to={`/campaigns/${campaign.UUID}/review`}>
              {campaign.Name}
            </StyledLink>
          ) : (
            campaign.Name
          ),
          campaign.State === "DRAFT" ? (
            <Tag closeable={false}>Draft</Tag>
          ) : (
            <Tag closeable={false} kind="positive">
              Sent
            </Tag>
          ),
          stats && stats[campaign.UUID]
            ? formatStat(stats[campaign.UUID].Sent)
            : "--",
          stats && stats[campaign.UUID]
            ? formatStat(stats[campaign.UUID].Failed)
            : "--",
          stats && stats[campaign.UUID]
            ? formatStat(stats[campaign.UUID].Opens)
            : "--",
          stats && stats[campaign.UUID]
            ? formatStat(stats[campaign.UUID].Clicks)
            : "--",
          stats && stats[campaign.UUID]
            ? formatStat(stats[campaign.UUID].OpenRate)
            : "--",
          stats && stats[campaign.UUID]
            ? formatStat(stats[campaign.UUID].ClickRate)
            : "--",
        ])}
      />
    </PageLayout>
  );
};

function formatStat(number: ?number, percent: ?boolean) {
  if (number == null) {
    return "--";
  }
  const stat = number || 0;

  if (percent) {
    return `${(stat * 100.0).toLocaleString()}%`;
  }

  return stat.toLocaleString();
}

export default CampaignsPage;
