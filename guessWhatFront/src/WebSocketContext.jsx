// WebSocketContext.js
import React, { createContext, useEffect, useState } from "react";

const WebSocketContext = createContext(null);

export const WebSocketProvider = ({ children }) => {
  const [ws, setWebSocket] = useState(null);
  const [userInRoom, setUserInRoom] = useState([]);
  const [messageInRoom, setMessageInRoom] = useState([]);
  const [username, setUsername] = useState("");

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8000/ws");

    socket.addEventListener("open", () => {
      console.log("WebSocket connection established.");
    });

    socket.addEventListener("message", (event) => {
      const parsedData = JSON.parse(event.data);

      switch (parsedData.type) {
        case "room_join":
          console.log(parsedData);
          if ((parsedData?.data?.roomMessages).length > 0) {
            for (const msg of parsedData.data.roomMessages) {
              setMessageInRoom((prevMessages) => [...prevMessages, msg]);
            }
          }
          break;
        case "username":
          console.log("username", parsedData);
          setUsername(parsedData.username);
          break;
        case "room_info":
          setUserInRoom(parsedData.users);
          break;
        case "message_received":
          setMessageInRoom((prevMessages) => [
            ...prevMessages,
            parsedData.data,
          ]);
          break;
        case "end_game":
          setMessageInRoom((prevMessages) => [
            ...prevMessages,
            {
              username: "SERVER",
              content: "Le gagnant est " + parsedData.data.classement,
            },
          ]);
          console.log(parsedData);
          break;
        default:
          break;
      }
    });

    socket.addEventListener("close", () => {
      console.log("WebSocket connection closed.");
    });

    setWebSocket(socket);

    return () => {
      socket.close();
    };
  }, []);

  return (
    <WebSocketContext.Provider
      value={{ ws, userInRoom, messageInRoom, username }}
    >
      {children}
    </WebSocketContext.Provider>
  );
};

export default WebSocketContext;
