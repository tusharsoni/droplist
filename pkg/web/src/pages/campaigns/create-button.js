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

const CreateCampaignButton = (props: Props) => {
  const history = useHistory();
  const [isModalOpen, setIsModalOpen] = React.useState(false);
  const [name, setName] = React.useState("");
  const campaignsAPI = useFetch<Template>(`/campaigns`);

  const onCreate = React.useCallback(async () => {
    const campaign = await campaignsAPI.post({
      name,
    });
    if (campaignsAPI.response.ok) {
      setIsModalOpen(false);
      history.push(`/campaigns/${campaign.UUID}/edit`);
    }
  }, [campaignsAPI, history, name]);

  return (
    <>
      <Button
        kind={KIND.secondary}
        size={SIZE.compact}
        onClick={() => setIsModalOpen(true)}
      >
        Create Campaign
      </Button>

      <Modal
        onClose={() => setIsModalOpen(false)}
        isOpen={isModalOpen}
        role={ROLE.dialog}
        closeable={!campaignsAPI.loading}
        autoFocus
        unstable_ModalBackdropScroll={true}
      >
        <ModalHeader>Create Campaign</ModalHeader>
        <ModalBody>
          <Input
            placeholder={"Campaign Name"}
            value={name}
            onChange={(e) => {
              setName(e.target.value);
            }}
          />
        </ModalBody>
        <ModalFooter>
          {campaignsAPI.error && (
            <Notification
              kind={NotificationKind.negative}
              overrides={{
                Body: { style: { textAlign: "left", width: "auto" } },
              }}
            >
              Failed to create the campaign. Please try again.
            </Notification>
          )}
          <ModalButton
            kind={KIND.tertiary}
            disabled={campaignsAPI.loading}
            onClick={() => setIsModalOpen(false)}
          >
            Cancel
          </ModalButton>
          <ModalButton isLoading={campaignsAPI.loading} onClick={onCreate}>
            Create Campaign
          </ModalButton>
        </ModalFooter>
      </Modal>
    </>
  );
};

export default CreateCampaignButton;
