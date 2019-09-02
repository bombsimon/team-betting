import React, { useState, useEffect } from 'react'

import { Competition } from '../Competition'
import { AddCompetitor, Competitor } from '../Competitor'
import { Bet } from '../Bet'
import HttpService from '../HttpClient'

export default function CompetitionPage(props) {
  const [state, setCompetition] = useState({ competition: {}, loading: true })
  const [betsPerCompetitor, setBetsPerCompetitor] = useState({})

  useEffect(() => {
    const getCompetition = async () => {
      const apiResult = await HttpService.GetCompetition(props.match.params.code)

      setCompetition(state => ({
        ...state,
        competition: apiResult
      }))

      // TODO: This should either be a list of everyones bets or just the
      // current users bet, although there's no user state implemented yet.
      let bpc = {}
      apiResult.bets.map((item, key) =>
        bpc[item.competitor.id] = item
      )

      setBetsPerCompetitor(bpc)

      setLoading(false)
    }

    getCompetition()
   }, [props.match.params.code])

  const setLoading = bool => {
    setCompetition(state => ({
      ...state,
      loading: bool
    }))
  }

  const updateCompetitors = competitor => {
    setCompetition(state => ({
      ...state,
      competition: {
        ...state.competition,
        competitors: [...state.competition.competitors, competitor]
      }
    }))
  }

  const updateBetPerCompetitor = bet => {
    setBetsPerCompetitor(state => ({
      ...betsPerCompetitor,
      [bet.competitor_id]: bet
    }))
  }

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
      {state.competition.competitors.map((competitor) =>
        <div
          key={competitor.id} 
          style={{
            display: 'flex',
            borderBottom: '1px solid #ccc',
            padding: '20px'
          }}
        >
          <div style={{float: 'left', flex: '50%'}}>
            <Competitor competitor={competitor} />
          </div>
          <div style={{flex: '50%'}}>
            <Bet
              competitorId={competitor.id}
              competition={state.competition}
              bets={betsPerCompetitor[competitor.id]}
              onAddedBet={updateBetPerCompetitor} 
              selectInputs={false}
            />
          </div>
        </div>
      )}
      <hr />
    </div>
  )
}

// vim: set ts=2 sw=2 et:
