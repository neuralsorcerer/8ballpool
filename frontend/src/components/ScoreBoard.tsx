import React from "react";

interface ScoreBoardProps {
  scores: [number, number];
  currentTurn: number;
  winningPlayer: number;
  gameOver: boolean;
  onRestart: () => void;
}

const ScoreBoard: React.FC<ScoreBoardProps> = ({
  scores,
  currentTurn,
  winningPlayer,
  gameOver,
  onRestart,
}) => {
  const player1Class = gameOver
    ? winningPlayer === 1
      ? "bg-green-300"
      : ""
    : currentTurn === 1
    ? "bg-yellow-300"
    : "";

  const player2Class = gameOver
    ? winningPlayer === 2
      ? "bg-green-300"
      : ""
    : currentTurn === 2
    ? "bg-yellow-300"
    : "";
  return (
    <div className="flex flex-col items-center space-y-2">
      <div className={`px-4 py-1 rounded-lg ${player1Class}`}>
        <span className="font-bold text-black">Player 1: {scores[0]}</span>
      </div>
      <div className={`px-4 py-1 rounded-lg ${player2Class}`}>
        <span className="font-bold text-black">Player 2: {scores[1]}</span>
      </div>
      {!gameOver && (
        <div className="text-black font-bold">
          Player {currentTurn}&apos;s turn
        </div>
      )}
      {gameOver && (
        <div className="text-red-600 font-bold">
          Game Over! Player {winningPlayer} wins!
        </div>
      )}
      {gameOver && (
        <button
          className="mt-2 px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-700"
          onClick={onRestart}
        >
          Restart
        </button>
      )}
    </div>
  );
};

export default ScoreBoard;
