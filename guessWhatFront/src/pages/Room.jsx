import React, { useEffect, useState, useContext } from "react";
import Chat from "../components/Chat";
import Cards from "../components/Cards";
import Canvas from "../components/Canvas";
import Drawing from "../components/drawing";
import WebSocketContext from "../WebSocketContext";

export default function Room() {
  const { ws, username } = useContext(WebSocketContext);
  const [drawer, setDrawer] = useState("");

  console.log("DRAWER : ", drawer, " | ", "USERNAME : ", username);
  useEffect(() => {
    if (ws) {
      ws.addEventListener("message", (event) => {
        const parsedData = JSON.parse(event.data);

        if (parsedData.type === "game_started") {
          setDrawer(parsedData.data.drawer);
        } else if (parsedData.type === "room_created") {
        }
      });
    }
  }, [ws]);

  return (
    <div
      className="Room"
      style={{
        display: "flex",
        justifyContent: "space-evenly",
        userSelect: "none",
      }}
    >
      <Cards />
      {username == drawer ? <Canvas /> : <Drawing />}
      <Chat />
    </div>
  );
}
