// @flow
import * as React from "react";
import useFetch from "use-http";
import type { Template } from "../../lib/types/content";
import { StyledSpinnerNext as Spinner } from "baseui/spinner";
import { useStyletron } from "baseui";
import { Label4, LabelSmall } from "baseui/typography";
import { Spacer20, Spacer8 } from "../../style-guide/spacer";
import { Checkbox } from "baseui/checkbox";
import { Button, SIZE, KIND } from "baseui/button";
import { ChevronLeft, ChevronRight } from "baseui/icon";
import { KIND as NotificationKind, Notification } from "baseui/notification";
import TemplatePreview from "../template-preview";

type Props = {
  initialTemplateUUID?: ?string,
  onSelect: (template: Template) => void,
};

const TemplatePicker = (props: Props) => {
  const [selectedTemplate, setSelectedTemplate] = React.useState<?Template>(
    null
  );
  const [page, setPage] = React.useState(0);
  const [css] = useStyletron();
  const { loading, error, data: templates, ...templatesAPI } = useFetch<
    Template[]
  >("/content/templates");

  React.useEffect(() => {
    async function loadData() {
      const templates = await templatesAPI.get();

      console.log("=========>", props.initialTemplateUUID);
      if (!templatesAPI.response.ok || !props.initialTemplateUUID) {
        return;
      }

      const initialTemplate = templates.find(
        (t) => t.UUID === props.initialTemplateUUID
      );

      if (!initialTemplate) {
        return;
      }

      const initialPage = Math.floor(templates.indexOf(initialTemplate) / 2);

      setPage(initialPage);
      setSelectedTemplate(initialTemplate);
    }

    loadData();
  }, []); /* eslint-disable-line react-hooks/exhaustive-deps */

  if (loading) {
    return (
      <div
        className={css({
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          minHeight: "636px",
        })}
      >
        <Spinner />
      </div>
    );
  }

  if (error) {
    return (
      <Notification kind={NotificationKind.negative}>
        Failed to load templates
      </Notification>
    );
  }

  if (!templates) {
    return null;
  }

  return (
    <div className={css({ display: "flex", minHeight: "636px" })}>
      <div className={css({ display: "flex", alignItems: "center" })}>
        <div>
          <Button
            size={SIZE.mini}
            kind={KIND.secondary}
            disabled={page === 0}
            onClick={() => {
              setPage(page - 1);
            }}
          >
            <ChevronLeft size={32} />
          </Button>
        </div>
      </div>
      {templates.slice(page * 2, page * 2 + 2).map((template: Template) => (
        <div
          key={template.UUID}
          className={css({
            minWidth: "400px",
            width: "400px",
            padding: "20px",
          })}
        >
          <div
            className={css({
              display: "flex",
              justifyContent: "space-between",
            })}
          >
            <div
              className={css({ textOverflow: "ellipsis", overflow: "hidden" })}
            >
              <LabelSmall>{template.Name}</LabelSmall>
              <Spacer8 />
              <Label4
                overrides={{
                  Block: {
                    style: {
                      whiteSpace: "nowrap",
                    },
                  },
                }}
              >
                {template.Subject}
              </Label4>
            </div>
            <div>
              <Checkbox
                checked={Boolean(
                  selectedTemplate && selectedTemplate.UUID === template.UUID
                )}
                onChange={() => {
                  setSelectedTemplate(template);
                  props.onSelect(template);
                }}
              />
            </div>
          </div>
          <Spacer8 />
          <TemplatePreview template={template} />
          <Spacer20 />
        </div>
      ))}
      <div className={css({ display: "flex", alignItems: "center" })}>
        <div>
          <Button
            size={SIZE.mini}
            kind={KIND.secondary}
            disabled={(page + 1) * 2 >= templates.length}
            onClick={() => {
              setPage(page + 1);
            }}
          >
            <ChevronRight size={32} />
          </Button>
        </div>
      </div>
    </div>
  );
};

export default TemplatePicker;
