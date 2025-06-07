package engine

import "math"

type Ball struct {
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Velocity float64 `json:"velocity"`
	Angle    float64 `json:"angle"`
	Color    string  `json:"color"`
	Potted   bool    `json:"potted"`
	PottedBy int     `json:"pottedBy"`
}

const ballRadius = 10
const friction = 0.98
const restitution = 0.9
const tableWidth = 800
const tableHeight = 400

var pockets = [][2]float64{
	{0, 0}, {tableWidth / 2, 0}, {tableWidth, 0},
	{0, tableHeight}, {tableWidth / 2, tableHeight}, {tableWidth, tableHeight},
}

const pocketRadius = 20

func (b *Ball) UpdatePosition() {
	if b.Potted {
		return
	}

	b.X += b.Velocity * math.Cos(b.Angle)
	b.Y += b.Velocity * math.Sin(b.Angle)
	b.Velocity *= friction

	if b.X < ballRadius {
		b.X = ballRadius
		b.Angle = math.Pi - b.Angle
	} else if b.X > tableWidth-ballRadius {
		b.X = tableWidth - ballRadius
		b.Angle = math.Pi - b.Angle
	}

	if b.Y < ballRadius {
		b.Y = ballRadius
		b.Angle = -b.Angle
	} else if b.Y > tableHeight-ballRadius {
		b.Y = tableHeight - ballRadius
		b.Angle = -b.Angle
	}

	if b.Velocity < 0.01 {
		b.Velocity = 0
	}

	for _, pocket := range pockets {
		if math.Sqrt(math.Pow(b.X-pocket[0], 2)+math.Pow(b.Y-pocket[1], 2)) < pocketRadius {
			b.Potted = true
			b.Velocity = 0
			break
		}
	}
}

func CheckCollision(b1, b2 *Ball) bool {
	dx := b1.X - b2.X
	dy := b1.Y - b2.Y
	distance := math.Sqrt(dx*dx + dy*dy)
	return distance < 2*ballRadius
}

func HandleCollision(b1, b2 *Ball) {
	dx := b1.X - b2.X
	dy := b1.Y - b2.Y
	distance := math.Sqrt(dx*dx + dy*dy)
	if distance == 0 {
		distance = 0.01
	}
	nx := dx / distance
	ny := dy / distance

	v1x := b1.Velocity * math.Cos(b1.Angle)
	v1y := b1.Velocity * math.Sin(b1.Angle)
	v2x := b2.Velocity * math.Cos(b2.Angle)
	v2y := b2.Velocity * math.Sin(b2.Angle)

	p := 2 * ((v1x-v2x)*nx + (v1y-v2y)*ny) / 2

	v1x -= p * nx
	v1y -= p * ny
	v2x += p * nx
	v2y += p * ny

	v1x *= restitution
	v1y *= restitution
	v2x *= restitution
	v2y *= restitution

	b1.Velocity = math.Hypot(v1x, v1y)
	b1.Angle = math.Atan2(v1y, v1x)
	b2.Velocity = math.Hypot(v2x, v2y)
	b2.Angle = math.Atan2(v2y, v2x)

	overlap := 2*ballRadius - distance
	b1.X += nx * overlap / 2
	b1.Y += ny * overlap / 2
	b2.X -= nx * overlap / 2
	b2.Y -= ny * overlap / 2
}