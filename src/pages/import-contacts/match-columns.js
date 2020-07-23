// @flow
import * as React from "react";
import CSVParser from "papaparse";
import { Card, StyledAction, StyledBody } from "baseui/card";
import { StyledSpinnerNext as Spinner } from "baseui/spinner";
import { Select } from "baseui/select";
import { Label1 } from "baseui/typography";
import { Spacer20, Spacer40 } from "../../style-guide/spacer";
import { useStyletron } from "baseui";
import { Notification, KIND as NotificationKind } from "baseui/notification";
import { Button } from "baseui/button";
import { ListItem, ListItemLabel } from "baseui/list";
import useFetch from "use-http";
import type { CreateContactResult } from "../../lib/types/audience";

type Props = {
  csvFile: File,
  onBack: () => void,
  onUpload: (results: CreateContactResult[]) => void,
};

const MatchColumns = (props: Props) => {
  const uploadContactsAPI = useFetch<CreateContactResult[]>(
    "/audience/contacts"
  );
  const [css] = useStyletron();
  const [loading, setLoading] = React.useState(true);
  const [columns, setColumns] = React.useState<string[]>([]);
  const [readErr, setReadErr] = React.useState(false);
  const [csvDataPreview, setCSVDataPreview] = React.useState<
    { [string]: string }[]
  >([]);
  const [emailColumn, setEmailColumn] = React.useState<?string>(null);
  const [firstNameColumn, setFirstNameColumn] = React.useState<?string>(null);
  const [lastNameColumn, setLastNameColumn] = React.useState<?string>(null);
  const [uploading, setUploading] = React.useState(false);
  const [uploadErr, setUploadErr] = React.useState(false);

  React.useEffect(() => {
    CSVParser.parse(props.csvFile, {
      header: true,
      preview: 5,
      skipEmptyLines: true,
      complete: ({ data, errors, meta }) => {
        if (errors && errors.length) {
          setReadErr(true);
          setLoading(false);
          return;
        }

        setColumns(meta.fields);
        setCSVDataPreview(data);
        setLoading(false);
      },
    });
  }, [props.csvFile]);

  const uploadContacts = React.useCallback(() => {
    setUploadErr(false);
    setUploading(true);

    CSVParser.parse(props.csvFile, {
      header: true,
      skipEmptyLines: true,
      complete: async ({ data, errors }) => {
        if (errors && errors.length) {
          setReadErr(true);
          setUploading(false);
          return;
        }

        const contacts = data.map((row) => ({
          email: row[emailColumn],
          params: JSON.stringify({
            FirstName: firstNameColumn ? row[firstNameColumn] : null,
            LastName: lastNameColumn ? row[lastNameColumn] : null,
          }),
        }));

        await uploadContactsAPI.post({
          contacts,
        });

        if (
          !uploadContactsAPI.response.ok ||
          !uploadContactsAPI.response.data
        ) {
          setUploadErr(true);
          setUploading(false);
          return;
        }

        props.onUpload(uploadContactsAPI.response.data);
      },
    });
  }, [emailColumn, firstNameColumn, lastNameColumn, props, uploadContactsAPI]);

  if (loading) {
    return <Spinner />;
  }

  if (readErr) {
    return (
      <>
        <Notification
          kind={NotificationKind.negative}
          overrides={{
            Body: { style: { width: "auto" } },
          }}
        >
          Failed to read the file. Make sure it is a valid CSV file and try
          again.
        </Notification>
        <Spacer20 />
        <Button onClick={props.onBack}>Try Again</Button>
      </>
    );
  }

  return (
    <>
      <Label1>Match columns from the uploaded file to your audience</Label1>
      <Spacer40 />
      <div
        className={css({
          display: "flex",
          overflowX: "auto",
          paddingBottom: "20px",
        })}
      >
        <MatchColumnCard
          title={"Email Address"}
          columns={columns}
          dataPreview={csvDataPreview}
          onMatch={setEmailColumn}
        />
        <MatchColumnCard
          title={"First Name"}
          columns={columns}
          dataPreview={csvDataPreview}
          onMatch={setFirstNameColumn}
        />
        <MatchColumnCard
          title={"Last Name"}
          columns={columns}
          dataPreview={csvDataPreview}
          onMatch={setLastNameColumn}
        />
      </div>
      <Spacer40 />

      {uploadErr && (
        <Notification
          kind={NotificationKind.negative}
          overrides={{
            Body: { style: { width: "auto" } },
          }}
        >
          Failed to upload contacts. Please try again.
        </Notification>
      )}

      <Button
        disabled={!emailColumn}
        isLoading={uploading}
        onClick={uploadContacts}
      >
        Upload Contacts
      </Button>
    </>
  );
};

type MatchColumnCardProps = {
  title: string,
  columns: string[],
  dataPreview: { [string]: string }[],
  onMatch: (column: string) => void,
};

const MatchColumnCard = (props: MatchColumnCardProps) => {
  const options = props.columns.map((col) => ({ id: col, label: col }));
  const [column, setColumn] = React.useState([]);

  return (
    <Card
      title={props.title}
      overrides={{
        Root: {
          style: { flexBasis: "328px", flexShrink: 0, marginRight: "20px" },
        },
      }}
    >
      <StyledAction>
        <Select
          options={options}
          value={column}
          placeholder="Select column"
          onChange={(params) => {
            // $FlowFixMe: BaseUI Flow errors
            setColumn(params.value);

            // $FlowFixMe: BaseUI Flow errors
            params.value.length && props.onMatch(params.value[0].id);
          }}
        />
      </StyledAction>
      {column.length ? (
        <StyledBody>
          {props.dataPreview.map((row, idx) => (
            <ListItem key={`row-${idx}`}>
              <ListItemLabel>{row[column[0].id]}</ListItemLabel>
            </ListItem>
          ))}
        </StyledBody>
      ) : null}
    </Card>
  );
};

export default MatchColumns;
