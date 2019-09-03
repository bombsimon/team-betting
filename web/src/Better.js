import React, { useState } from "react";

import Generic from "./Generic";
import HttpService from "./HttpClient";

export function SaveBetter({ current }) {
  const initialBetterState = {
    name: current === undefined ? "" : current.name,
    email: current === undefined ? "" : current.email,
    image: current === undefined ? "" : current.image,
    confirmed: current === undefined ? false : current.confirmed
  };

  console.log(current);

  const initialButtonState = {
    label: current === undefined ? "Add" : "Update",
    disabled: true
  };

  const [better, setBetter] = useState(initialBetterState);
  const [button, setButton] = useState(initialButtonState);

  const onSubmit = event => {
    event.preventDefault();

    if (button.disabled === true) {
      return;
    }

    (async () => {
      const apiResult = await HttpService.Request({
        method: "post",
        url: "/better",
        data: better
      });

      if (apiResult !== undefined) {
        setBetter(initialBetterState);
      }
    })();
  };

  const handleInputChange = event => {
    const { name, value } = event.target;

    // Assume valid - set button to true.
    Generic.SetBoolKey(setButton, "disabled", false);

    setBetter(state => ({
      ...state,
      [name]: value
    }));
  };

  const handleImageClick = event => {
    const { alt } = event.target;

    setBetter(state => ({
      ...state,
      image: alt
    }));
  };

  const images = [];
  for (let i = 1; i <= 5; i++) {
    const filename = `avatar${i}.png`;
    const imgStyle = {
      cursor: "pointer",
      width: "64px",
      padding: "10px",
      border: filename === better.image ? "2px solid black" : ""
    };

    images.push(
      <img
        key={filename}
        alt={filename}
        src={"avatar/" + filename}
        style={imgStyle}
        onClick={handleImageClick}
      />
    );
  }

  return (
    <form onSubmit={onSubmit}>
      <Generic.FormGroupInput
        id="name"
        name="Name"
        value={better.name}
        onChange={handleInputChange}
      />

      <Generic.FormGroupInput
        id="email"
        name="Email"
        value={better.email}
        onChange={handleInputChange}
      />

      <div className="form-group">{images}</div>

      <button
        className={
          "btn btn-lg btn-primary" + (button.disabled ? " disabled" : "")
        }
      >
        {button.label}
      </button>
    </form>
  );
}

export function SendLoginEmail() {
  const [email, setEmail] = useState("");

  const handleInputChange = event => {
    const { value } = event.target;

    setEmail(value);
  };

  const onSubmit = event => {
    event.preventDefault();

    if (email === "") {
      return;
    }
  };

  return (
    <div>
      <h1>Been here before?</h1>
      <p className="lead">
        Just write your e-mail and we'll send you a sign in link!
      </p>
      <form onSubmit={onSubmit}>
        <Generic.FormGroupInput
          id="email"
          name="Email"
          onChange={handleInputChange}
        />

        <button className="btn btn-lg btn-primary">Send!</button>
      </form>
    </div>
  );
}

const BetterService = {
  SendLoginEmail,
  SaveBetter
};

export default BetterService;

// vim: set sw=2 ts=2 et:
