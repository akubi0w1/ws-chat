import React from 'react';

const MessageForm = (props) => {
    return (
        <form className="send-form" onSubmit={props.handleSubmit} method="POST">
            <input type="text" value={props.message} placeholder="Type message..." onChange={props.handleChange} />
            <button type="submit" class="btn btn-lg btn-accent">POST</button>
        </form>
    )
}

export default MessageForm;