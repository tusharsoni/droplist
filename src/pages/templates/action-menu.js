// @flow
import * as React from "react";
import type { Template } from "../../lib/types/content";
import { Link } from "react-router-dom";
import { Button, KIND, SIZE } from "baseui/button";
import { PLACEMENT as PopoverPlacement, StatefulPopover } from "baseui/popover";
import { StatefulMenu } from "baseui/menu";
import { ChevronDown } from "baseui/icon";
import { useStyletron } from "baseui";
import {
  Modal,
  ModalBody,
  ModalButton,
  ModalFooter,
  ModalHeader,
  ROLE,
} from "baseui/modal";
import { Input } from "baseui/input";
import useFetch from "use-http";
import { KIND as NotificationKind, Notification } from "baseui/notification";

type Props = {
  template: Template,
  onUpdate: () => void,
};

const TemplateActionMenu = (props: Props) => {
  const [css] = useStyletron();
  const [isRenameModalOpen, setIsRenameModalOpen] = React.useState(false);
  const [name, setName] = React.useState(props.template.Name);
  const updateTemplateAPI = useFetch<Template>(
    `/content/templates/${props.template.UUID}`
  );

  const onRename = React.useCallback(async () => {
    await updateTemplateAPI.post({
      name,
      subject: props.template.Subject,
      preview_text: props.template.PreviewText,
      html_body: props.template.HTMLBody,
    });
    if (updateTemplateAPI.response.ok) {
      setIsRenameModalOpen(false);
      props.onUpdate();
    }
  }, [name, props, updateTemplateAPI]);

  return (
    <>
      <Link
        className={css({ textDecoration: "none" })}
        to={`/templates/${props.template.UUID}/edit`}
      >
        <Button size={SIZE.compact} kind={KIND.secondary}>
          Edit
        </Button>
      </Link>
      <StatefulPopover
        focusLock
        placement={PopoverPlacement.bottomLeft}
        content={({ close }) => (
          <StatefulMenu
            items={[
              { id: "rename", label: "Rename" },
              { id: "delete", label: "Delete" },
            ]}
            onItemSelect={({ item }) => {
              if (item.id === "rename") {
                setIsRenameModalOpen(true);
              }
              close();
            }}
          />
        )}
      >
        <Button size={SIZE.compact} kind={KIND.secondary}>
          <ChevronDown />
        </Button>
      </StatefulPopover>

      <Modal
        onClose={() => setIsRenameModalOpen(false)}
        isOpen={isRenameModalOpen}
        size={SIZE.default}
        role={ROLE.dialog}
        closeable={!updateTemplateAPI.loading}
        autoFocus
      >
        <ModalHeader>Rename Template</ModalHeader>
        <ModalBody>
          <Input
            value={name}
            onChange={(e) => {
              setName(e.target.value);
            }}
          />
        </ModalBody>
        <ModalFooter>
          {updateTemplateAPI.error && (
            <Notification
              kind={NotificationKind.negative}
              overrides={{
                Body: { style: { textAlign: "left", width: "auto" } },
              }}
            >
              Failed to save changes. Please try again.
            </Notification>
          )}
          <ModalButton
            kind={KIND.tertiary}
            disabled={updateTemplateAPI.loading}
            onClick={() => setIsRenameModalOpen(false)}
          >
            Cancel
          </ModalButton>
          <ModalButton isLoading={updateTemplateAPI.loading} onClick={onRename}>
            Save Changes
          </ModalButton>
        </ModalFooter>
      </Modal>
    </>
  );
};

export default TemplateActionMenu;
