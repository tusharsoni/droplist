// @flow

import * as React from "react";
import { useStyletron } from "baseui";

export const Spacer8 = () => {
  const [css] = useStyletron();

  return <div className={css({ height: "8px", width: "8px" })} />;
};

export const Spacer20 = () => {
  const [css] = useStyletron();

  return <div className={css({ height: "20px", width: "20px" })} />;
};

export const Spacer40 = () => {
  const [css] = useStyletron();

  return <div className={css({ height: "40px", width: "40px" })} />;
};
