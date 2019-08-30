import React, { useState, useEffect } from 'react'

import HttpService from '../HttpClient'
import { CompetitionLink } from '../Competition'


export default function HomePage() {
  const initialCompetitionsState = {
    competitions: {},
    loading: true,
  }

  const [competitions, setCompetitions] = useState(initialCompetitionsState)

  useEffect(() => {
    const getCompetitions = async () => {
      const apiResult = await HttpService.Request({
        method: 'get',
        url: '/competition'
      })

      setCompetitions(apiResult)
    }

    getCompetitions()
  }, []) 

  return competitions.loading ? (
    <div>Loading...</div>
  ) : (
    <div className="container">
      {competitions.map((competition) =>
        <CompetitionLink key={competition.id} competition={competition} />
      )}
      <small>{competitions.status} {competitions.statusText}</small>
    </div>
  )
}

// vim: set ts=2 sw=2 et
