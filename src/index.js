// @flow
import React from "react";
import ReactDOM from "react-dom";
import { Client as Styletron } from "styletron-engine-atomic";
import { Provider as StyletronProvider } from "styletron-react";
import { LightTheme, BaseProvider } from "baseui";
import AppRouter from "./router";
import "./index.css";

const engine = new Styletron();

ReactDOM.render(
  <React.StrictMode>
    <StyletronProvider value={engine}>
      <BaseProvider theme={LightTheme}>
        <AppRouter />
      </BaseProvider>
    </StyletronProvider>
  </React.StrictMode>,
  // $FlowFixMe: root element exists
  document.getElementById("root")
);
