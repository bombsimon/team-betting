import React from "react";

export function FormGroupInput({ id, type, name, value, onChange }) {
  return (
    <div className="form-group">
      <label htmlFor={id}>{name}</label>
      <input
        type={type === undefined ? "text" : type}
        name={id}
        className="form-control"
        id={id}
        value={value === null ? "" : value}
        aria-describedby={name + "Help"}
        placeholder={name}
        onChange={onChange}
      />
    </div>
  );
}

export function FormGroupSelect({ name, label, value, options, onChange }) {
  return (
    <div className="form-group">
      <label htmlFor={name}>{label}</label>
      <select
        className="form-control"
        name={name}
        value={value === null ? "" : value}
        onChange={onChange}
      >
        {options}
      </select>
    </div>
  );
}

export function SmallDate({ date }) {
  return (
    <small>
      Created at{" "}
      {new Intl.DateTimeFormat("en-US", {
        year: "numeric",
        month: "long",
        day: "2-digit"
      }).format(new Date(date))}
    </small>
  );
}

export function SetBoolKey(fn, keyName, value) {
  fn(state => ({
    ...state,
    [keyName]: value
  }));
}

const Generic = {
  SetBoolKey,
  SmallDate,
  FormGroupInput,
  FormGroupSelect
};

export default Generic;

// vim: set sw=2 ts=2 et:
