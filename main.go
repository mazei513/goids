package main

import (
	"image"
	_ "image/png"
	"math"
	"math/rand/v2"
	"os"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	windowSize    = 1000
	screenSize    = 6000
	maxSpeed      = 50
	minSpeed      = 5
	nBoids        = 6500
	vision        = 60
	repulseVision = 20.0
	repulseDampen = 30.0
	attract       = 100.0
	follow        = 16.0
	nPar          = 16
)

type Vec struct{ x, y float64 }
type Boid struct{ x, y, dx, dy float64 }
type Game struct {
	vs     []Vec
	boids  []*Boid
	colors []ebiten.ColorScale
}

func mag(x, y float64) float64 { return math.Sqrt(x*x + y*y) }

func (g *Game) Update() error {
	// Boid update
	wg := &sync.WaitGroup{}
	wg.Add(nPar)
	for i := range nPar {
		nChunk := nBoids / nPar
		start := nChunk * i
		end := nChunk * (i + 1)
		if i == nPar-1 {
			end = nBoids
		}
		go func() { Calc(g.boids, g.boids[start:end], g.vs[start:end], start); wg.Done() }()
	}

	wg.Wait()

	for i, b := range g.boids {
		// Update vector
		b.dx += g.vs[i].x
		b.dy += g.vs[i].y

		// Limit max speed
		dMag := mag(b.dx, b.dy)
		if dMag > maxSpeed {
			b.dx = (b.dx / dMag) * maxSpeed
			b.dy = (b.dy / dMag) * maxSpeed
		} else if dMag < minSpeed {
			b.dx = b.dx * minSpeed / dMag * (1 + rand.Float64()*2)
			b.dy = b.dy * minSpeed / dMag * (1 + rand.Float64()*2)
		}

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

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	offset := -float64(boidImage.Bounds().Dx()) / 2
	const halfPi = math.Pi / 2
	for i, b := range g.boids {
		opts.ColorScale = g.colors[i]
		opts.GeoM.Translate(offset, offset)
		opts.GeoM.Rotate(math.Atan2(b.dy, b.dx) + halfPi)
		opts.GeoM.Translate(b.x, b.y)
		screen.DrawImage(boidImage, opts)
		opts.GeoM.Reset()
	}
	// TODO add marker when TPS <59.9
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

	ebiten.SetWindowSize(windowSize, windowSize)
	ebiten.SetWindowTitle("goids")
	boids := make([]*Boid, nBoids)
	colors := make([]ebiten.ColorScale, nBoids)
	for i := range boids {
		dx := rand.Float64() - 0.5
		dy := rand.Float64() - 0.5
		dMag := mag(dx, dy)
		dx = (dx / dMag) * minSpeed
		dy = (dy / dMag) * minSpeed
		boids[i] = &Boid{(rand.Float64()) * screenSize, (rand.Float64()) * screenSize, dx, dy}
		colors[i].Scale(rand.Float32(), rand.Float32(), rand.Float32(), 1)
	}
	err = ebiten.RunGame(&Game{make([]Vec, nBoids), boids, colors})
	if err != nil {
		panic(err)
	}
}

func Calc(all, boids []*Boid, out []Vec, offset int) {
	for i, b := range boids {
		v1x, v1y := 0.0, 0.0
		v2x, v2y := 0.0, 0.0
		v3x, v3y := 0.0, 0.0
		n := 0.0
		for j, b2 := range all {
			dx, dy := b2.x-b.x, b2.y-b.y
			dMag := mag(dx, dy)
			if offset+i == j || dMag > vision {
				continue
			}

			// Rule 1
			v1x += b2.x
			v1y += b2.y

			// Rule 2
			if dMag < repulseVision {
				v2x -= dx
				v2y -= dy
			}

			// Rule 3
			v3x += b2.dx
			v3y += b2.dy

			n += 1
		}

		if n == 0 {
			v1x, v1y, v3x, v3y = 0, 0, 0, 0
		} else {
			// Rule 1
			v1x /= n
			v1y /= n
			v1x = (v1x - b.x) / attract
			v1y = (v1y - b.y) / attract

			// Rule 3
			v3x /= n
			v3y /= n
			v3x = (v3x - b.dx) / follow
			v3y = (v3y - b.dy) / follow
		}

		// Store vector
		out[i].x = v1x + v2x/repulseDampen + v3x
		out[i].y = v1y + v2y/repulseDampen + v3y
	}
}
