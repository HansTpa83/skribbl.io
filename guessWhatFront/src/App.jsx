import WebSocketContext from "./WebSocketContext";
import React, { useEffect, useContext } from "react";
import { BrowserRouter, Routes, Route, useNavigate } from "react-router-dom";

// components
import Navbar from "./components/Navbar";
import Room from "./pages/Room";
import Home from "./pages/Home";

const App = () => {
  const { ws } = useContext(WebSocketContext);
  const navigate = useNavigate();
  useEffect(() => {
    if (ws) {
      ws.addEventListener("message", (event) => {
        const dataParsed = JSON.parse(event.data);
        switch (dataParsed.type) {
          case "room_created":
            navigate(`/room/${dataParsed.data.roomId}`);
            break;
          case "room_join":
            navigate(`/room/${dataParsed.data.roomId}`);
            break;
          default:
            break;
        }
      });
    }
  }, [ws]);

  return (
    <div className="App">
      <Navbar />
      <Routes>
        <Route path="/room/:roomId" element={<Room />} />
        <Route path="/" element={<Home />} />
        {/* <Route path="/room" element={<Room socket={ws} />} /> */}
      </Routes>
      {/* <Chat /> */}
    </div>
  );
};

export default App;
