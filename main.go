package main

import (
	"image"
	"math"
	"math/rand/v2"
	"os"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	_ "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenSize = 4000
	screenHalf = screenSize / 2.0
	maxSpeed   = 50
	nBoids     = 6500
	vision     = 60
	repulse    = 30.0
	attract    = 200.0
	follow     = 8.0
	nPar       = 16
)

type Boid struct{ x, y, dx, dy float64 }
type Game struct{ boids []*Boid }

func mag(x, y float64) float64 { return math.Sqrt(x*x + y*y) }

func (g *Game) Update() error {
	// Boid update
	dxs := make([]float64, nBoids)
	dys := make([]float64, nBoids)
	wg := &sync.WaitGroup{}
	wg.Add(nPar)
	for i := range nPar {
		nChunk := nBoids / nPar
		start := nChunk * i
		end := nChunk * (i + 1)
		if i == nPar-1 {
			end = nBoids
		}
		go func(boids []*Boid, outX, outY []float64, offset int) {
			for i, b := range boids {
				v1x, v1y := 0.0, 0.0
				v2x, v2y := 0.0, 0.0
				v3x, v3y := 0.0, 0.0
				n := 0.0
				for j, b2 := range g.boids {
					dx, dy := b2.x-b.x, b2.y-b.y
					dMag := mag(dx, dy)
					if offset+i == j || dMag > vision {
						continue
					}

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
				outX[i] = v1x + v2x/70 + v3x
				outY[i] = v1y + v2y/70 + v3y
			}
			wg.Done()
		}(g.boids[start:end], dxs[start:end], dys[start:end], start)
	}

	wg.Wait()

	for i, b := range g.boids {
		// Update vector
		b.dx += dxs[i]
		b.dy += dys[i]

		// Limit max speed
		dMag := mag(b.dx, b.dy)
		if dMag > maxSpeed {
			b.dx = (b.dx / dMag) * maxSpeed
			b.dy = (b.dy / dMag) * maxSpeed
		} else if dMag < 5 {
			b.dx = b.dx * 3
			b.dy = b.dy * 3
		}

		// Update movement
		b.x += b.dx
		b.y += b.dy

		// Bounce off corner
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
	opts := &ebiten.DrawImageOptions{}
	offset := -float64(boidImage.Bounds().Dx()) / 2
	const halfPi = math.Pi / 2
	for i, b := range g.boids {
		c := float32(i) / nBoids
		opts.ColorScale.Scale(1, c, 1, 1)
		opts.GeoM.Translate(offset, offset)
		opts.GeoM.Rotate(math.Atan2(b.dy, b.dx) + halfPi)
		opts.GeoM.Translate(b.x+screenHalf, b.y+screenHalf)
		screen.DrawImage(boidImage, opts)
		opts.ColorScale.Reset()
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

	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetWindowTitle("goids")
	boids := make([]*Boid, nBoids)
	for i := range boids {
		dx := rand.Float64() - 0.5
		dy := rand.Float64() - 0.5
		dMag := mag(dx, dy)
		dx = (dx / dMag) * maxSpeed
		dy = (dy / dMag) * maxSpeed
		boids[i] = &Boid{(rand.Float64() - 0.5) * screenSize, (rand.Float64() - 0.5) * screenSize, dx, dy}
	}
	err = ebiten.RunGame(&Game{boids})
	if err != nil {
		panic(err)
	}
}
