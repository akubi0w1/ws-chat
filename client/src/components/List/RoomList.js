import React from 'react';
import { Link } from 'react-router-dom';

const RoomList = (props) => {
    return (
        <ul className="list">
            {
                props.roomList.map(room => (
                    <ListItem key={room.id} id={room.id} name={room.name}/>
                ))
            }
        </ul>
    )
};

const ListItem = (props) => {
    return (
        <li className="list-item">
            <Link to={`/rooms/${props.id}/dialog`} className="list-link">{props.name}</Link>
        </li>
    )
};

export default RoomList;