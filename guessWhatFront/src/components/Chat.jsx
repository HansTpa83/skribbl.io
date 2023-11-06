import React, { useState, useEffect, useContext } from "react";
import WebSocketContext from "../WebSocketContext";
import { useParams } from "react-router-dom";

const Chat = () => {
  const { ws, messageInRoom } = useContext(WebSocketContext);
  const [message, setMessage] = useState("");
  let { roomId } = useParams();
  console.log("messageInRoom : ", messageInRoom);
  useEffect(() => {}, [messageInRoom]);

  const sendMessage = () => {
    ws.send(
      JSON.stringify({
        type: "chat",
        data: {
          content: message,
          room: roomId,
        },
      })
    );
  };

  return (
    <div style={{ width: "33vw", textAlign: "center" }}>
      <div style={{ height: "70vh", overflow: "auto" }}>
        {messageInRoom.map((msg, index) => {
          console.log("msg : ", msg);
          return (
            <div key={index}>
              {msg.username} : {msg.content}
            </div>
          );
        })}
      </div>
      <input
        type="text"
        value={message}
        onChange={(e) => setMessage(e.target.value)}
      />
      <button onClick={sendMessage}>Envoyer</button>
    </div>
  );
};

export default Chat;
