import React from "react";
import { Link } from 'react-router-dom'

export function Competition(props) {
  const { data } = props

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
  Competition, CompetitionLink
}

export default CompetitionData;

// vim: set sw=2 ts=2 et:
