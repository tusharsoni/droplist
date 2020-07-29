// @flow
import * as React from "react";
import { Button, KIND, SIZE } from "baseui/button";
import {
  Modal,
  ModalBody,
  ModalButton,
  ModalFooter,
  ModalHeader,
} from "baseui/modal";
import { KIND as NotificationKind, Notification } from "baseui/notification";
import useFetch from "use-http";
import { Input } from "baseui/input";
import { FormControl } from "baseui/form-control";

type Props = {
  campaignUUID: string,
};

const SendTestButton = (props: Props) => {
  const [showModal, setShowModal] = React.useState(false);
  const testCampaignAPI = useFetch(`/campaigns/${props.campaignUUID}/test`);
  const [email, setEmail] = React.useState("");

  return (
    <>
      <Button
        kind={KIND.secondary}
        size={SIZE.compact}
        onClick={() => setShowModal(true)}
      >
        Send Test Email
      </Button>
      <Modal
        onClose={() => setShowModal(false)}
        isOpen={showModal}
        closeable={!testCampaignAPI.loading}
        unstable_ModalBackdropScroll={true}
      >
        <ModalHeader>Send Test Email</ModalHeader>
        <ModalBody>
          <FormControl
            label="Contact Email"
            caption="This contact must be in your audience and be subscribed to receive emails"
          >
            <Input
              placeholder={"jane@doe.com"}
              value={email}
              onChange={(e) => {
                setEmail(e.target.value);
              }}
            />
          </FormControl>

          {testCampaignAPI.response.ok && (
            <Notification kind={NotificationKind.positive}>
              The test email has been queued for sending
            </Notification>
          )}

          {testCampaignAPI.error && (
            <Notification kind={NotificationKind.negative}>
              Failed to send a test email
            </Notification>
          )}
        </ModalBody>
        <ModalFooter>
          <ModalButton
            isLoading={testCampaignAPI.loading}
            disabled={testCampaignAPI.loading}
            onClick={() => {
              testCampaignAPI.post({ emails: [email] });
            }}
          >
            Send
          </ModalButton>
        </ModalFooter>
      </Modal>
    </>
  );
};

export default SendTestButton;
