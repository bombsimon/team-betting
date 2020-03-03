import React, { useState, useEffect } from "react";

import {
  sortableContainer,
  sortableElement,
  sortableHandle,
  arrayMove
} from "react-sortable-hoc";
import { Competition } from "../Competition";
import { AddCompetitor, Competitor } from "../Competitor";
import { Bet } from "../Bet";
import HttpService from "../HttpClient";

export default function CompetitionPage(props) {
  const [state, setCompetition] = useState({ competition: {}, loading: true });
  const [betsPerCompetitor, setBetsPerCompetitor] = useState({});
  const {
    match: {
      params: { code }
    }
  } = props;

  const setLoading = bool => {
    setCompetition(prev => ({
      ...prev,
      loading: bool
    }));
  };

  useEffect(() => {
    const getCompetition = async () => {
      const apiResult = await HttpService.Request({
        headers: {
          Authorization: `Bearer ${localStorage.getItem("authorization")}`
        },
        method: "get",
        url: `/competition/${code}`
      });

      setCompetition(prev => ({
        ...prev,
        competition: apiResult
      }));

      // TODO: This should either be a list of everyones bets or just the
      // current users bet, although there's no user state implemented yet.
      const bpc = {};
      apiResult.bets.map(item => {
        bpc[item.competitor.id] = item;
        return item;
      });

      setBetsPerCompetitor(bpc);

      setLoading(false);
    };

    getCompetition();
  }, [code]);

  const updateCompetitors = competitor => {
    setCompetition(prev => ({
      ...prev,
      competition: {
        ...prev.competition,
        competitors: [...prev.competition.competitors, competitor]
      }
    }));
  };

  const updateBetPerCompetitor = bet => {
    setBetsPerCompetitor({
      ...betsPerCompetitor,
      [bet.competitor_id]: bet
    });
  };

  const onSortEnd = ({ oldIndex, newIndex }) => {
    setCompetition(prev => ({
      ...prev,
      competition: {
        ...prev.competition,
        competitors: arrayMove(
          state.competition.competitors,
          oldIndex,
          newIndex
        )
      }
    }));
  };

  const DragHandle = sortableHandle(() => <div className="DragHandle" />);

  const SortableItem = sortableElement(({ value }) => (
    <li className="SortableItem">
      <DragHandle />
      {value.name} - {value.description}
    </li>
  ));

  const SortableContainer = sortableContainer(({ children }) => {
    return <ul className="SortableList">{children}</ul>;
  });

  return state.loading ? (
    <div>Loading...</div>
  ) : (
    <div className="container">
      <Competition competition={state.competition} />
      <hr />
      <AddCompetitor
        competitionId={state.competition.id}
        onAddedCompetitor={updateCompetitors}
      />
      <hr />
      <h1>Competitors for competition</h1>
      {/*
      {state.competition.competitors.map(competitor => (
        <div
          key={competitor.id}
          style={{
            display: "flex",
            borderBottom: "1px solid #ccc",
            padding: "20px"
          }}
        >
          <div style={{ float: "left", flex: "50%" }}>
            <Competitor competitor={competitor} />
          </div>
          <div style={{ flex: "50%" }}>
            <Bet
              competitorId={competitor.id}
              competition={state.competition}
              bets={betsPerCompetitor[competitor.id]}
              onAddedBet={updateBetPerCompetitor}
              selectInputs
            />
          </div>
        </div>
      ))}
      */}
      <SortableContainer onSortEnd={onSortEnd} useDragHandle>
        {state.competition.competitors.map((value, index) => (
          <SortableItem
            helperClass="SortableHelper"
            key={`item-${value.id}`}
            index={index}
            value={value}
          />
        ))}
      </SortableContainer>
      <hr />
    </div>
  );
}

// vim: set ts=2 sw=2 et:
