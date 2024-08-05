import React, { useEffect, useState } from "react";
import Game from "./components/Game";

const App: React.FC = () => {
  const [ws, setWs] = useState<WebSocket | null>(null);

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8080/ws");
    socket.onopen = () => console.log("WebSocket connection established");
    socket.onclose = (event) =>
      console.log("WebSocket connection closed: ", event);
    socket.onerror = (error) => console.log("WebSocket error: ", error);

    setWs(socket);

    return () => {
      socket.close();
    };
  }, []);

  return <Game ws={ws} />;
};

export default App;
