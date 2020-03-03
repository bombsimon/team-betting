import React, { useState, useEffect } from "react";
import { Redirect } from "react-router-dom";

import { SendLoginEmail, SaveBetter } from "../Better";

import HttpService from "../HttpClient";

export default function RegisterPage({ location, flash }) {
  const params = new URLSearchParams(location.search);
  const linkData = params.get("signin");

  const [validLink, setValidLink] = useState(false);

  useEffect(() => {
    (async () => {
      try {
        if (!linkData) {
          return;
        }

        const result = await HttpService.Request({
          method: "post",
          url: "/email/verify",
          data: {
            encoding: linkData
          }
        });

        localStorage.setItem("authorization", result.jwt);
        setValidLink(true);
      } catch (error) {
        // HTTP 400?
        flash({
          level: "danger",
          message: error.data.error
        });
      }
    })();
  }, [flash, linkData, setValidLink]);

  const betterSaved = () => {
    setValidLink(true);
  };

  return validLink ? (
    <Redirect to="/" />
  ) : (
    <div className="container">
      <h1>Register here</h1>
      <p className="lead">
        If you want to save all your contributions in all competitions enter
        your email address. You don&apos;t have to do anything now but the next
        time you get here we can send you an email to keep your progress!
      </p>
      <SaveBetter flash={flash} onSave={betterSaved} />

      <hr />

      <h1>Been here before?</h1>
      <p className="lead">
        Just write your e-mail and we&apos;ll send you a sign in link!
      </p>
      <SendLoginEmail flash={flash} />
    </div>
  );
}

// vim: set ts=2 sw=2 et:
