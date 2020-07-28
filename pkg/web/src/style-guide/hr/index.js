// @flow
import * as React from "react";
import { useStyletron } from "baseui";

const HR = () => {
  const [css] = useStyletron();

  return (
    <hr
      className={css({
        border: "none",
        borderBottom: "1px solid #ddd",
        height: "1px",
      })}
    />
  );
};

export default HR;
