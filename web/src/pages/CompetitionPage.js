import React, { useState, useEffect } from 'react'
import axios from 'axios'

export default function CompetitionPage(props) {
  const initialCompetitionState = {
    competition: {},
    loading: true,
  }

  const [competition, setCompetition] = useState(initialCompetitionState)

  useEffect(() => {
    const getCompetition = async () => {
      let apiResult = null;

      try {
        const { data } = await axios(`http://localhost:5000/competition/${props.match.params.code}`)

        apiResult = data
      } catch (err) {
        apiResult = err.response;
      } finally {
        console.log(apiResult)
      }

      setCompetition(apiResult)
    }

    // Invoke the async function if competition not passed.
    const passedCompetition = props.location.state
    if ( passedCompetition === undefined ) {
      getCompetition()
    }
    else {
      setCompetition(passedCompetition.competition)
    }
  }, [props.location.state, props.match.params.code]) // Don't forget the `[]`, which will prevent useEffect from running in an infinite loop

  return competition.loading ? (
      <div>Loading...</div>
  ) : (
      <div className="container">
          <h1>{competition.name}</h1>
          <p>{competition.description}</p>
          <small>{competition.status} {competition.statusText}</small>
      </div>
  )
}

// vim: set ts=2 sw=2 et:
