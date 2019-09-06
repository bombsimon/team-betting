import React, { useState } from "react";

import { FormGroupInput, SetBoolKey } from "./Generic";
import HttpService from "./HttpClient";

export function SaveBetter({ current, flash, onSave }) {
  const initialBetterState = {
    name: current === undefined ? "" : current.name,
    email: current === undefined ? "" : current.email,
    image: current === undefined ? "" : current.image
  };

  const initialButtonState = {
    label: current === undefined ? "Register" : "Update",
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
      try {
        const result = await HttpService.Request({
          headers: {
            Authorization: `Bearer ${localStorage.getItem("authorization")}`
          },
          method: "post",
          url: "/better",
          data: better
        });

        setBetter(initialBetterState);

        localStorage.setItem("authorization", result.jwt);

        onSave();
      } catch (error) {
        // HTTP 400?
        flash({
          level: "danger",
          message: error.data.error
        });
      }
    })();
  };

  const handleInputChange = event => {
    const { name, value } = event.target;

    // Assume valid - set button to true.
    SetBoolKey(setButton, "disabled", false);

    setBetter(state => ({
      ...state,
      [name]: value
    }));
  };

  const handleImageClick = event => {
    const { alt } = event.target;

    setBetter(state => ({
      ...state,
      image: state.image === alt ? undefined : alt
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
      <button
        key={filename}
        type="button"
        onClick={handleImageClick}
        style={{ background: "none", border: "none", outline: 0 }}
      >
        <img alt={filename} src={`avatar/${filename}`} style={imgStyle} />
      </button>
    );
  }

  return (
    <form onSubmit={onSubmit}>
      <FormGroupInput
        id="name"
        name="Name"
        value={better.name}
        onChange={handleInputChange}
      />

      <FormGroupInput
        id="email"
        name="Email"
        value={better.email}
        onChange={handleInputChange}
      />

      <div className="form-group">{images}</div>

      <div className="form-group">
        <button
          type="submit"
          className={`btn btn-lg btn-primary${
            button.disabled ? " disabled" : ""
          }`}
        >
          {button.label}
        </button>
      </div>
    </form>
  );
}

export function SendLoginEmail({ flash }) {
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

    (async () => {
      try {
        await HttpService.Request({
          method: "post",
          url: "/email/send",
          data: {
            email
          }
        });

        setEmail("");
      } catch (error) {
        flash({
          level: "danger",
          message: error.data.error
        });
      }
    })();
  };

  return (
    <div>
      <form onSubmit={onSubmit}>
        <FormGroupInput
          value={email}
          id="email"
          name="Email"
          onChange={handleInputChange}
        />

        <button type="submit" className="btn btn-lg btn-primary">
          Send!
        </button>
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
