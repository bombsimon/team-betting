import React, { useState, useEffect } from 'react'

import HttpService from '../HttpClient'
import { AddCompetition, CompetitionLink } from '../Competition'


export default function HomePage() {
  const initialCompetitionsState = {
    competitions: {},
    competitionToAdd: {
      created_by_id: 1,
      name: '',
      description: '',
      min_score: null,
      max_score: null
    },
    loading: true,
  }

  const [competitions, setCompetitions] = useState(initialCompetitionsState)

  useEffect(() => {
    const getCompetitions = async () => {
      const apiResult = await HttpService.Request({
        method: 'get',
        url: '/competition'
      })

      setCompetitions(state => ({
        ...state,
        competitions: apiResult,
        loading: false
      }))
    }

    getCompetitions()
  }, [setCompetitions])

  const handleInputChange = event => {
    const { name } = event.target
    let { value } = event.target

    const integers = [
      'min_score', 'max_score'
    ]

    if (integers.includes(name)) {
      if (Number.isNaN(Number(value))) {
        return
      }

      value = Number(value)
    }

    setCompetitions(state => ({
      ...state,
      competitionToAdd: {
        ...state.competitionToAdd,
        [name]: value
      }
    }))
  }

  const onSubmit = event => {
    event.preventDefault()

    ;(async () => {
      const apiResult = await HttpService.Request({
        method: 'post',
        url: '/competition',
        data: competitions.competitionToAdd
      })

      if (apiResult !== undefined) {
        setCompetitions(state => ({
          ...state,
          competitions: [...state.competitions, apiResult]
        }))
      }
    })();
  }

  return competitions.loading ? (
    <div>Loading...</div>
  ) : (
    <div className="container">
      <AddCompetition onSubmit={onSubmit} onChange={handleInputChange}/>
      {competitions.competitions.map((competition) =>
        <CompetitionLink key={competition.id} competition={competition} />
      )}
      <small>{competitions.status} {competitions.statusText}</small>
    </div>
  )
}

// vim: set ts=2 sw=2 et:
