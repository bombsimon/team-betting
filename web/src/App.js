import React from "react";
import { Route, Switch } from "react-router-dom";

import "./index.css";

import CompetitionPage from "./pages/CompetitionPage";
import CompetitionsPage from "./pages/CompetitionsPage";
import RegisterPage from "./pages/RegisterPage";

export default function App() {
  return (
    <Switch>
      <Route exact path="/" component={RegisterPage} />
      <Route path="/list" component={CompetitionsPage} />
      <Route path="/:code" component={CompetitionPage} />
    </Switch>
  );
}

// vim: set ts=2 sw=2 et:
