// @flow

import * as React from "react";
import { Redirect, Route, Switch } from "react-router-dom";
import AudiencePage from "./pages/audience";
import TemplatesPage from "./pages/templates";
import CampaignsPage from "./pages/campaigns";
import ImportContactsPage from "./pages/import-contacts";
import EditTemplatePage from "./pages/edit-template";
import ReviewCampaignPage from "./pages/review-campaign";
import EditCampaignPage from "./pages/edit-campaign";
import ProfilePage from "./pages/profile";
import LoginPage from "./pages/login";
import LogoutPage from "./pages/logout";

const AppRouter = () => (
  <Switch>
    <Route exact path="/">
      <Redirect to="/campaigns" />
    </Route>

    <Route exact path="/login" component={LoginPage} />
    <Route exact path="/logout" component={LogoutPage} />

    <Route exact path="/profile" component={ProfilePage} />

    <Route exact path="/audience" component={AudiencePage} />
    <Route
      exact
      path="/audience/contacts/import"
      component={ImportContactsPage}
    />

    <Route exact path="/templates" component={TemplatesPage} />
    <Route exact path="/templates/:uuid/edit" component={EditTemplatePage} />

    <Route exact path="/campaigns" component={CampaignsPage} />
    <Route exact path="/campaigns/:uuid/edit" component={EditCampaignPage} />
    <Route
      exact
      path="/campaigns/:uuid/review"
      component={ReviewCampaignPage}
    />
  </Switch>
);

export default AppRouter;
