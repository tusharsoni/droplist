// @flow

import { Provider } from "use-http";
import * as React from "react";
import { useHistory } from "react-router-dom";
import { getSession } from "../auth";

type Props = {
  children: React.Node,
};

const whitelistedPaths = ["/auth/email-otp/signup", "/auth/email-otp/login"];

const HTTPProvider = (props: Props) => {
  const history = useHistory();

  const options = {
    cachePolicy: "no-cache",
    interceptors: {
      request: async ({ options, path }) => {
        const session = getSession();
        const shouldSendAuthHeaders = whitelistedPaths.indexOf(path) < 0;

        if (shouldSendAuthHeaders && !session) {
          history.push("/login");
          return options;
        }

        if (shouldSendAuthHeaders && session) {
          const username = session.userUUID;
          const password = session.sessionToken;
          const token = new Buffer(username + ":" + password).toString(
            "base64"
          );

          options.headers.Authorization = `Basic ${token}`;
        }

        return options;
      },
    },
  };

  return (
    <Provider url="/api" options={options}>
      {props.children}
    </Provider>
  );
};

export default HTTPProvider;
