import React, { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import axios from 'axios'

export default function HomePage() {
  const initialCompetitionsState = {
    competitions: {},
    loading: true,
  }

  const [competitions, setCompetitions] = useState(initialCompetitionsState)

  useEffect(() => {
    const getCompetitions = async () => {
      let apiResult = null;

      try {
        const { data } = await axios(`http://localhost:5000/competition`)

        apiResult = data
      } catch (err) {
        apiResult = err.response;
      } finally {
        console.log(apiResult)
      }

      setCompetitions(apiResult)
    }

    // Invoke the async function
    getCompetitions()
  }, []) 

  const cLinks = []
  for (let i = 0; i <= competitions.length; i++) {
    const value = competitions[i]

    if (value === undefined) {
      continue
    }

      cLinks.push(
        <div>
          <Link
              key="{i}"
              to={{
                  pathname: "/" + value.id,
                  state: {
                    competition: value
                  }
              }}
          >
            <h1>{value.name}</h1>
          </Link>
          <p>{value.description}</p>
        </div>
      )
  }

  return competitions.loading ? (
      <div>Loading...</div>
  ) : (
      <div className="container">
          {cLinks}
          <small>{competitions.status} {competitions.statusText}</small>
      </div>
  )
}
