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
	Erase TileType = 99
)

type Size struct {
	W, H int
}

var tilePrices = map[TileType]int{
	Road:     2,
	House:    50,
	Hospital: 20,
	Farm:     20,
}

var tileSizes = map[TileType]Size{
	Road:     {32, 32},
	House:    {64, 64},
	Hospital: {96, 96},
	Farm:     {48, 48},
}

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
	// Build or erase with mouse click
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		tx, ty := x/tileSize, y/tileSize
		if tx >= 0 && tx < gridWidth && ty >= 0 && ty < gridHeight {
			if g.buildMode == Erase {
				g.refundAfterErase(ty, tx)
				g.grid[ty][tx] = Empty
			} else {
				price := tilePrices[g.buildMode]
				if g.balance >= price {
					g.grid[ty][tx] = g.buildMode
					g.balance -= price
				}
			}
		}
	}

	// Tool selection (keys 0–4)
	switch {
	case ebiten.IsKeyPressed(ebiten.Key0):
		g.buildMode = Erase
	case ebiten.IsKeyPressed(ebiten.Key1):
		g.buildMode = Road
	case ebiten.IsKeyPressed(ebiten.Key2):
		g.buildMode = House
	case ebiten.IsKeyPressed(ebiten.Key3):
		g.buildMode = Hospital
	case ebiten.IsKeyPressed(ebiten.Key4):
		g.buildMode = Farm
	}

	return nil
}

func (g *Game) refundAfterErase(ty, tx int) {
	tileType := g.grid[ty][tx]
	g.balance += tilePrices[tileType]
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{100, 180, 100, 255})

	// Draw tiles
	for y := 0; y < gridHeight; y++ {
		for x := 0; x < gridWidth; x++ {
			tile := g.grid[y][x]
			if tile == Empty {
				continue
			}

			size := tileSizes[tile]
			var c color.Color
			switch tile {
			case Road:
				c = color.RGBA{120, 120, 120, 255}
			case House:
				c = color.RGBA{255, 180, 180, 255}
			case Hospital:
				c = color.RGBA{255, 0, 0, 255}
			case Farm:
				c = color.RGBA{0, 200, 0, 255}
			}

			px := float64(x * tileSize)
			py := float64(y * tileSize)
			ebitenutil.DrawRect(screen, px, py, float64(size.W), float64(size.H), c)
		}
	}

	// UI
	toolName := map[TileType]string{
		Road:     "Road",
		House:    "House",
		Hospital: "Hospital",
		Farm:     "Farm",
		Erase:    "Erase",
	}[g.buildMode]

	ebitenutil.DebugPrintAt(screen, "Tools: 0-Erase 1-Road 2-House 3-Hospital 4-Farm", 10, 10)
	ebitenutil.DebugPrintAt(screen, "Selected Tool: "+toolName, 10, 30)
	ebitenutil.DebugPrintAt(screen, "Balance: $"+strconv.Itoa(g.balance), 10, 50)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowTitle("City Builder Mayor Game")
	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
