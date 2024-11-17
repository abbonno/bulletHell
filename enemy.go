package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Enemy struct {
	X     float64     // Posición X
	Y     float64     // Posición Y
	Speed float64     // Velocidad de seguimiento
	Color color.Color // Color del enemigo
}

func (e *Enemy) update(iconX, iconY float64) {
	// Calcula la dirección hacia el icono
	dx := iconX - e.X
	dy := iconY - e.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance > 0 {
		e.X += e.Speed * (dx / distance) // Movimiento normalizado
		e.Y += e.Speed * (dy / distance)
	}
}

func (e *Enemy) draw(screen *ebiten.Image) {
	// Dibuja al enemigo como un círculo
	ebitenutil.DrawCircle(screen, e.X, e.Y, bulletRadius, e.Color)
}
