import React from "react";
import { Link } from 'react-router-dom'

import Generic from './Generic'

export function AddCompetition({ onSubmit, onChange }) {
  return (
    <form onSubmit={onSubmit}>
      <Generic.FormGroupInput id="name" name="Name" onChange={onChange} />
      <Generic.FormGroupInput id="description" name="Description" onChange={onChange} />
      <Generic.FormGroupInput id="min_score" name="Minimum score" onChange={onChange} />
      <Generic.FormGroupInput id="max_score" name="Maximum score" onChange={onChange} />

      <button>Add</button>
    </form>
  )
}

export function Competition({ data }) {
  const fields = [
    'name', 'description', 'code',
    'created_at', 'updated_at',
    'min_score', 'max_score'
  ]

  const tableBody = fields.map((key) =>
    <tr key={key}>
      <td>{key}</td>
      <td>{data[key]}</td>
    </tr>
  );

  return (
    <div>
      <h1>{data.name}</h1>
      <table style={{width: '100%'}}>
        <thead>
          <tr>
            <th>Key</th>
            <th>Value</th>
          </tr>
        </thead>
        <tbody>
          {tableBody}
        </tbody>
      </table>
    </div>
  )
}

export function CompetitionLink(props) {
  return (
    <div>
      <Link
        key="{i}"
        to={"/" + props.competition.id}
      >
        <h1>{props.competition.name}</h1>
      </Link>
      <p>{props.competition.description}</p>
    </div>
  )
}

const CompetitionData = {
  Competition, CompetitionLink, AddCompetition
}

export default CompetitionData;

// vim: set sw=2 ts=2 et:
