import React, { useState } from "react";
import { Route, Switch } from "react-router-dom";

import "./index.css";

import CompetitionPage from "./pages/CompetitionPage";
import CompetitionsPage from "./pages/CompetitionsPage";
import RegisterPage from "./pages/RegisterPage";

export default function App() {
  const [alert, setAlert] = useState({});

  const alertBlock =
    alert && alert.message && alert.level ? (
      <div
        className={`alert alert-${alert.level} alert-dismissible fade show`}
        role="alert"
      >
        {alert.message}
        <button
          type="button"
          className="close"
          data-dismiss="alert"
          aria-label="Close"
        >
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
    ) : (
      ""
    );

  return (
    <>
      {alertBlock}

      <Switch>
        <Route
          exact
          path="/"
          render={props => <RegisterPage {...props} flash={setAlert} />}
        />
        <Route
          path="/list"
          render={props => <CompetitionsPage {...props} flash={setAlert} />}
        />
        <Route
          path="/:code"
          render={props => <CompetitionPage {...props} flash={setAlert} />}
        />
      </Switch>
    </>
  );
}

// vim: set ts=2 sw=2 et:
