import React, { useState } from "react";

import Generic from "./Generic";
import HttpService from "./HttpClient";

export function AddCompetitor({ competitionId, onAddedCompetitor }) {
  const initialCompetitorState = {
    created_by_id: 1,
    name: "",
    description: "",
    competition_id: competitionId
  };

  const [competitor, setCompetitor] = useState(initialCompetitorState);

  const handleInputChange = event => {
    const { name, value } = event.target;

    setCompetitor(state => ({
      ...state,
      [name]: value
    }));
  };

  const onSubmit = event => {
    event.preventDefault();
    (async () => {
      const apiResult = await HttpService.Request({
        headers: {
          Authorization: `Bearer ${localStorage.getItem("authorization")}`
        },
        method: "post",
        url: "/competitor",
        data: competitor
      });

      if (apiResult !== undefined) {
        onAddedCompetitor(apiResult);
      }
    })();

    setCompetitor(initialCompetitorState);
  };

  return (
    <div>
      <h1>Add competitor to competition</h1>
      <form onSubmit={onSubmit}>
        <Generic.FormGroupInput
          value={competitor.name}
          id="name"
          name="Name"
          onChange={handleInputChange}
        />
        <Generic.FormGroupInput
          value={competitor.description}
          id="description"
          name="Description"
          onChange={handleInputChange}
        />

        <button type="submit" className="btn btn-lg btn-primary">
          Add
        </button>
      </form>
    </div>
  );
}

export function Competitor({ competitor }) {
  return (
    <div>
      <h2>{competitor.name}</h2>
      <p className="lead">{competitor.description}</p>
      <Generic.SmallDate date={competitor.created_at} />
    </div>
  );
}

const CompetitorData = {
  AddCompetitor,
  Competitor
};

export default CompetitorData;

// vim: set sw=2 ts=2 et:
