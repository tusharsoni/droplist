// @flow
import * as React from "react";
import type { Template } from "../../lib/types/content";
import { useStyletron } from "baseui";
import useFetch from "use-http";
import { KIND as NotificationKind, Notification } from "baseui/notification";

type Props = {
  template: Template,
};

const defaultStyleBlock = `
<style type="text/css">
  body { font-family: sans-serif; }
</style>
`;

const TemplatePreview = (props: Props) => {
  const [css] = useStyletron();
  const { loading, error, data: preview } = useFetch<{ html: string }>(
    `/content/templates/${props.template.UUID}/preview`,
    {},
    []
  );

  if (loading) {
    return null;
  }

  if (error) {
    return (
      <Notification kind={NotificationKind.negative}>
        Failed to preview this email
      </Notification>
    );
  }

  return (
    <iframe
      title={`${props.template.Name} Preview`}
      className={css({
        width: "100%",
        height: "500px",
        border: "1px solid #ddd",
      })}
      srcDoc={preview.html + defaultStyleBlock}
    />
  );
};

export default TemplatePreview;
