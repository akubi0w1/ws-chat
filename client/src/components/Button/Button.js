import React from 'react';

export const AccentButton = (props) => {
    return (
        <button type={props.type} onClick={props.handleClick} className="btn btn-accent">
            {props.label}
        </button>
    )
}

export const AccentButtonLarge = (props) => {
    return (
        <button type={props.type} onClick={props.handleClick} className="btn btn-lg btn-accent">
            {props.label}
        </button>
    )
}
