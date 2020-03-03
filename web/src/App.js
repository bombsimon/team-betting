import React, { useState } from "react";
import { Redirect, Route, Switch } from "react-router-dom";

import "./index.css";
import "./App.scss";

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

  const isUserLoggedIn = () => {
    const jwt = localStorage.getItem("authorization");

    return Boolean(jwt);
  };

  return (
    <>
      {alertBlock}

      <Switch>
        <Route
          path="/login"
          render={props =>
            isUserLoggedIn() ? (
              <Redirect to="/" />
            ) : (
              <RegisterPage {...props} flash={setAlert} />
            )
          }
        />
        <Route
          path="/"
          exact
          render={props =>
            isUserLoggedIn() ? (
              <CompetitionsPage {...props} flash={setAlert} />
            ) : (
              <Redirect to="/login" />
            )
          }
        />
        <Route
          path="/:code"
          render={props =>
            isUserLoggedIn() ? (
              <CompetitionPage {...props} flash={setAlert} />
            ) : (
              <Redirect to="/login" />
            )
          }
        />
      </Switch>
    </>
  );
}

// vim: set ts=2 sw=2 et:
