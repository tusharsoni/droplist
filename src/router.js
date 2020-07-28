// @flow

import * as React from "react";
import { BrowserRouter, Redirect, Route, Switch } from "react-router-dom";
import AudiencePage from "./pages/audience";
import TemplatesPage from "./pages/templates";
import CampaignsPage from "./pages/campaigns";
import ImportContactsPage from "./pages/import-contacts";
import EditTemplatePage from "./pages/edit-template";
import CreateCampaignPage from "./pages/create-campaign";
import ReviewCampaignPage from "./pages/review-campaign";
import EditCampaignPage from "./pages/edit-campaign";

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

      <Route exact path="/campaigns" component={CampaignsPage} />
      <Route exact path="/campaigns/create" component={CreateCampaignPage} />
      <Route exact path="/campaigns/:uuid/edit" component={EditCampaignPage} />
      <Route
        exact
        path="/campaigns/:uuid/review"
        component={ReviewCampaignPage}
      />
    </Switch>
  </BrowserRouter>
);

export default AppRouter;
