import React from 'react';
import { Link } from 'react-router-dom';

import RoomList from '../components/List/RoomList';

const ROOM_LIST = [
    {
        id: "1",
        label: "room name1",
    },
    {
        id: "2",
        label: "room name2",
    },
    {
        id: "3",
        label: "room name3",
    },
]

const Home = () => (
    <div className="container">
        <div className="title">
            <label>Chat Rooms</label>
            <div>
                <Link to={`/rooms/new`} className="btn btn-lg btn-accent">Create Room</Link>
            </div>
        </div>
        <div className="content">
            <RoomList
                roomList={ROOM_LIST}
            />
        </div>
    </div>
);

export default Home;