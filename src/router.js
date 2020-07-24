// @flow

import * as React from "react";
import { BrowserRouter, Redirect, Route, Switch } from "react-router-dom";
import AudiencePage from "./pages/audience";
import TemplatesPage from "./pages/templates";
import CampaignsPage from "./pages/campaigns";
import ImportContactsPage from "./pages/import-contacts";
import EditTemplatePage from "./pages/edit-template";

const AppRouter = () => (
  <BrowserRouter>
    <Switch>
      <Route exact path="/">
        <Redirect to="/campaigns" />
      </Route>

      <Route exact path="/audience" component={AudiencePage} />
      <Route
        exact
        path="/audience/contacts/import"
        component={ImportContactsPage}
      />

      <Route exact path="/templates" component={TemplatesPage} />
      <Route exact path="/templates/:uuid/edit" component={EditTemplatePage} />

      <Route path="/campaigns" component={CampaignsPage} />
    </Switch>
  </BrowserRouter>
);

export default AppRouter;
