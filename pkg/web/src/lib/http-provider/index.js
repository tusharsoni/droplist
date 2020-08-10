// @flow

import { Provider } from "use-http";
import * as React from "react";
import { useHistory, useLocation } from "react-router-dom";
import { clearSession, getSession } from "../auth";

type Props = {
  children: React.Node,
};

const whitelistedPaths = ["/auth/email-otp/signup", "/auth/email-otp/login"];

const HTTPProvider = (props: Props) => {
  const history = useHistory();
  const location = useLocation();
  const loginURL = `/login?to=${encodeURIComponent(
    location.pathname + location.search
  )}`;

  const options = {
    cachePolicy: "no-cache",
    interceptors: {
      request: async ({ options, path }) => {
        const session = getSession();
        const shouldSendAuthHeaders = whitelistedPaths.indexOf(path) < 0;

        if (shouldSendAuthHeaders && !session) {
          history.push(loginURL);
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
      response: async ({ response }) => {
        if (response.status === 401) {
          clearSession();
          history.push(loginURL);
        }

        return response;
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
