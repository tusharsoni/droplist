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
  onCreate: () => void,
};

const CreateContactButton = (props: Props) => {
  const [showModal, setShowModal] = React.useState(false);
  const contactsAPI = useFetch(`/audience/contacts`);
  const [email, setEmail] = React.useState("");
  const [firstName, setFirstName] = React.useState("");
  const [lastName, setLastName] = React.useState("");

  return (
    <>
      <Button
        kind={KIND.secondary}
        size={SIZE.compact}
        onClick={() => setShowModal(true)}
      >
        Add Contact
      </Button>
      <Modal
        onClose={() => setShowModal(false)}
        isOpen={showModal}
        closeable={!contactsAPI.loading}
        unstable_ModalBackdropScroll={true}
      >
        <ModalHeader>Add Contact</ModalHeader>
        <ModalBody>
          <FormControl label="Email">
            <Input
              value={email}
              onChange={(e) => {
                setEmail(e.target.value);
              }}
            />
          </FormControl>
          <FormControl label="First Name">
            <Input
              value={firstName}
              onChange={(e) => {
                setFirstName(e.target.value);
              }}
            />
          </FormControl>
          <FormControl label="Last Name">
            <Input
              value={lastName}
              onChange={(e) => {
                setLastName(e.target.value);
              }}
            />
          </FormControl>

          {contactsAPI.response.ok && (
            <Notification kind={NotificationKind.positive}>
              Successfully added contact
            </Notification>
          )}

          {contactsAPI.error && (
            <Notification kind={NotificationKind.negative}>
              Failed to add contact
            </Notification>
          )}
        </ModalBody>
        <ModalFooter>
          <ModalButton
            isLoading={contactsAPI.loading}
            disabled={contactsAPI.loading}
            onClick={async () => {
              await contactsAPI.post({
                contacts: [
                  {
                    email: email.toLowerCase().trim(),
                    params: JSON.stringify({
                      FirstName: firstName,
                      LastName: lastName,
                    }),
                  },
                ],
              });
              if (contactsAPI.response.ok) {
                setShowModal(false);
                props.onCreate();
              }
            }}
          >
            Add Contact
          </ModalButton>
        </ModalFooter>
      </Modal>
    </>
  );
};

export default CreateContactButton;
