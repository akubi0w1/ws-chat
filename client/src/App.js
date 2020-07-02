import React from 'react';
import './assets/css/main.css';

import { BrowserRouter, Route } from 'react-router-dom';

import Home from './pages/Home';
import CreateRoom from './pages/CreateRoom';
import ChatRoom from './pages/ChatRoom';

function App() {
  return (
    <div className="container-wrapper">
        <BrowserRouter>
            <Route exact path="/" component={Home}/>
            <Route exact path="/rooms/new" component={CreateRoom}/>
            <Route exact path="/rooms/:id/dialog" component={ChatRoom} />
        </BrowserRouter>
    </div>
  );
}

export default App;
