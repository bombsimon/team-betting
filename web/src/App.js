import React from 'react'
import { Route, Switch } from 'react-router-dom'

import './index.css'

import CompetitionPage from './pages/CompetitionPage'
import CompetitionsPage from './pages/CompetitionsPage'

export default function App() {
  return (
    <Switch>
      <Route exact path="/" component={CompetitionsPage} />
      <Route path="/:code" component={CompetitionPage} />
    </Switch>
  )
}
