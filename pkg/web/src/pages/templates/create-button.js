// @flow
import * as React from "react";
import { Button, KIND, SIZE } from "baseui/button";
import {
  Modal,
  ModalBody,
  ModalButton,
  ModalFooter,
  ModalHeader,
  ROLE,
} from "baseui/modal";
import { KIND as NotificationKind, Notification } from "baseui/notification";
import useFetch from "use-http";
import type { Template } from "../../lib/types/content";
import { Input } from "baseui/input";
import { useHistory } from "react-router-dom";

type Props = {};

const CreateTemplateButton = (props: Props) => {
  const history = useHistory();
  const [isModalOpen, setIsModalOpen] = React.useState(false);
  const [name, setName] = React.useState("");
  const createTemplatesAPI = useFetch<Template>(`/content/templates`);

  const onCreate = React.useCallback(async () => {
    const template = await createTemplatesAPI.post({
      name,
    });
    if (createTemplatesAPI.response.ok) {
      setIsModalOpen(false);
      history.push(`/templates/${template.UUID}/edit`);
    }
  }, [createTemplatesAPI, history, name]);

  return (
    <>
      <Button
        kind={KIND.secondary}
        size={SIZE.compact}
        onClick={() => setIsModalOpen(true)}
      >
        Create Template
      </Button>

      <Modal
        onClose={() => setIsModalOpen(false)}
        isOpen={isModalOpen}
        role={ROLE.dialog}
        closeable={!createTemplatesAPI.loading}
        autoFocus
        unstable_ModalBackdropScroll={true}
      >
        <ModalHeader>Create Template</ModalHeader>
        <ModalBody>
          <Input
            placeholder={"Template Name"}
            value={name}
            onChange={(e) => {
              setName(e.target.value);
            }}
          />
        </ModalBody>
        <ModalFooter>
          {createTemplatesAPI.error && (
            <Notification
              kind={NotificationKind.negative}
              overrides={{
                Body: { style: { textAlign: "left", width: "auto" } },
              }}
            >
              Failed to create the template. Please try again.
            </Notification>
          )}
          <ModalButton
            kind={KIND.tertiary}
            disabled={createTemplatesAPI.loading}
            onClick={() => setIsModalOpen(false)}
          >
            Cancel
          </ModalButton>
          <ModalButton
            isLoading={createTemplatesAPI.loading}
            onClick={onCreate}
          >
            Create Template
          </ModalButton>
        </ModalFooter>
      </Modal>
    </>
  );
};

export default CreateTemplateButton;
