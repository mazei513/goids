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

const screenSize = 800
const maxSpeed = 10
const nBoids = 10

type Boid struct{ x, y, dx, dy float64 }

type Game struct {
	boids []*Boid
}

func (g *Game) Update() error {
	for i, b := range g.boids {
		// Rule 1
		centerX, centerY := 0.0, 0.0
		for j, b2 := range g.boids {
			if i == j {
				continue
			}
			centerX += b2.x
			centerY += b2.y
		}
		centerX /= float64(len(g.boids)) - 1
		centerY /= float64(len(g.boids)) - 1
		v1x, v1y := (centerX-b.x)/100, (centerY-b.y)/100

		// Rule 2
		v2x, v2y := 0.0, 0.0
		for j, b2 := range g.boids {
			if i == j {
				continue
			}
			dx, dy := b.x-b2.x, b.y-b2.y
			if math.Sqrt(dx*dx+dy*dy) >= 16 {
				continue
			}
			v2x -= dx
			v2y -= dy
		}

		// Update vector
		b.dx += v1x + v2x
		b.dy += v1y + v2y
		fmt.Println(v1x, v1y, v2x, v2y, b.dx, b.dy)

		// Update movement
		b.x += b.dx
		b.y += b.dy

		// Bounce off corner
		if b.x < 0 || b.x > screenSize {
			b.dx = -b.dx
			b.x += b.dx * 2
		}
		if b.y < 0 || b.y > screenSize {
			b.dy = -b.dy
			b.y += b.dy * 2
		}
	}
	fmt.Println()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i, b := range g.boids {
		opts := &ebiten.DrawImageOptions{}
		c := float32(i) / float32(len(g.boids))
		opts.ColorScale.Scale(1, c, 1, 1)
		opts.GeoM.Translate(-float64(boidImage.Bounds().Dx())/2, -float64(boidImage.Bounds().Dy())/2)
		opts.GeoM.Rotate(math.Atan2(b.dy, b.dx) + math.Pi/2)
		opts.GeoM.Translate(b.x, b.y)
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

	ebiten.SetWindowSize(screenSize, screenSize)
	ebiten.SetWindowTitle("goids")
	boids := make([]*Boid, nBoids)
	for i := range boids {
		dx := rand.Float64()
		boids[i] = &Boid{rand.Float64() * screenSize, rand.Float64() * screenSize, dx, 1 - dx}
	}
	err = ebiten.RunGame(&Game{boids: boids})
	if err != nil {
		panic(err)
	}
}
