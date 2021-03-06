// @flow
import React from "react";
import ReactDOM from "react-dom";
import { Client as Styletron } from "styletron-engine-atomic";
import { Provider as StyletronProvider } from "styletron-react";
import { BaseProvider } from "baseui";
import AppRouter from "./router";
import "./index.css";
import HTTPProvider from "./lib/http-provider";
import theme from "./style-guide/theme";
import { BrowserRouter } from "react-router-dom";

const engine = new Styletron();

ReactDOM.render(
  <React.StrictMode>
    <StyletronProvider value={engine}>
      <BaseProvider
        theme={theme}
        overrides={{ AppContainer: { style: { height: "100%" } } }}
      >
        <BrowserRouter>
          <HTTPProvider>
            <AppRouter />
          </HTTPProvider>
        </BrowserRouter>
      </BaseProvider>
    </StyletronProvider>
  </React.StrictMode>,
  // $FlowFixMe: root element exists
  document.getElementById("root")
);
