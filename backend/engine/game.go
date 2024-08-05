package engine

type Game struct {
	Balls         []*Ball `json:"balls"`
	CurrentTurn   int     `json:"currentTurn"`
	Scores        [2]int  `json:"scores"`
	PottedBalls   []*Ball `json:"pottedBalls"`
	CanShoot      bool    `json:"canShoot"`
	GameOver      bool    `json:"gameOver"`
	WinningPlayer int     `json:"winningPlayer"`
}

func NewGame() *Game {
	colors := []string{"yellow", "blue", "red", "purple", "orange", "green", "brown", "black",
		"yellow", "blue", "red", "purple", "orange", "green", "brown"}
	balls := make([]*Ball, 0, len(colors)+1) // +1 for the cue ball
	positions := []struct{ x, y float64 }{
		{400, 200}, {430, 185}, {430, 215},
		{460, 170}, {460, 200}, {460, 230},
		{490, 155}, {490, 185}, {490, 215}, {490, 245},
		{520, 140}, {520, 170}, {520, 200}, {520, 230}, {520, 260},
	}
	for i, color := range colors {
		balls = append(balls, &Ball{
			X:        positions[i%len(positions)].x,
			Y:        positions[i%len(positions)].y,
			Velocity: 0,
			Angle:    0,
			Color:    color,
			Potted:   false,
		})
	}

	balls = append(balls, &Ball{
		X:        200,
		Y:        200,
		Velocity: 0,
		Angle:    0,
		Color:    "white",
		Potted:   false,
	})
	return &Game{Balls: balls, CurrentTurn: 1, Scores: [2]int{0, 0}, PottedBalls: []*Ball{}, CanShoot: true, GameOver: false, WinningPlayer: 0}
}

func (g *Game) Update() {
	if g.GameOver {
		return
	}

	moving := false
	for _, ball := range g.Balls {
		if !ball.Potted {
			ball.UpdatePosition()
			if ball.Velocity != 0 {
				moving = true
			}
		} else if !contains(g.PottedBalls, ball) {
			g.PottedBalls = append(g.PottedBalls, ball)
			if ball.Color == "black" {
				if len(g.PottedBalls) < 15 {
					g.GameOver = true
					g.WinningPlayer = 3 - g.CurrentTurn
				} else {
					g.GameOver = true
					g.WinningPlayer = g.CurrentTurn
				}
			} else if ball.Color != "white" {
				g.Scores[g.CurrentTurn-1]++
			} else {
				ball.Potted = false
				ball.X = 100
				ball.Y = 200
				ball.Velocity = 0
				ball.Angle = 0
			}
		}
	}

	for i := 0; i < len(g.Balls); i++ {
		for j := i + 1; j < len(g.Balls); j++ {
			if CheckCollision(g.Balls[i], g.Balls[j]) {
				HandleCollision(g.Balls[i], g.Balls[j])
			}
		}
	}

	g.CanShoot = !moving

	if !moving && !g.GameOver {
		g.CurrentTurn = 3 - g.CurrentTurn
	}
}

func (g *Game) ShootBall(angle, power float64) {
	if !g.CanShoot || g.GameOver {
		return
	}
	for _, ball := range g.Balls {
		if ball.Color == "white" && !ball.Potted {
			ball.Velocity = power
			ball.Angle = angle
			break
		}
	}
	g.CanShoot = false
}

func contains(s []*Ball, e *Ball) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
