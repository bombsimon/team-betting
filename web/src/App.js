import React from 'react'
import { Route, Switch } from 'react-router-dom'

import './index.css'

import CompetitionPage from './pages/CompetitionPage'
import HomePage from './pages/HomePage'

export default function App() {
  return (
    <Switch>
      <Route exact path="/" component={HomePage} />
      <Route path="/:code" component={CompetitionPage} />
    </Switch>
  )
}
