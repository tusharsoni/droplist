// @flow
import * as React from "react";
import { Table } from "baseui/table-semantic";
import { Tag } from "baseui/tag";
import { DateTime } from "luxon";
import type { AudienceSummary, Contact } from "../../lib/types/audience";
import { Checkbox } from "baseui/checkbox";
import AudienceTableHeader from "./table-header";
import { Spacer40 } from "../../style-guide/spacer";

type Props = {
  summary: AudienceSummary,
  contacts: Contact[],
  onRefresh: () => void,
};

const AudienceTable = (props: Props) => {
  const [checkAll, setCheckAll] = React.useState(false);
  const [selectedContacts, setSelectedContacts] = React.useState<Contact[]>([]);

  return (
    <>
      <AudienceTableHeader
        totalContacts={props.summary.TotalContacts}
        subscribedContacts={props.summary.SubscribedContacts}
        selectedAll={checkAll}
        selectedContacts={selectedContacts}
        onDelete={props.onRefresh}
      />
      <Spacer40 />
      <Table
        columns={[
          <Checkbox
            checked={checkAll}
            onChange={(e) => {
              setCheckAll(e.target.checked);
              setSelectedContacts([]);
            }}
          />,
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
            <Checkbox
              checked={checkAll || selectedContacts.indexOf(contact) >= 0}
              onChange={(e) => {
                if (checkAll && !e.target.checked) {
                  setCheckAll(false);
                  setSelectedContacts(
                    props.contacts.filter((c) => c.UUID !== contact.UUID)
                  );
                  return;
                }

                e.target.checked
                  ? setSelectedContacts([...selectedContacts, contact])
                  : setSelectedContacts(
                      selectedContacts.filter((c) => c.UUID !== contact.UUID)
                    );
              }}
            />,
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
        overrides={{
          TableBodyCell: {
            style: { verticalAlign: "middle" },
          },
        }}
      />
    </>
  );
};

export default AudienceTable;
