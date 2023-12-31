import { useRef, useEffect } from "react";

export function useOnDraw(onDraw /*, sendCanvasUpdate*/) {
  const canvasRef = useRef(null);

  const isDrawingRef = useRef(false);

  const mouseMoveListenerRef = useRef(null);
  const mouseDownListenerRef = useRef(null);
  const mouseUpListenerRef = useRef(null);

  const prevPointRef = useRef(null);

  useEffect(() => {
    return () => {
      if (mouseMoveListenerRef.current) {
        window.removeEventListener("mousemove", mouseMoveListenerRef.current);
      }
      if (mouseUpListenerRef.current) {
        window.removeEventListener("mouseup", mouseUpListenerRef.current);
      }
    };
  }, []);

  function setCanvasRef(ref) {
    if (!ref) return;

    if (canvasRef.current) {
      canvasRef.current.removeEventListener(
        "mousedown",
        mouseDownListenerRef.current
      );
    }

    canvasRef.current = ref;
    initMouseMoveListener();
    initMouseDownListener();
    initMouseUpListener();
  }

  function initMouseMoveListener() {
    const mouseMoveListener = (e) => {
      if (isDrawingRef.current) {
        const point = computePointInCanvas(e.clientX, e.clientY);
        const ctx = canvasRef.current.getContext("2d");
        if (onDraw) onDraw(ctx, point, prevPointRef.current);
        prevPointRef.current = point;
      }
    };
    mouseMoveListenerRef.current = mouseMoveListener;
    window.addEventListener("mousemove", mouseMoveListener);
  }

  function initMouseDownListener() {
    if (!canvasRef.current) return;

    const listener = () => {
      isDrawingRef.current = true;
      // sendCanvasUpdate(); // Envoyer les modifications au socket
    };
    mouseDownListenerRef.current = listener;

    canvasRef.current.addEventListener("mousedown", listener);
  }

  function initMouseUpListener() {
    if (!canvasRef.current) return;

    const listener = () => {
      isDrawingRef.current = false;
      prevPointRef.current = null;
    };
    mouseUpListenerRef.current = listener;

    window.addEventListener("mouseup", listener);
  }

  function computePointInCanvas(clientX, clientY) {
    if (!canvasRef.current) return null;
    const canvasRect = canvasRef.current.getBoundingClientRect();

    return {
      x: clientX - canvasRect.left,
      y: clientY - canvasRect.top,
    };
  }

  return setCanvasRef;
}
