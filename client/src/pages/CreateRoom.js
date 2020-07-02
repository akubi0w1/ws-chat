import React, { useState } from 'react';
import { Redirect, Link } from 'react-router-dom';
import CreateRoomForm from '../components/Form/CreateRoomForm';

const CreateRoom = (props) => {
    const [roomName, setRoomName] = useState("");
    const [isRedirect, setIsRedirect] = useState(false);

    const handleChange = e => {
        setRoomName(e.target.value)
    }
    const handleSubmit = e => {
        e.preventDefault()
        console.log(roomName)
        // TODO: ヴァリデーション
        setRoomName("");
        setIsRedirect(true);
    }

    return (
        <>
        {
            isRedirect
            ? <Redirect to={`/`}/>
            : <>
                <div className="container">
                    <div class="title">
                        <label>Create Room</label>
                        <div>
                            <Link to={`/`} className="btn btn-sub">Back Room List</Link>
                        </div>
                    </div>

                    <div class="content">
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