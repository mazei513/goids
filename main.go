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
	for _, b := range g.boids {
		b.x += b.dx * maxSpeed
		if b.x < 0 || b.x > screenSize {
			b.dx = -b.dx
			b.x += b.dx * maxSpeed * 2
		}
		b.y += b.dy * maxSpeed
		if b.y < 0 || b.y > screenSize {
			b.dy = -b.dy
			b.y += b.dy * maxSpeed * 2
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, b := range g.boids {
		opts := &ebiten.DrawImageOptions{}
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
