import React, { useState } from "react";
import { Link } from "react-router-dom";

import Generic from "./Generic";
import HttpService from "./HttpClient";

export function AddCompetition(props) {
  const initialState = {
    created_by_id: 1, // TODO: Get from JWT/localStore
    name: "",
    description: "",
    min_score: "",
    max_score: ""
  };

  const [competition, setCompetition] = useState(initialState);

  const handleInputChange = event => {
    const { name, value } = event.target;
    let finalValue = value;

    if (["min_score", "max_score"].includes(name)) {
      if (Number.isNaN(Number(value))) {
        finalValue = competition[name];
      } else if (value === "") {
        finalValue = "";
      } else {
        finalValue = Number(value);
      }
    }

    setCompetition(state => ({
      ...state,
      [name]: finalValue
    }));
  };

  const onSubmit = event => {
    event.preventDefault();

    const request = {
      ...competition,
      min_score: competition.min_score === "" ? null : competition.min_score,
      max_score: competition.max_score === "" ? null : competition.max_score
    };

    (async () => {
      const apiResult = await HttpService.Request({
        method: "post",
        url: "/competition",
        data: request
      });

      if (apiResult !== undefined && apiResult.error === undefined) {
        props.onAddedCompetition(apiResult);
      }
    })();

    setCompetition(initialState);
  };

  return (
    <form onSubmit={onSubmit}>
      <Generic.FormGroupInput
        value={competition.name}
        id="name"
        name="Name"
        onChange={handleInputChange}
      />
      <Generic.FormGroupInput
        value={competition.description}
        id="description"
        name="Description"
        onChange={handleInputChange}
      />
      <Generic.FormGroupInput
        value={competition.min_score}
        id="min_score"
        name="Minimum score"
        onChange={handleInputChange}
      />
      <Generic.FormGroupInput
        value={competition.max_score}
        id="max_score"
        name="Maximum score"
        onChange={handleInputChange}
      />

      <button className="btn btn-lg btn-primary">Add</button>
    </form>
  );
}

export function Competition({ competition }) {
  return (
    <div>
      <h1>{competition.name}</h1>
      <p className="lead">{competition.description}</p>
      <Generic.SmallDate date={competition.created_at} />
    </div>
  );
}

export function CompetitionLink({ competition }) {
  return (
    <div>
      <Link to={"/" + competition.id}>
        <h1>{competition.name}</h1>
      </Link>
      <p>{competition.description}</p>
    </div>
  );
}

export function Competitions({ competitions }) {
  return (
    <div>
      {competitions.map(competition => (
        <CompetitionLink key={competition.id} competition={competition} />
      ))}
    </div>
  );
}

const CompetitionData = {
  AddCompetition,
  Competition,
  CompetitionLink,
  Competitions
};

export default CompetitionData;

// vim: set sw=2 ts=2 et:
