import React, { useState, useEffect } from 'react';
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
    const [conn, setConn] = useState(null);

    useEffect(() => {
        console.log("effect!")
        // connect websocket
        if(window["WebSocket"]) {
            setConn(new WebSocket("ws://" + document.location.hostname + ":8080/ws/" + props.match.params.id))
            // setConn(new WebSocket("ws://" + document.location.hostname + ":8080/ws/" + props.match.params.id))
            // conn = new WebSocket("ws://" + document.location.hostname + ":8080/ws/" + props.match.params.id)
            // if (!conn) {
            //     console.log("conn is null")
            //     return
            // }
            // なんかすぐに反映されんくて、nullにoncloseはつけれんぞって言われるので、別のeffectにする...
            // カスタムhook作ると良いのか？
            // conn.onclose = function(e) {
            //     setMessages([
            //         {
            //             user: { name: "system" },
            //             body: "connection closed",
            //         },
            //         ...messages
            //     ])
            // };
            // conn.onmessage = function(e) {
            //     console.log("onmessage")
            //     console.log(e)
            // };
        } else {
            setMessages([
                {
                    user: { name: "system" },
                    body: "Your browser does not support websockets",
                },
                ...messages
            ])
        }
        return () => conn.close();
    }, [])

    // connに変化があった時だけ呼ばれる
    useEffect(() => {
        if (!conn) {
            console.log("conn is null")
            return
        }
        conn.onclose = function (e) {
            setMessages([
                {
                    user: { name: "system" },
                    body: "connection closed",
                },
                ...messages
            ])
        };
        conn.onmessage = function (e) {
            var msg = JSON.parse(e.data)
            setMessages([
                {
                    id: msg.id,
                    user: { name: msg.sender.name },
                    body: msg.body,
                },
                ...messages
            ])
        };
    }, [conn, messages])

    const handleChange = e => {
        setMessage(e.target.value);
    }

    const handleSubmit = e => {
        e.preventDefault();
        console.log("submit: ", message);
        // TODO: バリデーション

        // ---- websocket
        if (!conn) {
            console.log("failed to get websocket connection")
            return false;
        }
        if (!message) {
            console.log("message is empty")
            return false;
        }
        conn.send(JSON.stringify({
            sender_id: 1, // TODO: 固定値きもい
            body: message
        }))
        // ----
        setMessage("");

    }

    return (
        <div className="container">
            <div className="title">
                {/* TODO: props...? */}
                <label>RoomName</label>
                <div>
                    <a href="/" className="btn btn-sub">Leave Room</a>
                </div>
            </div>

            <div className="content">
                <MessageForm
                    message={message}
                    handleChange={handleChange}
                    handleSubmit={handleSubmit}
                />
            </div>

            <div className="content">
                <ChatDialog
                    messages={messages}
                />
            </div>
        </div>
    )
}

export default ChatRoom;