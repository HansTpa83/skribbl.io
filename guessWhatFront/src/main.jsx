import "./index.css";
import React from "react";
import App from "./App.jsx";
import ReactDOM from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import { WebSocketProvider } from "./WebSocketContext";

ReactDOM.createRoot(document.getElementById("root")).render(
  <WebSocketProvider>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </WebSocketProvider>
);
