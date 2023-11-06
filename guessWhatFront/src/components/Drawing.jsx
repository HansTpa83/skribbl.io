import React, { useState, useEffect, useContext } from "react";
import WebSocketContext from "../WebSocketContext";

export default function Drawing() {
  const { ws } = useContext(WebSocketContext);
  const [drawing, setDrawing] = useState([]);
  let key = 0;

  useEffect(() => {
    if (ws) {
      ws.addEventListener("message", (event) => {
        const parsedData = JSON.parse(event.data);

        if (parsedData.type === "canvas_update") {
          const points = parsedData.data.points;
          setDrawing(points);
        } else if (parsedData.type === "game_started") {
          console.log("game_started : ", parsedData.data);
        } else if (parsedData.type === "reset_canvas") {
          setDrawing([]);
        }
      });
    }
  }, [ws]);

  useEffect(() => {
    const canvas = document.getElementById("canvas-viewer");
    const context = canvas.getContext("2d");

    // Efface le canvas
    context.clearRect(0, 0, canvas.width, canvas.height);

    // Redessine les dessins existants
    drawing.forEach((line) => {
      drawLine(line.start, line.end, context, "#000000", 5);
    });
  }, [drawing]);

  function drawLine(start, end, ctx, color, width) {
    ctx.beginPath();
    ctx.lineWidth = width;
    ctx.strokeStyle = color;
    ctx.moveTo(start.x, start.y);
    ctx.lineTo(end.x, end.y);
    ctx.stroke();

    ctx.fillStyle = "#000000";
    ctx.beginPath();
    ctx.arc(start.x, start.y, 2, 0, 2 * Math.PI);
    ctx.fill();
  }

  return (
    <canvas
      key={key}
      id="canvas-viewer"
      width={500}
      height={500}
      style={{ border: "1px solid black" }}
    />
  );
}
