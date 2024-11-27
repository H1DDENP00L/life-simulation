package main

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/exp/rand"
)

const (
	screenHeight = 640
	screenWidth  = 480
	cellSize     = 10
	rows         = screenHeight / cellSize
	cols         = screenWidth / cellSize
)

type Game struct {
	// состояние клеток: true — живая, false — мертвая
	cells [rows][cols]bool
	// счетчик поколений
	generations int
}

func (g *Game) Randomsize() {
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			g.cells[y][x] = rand.Intn(2) == 0
		}
	}
}

var lastUpdate time.Time

func (g *Game) Update() error {
	if time.Since(lastUpdate) > 100*time.Millisecond {
		g.NextGeneration()
		lastUpdate = time.Now()
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.Randomsize()
	}
	g.generations++
	return nil
}

func (g *Game) NextGeneration() {
	var next [rows][cols]bool
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			aliveNeighbors := g.countAliveNeighboors(x, y)
			if g.cells[y][x] {
				// живая клетка остается живой с 2 или с 3 соседями
				next[y][x] = aliveNeighbors == 2 || aliveNeighbors == 3
			} else {
				//мертвая клетка оживает с 3 соседями
				next[y][x] = aliveNeighbors == 3
			}
		}

	}
	g.cells = next
}

func (g *Game) countAliveNeighboors(x, y int) int {
	count := 0
	// перебор соседей
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			nx, ny := x+dx, y+dy

			//пропускаем себя
			if dx == 0 && dy == 0 {
				continue
			}
			//учитываем только видимые клетки
			if nx >= 0 && ny >= 0 && nx < cols && ny < rows && g.cells[ny][nx] {
				count++
			}

		}
	}
	return count
}

func (g *Game) Draw(screen *ebiten.Image) {

	// отрисовка клеток
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			if g.cells[y][x] {
				rect := ebiten.NewImage(cellSize, cellSize)
				rect.Fill(color.RGBA{255, 255, 255, 10}) // живая клетка
				opts := &ebiten.DrawImageOptions{}
				opts.GeoM.Translate(float64(x*cellSize), float64(y*cellSize))
				screen.DrawImage(rect, opts)
			}
		}

	}
	ebitenutil.DebugPrint(screen, "Press SPACE to randomize\n\n")

}

func (g *Game) Layout(outsideWidth, outsideHeigth int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("That's life... That's what all the people say")

	game := &Game{}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
