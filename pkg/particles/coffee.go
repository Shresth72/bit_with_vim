package particles

import (
	"math"
	"math/rand"
	"time"
)

type Coffee struct {
	ParticleSystem
}

func reset(p *Particle, params *ParticleParams) {
	p.Lifetime = int64(math.Floor(float64(params.MaxLife) * rand.Float64()))
	p.Speed = params.MaxSpeed * rand.Float64()

	maxX := math.Floor(float64(params.X) / 2)
	x := math.Max(-maxX, math.Min(rand.NormFloat64() * params.XScale, maxX))
	p.X = x + maxX

	p.Y = 0
}

func nextPos(particle *Particle, deltaMS int64) {
	particle.Lifetime -= deltaMS
	if particle.Lifetime <= 0 {
		return
	}

  percent := (float64(deltaMS) / 1000.0)
	particle.Y += particle.Speed * percent
}

func NewCoffee(width, height int, scale float64) Coffee {
	// assert.Assert(width % 2 == 1, "width of particle system must be odd")
  startTime := time.Now().UnixMilli()
  ascii := func(row, col int, counts [][]int) string {
    count := counts[row][col]
    if count == 0 {
      return " "
    }
    if count == 1 {
      return "."
    }

    direction := row + 
      int(((time.Now().UnixMilli() - startTime) / 1000) % 2)

    if direction % 2 == 0 {
      return "}"
    }
    return "{"
  }

	return Coffee{
		ParticleSystem: NewParticleSystem(
			ParticleParams{
				MaxLife:       6000,
				MaxSpeed:      1.5,
				ParticleCount: 700,

				reset:        reset,
				ascii:        ascii,
				nextPosition: nextPos,

        XScale: scale,
        X: width,
        Y: height,
			},
		),
	}
}
