// @flow

import { Provider } from "use-http";
import * as React from "react";

type Props = {
  children: React.Node,
};

const HTTPProvider = (props: Props) => {
  const options = {
    cachePolicy: "no-cache",
    interceptors: {
      request: async ({ options }) => {
        const username = "490a8945-be70-4324-969a-8b9475eb970e";
        const password =
          "fKUkdvYo0PGXtMQZsvC5Wj9N02zAPN7FJdyf3g89KPB13X6GZYZmWxDj2qnl3TjA";
        const token = new Buffer(username + ":" + password).toString("base64");

        options.headers.Authorization = `Basic ${token}`;

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
