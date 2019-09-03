import React from "react";
import { Link } from "react-router-dom";

import { SendLoginEmail, SaveBetter } from "../Better";

export default function RegisterPage() {
  return (
    <div className="container">
      <h1>Register here</h1>
      <SaveBetter />

      <hr />

      <SendLoginEmail />

      <hr />
      <Link to="/list">
        <h1>Show competitions</h1>
      </Link>
    </div>
  );
}

// vim: set ts=2 sw=2 et:
