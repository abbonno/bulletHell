package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/nfnt/resize"
)

// Parámetros de la pantalla y el icono
const (
	screenWidth  = 1400
	screenHeight = 600
	iconSize     = 50
	bulletRadius = 5
	moveSpeed    = 5
	background   = ""
	icon         = "./public/img/icon.png"
	patterns     = "./public/json/patterns.json"
)

var tickCount int

type Game struct {
	iconImage       *ebiten.Image
	bgImage         *ebiten.Image
	bgColor         color.Color
	iconX           float64
	iconY           float64
	bullets         []Bullet
	patterns        []BulletPattern
	enemy           Enemy
	health          int // Vida del jugador
	score           int // Puntos
	start           time.Time
	gameTime        float64 // Tiempo de juego
	invincible      bool    // Estado de invulnerabilidad
	invincibleTimer int     // Duración de la invulnerabilidad en ticks
	coins           []Coin  // Monedas generadas aleatoriamente
}

func (g *Game) drawUI(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("Vida: %d\nPuntos: %d\nTiempo: %.1fs", g.health, g.score, g.gameTime))
}

func (g *Game) handleCollisions() {
	if g.invincible {
		g.invincibleTimer--
		if g.invincibleTimer <= 0 {
			g.invincible = false
		}
		return
	}

	iconCenterX := g.iconX + iconSize/2
	iconCenterY := g.iconY + iconSize/2

	for _, bullet := range g.bullets {
		if !bullet.Active {
			continue
		}

		dx := bullet.X - iconCenterX
		dy := bullet.Y - iconCenterY
		distance := math.Sqrt(dx*dx + dy*dy)

		if distance <= bulletRadius+iconSize/2 {
			g.health--
			g.invincible = true
			g.invincibleTimer = 60 // 1 segundo de invulnerabilidad si TPS=60
			if g.health <= 0 {
				log.Println("Game Over")
				return
			}
			break
		}
	}
}

// Inicializa el juego
func (g *Game) initGame() {
	// Cargar la imagen del icono
	file, err := os.Open(icon)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	resizedImg := resize.Resize(iconSize, iconSize, img, resize.Lanczos3)
	g.iconImage = ebiten.NewImageFromImage(resizedImg)

	// Cargar el fondo
	g.bgImage, _, err = ebitenutil.NewImageFromFile(background)
	if err != nil {
		log.Println("No se pudo cargar la imagen de fondo, usando color de fondo.")
		g.bgImage = nil
		g.bgColor = color.RGBA{50, 50, 150, 255}
	}

	// Posición inicial del icono
	g.iconX = screenWidth/2 - iconSize/2
	g.iconY = screenHeight/2 - iconSize/2

	// Interfaz
	g.health = 3
	g.score = 0
	g.gameTime = 0
	g.start = time.Now()
	g.invincible = false
	g.invincibleTimer = 0
	g.coins = []Coin{}

	// Inicializar enemigo
	g.enemy = Enemy{
		X:     0, // Esquina superior izquierda
		Y:     0,
		Speed: 1.5,                        // Velocidad de seguimiento
		Color: color.RGBA{255, 0, 0, 255}, // Rojo
	}

	// Cargar patrones de balas
	g.patterns, err = loadPatterns(patterns)
	if err != nil {
		log.Fatal(err)
	}
}

// Actualiza el estado del juego en cada tick
func (g *Game) Update() error {
	tickCount++
	// Mover el icono con las teclas de flecha
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.iconY -= moveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.iconY += moveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.iconX -= moveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.iconX += moveSpeed
	}

	// Interfaz
	g.handleCollisions()
	g.gameTime = time.Since(g.start).Seconds()
	g.generateCoins()
	g.handleCoinCollisions()

	// Crear nuevas balas en intervalos
	if tickCount%60 == 0 { // Cada 60 ticks (1 segundo si TPS=60)
		g.bullets = append(g.bullets, createBulletsFromPattern(g.patterns, g.enemy.X, g.enemy.Y)...)

	}

	// Actualizar enemigo
	g.enemy.update(g.iconX, g.iconY)

	// Actualizar balas existentes
	updateBullets(&g.bullets, g.iconX, g.iconY)

	return nil
}

// Dibuja el juego en la ventana
func (g *Game) Draw(screen *ebiten.Image) {
	// Dibuja el fondo
	if g.bgImage != nil {
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(g.bgImage, op)
	} else {
		screen.Fill(g.bgColor)
	}

	// Dibuja el icono
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.iconX, g.iconY)
	screen.DrawImage(g.iconImage, op)

	// Dibuja la interfaz
	g.drawUI(screen)

	// Dibuja las monedas
	for _, coin := range g.coins {
		ebitenutil.DrawCircle(screen, coin.X, coin.Y, bulletRadius, coin.Color)
	}

	// Dibuja al enemigo
	g.enemy.draw(screen)

	// Dibuja las balas
	drawBullets(screen, g.bullets)
}

// Define el tamaño de la ventana
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{}
	game.initGame()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Mover Icono con Flechas")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
