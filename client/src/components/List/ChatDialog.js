import React from 'react';

const ChatDialog = (props) => {
    return (
        <ul class="dialog">
            {
                props.messages.map(msg => (
                    <Message
                        userName={msg.user.name}
                        body={msg.body}
                    />
                ))
            }
        </ul>
    )
}

const Message = (props) => {
    return (
        <li class="message">
            <div class="sender">
                <div>
                    <div class="avatar"></div>
                </div>
                <span class="name">{props.userName}</span>
            </div>
            <div class="body">{props.body}</div>
        </li>
    )
}

export default ChatDialog;