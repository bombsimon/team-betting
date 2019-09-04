import React from "react";

import { SendLoginEmail, SaveBetter } from "../Better";

export default function RegisterPage({ flash }) {
  return (
    <div className="container">
      <h1>Been here before?</h1>
      <p className="lead">
        Just write your e-mail and we&apos;ll send you a sign in link!
      </p>
      <SendLoginEmail flash={flash} />

      <hr />

      <h1>Register here</h1>
      <p className="lead">
        If you wan&apos;t to save all your contributions in all competitions
        enter your email address. You don&apos;t have to do anything now but the
        next time you get here we can send you an email to keep your progress!
      </p>
      <SaveBetter flash={flash} />
    </div>
  );
}

// vim: set ts=2 sw=2 et:
