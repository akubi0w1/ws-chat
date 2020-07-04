import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import axios from 'axios';

import RoomList from '../components/List/RoomList';

const Home = () => {
    const [roomList, setRoomList] = useState([]);
    const [sysMessage, setSysMessage] = useState("");

    useEffect(() => {
        axios.get("http://" + document.location.hostname + ":8080/rooms")
            .then((res) => {
                setRoomList(res.data.rooms)
            })
            .catch((err) => setSysMessage("sorry. failed to connection api｡ﾟ(ﾟ´ω`ﾟ)ﾟ｡"))
    }, [])

    return (
        <div className="container">
            <div className="title">
                <label>Chat Rooms</label>
                <div>
                    <Link to={`/rooms/new`} className="btn btn-lg btn-accent">Create Room</Link>
                </div>
            </div>
            {
                sysMessage
                    ? <div className="content"><b>{sysMessage}</b></div>
                : <></>
            }
            <div className="content">
                <RoomList
                    roomList={roomList}
                />
            </div>
        </div>
    )
};

export default Home;