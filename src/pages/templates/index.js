// @flow

import React from "react";
import PageLayout from "../../style-guide/page-layout";
import useFetch from "use-http";
import { StyledSpinnerNext as Spinner } from "baseui/spinner";
import type { Template } from "../../lib/types/content";
import { useStyletron } from "baseui";
import { Display3, Label1 } from "baseui/typography";
import { Spacer40, Spacer8 } from "../../style-guide/spacer";
import TemplateActionMenu from "./action-menu";
import CreateTemplateButton from "./create-button";
import { ListItem, ListItemLabel } from "baseui/list";
import ContentSvg from "../../style-guide/illustrations/content";

const TemplatesPage = () => {
  const [css] = useStyletron();
  const { get, loading, error, data: templates } = useFetch<Template[]>(
    "/content/templates",
    {},
    []
  );

  if (loading) {
    return (
      <PageLayout>
        <Spinner />
      </PageLayout>
    );
  }

  if (error) {
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
          <Display3>Templates</Display3>
        </div>
        {templates.length > 0 && <CreateTemplateButton />}
      </div>
      <Spacer40 />

      {templates.length === 0 && (
        <div
          className={css({
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            justifyContent: "center",
          })}
        >
          <ContentSvg height={250} />
          <Spacer40 />
          <Label1>You don't have any templates yet</Label1>
          <Spacer8 />
          <CreateTemplateButton />
        </div>
      )}

      <ul className={css({ padding: 0, margin: 0 })}>
        {templates.map((template: Template) => (
          <ListItem
            key={template.UUID}
            endEnhancer={() => (
              <TemplateActionMenu template={template} onUpdate={get} />
            )}
          >
            <ListItemLabel description={template.Subject}>
              {template.Name}
            </ListItemLabel>
          </ListItem>
        ))}
      </ul>
    </PageLayout>
  );
};

export default TemplatesPage;
