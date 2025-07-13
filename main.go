package main

import (
	"image/color"
	"math/rand"

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

type Ghost struct {
	X     float64
	Y     float64
	Speed float64
	DirX  float64
	DirY  float64
}

func (g *Ghost) Update(maze [][]int) {
	newX := g.X + g.DirX*g.Speed
	newY := g.Y + g.DirY*g.Speed
	
	if g.isColliding(newX, newY, maze) {
		g.chooseRandomDirection(maze)
	} else {
		g.X = newX
		g.Y = newY
		
		if g.isAtIntersection(maze) {
			g.chooseRandomDirection(maze)
		}
	}
}

func (g *Ghost) isColliding(x, y float64, maze [][]int) bool {
	radius := float64(TileSize) / 3
	
	checkPoints := []struct{ px, py float64 }{
		{x - radius, y - radius},
		{x + radius, y - radius},
		{x - radius, y + radius},
		{x + radius, y + radius},
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

func (g *Ghost) isAtIntersection(maze [][]int) bool {
	tileX := int(g.X / TileSize)
	tileY := int(g.Y / TileSize)
	
	if tileY < 0 || tileY >= len(maze) || tileX < 0 || tileX >= len(maze[0]) {
		return false
	}
	
	centerX := float64(tileX*TileSize + TileSize/2)
	centerY := float64(tileY*TileSize + TileSize/2)
	
	if abs(g.X-centerX) < 5 && abs(g.Y-centerY) < 5 {
		directions := [][]float64{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
		validDirections := 0
		
		for _, dir := range directions {
			nextTileX := tileX + int(dir[0])
			nextTileY := tileY + int(dir[1])
			
			if nextTileY >= 0 && nextTileY < len(maze) && nextTileX >= 0 && nextTileX < len(maze[0]) {
				if maze[nextTileY][nextTileX] != 1 {
					validDirections++
				}
			}
		}
		
		return validDirections > 2
	}
	
	return false
}

func (g *Ghost) chooseRandomDirection(maze [][]int) {
	tileX := int(g.X / TileSize)
	tileY := int(g.Y / TileSize)
	
	directions := [][]float64{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	var validDirections [][]float64
	
	for _, dir := range directions {
		nextTileX := tileX + int(dir[0])
		nextTileY := tileY + int(dir[1])
		
		if nextTileY >= 0 && nextTileY < len(maze) && nextTileX >= 0 && nextTileX < len(maze[0]) {
			if maze[nextTileY][nextTileX] != 1 {
				if dir[0] != -g.DirX || dir[1] != -g.DirY {
					validDirections = append(validDirections, dir)
				}
			}
		}
	}
	
	if len(validDirections) > 0 {
		chosen := validDirections[rand.Intn(len(validDirections))]
		g.DirX = chosen[0]
		g.DirY = chosen[1]
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

type Game struct {
	maze   [][]int
	player Player
	ghost  Ghost
	Score  int
}

func (g *Game) Update() error {
	g.player.Update(g.maze)
	g.ghost.Update(g.maze)
	g.checkItemCollection()
	return nil
}

func (g *Game) checkItemCollection() {
	tileX := int(g.player.X / TileSize)
	tileY := int(g.player.Y / TileSize)
	
	if tileY >= 0 && tileY < len(g.maze) && tileX >= 0 && tileX < len(g.maze[0]) {
		if g.maze[tileY][tileX] == 2 {
			g.maze[tileY][tileX] = 0
			g.Score += 10
		} else if g.maze[tileY][tileX] == 3 {
			g.maze[tileY][tileX] = 0
			g.Score += 50
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	for y, row := range g.maze {
		for x, tile := range row {
			switch tile {
			case 1:
				vector.DrawFilledRect(screen, float32(x*TileSize), float32(y*TileSize), TileSize, TileSize, color.RGBA{R: 0, G: 0, B: 255, A: 255}, false)
			case 2:
				centerX := float32(x*TileSize + TileSize/2)
				centerY := float32(y*TileSize + TileSize/2)
				vector.DrawFilledCircle(screen, centerX, centerY, 2, color.RGBA{R: 255, G: 255, B: 255, A: 255}, false)
			case 3:
				centerX := float32(x*TileSize + TileSize/2)
				centerY := float32(y*TileSize + TileSize/2)
				vector.DrawFilledCircle(screen, centerX, centerY, 5, color.RGBA{R: 255, G: 255, B: 255, A: 255}, false)
			}
		}
	}
	
	vector.DrawFilledCircle(screen, float32(g.player.X), float32(g.player.Y), TileSize/3, color.RGBA{R: 0xff, G: 0xff, B: 0, A: 0xff}, false)
	vector.DrawFilledCircle(screen, float32(g.ghost.X), float32(g.ghost.Y), TileSize/3, color.RGBA{R: 255, G: 0, B: 0, A: 255}, false)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 32 * TileSize, 17 * TileSize
}

func main() {
	game := &Game{
		maze: [][]int{
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			{1, 3, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 3, 1},
			{1, 2, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 2, 1},
			{1, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 1},
			{1, 2, 1, 1, 1, 1, 2, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 2, 1, 1, 1, 1, 2, 1, 1},
			{1, 2, 2, 2, 2, 2, 2, 1, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 1, 2, 2, 2, 2, 2, 2, 2, 1},
			{1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1},
			{1, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 1},
			{1, 1, 1, 1, 1, 1, 2, 1, 2, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 1, 2, 1, 2, 1, 1, 1, 1, 1, 1, 1},
			{1, 2, 2, 2, 2, 2, 2, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 2, 2, 2, 2, 2, 2, 2, 1},
			{1, 1, 1, 1, 1, 1, 2, 1, 2, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 2, 1, 2, 1, 1, 1, 1, 1, 1, 1},
			{1, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 1},
			{1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1},
			{1, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 1},
			{1, 2, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 2, 1},
			{1, 3, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 3, 1},
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
		player: Player{
			X:     TileSize + TileSize/2,
			Y:     TileSize + TileSize/2,
			Speed: 2.0,
		},
		ghost: Ghost{
			X:     TileSize*16 + TileSize/2,
			Y:     TileSize*9 + TileSize/2,
			Speed: 1.5,
			DirX:  1.0,
			DirY:  0.0,
		},
	}
	
	ebiten.SetWindowTitle("PackMan Game")
	ebiten.SetWindowSize(32*TileSize, 17*TileSize)
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}