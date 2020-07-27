// @flow

import React from "react";
import PageLayout from "../../style-guide/page-layout";
import { Display3 } from "baseui/typography";
import { useStyletron } from "baseui";
import { Button, KIND, SIZE } from "baseui/button";
import { Link } from "react-router-dom";

const CampaignsPage = () => {
  const [css] = useStyletron();

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
    </PageLayout>
  );
};

export default CampaignsPage;
