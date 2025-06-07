# 8 Ball Pool Game

This is a simple 8 Ball Pool game built using Go for the backend and React with Konva for the frontend. The game supports multiplayer mode on the same machine and features a realistic cue stick and ball physics.

## Features

- Realistic ball and cue stick physics
- Multiplayer support on the same machine
- Real-time updates using WebSocket
- Visual power meter and highlighted scoreboard (current turn and winner)

## Prerequisites

- Go (1.16 or higher)
- Node.js (14.x or higher)

## Getting Started

### Backend

1. **Navigate to the backend directory**:
   ```sh
   cd backend
   ```
2. **Install Go modules**:
   ```sh
   go mod tidy
   ```
3. **Run the backend server**:
   ```sh
   go run main.go
   ```

The backend server should now be running on `http://localhost:8080`.

- `main.go`: The main entry point for the Go server, handling WebSocket connections and game state updates.
- `engine/`: Contains the game logic and physics calculations.

### Frontend

1. **Navigate to the frontend directory**:
   ```sh
   cd frontend
   ```
2. **Install npm dependencies**:
   ```sh
   npm install
   ```
3. **Run the frontend development server**:
   ```sh
   npm run dev
   ```

The frontend application should now be running on `http://localhost:3000`.

- `src/App.tsx`: Main React component.
- `src/Game.tsx`: React component that handles the game interface using Konva for rendering the game objects.

## Gameplay

- Use the mouse to aim the cue stick.
- Click and drag to set the power and release to shoot.
- Players take turns to pot their respective balls.
- Pot the black ball to win, but only after potting all your other balls.

## License

This project is licensed under the MIT License.
