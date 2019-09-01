import React, { useState } from "react";

import Generic from './Generic'
import HttpService from './HttpClient'

export function Bet({ competitionId, competitorId, betData, minScore, maxScore, maxPos, selectInputs}) {
  let initialBetState = {
    placing: '',
    score: '',
    note: '',
    btnDisabled: true
  }

  if (betData !== undefined) {
    initialBetState = {
      placing: initialBetState.placing === null ? '' : betData.placing,
      score: initialBetState.score === null ? '' : betData.score,
      note: initialBetState.note === null ? '' : betData.note,
      btnDisabled: true
    }
  }

  const [bet, setBet] = useState(initialBetState)

  const handleInputChange = event => {
    const { name, value } = event.target

    setBet(state => {
      let newState =  {
        ...state,
        [name]: [ "placing", "score" ].includes(name) ? Number(value) : value,
        btnDisabled: false
      }

      if (
        newState.score < minScore ||
        newState.score > maxScore ||
        newState.placing < 1 ||
        newState.placing > maxPos
      ) {
        console.log(`Invalud value for ${name} - cannot use ${value}`)
        newState.btnDisabled = true
      }

      return newState
    })
  }

  const onSubmit = event => {
    event.preventDefault()

    const request = {
      competition_id: competitionId,
      competitor_id: competitorId,
      better_id: 1,
      placing: Number(bet.placing),
      score: Number(bet.score),
      note: bet.note
    }

    ;(async () => {
      const apiResult = await HttpService.Request({
        method: 'put',
        url: '/bet',
        data: request
      })

      if (apiResult !== undefined) {
        setBet(state => ({
          ...state,
          apiResult
        }))
      }
    })();

    // TODO: Super class to update all bets list?
  }

  let numberInputs = null;

  if (selectInputs) {
    let placingOptions = []
    let scoreOptions = []

    for (var i = 1; i <= maxPos; i++) {
      placingOptions.push(<option value={i}>{i}</option>)
    }

    for (var j = minScore; j <= maxScore; j++) {
      scoreOptions.push(<option value={j}>{j}</option>)
    }

    numberInputs = (
      <>
      <Generic.FormGroupSelect
        name="placing"
        label="Placing"
        options={placingOptions}
        value={bet.placing}
        onChange={handleInputChange}
      />

      <Generic.FormGroupSelect
        name="score"
        label="Score"
        options={scoreOptions}
        value={bet.score}
        onChange={handleInputChange}
      />
      </>
    )
  }
  else {
    numberInputs = (
      <>
        <Generic.FormGroupInput
          type="number"
          id="placing"
          name="Placing"
          value={bet.placing}
          onChange={handleInputChange}
        />
        <Generic.FormGroupInput
          type="number"
          id="score"
          name="Score"
          value={bet.score}
          onChange={handleInputChange}
        />
      </>
    )
  }

  return (
    <form onSubmit={onSubmit}>
      {numberInputs}
      <Generic.FormGroupInput
        id="note"
        name="Note"
        value={bet.note}
        onChange={handleInputChange}
      />

    <button
      className={"btn btn-lg btn-primary" + (bet.btnDisabled ? " disabled" : "")}
    >
      {betData === undefined ? 'Add' : 'Update'}
    </button>
    </form>
  )
}

const BetData = {
    Bet
}

export default BetData;

// vim: set sw=2 ts=2 et:
