// @flow

import * as React from "react";
import { useStyletron } from "baseui";

export const Spacer = (props: { size: number }) => {
  const [css] = useStyletron();
  const size = `${props.size}px`;

  return <div className={css({ height: size, width: size })} />;
};

export const Spacer8 = () => <Spacer size={8} />;

export const Spacer20 = () => <Spacer size={20} />;

export const Spacer40 = () => <Spacer size={40} />;
