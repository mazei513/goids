package main

import (
	"fmt"
	"image"
	"math"
	"math/rand/v2"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const screenSize = 4000
const maxSpeed = 20
const nBoids = 2000
const repulse = 10.0
const attract = 400.0
const follow = 8.0

type Boid struct{ x, y, dx, dy float64 }

type Game struct {
	boids []*Boid
}

func mag(dx, dy float64) float64 {
	return math.Sqrt(dx*dx + dy*dy)
}

func (g *Game) Update() error {
	// Boid update
	for i, b := range g.boids {
		v1x, v1y := 0.0, 0.0
		v2x, v2y := 0.0, 0.0
		v3x, v3y := 0.0, 0.0
		n := 0.0
		for j, b2 := range g.boids {
			dx, dy := b2.x-b.x, b2.y-b.y
			if i == j {
				continue
			}
			dMag := mag(dx, dy)
			if dMag > 200 {
				continue
			}
			n += 1
			// Rule 1
			v1x += b2.x
			v1y += b2.y

			// Rule 2
			if dMag < repulse {
				v2x -= dx
				v2y -= dy
			}

			// Rule 3
			v3x += b2.dx
			v3y += b2.dy
		}

		if n == 0 {
			v1x, v1y, v3x, v3y = 0, 0, 0, 0
		} else {
			// Rule 1
			v1x /= n
			v1y /= n
			v1x, v1y = (v1x-b.x)/attract, (v1y-b.y)/attract

			// Rule 3
			v3x /= n
			v3y /= n
			v3x = (v3x - b.dx) / follow
			v3y = (v3y - b.dy) / follow
		}

		// Update vector
		b.dx += v1x + v2x + v3x
		b.dy += v1y + v2y + v3y

		// Limit max speed
		dMag := mag(b.dx, b.dy)
		if dMag > maxSpeed {
			b.dx = (b.dx / dMag) * maxSpeed
			b.dy = (b.dy / dMag) * maxSpeed
		}

		// Update movement
		b.x += b.dx
		b.y += b.dy

		// Bounce off corner
		screenHalf := screenSize / 2.0
		if b.x < -screenHalf || b.x > screenHalf {
			b.dx = -b.dx
			b.x += b.dx * 2
		}
		if b.y < -screenHalf || b.y > screenHalf {
			b.dy = -b.dy
			b.y += b.dy * 2
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i, b := range g.boids {
		opts := &ebiten.DrawImageOptions{}
		c := float32(i) / float32(len(g.boids))
		opts.ColorScale.Scale(1, c, 1, 1)
		opts.GeoM.Translate(-float64(boidImage.Bounds().Dx())/2, -float64(boidImage.Bounds().Dy())/2)
		opts.GeoM.Rotate(math.Atan2(b.dy, b.dx) + math.Pi/2)
		opts.GeoM.Translate(b.x+screenSize/2, b.y+screenSize/2)
		screen.DrawImage(boidImage, opts)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f, TPS: %f", ebiten.ActualFPS(), ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenSize, screenSize
}

var boidImage *ebiten.Image

func main() {
	r, err := os.Open("boid.png")
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(r)
	if err != nil {
		panic(err)
	}
	boidImage = ebiten.NewImageFromImage(img)

	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetWindowTitle("goids")
	boids := make([]*Boid, nBoids)
	for i := range boids {
		dx := rand.Float64()
		dy := rand.Float64()
		dMag := mag(dx, dy)
		dx = (dx / dMag) * maxSpeed
		dy = (dy / dMag) * maxSpeed
		boids[i] = &Boid{(rand.Float64() - 0.5) * screenSize, (rand.Float64() - 0.5) * screenSize, dx, dy}
	}
	err = ebiten.RunGame(&Game{boids: boids})
	if err != nil {
		panic(err)
	}
}
