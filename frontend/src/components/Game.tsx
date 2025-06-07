import Konva from "konva";
import React, { useRef, useEffect, useState, useCallback } from "react";
import { Stage, Layer, Circle, Line } from "react-konva";
import ScoreBoard from "./ScoreBoard";

interface Ball {
  x: number;
  y: number;
  color: string;
  potted: boolean;
  pottedBy?: number;
}

interface GameProps {
  ws: WebSocket | null;
}

const pockets = [
  { x: 0, y: 0 },
  { x: 400, y: 0 },
  { x: 800, y: 0 },
  { x: 0, y: 400 },
  { x: 400, y: 400 },
  { x: 800, y: 400 },
];

const Game: React.FC<GameProps> = ({ ws }) => {
  const stageRef = useRef<Konva.Stage | null>(null);
  const ballsRef = useRef<Ball[]>([]);
  const cueStickRef = useRef({ x: 0, y: 0, angle: 0, power: 0 });
  const [balls, setBalls] = useState<Ball[]>([]);
  const [currentTurn, setCurrentTurn] = useState(1);
  const [scores, setScores] = useState<[number, number]>([0, 0]);
  const [pottedBalls, setPottedBalls] = useState<Ball[]>([]);
  const [isDragging, setIsDragging] = useState(false);
  const [canShoot, setCanShoot] = useState(true);
  const [gameOver, setGameOver] = useState(false);
  const [winningPlayer, setWinningPlayer] = useState(0);
  const [power, setPower] = useState(0);

  const handleWsMessage = useCallback((event: MessageEvent) => {
    const data = JSON.parse(event.data);
    ballsRef.current = data.balls;
    setBalls(data.balls);
    setCurrentTurn(data.currentTurn);
    setScores(data.scores);
    setPottedBalls(data.pottedBalls);
    setCanShoot(data.canShoot);
    setGameOver(data.gameOver);
    setWinningPlayer(data.winningPlayer);
  }, []);

  useEffect(() => {
    if (ws) {
      ws.onopen = () => console.log("WebSocket connected");
      ws.onmessage = handleWsMessage;
      ws.onerror = (error) => console.log("WebSocket error: ", error);
      ws.onclose = () => console.log("WebSocket closed");
    }
  }, [ws, handleWsMessage]);

  const handleMouseMove = (e: any) => {
    const stage = stageRef.current;
    if (stage && ballsRef.current.length > 0 && canShoot && !gameOver) {
      const pointer = stage.getPointerPosition();
      if (pointer) {
        const whiteBall = ballsRef.current.find(
          (ball) => ball.color === "white"
        );
        if (whiteBall) {
          cueStickRef.current = {
            ...cueStickRef.current,
            x: pointer.x,
            y: pointer.y,
            angle: Math.atan2(pointer.y - whiteBall.y, pointer.x - whiteBall.x),
            power: isDragging
              ? Math.min(
                  100,
                  Math.max(
                    0,
                    Math.hypot(
                      pointer.x - whiteBall.x,
                      pointer.y - whiteBall.y
                    ) - 50
                  )
                )
              : cueStickRef.current.power,
          };
          setPower(cueStickRef.current.power);
        }
      }
    }
  };

  const handleMouseDown = () => {
    if (canShoot && !gameOver) {
      setIsDragging(true);
      setPower(0);
    }
  };

  const handleMouseUp = () => {
    if (isDragging && ws) {
      const { angle, power } = cueStickRef.current;
      const scaledPower = power / 10; // Scale power down
      ws.send(JSON.stringify({ type: "shoot", angle, power: scaledPower }));
    }
    setIsDragging(false);
    setPower(0);
  };

  const restartGame = () => {
    if (ws) {
      ws.send(JSON.stringify({ type: "restart" }));
    }
  };

  return (
    <div className="flex flex-col items-center mt-4 space-y-4">
      <div className="bg-green-700 rounded-lg shadow-lg p-4">
        <Stage
          width={800}
          height={400}
          ref={stageRef}
          onMouseMove={handleMouseMove}
          onMouseDown={handleMouseDown}
          onMouseUp={handleMouseUp}
          className="bg-green-600 rounded-lg shadow-inner"
        >
          <Layer>
            {pockets.map((pocket, i) => (
              <Circle
                key={i}
                x={pocket.x}
                y={pocket.y}
                radius={20}
                fill="black"
              />
            ))}
            {balls.map((ball, i) => (
              <Circle
                key={i}
                x={ball.x}
                y={ball.y}
                radius={10}
                fill={ball.color}
                visible={!ball.potted}
              />
            ))}
            {pottedBalls.map((ball, i) => (
              <Circle
                key={i}
                x={820}
                y={20 + i * 20}
                radius={10}
                fill={ball.color}
              />
            ))}
            {balls.length > 0 &&
              !balls.find((ball) => ball.color === "white")?.potted && (
                <>
                  <Line
                    points={[
                      balls.find((ball) => ball.color === "white")!.x,
                      balls.find((ball) => ball.color === "white")!.y,
                      cueStickRef.current.x,
                      cueStickRef.current.y,
                    ]}
                    stroke="#8B4513"
                    strokeWidth={6}
                    lineCap="round"
                    lineJoin="round"
                  />
                  {isDragging && (
                    <Line
                      points={[
                        balls.find((ball) => ball.color === "white")!.x,
                        balls.find((ball) => ball.color === "white")!.y,
                        balls.find((ball) => ball.color === "white")!.x +
                          cueStickRef.current.power *
                            Math.cos(cueStickRef.current.angle),
                        balls.find((ball) => ball.color === "white")!.y +
                          cueStickRef.current.power *
                            Math.sin(cueStickRef.current.angle),
                      ]}
                      stroke="red"
                      strokeWidth={2}
                      lineCap="round"
                      lineJoin="round"
                    />
                  )}
                </>
              )}
          </Layer>
        </Stage>
      </div>
      <div className="w-64 h-4 bg-gray-300 rounded overflow-hidden relative">
        <div className="h-full bg-red-500" style={{ width: `${power}%` }}></div>
        <span className="absolute inset-0 text-xs text-center text-black">
          Power: {Math.round(power)}%
        </span>
      </div>
      <ScoreBoard
        scores={scores}
        currentTurn={currentTurn}
        winningPlayer={winningPlayer}
        gameOver={gameOver}
        onRestart={restartGame}
      />
      <div className="flex space-x-8 mt-4">
        <div className="text-center">
          <div className="text-lg font-bold text-black">
            Player 1 Potted Balls:
          </div>
          <div className="flex space-x-2 mt-2">
            {pottedBalls
              .filter((ball) => ball.color !== "white" && ball.pottedBy === 1)
              .map((ball, i) => (
                <div key={i} className="relative">
                  <div className="w-8 h-8 bg-gray-800 rounded-full flex items-center justify-center">
                    <div
                      className="w-6 h-6 bg-gray-500 rounded-full"
                      style={{ backgroundColor: ball.color }}
                    ></div>
                  </div>
                </div>
              ))}
          </div>
        </div>
        <div className="text-center">
          <div className="text-lg font-bold text-black">
            Player 2 Potted Balls:
          </div>
          <div className="flex space-x-2 mt-2">
            {pottedBalls
              .filter((ball) => ball.color !== "white" && ball.pottedBy === 2)
              .map((ball, i) => (
                <div key={i} className="relative">
                  <div className="w-8 h-8 bg-gray-800 rounded-full flex items-center justify-center">
                    <div
                      className="w-6 h-6 bg-gray-500 rounded-full"
                      style={{ backgroundColor: ball.color }}
                    ></div>
                  </div>
                </div>
              ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Game;
