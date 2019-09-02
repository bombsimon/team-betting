import React from 'react'

import { SendLoginEmail, SaveBetter } from '../Better'

export default function RegisterPage() {
  return (
    <div className="container">
      <h1>Register here</h1>
      <SaveBetter />

      <hr />

      <SendLoginEmail />

      <hr />
      <a href="/list"><h1>Show competitions</h1></a>
    </div>
  )
}

// vim: set ts=2 sw=2 et:
