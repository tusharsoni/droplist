// @flow
import * as React from "react";
import { Label1 } from "baseui/typography";
import { Spacer40 } from "../../style-guide/spacer";
import { Button } from "baseui/button";
import { Link } from "react-router-dom";
import { useStyletron } from "baseui";
import type { CreateContactResult } from "../../lib/types/audience";

type Props = {
  results: CreateContactResult[],
};

const Success = (props: Props) => {
  const [css] = useStyletron();
  const success = props.results.filter((r) => r.success).length;
  const failed = props.results.length - success;

  if (success === 0) {
    return (
      <>
        <Label1>Failed to upload your contacts. Please try again.</Label1>
        <Spacer40 />
        <Link
          className={css({ textDecoration: "none" })}
          to={"/audience/contacts/import"}
        >
          <Button>Import Contacts</Button>
        </Link>
      </>
    );
  }

  return (
    <>
      <Label1>
        Your contacts have been uploaded successfully!
        {failed === 0
          ? null
          : failed === 1
          ? " 1 contact failed."
          : ` ${failed} contacts failed.`}
      </Label1>
      <Spacer40 />
      <Link className={css({ textDecoration: "none" })} to={"/audience"}>
        <Button>Audience</Button>
      </Link>
    </>
  );
};

export default Success;
