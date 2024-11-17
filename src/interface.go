package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Button struct {
	X, Y, W, H float64
	Text       string
	OnClick    func()
}

func (b *Button) isMouseOver(x, y float64) bool {
	return x >= b.X && x <= b.X+b.W && y >= b.Y && y <= b.Y+b.H
}

func (b *Button) handleClick(x, y float64) {
	if b.isMouseOver(x, y) && b.OnClick != nil {
		b.OnClick()
	}
}

func (b *Button) draw(screen *ebiten.Image) {
	// Dibuja el botón como un rectángulo
	ebitenutil.DrawRect(screen, b.X, b.Y, b.W, b.H, color.RGBA{0, 0, 255, 255}) // Azul
	// Dibuja el texto del botón
	textX := b.X + b.W/4
	textY := b.Y + b.H/4
	ebitenutil.DebugPrintAt(screen, b.Text, int(textX), int(textY))
}
