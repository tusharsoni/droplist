// @flow
import * as React from "react";
import { Table } from "baseui/table-semantic";
import { Tag } from "baseui/tag";
import { DateTime } from "luxon";
import type { Contact } from "../../lib/types/audience";

type Props = {
  contacts: Contact[],
};

const AudienceTable = (props: Props) => {
  return (
    <Table
      columns={[
        "Email Address",
        "First Name",
        "Last Name",
        "Status",
        "Date Added",
        "Last Changed",
      ]}
      data={props.contacts.map((contact) => {
        const params = JSON.parse(contact.Params);

        return [
          contact.Email,
          params["FirstName"],
          params["LastName"],
          contact.Status === "SUBSCRIBED" ? (
            <Tag closeable={false} kind="positive">
              Subscribed
            </Tag>
          ) : (
            <Tag closeable={false} kind="negative">
              Unsubscribed
            </Tag>
          ),
          DateTime.fromISO(contact.CreatedAt).toLocaleString(
            DateTime.DATETIME_SHORT
          ),
          DateTime.fromISO(contact.UpdatedAt).toLocaleString(
            DateTime.DATETIME_SHORT
          ),
        ];
      })}
    />
  );
};

export default AudienceTable;
