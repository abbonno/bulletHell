package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Bullet struct {
	X      float64     // Posición X
	Y      float64     // Posición Y
	Angle  float64     // Dirección en radianes
	Speed  float64     // Velocidad
	Active bool        // Si la bala está activa
	Color  color.Color // Color de la bala
}

func drawBullets(screen *ebiten.Image, bullets []Bullet) {
	for _, bullet := range bullets {
		if !bullet.Active {
			continue
		}
		ebitenutil.DrawCircle(screen, bullet.X, bullet.Y, bulletRadius, bullet.Color)
	}
}

func updateBullets(bullets *[]Bullet, iconX, iconY float64) {
	iconCenterX := iconX + iconSize/2
	iconCenterY := iconY + iconSize/2

	for i := range *bullets {
		bullet := &(*bullets)[i]

		// Solo actualiza balas activas
		if !bullet.Active {
			continue
		}

		// Actualiza la posición
		bullet.X += math.Cos(bullet.Angle) * bullet.Speed
		bullet.Y += math.Sin(bullet.Angle) * bullet.Speed

		// Detecta colisión con el icono
		dx := bullet.X - iconCenterX
		dy := bullet.Y - iconCenterY
		distance := math.Sqrt(dx*dx + dy*dy)
		if distance <= bulletRadius+iconSize/2 {
			bullet.Color = color.RGBA{255, 0, 0, 255} // Cambia a rojo
			//bullet.Active = false                     // Desactiva la bala tras la colisión
		}

		// Desactiva balas fuera de la pantalla
		if bullet.X < 0 || bullet.X > screenWidth || bullet.Y < 0 || bullet.Y > screenHeight {
			bullet.Active = false
		}
	}
}
