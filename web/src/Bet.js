import React, { useState } from "react";

import { FormGroupInput, FormGroupSelect, SetBoolKey } from "./Generic";
import HttpService from "./HttpClient";

export function Bet({ competitorId, competition, bets, selectInputs }) {
  const initialBetState = {
    competitor_id: competitorId,
    competition_id: competition.id,
    better_id: 1,
    placing: bets === undefined ? 1 : bets.placing,
    score: bets === undefined ? competition.min_score : bets.score,
    note: bets === undefined ? "" : bets.note
  };

  const initialButtonState = {
    label: bets === undefined ? "Add" : "Update",
    disabled: true
  };

  const [bet, setBet] = useState(initialBetState);
  const [button, setButton] = useState(initialButtonState);

  const handleInputChange = event => {
    const { name, value } = event.target;

    // Assume valid - set button to true.
    SetBoolKey(setButton, "disabled", false);

    // Update the state with the new value. If placing or score and not empty,
    // convert to number.
    const newState = {
      ...bet,
      [name]:
        ["placing", "score"].includes(name) &&
        value !== "" &&
        !Number.isNaN(Number(value))
          ? Number(value)
          : value
    };

    setBet(newState);

    if (
      newState.score < competition.min_score ||
      newState.score > competition.max_score ||
      newState.placing < 1 ||
      newState.placing > competition.competitors.length
    ) {
      // The new data on the state is not valid, disable button again.
      SetBoolKey(setButton, "disabled", true);
    }
  };

  const onSubmit = event => {
    event.preventDefault();

    if (button.disabled) {
      return;
    }

    (async () => {
      const apiResult = await HttpService.Request({
        headers: {
          Authorization: `Bearer ${localStorage.getItem("authorization")}`
        },
        method: "put",
        url: "/bet",
        data: {
          ...bet,
          placing: bet.placing === "" ? null : bet.placing,
          score: bet.score === "" ? null : bet.score
        }
      });

      if (apiResult !== undefined) {
        setBet(apiResult);
        setButton({
          label: "Update",
          disabled: true
        });
      }
    })();

    // TODO: Super class to update all bets list?
  };

  let numberInputs = null;

  if (selectInputs) {
    const placingOptions = [];
    const scoreOptions = [];

    for (let i = 1; i <= competition.competitors.length; i++) {
      placingOptions.push(
        <option key={i} value={i}>
          {i}
        </option>
      );
    }

    for (let j = competition.min_score; j <= competition.max_score; j++) {
      scoreOptions.push(
        <option key={j} value={j}>
          {j}
        </option>
      );
    }

    numberInputs = (
      <>
        <FormGroupSelect
          name="placing"
          label="Placing"
          options={placingOptions}
          value={bet.placing}
          onChange={handleInputChange}
        />

        <FormGroupSelect
          name="score"
          label="Score"
          options={scoreOptions}
          value={bet.score}
          onChange={handleInputChange}
        />
      </>
    );
  } else {
    numberInputs = (
      <>
        <FormGroupInput
          type="number"
          id="placing"
          name="Placing"
          value={bet.placing}
          onChange={handleInputChange}
        />
        <FormGroupInput
          type="number"
          id="score"
          name="Score"
          value={bet.score}
          onChange={handleInputChange}
        />
      </>
    );
  }

  return (
    <form onSubmit={onSubmit}>
      {numberInputs}
      <FormGroupInput
        id="note"
        name="Note"
        value={bet.note}
        onChange={handleInputChange}
      />

      <button
        type="submit"
        className={`btn btn-lg btn-primary${
          button.disabled ? " disabled" : ""
        }`}
      >
        {button.label}
      </button>
    </form>
  );
}

const BetData = {
  Bet
};

export default BetData;

// vim: set sw=2 ts=2 et:
