// @flow
import * as React from "react";
import { useStyletron } from "baseui";
import { Button, KIND, SIZE as ButtonSIZE } from "baseui/button";
import { Spacer20 } from "../../style-guide/spacer";
import { KIND as NotificationKind, Notification } from "baseui/notification";
import useFetch from "use-http";
import type { Template } from "../../lib/types/content";

const PreviewModes = {
  Mobile: "MOBILE",
  Desktop: "DESKTOP",
};

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
  const [previewMode, setPreviewMode] = React.useState(PreviewModes.Desktop);

  const { error: previewError, data: preview } = useFetch<{ html: string }>(
    `/content/templates/${props.template.UUID}/preview`,
    {},
    [props.template.HTMLBody]
  );

  return (
    <div
      className={css({
        display: "flex",
        flexDirection: "column",
        height: "100%",
      })}
    >
      <div className={css({ textAlign: "center" })}>
        <Button
          onClick={() => setPreviewMode(PreviewModes.Desktop)}
          kind={KIND.tertiary}
          size={ButtonSIZE.mini}
          isSelected={previewMode === PreviewModes.Desktop}
        >
          Desktop
        </Button>
        <Button
          onClick={() => setPreviewMode(PreviewModes.Mobile)}
          kind={KIND.tertiary}
          size={ButtonSIZE.mini}
          isSelected={previewMode === PreviewModes.Mobile}
        >
          Mobile
        </Button>
      </div>
      <Spacer20 />
      {previewError && (
        <Notification
          kind={NotificationKind.negative}
          overrides={{
            Body: { style: { width: "auto" } },
          }}
        >
          Failed to preview the email. Please make sure are using valid snippets
          and try again.
        </Notification>
      )}

      {preview && (
        <div
          className={css({
            flex: 1,
            width: previewMode === PreviewModes.Mobile ? "360px" : "100%",
            marginLeft: "auto",
            marginRight: "auto",
            border: "1px solid #ddd",
          })}
        >
          <iframe
            title="Email Preview"
            className={css({
              width: "100%",
              height: "100%",
              border: "none",
            })}
            srcDoc={preview.html + defaultStyleBlock}
          />
        </div>
      )}
    </div>
  );
};

export default TemplatePreview;
