package engine

import "math"

type Ball struct {
    X        float64 `json:"x"`
    Y        float64 `json:"y"`
    Velocity float64 `json:"velocity"`
    Angle    float64 `json:"angle"`
    Color    string  `json:"color"`
    Potted   bool    `json:"potted"`
}

const ballRadius = 10
const friction = 0.98
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
    collisionAngle := math.Atan2(dy, dx)

    speed1 := b1.Velocity
    speed2 := b2.Velocity
    direction1 := b1.Angle
    direction2 := b2.Angle

    velocityX1 := speed1 * math.Cos(direction1-collisionAngle)
    velocityY1 := speed1 * math.Sin(direction1-collisionAngle)
    velocityX2 := speed2 * math.Cos(direction2-collisionAngle)
    velocityY2 := speed2 * math.Sin(direction2-collisionAngle)

    finalVelocityX1 := (velocityX1*(ballRadius-1) + (2 * ballRadius * velocityX2)) / (ballRadius + ballRadius)
    finalVelocityY1 := velocityY1
    finalVelocityX2 := (velocityX2*(ballRadius-1) + (2 * ballRadius * velocityX1)) / (ballRadius + ballRadius)
    finalVelocityY2 := velocityY2

    b1.Velocity = math.Sqrt(finalVelocityX1*finalVelocityX1 + finalVelocityY1*finalVelocityY1)
    b2.Velocity = math.Sqrt(finalVelocityX2*finalVelocityX2 + finalVelocityY2*finalVelocityY2)
    b1.Angle = math.Atan2(finalVelocityY1, finalVelocityX1) + collisionAngle
    b2.Angle = math.Atan2(finalVelocityY2, finalVelocityX2) + collisionAngle

    overlap := 2*ballRadius - distance
    b1.X += math.Cos(collisionAngle) * overlap / 2
    b1.Y += math.Sin(collisionAngle) * overlap / 2
    b2.X -= math.Cos(collisionAngle) * overlap / 2
    b2.Y -= math.Sin(collisionAngle) * overlap / 2
}
