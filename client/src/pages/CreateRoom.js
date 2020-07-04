import React, { useState } from 'react';
import { Redirect, Link } from 'react-router-dom';
import CreateRoomForm from '../components/Form/CreateRoomForm';
import axios from 'axios';

const CreateRoom = (props) => {
    const [roomName, setRoomName] = useState("");
    const [isRedirect, setIsRedirect] = useState(false);

    const handleChange = e => {
        setRoomName(e.target.value)
    }

    const handleSubmit = e => {
        e.preventDefault()
        if (!roomName) {
            console.log("roomName is empty")
            return
        }
        axios.post("http://" + document.location.hostname + ":8080/rooms", {
            name: roomName,
        }, {
            withCredentials: true,
        }).then(res => {
          setRoomName("");
          setIsRedirect(true);
        })
        .catch(err => console.log(err));
    }

    return (
        <>
        {
            isRedirect
            ? <Redirect to={`/`}/>
            : <>
                <div className="container">
                    <div className="title">
                        <label>Create Room</label>
                        <div>
                            <Link to={`/`} className="btn btn-sub">Back Room List</Link>
                        </div>
                    </div>

                    <div className="content">
                        <CreateRoomForm
                            roomName={roomName}
                            handleChange={handleChange}
                            handleSubmit={handleSubmit}
                        />
                    </div>
                </div>
            </>
        }
        </>
    )
}

export default CreateRoom;