import React, { useState } from 'react';
import { Link } from 'react-router-dom';

import MessageForm from '../components/Form/MessageForm';
import ChatDialog from '../components/List/ChatDialog';

// const MESSAGES = [
//     {
//         user: {
//             name: "name",
//         },
//         body: "本文"
//     },
//     {
//         user: {
//             name: "name",
//         },
//         body: "本文"
//     },
//     {
//         user: {
//             name: "name",
//         },
//         body: "本文"
//     },
// ]

const ChatRoom = (props) => {
    // TODO: websocket
    const [message, setMessage] = useState("");
    const [messages, setMessages] = useState([]);

    const handleChange = e => {
        setMessage(e.target.value);
    }

    const handleSubmit = e => {
        e.preventDefault();
        console.log("submit: ", message);
        // TODO: バリデーション
        setMessages([
            {
                user: {
                    name: "aaa",
                },
                body: message,
            },
            ...messages
        ])
        setMessage("");

    }

    return (
        <div className="container">
            <div class="title">
                {/* TODO: props...? */}
                <label>RoomName</label>
                <div>
                    <Link to={`/`} className="btn btn-sub">Leave Room</Link>
                </div>
            </div>

            <div class="content">
                <MessageForm
                    message={message}
                    handleChange={handleChange}
                    handleSubmit={handleSubmit}
                />
            </div>

            <div class="content">
                <ChatDialog
                    messages={messages}
                />
            </div>
        </div>
    )
}

export default ChatRoom;