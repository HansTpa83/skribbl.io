import React, { useState } from "react";
import "./css/nav.css";
import { BsFillDoorOpenFill, BsFillDoorClosedFill } from "react-icons/bs";
import { MdDraw } from "react-icons/md";

export default function Navbar() {
  const [isShown, setIsShown] = useState(false);
  return (
    <nav>
      <div>
        <MdDraw size={60} color="blue" />
      </div>
      <h1>Guess what</h1>
      <div
        style={{ width: "5vw" }}
        onMouseEnter={() => setIsShown(true)}
        onMouseLeave={() => setIsShown(false)}
      >
        {isShown ? (
          <BsFillDoorOpenFill size={60} color="red" />
        ) : (
          <BsFillDoorClosedFill size={60} color="red" />
        )}
      </div>
    </nav>
  );
}
