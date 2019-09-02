import React, { useState } from "react";

import Generic from './Generic'
import HttpService from './HttpClient'

export function SaveBetter({ better }) {
  let initialBetterState = {
    better: {
      name: '',
      email: '',
      image: '',
      confirmed: false,
    },
    btnDisabled: true,
  }

  if (better !== undefined) {
    initialBetterState = better
  }

  const [state, setState] = useState(initialBetterState)

  const onSubmit = event => {
    event.preventDefault()

    ;(async () => {
      const apiResult = await HttpService.Request({
        method: 'post',
        url: '/better',
        data: state.better
      })

      if (apiResult !== undefined) {
        setState(state => ({
          ...state,
          better: apiResult
        }))
      }
    })();
  }

  const handleInputChange = event => {
    const { name, value } = event.target

    setState(state => ({
      ...state,
      btnDisabled: false,
      better: {
        ...state.better,
        [name]: value,
      }
    }))
  }

  const handleImageClick = event => {
    const { alt } = event.target

    setState(state => ({
      ...state,
        better: {
          ...state.better,
          image: alt
        }
    }))
  }

  const images = []
  for (let i = 1; i <= 5; i++) {
    const filename = `avatar${i}.png`
    const imgStyle = {
      cursor: 'pointer',
      width: '64px',
      padding: '10px',
      border: (filename === state.better.image ? '2px solid black' : '' )
    }

    images.push(<img key={filename} alt={filename} src={'avatar/' + filename} style={imgStyle} onClick={handleImageClick} />)
  }

  return (
    <form onSubmit={onSubmit}>
      <Generic.FormGroupInput
        id="name"
        name="Name"
        value={state.better.name}
        onChange={handleInputChange}
      />

      <Generic.FormGroupInput
        id="email"
        name="Email"
        value={state.better.email}
        onChange={handleInputChange}
      />

      <div className="form-group">
        {images}
      </div>

      <button
        className={"btn btn-lg btn-primary" + (state.btnDisabled ? " disabled" : "")}
      >
        {better === undefined ? 'Add' : 'Update'}
      </button>
    </form>
  )
}

export function SendLoginEmail() {
  return (
    <div>
      <h1>Been here before?</h1>
      <p className="lead">Just write your e-mail and we'll send you a sign in link!</p>
      <form>
        <Generic.FormGroupInput
          id="email"
          name="Email"
        />

        <button className="btn btn-lg btn-primary">Send!</button>
      </form>
    </div>
  )
}

const BetterService = {
  SendLoginEmail, SaveBetter
}

export default BetterService;

// vim: set sw=2 ts=2 et:
