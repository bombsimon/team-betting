import React, { useState, useEffect } from "react";

import HttpService from "../HttpClient";
import { AddCompetition, Competitions } from "../Competition";

export default function CompetitionsPage() {
  const initialState = {
    competitions: null,
    loading: true
  };

  const [state, setState] = useState(initialState);

  useEffect(() => {
    const getCompetitions = async () => {
      const apiResult = await HttpService.Request({
        headers: {
          Authorization: `Bearer ${localStorage.getItem("authorization")}`
        },
        method: "get",
        url: "/competition"
      });

      setState(prev => ({
        ...prev,
        competitions: apiResult,
        loading: false
      }));
    };

    getCompetitions();
  }, [setState]);

  const onAddedCompetition = competition => {
    setState(prev => ({
      ...prev,
      competitions: [...state.competitions, competition]
    }));
  };

  return state.loading ? (
    <div>Loading...</div>
  ) : (
    <div className="container">
      <AddCompetition onAddedCompetition={onAddedCompetition} />

      <hr />

      <Competitions competitions={state.competitions} />
    </div>
  );
}

// vim: set ts=2 sw=2 et:
