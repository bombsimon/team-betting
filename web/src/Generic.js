import React from "react";

export function FormGroupInput({id, type, name, value, onChange}) {
    return (
      <div className="form-group">
        <label htmlFor={id}>{name}</label>
        <input
            type={type === undefined ? "text" : type}
            name={id}
            className="form-control"
            id={id}
            value={value === null ? undefined : value}
            aria-describedby={name + "Help"}
            placeholder={name}
            onChange={onChange}
        />
      </div>
    )
}

export function FormGroupSelect({ name, label, value, options, onChange }) {
    return (
        <div className="form-group">
        <label forHtml={name}>{label}</label>
        <select className="form-control" name={name} value={value} onChange={onChange}>
          {options}
        </select>
      </div>
    )
}

const Generic = {
    FormGroupInput, FormGroupSelect
}

export default Generic;
