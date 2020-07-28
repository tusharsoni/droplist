// @flow
import * as React from "react";
import { Button, KIND } from "baseui/button";
import type { Campaign } from "../../lib/types/campaign";
import useFetch from "use-http";
import {
  Modal,
  ModalBody,
  ModalButton,
  ModalFooter,
  ModalHeader,
} from "baseui/modal";
import { useHistory } from "react-router-dom";
import { KIND as NotificationKind, Notification } from "baseui/notification";

type Props = {
  campaign: Campaign,
};

const DeleteCampaignButton = (props: Props) => {
  const [showConfirmModal, setShowConfirmModal] = React.useState(false);
  const deleteCampaignAPI = useFetch(`/campaigns/${props.campaign.UUID}`);
  const history = useHistory();

  return (
    <>
      <Button kind={KIND.secondary} onClick={() => setShowConfirmModal(true)}>
        Delete
      </Button>
      <Modal
        onClose={() => setShowConfirmModal(false)}
        isOpen={showConfirmModal}
        closeable={!deleteCampaignAPI.loading}
        unstable_ModalBackdropScroll={true}
      >
        <ModalHeader>Delete Campaign</ModalHeader>
        <ModalBody>
          Are you sure? You cannot undo this action.
          {deleteCampaignAPI.error && (
            <Notification kind={NotificationKind.negative}>
              Failed to delete campaign
            </Notification>
          )}
        </ModalBody>
        <ModalFooter>
          <ModalButton
            isLoading={deleteCampaignAPI.loading}
            disabled={deleteCampaignAPI.loading}
            onClick={async () => {
              await deleteCampaignAPI.delete();

              if (deleteCampaignAPI.response.ok) {
                setShowConfirmModal(false);
                history.push("/campaigns");
              }
            }}
          >
            Confirm & Delete
          </ModalButton>
        </ModalFooter>
      </Modal>
    </>
  );
};

export default DeleteCampaignButton;
