// @flow
import * as React from "react";
import { Display3, Label1 } from "baseui/typography";
import { useStyletron } from "baseui";
import { Spacer20 } from "../../style-guide/spacer";
import type { Contact } from "../../lib/types/audience";
import { Button, KIND, SIZE } from "baseui/button";
import {
  Modal,
  ModalBody,
  ModalButton,
  ModalFooter,
  ModalHeader,
  ROLE,
} from "baseui/modal";
import useFetch from "use-http";

type Props = {
  totalContacts: number,
  subscribedContacts: number,
  selectedAll?: ?boolean,
  selectedContacts: Contact[],
};

const AudienceTableHeader = (props: Props) => {
  const deleteContactsAPI = useFetch("/audience/contacts");
  const [showConfirmModal, setShowConfirmModal] = React.useState(false);
  const [css] = useStyletron();
  const selectedCount = props.selectedAll
    ? props.totalContacts
    : props.selectedContacts.length;

  return (
    <>
      <div
        className={css({
          display: "flex",
          justifyContent: "space-between",
          alignItems: "flex-end",
        })}
      >
        <div>
          <Display3>Audience</Display3>
          <Spacer20 />

          <Label1>
            {props.totalContacts === 1
              ? "Your audience has 1 contact. "
              : `Your audience has ${props.totalContacts.toLocaleString()} contacts. `}
            {props.subscribedContacts === 1
              ? "1 is a subscriber."
              : `${props.subscribedContacts.toLocaleString()} of these are subscribers.`}
          </Label1>
        </div>
        <div>
          {selectedCount ? (
            <Button
              kind={KIND.secondary}
              size={SIZE.compact}
              onClick={() => setShowConfirmModal(true)}
            >
              {selectedCount === 1
                ? "Delete 1 contact"
                : `Delete ${selectedCount} contacts`}
            </Button>
          ) : null}
        </div>
      </div>

      <Modal
        onClose={() => setShowConfirmModal(false)}
        isOpen={showConfirmModal}
        closeable={deleteContactsAPI.loading}
      >
        <ModalHeader>
          {selectedCount === 1
            ? "Delete 1 contact"
            : `Delete ${selectedCount} contacts`}
        </ModalHeader>
        <ModalBody>Are you sure? You cannot undo this action.</ModalBody>
        <ModalFooter>
          <ModalButton
            isLoading={deleteContactsAPI.loading}
            onClick={async () => {
              await deleteContactsAPI.delete({
                delete_all: Boolean(props.selectedAll),
                contact_uuids: props.selectedContacts.map((c) => c.UUID),
              });

              if (!deleteContactsAPI.response.ok) {
                // todo: show error toast
                return;
              }

              setShowConfirmModal(false);
              window.location = "/audience";
            }}
          >
            Confirm & Delete
          </ModalButton>
        </ModalFooter>
      </Modal>
    </>
  );
};

export default AudienceTableHeader;
