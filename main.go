package main

import (
	"image/color"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 800
	screenHeight = 600
	tileSize     = 32
	gridWidth    = 20
	gridHeight   = 15
)

type TileType int

const (
	Empty TileType = iota
	Road
	House
	Hospital
	Farm
)

type Game struct {
	grid      [gridHeight][gridWidth]TileType
	balance   int
	buildMode TileType
}

func NewGame() *Game {
	return &Game{
		balance:   1000,
		buildMode: Road,
	}
}

func (g *Game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		tx, ty := x/tileSize, y/tileSize
		if tx >= 0 && tx < gridWidth && ty >= 0 && ty < gridHeight {
			if g.balance >= 100 {
				g.grid[ty][tx] = g.buildMode
				g.balance -= 1
			}
		}
	}

	// Press 1-4 to select build type
	if ebiten.IsKeyPressed(ebiten.Key1) {
		g.buildMode = Road
	}
	if ebiten.IsKeyPressed(ebiten.Key2) {
		g.buildMode = House
	}
	if ebiten.IsKeyPressed(ebiten.Key3) {
		g.buildMode = Hospital
	}
	if ebiten.IsKeyPressed(ebiten.Key4) {
		g.buildMode = Farm
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{200, 230, 255, 255}) // background

	for y := 0; y < gridHeight; y++ {
		for x := 0; x < gridWidth; x++ {
			tile := g.grid[y][x]
			var c color.Color
			switch tile {
			case Road:
				c = color.RGBA{128, 128, 128, 255}
			case House:
				c = color.RGBA{255, 200, 200, 255}
			case Hospital:
				c = color.RGBA{255, 0, 0, 255}
			case Farm:
				c = color.RGBA{0, 200, 0, 255}
			default:
				c = color.RGBA{255, 255, 255, 255}
			}
			ebitenutil.DrawRect(screen, float64(x*tileSize), float64(y*tileSize), tileSize, tileSize, c)
		}
	}

	// Top bar
	ebitenutil.DebugPrintAt(screen, "1:Road 2:House 3:Hospital 4:Farm", 10, 10)
	ebitenutil.DebugPrintAt(screen, "Balance: $"+strconv.Itoa(g.balance), 10, 30)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Return the desired screen width and height
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowTitle("City Builder Mayor Game")

	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
