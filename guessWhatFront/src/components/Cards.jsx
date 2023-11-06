import React, { useEffect, useState, useContext } from "react";
import { useParams } from "react-router-dom";
import WebSocketContext from "../WebSocketContext";

export default function Cards() {
  const { ws, userInRoom } = useContext(WebSocketContext);
  const [drawer, setDrawer] = useState([]);
  const [guessWord, setGuessWord] = useState("");
  const { roomId } = useParams();

  useEffect(() => {}, [userInRoom]);

  const startGame = () => {
    const data = {
      type: "start_game",
      data: { room: roomId },
    };
    ws.send(JSON.stringify(data));
  };

  useEffect(() => {
    if (ws) {
      ws.addEventListener("message", (event) => {
        const parsedData = JSON.parse(event.data);

        if (parsedData.type === "game_started") {
          setDrawer(parsedData.data.drawer);
          setGuessWord(parsedData?.data?.guessWord || "");
        }
      });
    }
  }, [ws]);

  return (
    <div>
      <h4>Infos : </h4>
      <div>
        <h4>Room code : {roomId}</h4>
      </div>
      <div>
        <h4>Users : </h4>
        {userInRoom.map((user, index) => {
          return (
            <div style={{ fontWeight: "bold" }} key={index}>
              {user}
            </div>
          );
        })}
      </div>
      <p>Mot à dessiner : {guessWord}</p>
      <p>Dessinateur : {drawer}</p>
      {localStorage.getItem("Admin") == "true" && (
        <div>
          <button onClick={startGame}>Start</button>
        </div>
      )}
    </div>
  );
}
