import React, { useState, useEffect } from "react";

import Generic from './Generic'
import HttpService from './HttpClient'

export function AddCompetitor({ competitionId, onCompetitorAdded }) {
  const initialCompetitorState = {
    name: '',
    description: '',
    competition_id: competitionId
  }

  const [competitor, setCompetitor] = useState(initialCompetitorState)

  // TODO: Maybe load existing competitors (not bound to competitionId to
  // support adding already existing competitor.
  useEffect(() => {
  }, [])

  const handleInputChange = event => {
    const { name, value } = event.target

    setCompetitor(state => ({
      ...state,
      [name]: value
    }))
  }

  const onSubmit = event => {
    event.preventDefault()

    ;(async () => {
      const apiResult = await HttpService.Request({
        method: 'post',
        url: '/competitor',
        data: {
          ...competitor,
          created_by_id: 1
        }
      })

      if (apiResult !== undefined) {
        onCompetitorAdded(apiResult)
      }
    })();
  }

  return (
    <form onSubmit={onSubmit}>
      <Generic.FormGroupInput id="name" name="Name" onChange={handleInputChange} />
      <Generic.FormGroupInput id="description" name="Description" onChange={handleInputChange} />

      <button className="btn btn-lg btn-primary">Add</button>
    </form>
  )
}

export function Competitor({ competitor, bet }) {
  return (
    <div>
      <h1>{competitor.name}</h1>
      <table className="table">
        <thead>
          <tr>
            <th>Key</th>
            <th>Value</th>
          </tr>
        </thead>
        <tbody>
          {Object.entries(competitor).map(([key, value]) =>
            <tr key={key}>
              <td>{key}</td>
              <td>{value}</td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  )
}

const CompetitorData = {
    AddCompetitor, Competitor
}

export default CompetitorData;

// vim: set sw=2 ts=2 et:
