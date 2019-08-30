import React from "react";

export function FormGroupInput({id, name, onChange}) {
    return (
      <div className="form-group">
        <label htmlFor={id}>{name}</label>
        <input
            type="text"
            name={id}
            className="form-control"
            id={id}
            aria-describedby={name + "Help"}
            placeholder={name}
            onChange={onChange}
        />
      </div>
    )
}

const Generic = {
    FormGroupInput
}

export default Generic;
