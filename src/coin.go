package main

import (
	"image/color"
	"math"
	"math/rand"
)

type Coin struct {
	X     float64     // Posición X
	Y     float64     // Posición Y
	Color color.Color // Color de la moneda
}

func (g *Game) generateCoins() {
	if len(g.coins) < 5 { // Máximo 5 monedas activas
		g.coins = append(g.coins, Coin{
			X:     float64(rand.Intn(screenWidth - int(iconSize))),
			Y:     float64(rand.Intn(screenHeight - int(iconSize))),
			Color: color.RGBA{255, 223, 0, 255}, // Amarillo
		})
	}
}

func (g *Game) handleCoinCollisions() {
	iconCenterX := g.iconX + iconSize/2
	iconCenterY := g.iconY + iconSize/2

	newCoins := []Coin{}
	for _, coin := range g.coins {
		dx := coin.X - iconCenterX
		dy := coin.Y - iconCenterY
		distance := math.Sqrt(dx*dx + dy*dy)

		if distance <= bulletRadius+iconSize/2 {
			g.score += 10 // Incrementa puntos
		} else {
			newCoins = append(newCoins, coin) // Mantiene monedas no recogidas
		}
	}
	g.coins = newCoins
}
