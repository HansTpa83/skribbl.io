import React, { useEffect, useState, useContext } from "react";
import WebSocketContext from "../WebSocketContext";

export default function Home() {
  const [inputs, setInputs] = useState({});
  const { ws } = useContext(WebSocketContext);

  const handleChange = (event) => {
    if (event.target.name === "room_join" && inputs.room_create) {
      delete inputs.room_create;
    }
    if (event.target.name === "room_create" && inputs.room_join) {
      delete inputs.room_join;
    }

    const name = event.target.name;
    const value = event.target.value;
    setInputs((values) => ({ ...values, [name]: value }));
  };

  const handleForm = (event) => {
    event.preventDefault();
    const keys = Object.keys(inputs);

    const message = JSON.stringify({
      type: keys.includes("room_create") ? "room_create" : "room_join",
      data: {
        username: inputs.username,
        room: inputs?.room_create || inputs?.room_join,
      },
    });

    ws.send(message);
  };

  return (
    <form
      className="Home"
      style={{ textAlign: "center" }}
      onSubmit={handleForm}
    >
      <h2>Rooms:</h2>
      <div style={{ display: "flex", justifyContent: "center" }}>
        <div style={{ marginRight: "2vw" }}>
          <h3>Join</h3>
          <input
            name="room_join"
            type="text"
            value={inputs.room_join || ""}
            onChange={handleChange}
          />
        </div>
        <div>
          <h3>Create</h3>
          <input
            name="room_create"
            type="text"
            value={inputs.room_create || ""}
            onChange={handleChange}
          />
        </div>
      </div>
      <div>
        <h2>Username</h2>
        <input
          name="username"
          type="text"
          value={inputs.username || ""}
          onChange={handleChange}
          required
        />
      </div>
      <button style={{ marginTop: "5vh", padding: "5px" }} type="submit">
        Submit
      </button>
    </form>
  );
}
