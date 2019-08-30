import React, { useState, useEffect } from 'react'

import { Competition } from '../Competition'
import HttpService from '../HttpClient'

export default function CompetitionPage(props) {
  const initialCompetitionState = {
    competition: {},
    loading: true,
  }

  const [competition, setCompetition] = useState(initialCompetitionState)

  useEffect(() => {
    const getCompetition = async () => {
      const apiResult = await HttpService.GetCompetition(props.match.params.code)

      setCompetition(apiResult)
    }

    getCompetition()
   }, [props.match.params.code])


  return competition.loading ? (
      <div>Loading...</div>
  ) : (
      <Competition data={competition} />
  )
}

// vim: set ts=2 sw=2 et:
