import React, { useState, useContext, useEffect, useRef } from "react";
import { useParams } from "react-router-dom";
import { useOnDraw } from "./Hooks";
import WebSocketContext from "../WebSocketContext";

export default function Canvas() {
  const { ws } = useContext(WebSocketContext);
  const [resetKey, setResetKey] = useState(0);
  const canvasRef = useOnDraw(onDraw);
  const { roomId } = useParams();
  useEffect(() => {}, [ws]);

  const handleResetClick = () => {
    setResetKey((prevKey) => prevKey + 1);
    ws.send(
      JSON.stringify({
        type: "reset_canvas",
        data: {
          room: roomId,
        },
      })
    );
  };

  function onDraw(ctx, point, prevPoint) {
    drawLine(prevPoint, point, ctx, "#000000", 5);

    if (ws) {
      const data = JSON.stringify({
        type: "canvas_update",
        data: {
          points: {
            start: prevPoint,
            end: point,
          },
          room: roomId,
        },
      });
      ws.send(data);
    }
  }

  function drawLine(start, end, ctx, color, width) {
    start = start ?? end;
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
    <div key={resetKey}>
      <canvas width={500} height={500} ref={canvasRef} />
      <button
        onClick={handleResetClick}
        style={{
          border: "1px solid black",
          display: "block",
          width: 500,
          height: 10,
        }}
      >
        Reset
      </button>
    </div>
  );
}
