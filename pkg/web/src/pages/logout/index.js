// @flow
import * as React from "react";
import { useHistory } from "react-router-dom";
import { clearSession } from "../../lib/auth";

const LogoutPage = () => {
  const history = useHistory();

  React.useEffect(() => {
    clearSession();
    history.push("/login");
  }, []); /* eslint-disable-line react-hooks/exhaustive-deps */

  return null;
};

export default LogoutPage;
