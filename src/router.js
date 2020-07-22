// @flow

import * as React from "react";
import { BrowserRouter, Redirect, Route, Switch } from "react-router-dom";
import AudiencePage from "./pages/audience";
import TemplatesPage from "./pages/templates";
import CampaignsPage from "./pages/campaigns";

const AppRouter = () => (
  <BrowserRouter>
    <Switch>
      <Route exact path="/">
        <Redirect to="/campaigns" />
      </Route>

      <Route path="/audience" component={AudiencePage} />
      <Route path="/templates" component={TemplatesPage} />
      <Route path="/campaigns" component={CampaignsPage} />
    </Switch>
  </BrowserRouter>
);

export default AppRouter;
