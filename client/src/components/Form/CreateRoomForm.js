import React from 'react';

const CreateRoomForm = props => {
    return (
        <form className="form" method="POST" onSubmit={props.handleSubmit}>
            <div>
                <label>room name</label>
                <input type="text" value={props.roomName} onChange={props.handleChange} />
            </div>
            <div>
                <button type="submit" className="btn btn-accent btn-lg">Create Room</button>
            </div>
        </form>
    )
}

export default CreateRoomForm;