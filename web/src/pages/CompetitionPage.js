import React, { useState, useEffect } from 'react'

import { Competition } from '../Competition'
import { AddCompetitor, Competitor } from '../Competitor'
import { Bet } from '../Bet'
import HttpService from '../HttpClient'

export default function CompetitionPage(props) {
  const initialCompetitionState = {
    competition: {},
    betsPerCompetitor: {},
    loading: true,
  }

  const [state, setCompetition] = useState(initialCompetitionState)

  useEffect(() => {
    const getCompetition = async () => {
      const apiResult = await HttpService.GetCompetition(props.match.params.code)

      let bpc = {}
      apiResult.bets.map((item, key) =>
        bpc[item.competitor.id] = item
      )

      setCompetition(state => ({
        ...state,
        competition: apiResult,
        betsPerCompetitor: bpc,
        loading: false
      }))
    }

    getCompetition()
   }, [props.match.params.code])

  const updateCompetitors = competitor => {
    let { competition }  = state

    competition.competitors = [...competition.competitors, competitor]

    setCompetition(state => ({
      ...state,
      competition: competition,
    }))
  }

  return state.loading ? (
      <div>Loading...</div>
  ) : (
    <div className="container">
      <h1>The competition</h1>
      <Competition data={state.competition} />
      <hr />

      <h1>All competitors</h1>
      <AddCompetitor
        competitionId={state.competition.id}
        onCompetitorAdded={updateCompetitors}
      />
      <hr />

      {state.competition.competitors.map((competitor) =>
        <div key={competitor.id} style={{display: 'flex'}}>
          <div style={{flex: '50%', paddingRight: '20px'}}>
            <Competitor competitor={competitor} />
          </div>
          <div style={{flex: '50%'}}>
            <h1>My bet</h1>
            <Bet
              competitionId={state.competition.id}
              competitorId={competitor.id}
              betData={state.betsPerCompetitor[competitor.id]}
              minScore={state.competition.min_score}
              maxScore={state.competition.max_score}
              maxPos={state.competition.competitors.length}
              selectInputs={true}
            />
          </div>
        </div>
      )}
      <hr />
    </div>
  )
}

// vim: set ts=2 sw=2 et:
