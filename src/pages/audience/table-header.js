// @flow
import * as React from "react";
import { Display3, Label1 } from "baseui/typography";
import { useStyletron } from "baseui";
import { Spacer20 } from "../../style-guide/spacer";

type Props = {
  totalContacts: number,
  subscribedContacts: number,
};

const AudienceTableHeader = (props: Props) => {
  const [css] = useStyletron();

  return (
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
            ? "1 of these is a subscriber."
            : `${props.subscribedContacts.toLocaleString()} of these are subscribers.`}
        </Label1>
      </div>
      <div>{/* todo: Actions such as import, export go here */}</div>
    </div>
  );
};

export default AudienceTableHeader;
