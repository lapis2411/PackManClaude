package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const TileSize = 30

type Player struct {
	X     float64
	Y     float64
	Speed float64
}

func (p *Player) Update(maze [][]int) {
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		newY := p.Y - p.Speed
		if !p.isColliding(p.X, newY, maze) {
			p.Y = newY
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		newY := p.Y + p.Speed
		if !p.isColliding(p.X, newY, maze) {
			p.Y = newY
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		newX := p.X - p.Speed
		if !p.isColliding(newX, p.Y, maze) {
			p.X = newX
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		newX := p.X + p.Speed
		if !p.isColliding(newX, p.Y, maze) {
			p.X = newX
		}
	}
}

func (p *Player) isColliding(x, y float64, maze [][]int) bool {
	radius := float64(TileSize) / 3
	
	// プレイヤーの円の境界4点をチェック
	checkPoints := []struct{ px, py float64 }{
		{x - radius, y - radius}, // 左上
		{x + radius, y - radius}, // 右上
		{x - radius, y + radius}, // 左下
		{x + radius, y + radius}, // 右下
	}
	
	for _, point := range checkPoints {
		tileX := int(point.px / TileSize)
		tileY := int(point.py / TileSize)
		
		if tileY < 0 || tileY >= len(maze) || tileX < 0 || tileX >= len(maze[0]) {
			return true
		}
		
		if maze[tileY][tileX] == 1 {
			return true
		}
	}
	
	return false
}

type Game struct {
	maze   [][]int
	player Player
}

func (g *Game) Update() error {
	g.player.Update(g.maze)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for y, row := range g.maze {
		for x, tile := range row {
			if tile == 1 {
				vector.DrawFilledRect(screen, float32(x*TileSize), float32(y*TileSize), TileSize, TileSize, color.RGBA{R: 0, G: 0, B: 255, A: 255}, false)
			}
		}
	}
	
	vector.DrawFilledCircle(screen, float32(g.player.X), float32(g.player.Y), TileSize/3, color.RGBA{R: 0xff, G: 0xff, B: 0, A: 0xff}, false)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 32 * TileSize, 17 * TileSize
}

func main() {
	game := &Game{
		maze: [][]int{
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 1, 1, 1, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 0, 1, 1, 1, 1, 0, 1, 1},
			{1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 0, 1, 0, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
		player: Player{
			X:     TileSize + TileSize/2,
			Y:     TileSize + TileSize/2,
			Speed: 2.0,
		},
	}
	
	ebiten.SetWindowTitle("PackMan Game")
	ebiten.SetWindowSize(32*TileSize, 17*TileSize)
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}