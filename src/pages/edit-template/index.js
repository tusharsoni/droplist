// @flow
import * as React from "react";
import PageLayout from "../../style-guide/page-layout";
import { useStyletron } from "baseui";
import useFetch from "use-http";
import type { Template } from "../../lib/types/content";
import { useParams } from "react-router-dom";
import TemplatePreview from "./preview";
import TemplateEditor from "./editor";
import { StyledSpinnerNext as Spinner } from "baseui/spinner";

const EditTemplatePage = () => {
  const { uuid: templateUUID } = useParams();
  const [css] = useStyletron();
  const [template, setTemplate] = React.useState<?Template>(null);

  const fetchTemplateAPI = useFetch<Template>(
    `/content/templates/${templateUUID}`
  );
  const updateTemplateAPI = useFetch<Template>(
    `/content/templates/${templateUUID}`
  );

  React.useEffect(() => {
    async function fetch() {
      const initialTemplate = await fetchTemplateAPI.get();

      if (fetchTemplateAPI.response.ok) {
        setTemplate(initialTemplate);
      }
    }

    fetch();
  }, []); /* eslint-disable-line react-hooks/exhaustive-deps */

  const onHTMLEdit = React.useCallback(
    async (template: Template) => {
      const updatedTemplate = await updateTemplateAPI.post({
        name: template.Name,
        subject: template.Subject,
        preview_text: template.PreviewText,
        html_body: template.HTMLBody,
      });

      if (updateTemplateAPI.response.ok) {
        setTemplate(updatedTemplate);
      }
    },
    [updateTemplateAPI]
  );

  if (fetchTemplateAPI.loading) {
    return (
      <PageLayout>
        <Spinner />
      </PageLayout>
    );
  }

  if (fetchTemplateAPI.error) {
    return <PageLayout>Failed to load this page. Please try again.</PageLayout>;
  }

  return (
    <PageLayout>
      <div
        className={css({
          display: "flex",
          height: "100%",
          paddingBottom: "20px",
        })}
      >
        <div
          className={css({
            flex: 1,
            marginRight: "20px",
          })}
        >
          {template && <TemplatePreview template={template} />}
        </div>

        {template && (
          <TemplateEditor
            template={template}
            saving={updateTemplateAPI.loading}
            saveError={Boolean(updateTemplateAPI.error)}
            onSave={onHTMLEdit}
          />
        )}
      </div>
    </PageLayout>
  );
};

export default EditTemplatePage;
