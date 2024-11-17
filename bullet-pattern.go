package main

import (
	"encoding/json"
	"image/color"
	"math"
	"os"
)

type BulletPattern struct {
	Angle float64 `json:"angle"`
	Speed float64 `json:"speed"`
}

func loadPatterns(fileName string) ([]BulletPattern, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var patterns []BulletPattern
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&patterns); err != nil {
		return nil, err
	}
	return patterns, nil
}

func createBulletsFromPattern(patterns []BulletPattern, centerX, centerY float64) []Bullet {
	bullets := []Bullet{}
	for _, pattern := range patterns {
		rad := pattern.Angle * (math.Pi / 180) // Convierte ángulo a radianes
		bullets = append(bullets, Bullet{
			X:      centerX,
			Y:      centerY,
			Angle:  rad,
			Speed:  pattern.Speed,
			Active: true,
			Color:  color.RGBA{255, 255, 255, 255}, // Blanco por defecto
		})
	}
	return bullets
}
