package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	TileSize = 30
	FrightenedDuration = 300 // フレーム数 (約5秒)
)

type GhostState int

const (
	Normal GhostState = 0
	Frightened GhostState = 1
)

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
	X               float64
	Y               float64
	Speed           float64
	DirX            float64
	DirY            float64
	State           GhostState
	FrightenedTimer int
	InitialX        float64
	InitialY        float64
}

func (g *Ghost) Update(maze [][]int, playerX, playerY float64) {
	if g.State == Frightened {
		g.FrightenedTimer--
		if g.FrightenedTimer <= 0 {
			g.State = Normal
		}
	}
	
	newX := g.X + g.DirX*g.Speed
	newY := g.Y + g.DirY*g.Speed
	
	if g.isColliding(newX, newY, maze) {
		g.chooseDirection(maze, playerX, playerY)
	} else {
		g.X = newX
		g.Y = newY
		
		if g.isAtIntersection(maze) {
			g.chooseDirection(maze, playerX, playerY)
		}
	}
}

func (g *Ghost) SetFrightened() {
	g.State = Frightened
	g.FrightenedTimer = FrightenedDuration
}

func (g *Ghost) ResetToInitialPosition() {
	g.X = g.InitialX
	g.Y = g.InitialY
	g.State = Normal
	g.FrightenedTimer = 0
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

func (g *Ghost) chooseDirection(maze [][]int, playerX, playerY float64) {
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
		var bestDirection []float64
		var bestDistance float64
		
		if g.State == Frightened {
			bestDistance = -1
		} else {
			bestDistance = float64(999999)
		}
		
		for _, dir := range validDirections {
			nextX := g.X + dir[0]*TileSize
			nextY := g.Y + dir[1]*TileSize
			
			dx := nextX - playerX
			dy := nextY - playerY
			distance := dx*dx + dy*dy
			
			if g.State == Frightened {
				if distance > bestDistance {
					bestDistance = distance
					bestDirection = dir
				}
			} else {
				if distance < bestDistance {
					bestDistance = distance
					bestDirection = dir
				}
			}
		}
		
		if bestDirection != nil {
			g.DirX = bestDirection[0]
			g.DirY = bestDirection[1]
		} else {
			chosen := validDirections[rand.Intn(len(validDirections))]
			g.DirX = chosen[0]
			g.DirY = chosen[1]
		}
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

type Scene interface {
	Update() Scene
	Draw(screen *ebiten.Image)
}

type GameScene struct {
	maze   [][]int
	player Player
	ghost  Ghost
	Score  int
}

func (gs *GameScene) Update() Scene {
	gs.player.Update(gs.maze)
	gs.ghost.Update(gs.maze, gs.player.X, gs.player.Y)
	gs.checkItemCollection()
	
	if gs.checkPlayerGhostCollision() {
		return &GameOverScene{}
	}
	
	if gs.checkStageClear() {
		return &StageClearScene{}
	}
	
	return gs
}

func (gs *GameScene) checkItemCollection() {
	tileX := int(gs.player.X / TileSize)
	tileY := int(gs.player.Y / TileSize)
	
	if tileY >= 0 && tileY < len(gs.maze) && tileX >= 0 && tileX < len(gs.maze[0]) {
		if gs.maze[tileY][tileX] == 2 {
			gs.maze[tileY][tileX] = 0
			gs.Score += 10
		} else if gs.maze[tileY][tileX] == 3 {
			gs.maze[tileY][tileX] = 0
			gs.Score += 50
			gs.ghost.SetFrightened()
		}
	}
}

func (gs *GameScene) checkPlayerGhostCollision() bool {
	playerRadius := float64(TileSize) / 3
	ghostRadius := float64(TileSize) / 3
	
	dx := gs.player.X - gs.ghost.X
	dy := gs.player.Y - gs.ghost.Y
	distance := dx*dx + dy*dy
	
	if distance < (playerRadius+ghostRadius)*(playerRadius+ghostRadius) {
		if gs.ghost.State == Frightened {
			gs.ghost.ResetToInitialPosition()
			gs.Score += 200
			return false
		} else {
			return true
		}
	}
	
	return false
}

func (gs *GameScene) checkStageClear() bool {
	for _, row := range gs.maze {
		for _, tile := range row {
			if tile == 2 || tile == 3 {
				return false
			}
		}
	}
	return true
}

func (gs *GameScene) Draw(screen *ebiten.Image) {
	for y, row := range gs.maze {
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
	
	vector.DrawFilledCircle(screen, float32(gs.player.X), float32(gs.player.Y), TileSize/3, color.RGBA{R: 0xff, G: 0xff, B: 0, A: 0xff}, false)
	
	var ghostColor color.RGBA
	if gs.ghost.State == Frightened {
		ghostColor = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	} else {
		ghostColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	}
	vector.DrawFilledCircle(screen, float32(gs.ghost.X), float32(gs.ghost.Y), TileSize/3, ghostColor, false)
	
	gs.drawScore(screen)
}

func (gs *GameScene) drawScore(screen *ebiten.Image) {
	scoreText := fmt.Sprintf("SCORE: %d", gs.Score)
	
	pixelSize := float32(2)
	letterSpacing := float32(12)
	startX := float32(10)
	startY := float32(10)
	
	letterPatterns := map[rune][][]bool{
		'S': {
			{false, true, true, true, false},
			{true, false, false, false, false},
			{false, true, true, true, false},
			{false, false, false, false, true},
			{true, true, true, true, false},
		},
		'C': {
			{false, true, true, true, false},
			{true, false, false, false, false},
			{true, false, false, false, false},
			{true, false, false, false, false},
			{false, true, true, true, false},
		},
		'O': {
			{false, true, true, true, false},
			{true, false, false, false, true},
			{true, false, false, false, true},
			{true, false, false, false, true},
			{false, true, true, true, false},
		},
		'R': {
			{true, true, true, true, false},
			{true, false, false, false, true},
			{true, true, true, true, false},
			{true, false, false, true, false},
			{true, false, false, false, true},
		},
		'E': {
			{true, true, true, true, true},
			{true, false, false, false, false},
			{true, true, true, true, false},
			{true, false, false, false, false},
			{true, true, true, true, true},
		},
		':': {
			{false, false, false, false, false},
			{false, true, false, false, false},
			{false, false, false, false, false},
			{false, true, false, false, false},
			{false, false, false, false, false},
		},
		' ': {
			{false, false, false, false, false},
			{false, false, false, false, false},
			{false, false, false, false, false},
			{false, false, false, false, false},
			{false, false, false, false, false},
		},
		'0': {
			{false, true, true, true, false},
			{true, false, false, false, true},
			{true, false, false, false, true},
			{true, false, false, false, true},
			{false, true, true, true, false},
		},
		'1': {
			{false, false, true, false, false},
			{false, true, true, false, false},
			{false, false, true, false, false},
			{false, false, true, false, false},
			{false, true, true, true, false},
		},
		'2': {
			{false, true, true, true, false},
			{false, false, false, false, true},
			{false, true, true, true, false},
			{true, false, false, false, false},
			{true, true, true, true, true},
		},
		'3': {
			{true, true, true, true, false},
			{false, false, false, false, true},
			{false, true, true, true, false},
			{false, false, false, false, true},
			{true, true, true, true, false},
		},
		'4': {
			{true, false, false, false, true},
			{true, false, false, false, true},
			{true, true, true, true, true},
			{false, false, false, false, true},
			{false, false, false, false, true},
		},
		'5': {
			{true, true, true, true, true},
			{true, false, false, false, false},
			{true, true, true, true, false},
			{false, false, false, false, true},
			{true, true, true, true, false},
		},
		'6': {
			{false, true, true, true, false},
			{true, false, false, false, false},
			{true, true, true, true, false},
			{true, false, false, false, true},
			{false, true, true, true, false},
		},
		'7': {
			{true, true, true, true, true},
			{false, false, false, false, true},
			{false, false, false, true, false},
			{false, false, true, false, false},
			{false, true, false, false, false},
		},
		'8': {
			{false, true, true, true, false},
			{true, false, false, false, true},
			{false, true, true, true, false},
			{true, false, false, false, true},
			{false, true, true, true, false},
		},
		'9': {
			{false, true, true, true, false},
			{true, false, false, false, true},
			{false, true, true, true, true},
			{false, false, false, false, true},
			{false, true, true, true, false},
		},
	}
	
	for i, char := range scoreText {
		if pattern, exists := letterPatterns[char]; exists {
			letterStartX := startX + float32(i)*letterSpacing
			
			for row := 0; row < 5; row++ {
				for col := 0; col < 5; col++ {
					if pattern[row][col] {
						x := letterStartX + float32(col)*pixelSize
						y := startY + float32(row)*pixelSize
						vector.DrawFilledRect(screen, x, y, pixelSize, pixelSize, color.RGBA{R: 255, G: 255, B: 255, A: 255}, false)
					}
				}
			}
		}
	}
}

type GameOverScene struct{}

func (gos *GameOverScene) Update() Scene {
	return gos
}

func (gos *GameOverScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	
	screenWidth := 32 * TileSize
	screenHeight := 17 * TileSize
	
	centerX := float32(screenWidth / 2)
	centerY := float32(screenHeight / 2)
	
	vector.DrawFilledRect(screen, centerX-150, centerY-40, 300, 80, color.RGBA{R: 255, G: 255, B: 255, A: 255}, false)
	vector.DrawFilledRect(screen, centerX-145, centerY-35, 290, 70, color.RGBA{R: 255, G: 0, B: 0, A: 255}, false)
	
	vector.DrawFilledRect(screen, centerX-130, centerY-20, 260, 40, color.RGBA{R: 0, G: 0, B: 0, A: 255}, false)
	
	letters := [][]bool{
		// G
		{true, true, true, true, true},
		{true, false, false, false, false},
		{true, false, true, true, true},
		{true, false, false, false, true},
		{true, true, true, true, true},
		
		// A
		{false, true, true, true, false},
		{true, false, false, false, true},
		{true, true, true, true, true},
		{true, false, false, false, true},
		{true, false, false, false, true},
		
		// M
		{true, false, false, false, true},
		{true, true, false, true, true},
		{true, false, true, false, true},
		{true, false, false, false, true},
		{true, false, false, false, true},
		
		// E
		{true, true, true, true, true},
		{true, false, false, false, false},
		{true, true, true, true, false},
		{true, false, false, false, false},
		{true, true, true, true, true},
		
		// (space)
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, false, false, false, false},
		
		// O
		{false, true, true, true, false},
		{true, false, false, false, true},
		{true, false, false, false, true},
		{true, false, false, false, true},
		{false, true, true, true, false},
		
		// V
		{true, false, false, false, true},
		{true, false, false, false, true},
		{false, true, false, true, false},
		{false, true, false, true, false},
		{false, false, true, false, false},
		
		// E
		{true, true, true, true, true},
		{true, false, false, false, false},
		{true, true, true, true, false},
		{true, false, false, false, false},
		{true, true, true, true, true},
		
		// R
		{true, true, true, true, false},
		{true, false, false, false, true},
		{true, true, true, true, false},
		{true, false, false, true, false},
		{true, false, false, false, true},
	}
	
	pixelSize := float32(3)
	letterSpacing := float32(20)
	startX := centerX - float32(9*int(letterSpacing))/2
	
	for letterIndex := 0; letterIndex < 9; letterIndex++ {
		letterStartX := startX + float32(letterIndex)*letterSpacing
		
		for row := 0; row < 5; row++ {
			for col := 0; col < 5; col++ {
				letterArrayIndex := letterIndex * 5 + row
				if letterArrayIndex < len(letters) && col < len(letters[letterArrayIndex]) && letters[letterArrayIndex][col] {
					x := letterStartX + float32(col)*pixelSize
					y := centerY - 10 + float32(row)*pixelSize
					vector.DrawFilledRect(screen, x, y, pixelSize, pixelSize, color.RGBA{R: 255, G: 255, B: 255, A: 255}, false)
				}
			}
		}
	}
}

type StageClearScene struct{}

func (scs *StageClearScene) Update() Scene {
	return scs
}

func (scs *StageClearScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	
	screenWidth := 32 * TileSize
	screenHeight := 17 * TileSize
	
	centerX := float32(screenWidth / 2)
	centerY := float32(screenHeight / 2)
	
	vector.DrawFilledRect(screen, centerX-180, centerY-40, 360, 80, color.RGBA{R: 0, G: 255, B: 0, A: 255}, false)
	vector.DrawFilledRect(screen, centerX-175, centerY-35, 350, 70, color.RGBA{R: 0, G: 0, B: 0, A: 255}, false)
	
	letters := [][]bool{
		// S
		{true, true, true, true, true},
		{true, false, false, false, false},
		{true, true, true, true, false},
		{false, false, false, false, true},
		{true, true, true, true, true},
		
		// T
		{true, true, true, true, true},
		{false, false, true, false, false},
		{false, false, true, false, false},
		{false, false, true, false, false},
		{false, false, true, false, false},
		
		// A
		{false, true, true, true, false},
		{true, false, false, false, true},
		{true, true, true, true, true},
		{true, false, false, false, true},
		{true, false, false, false, true},
		
		// G
		{true, true, true, true, true},
		{true, false, false, false, false},
		{true, false, true, true, true},
		{true, false, false, false, true},
		{true, true, true, true, true},
		
		// E
		{true, true, true, true, true},
		{true, false, false, false, false},
		{true, true, true, true, false},
		{true, false, false, false, false},
		{true, true, true, true, true},
		
		// (space)
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, false, false, false, false},
		
		// C
		{false, true, true, true, false},
		{true, false, false, false, false},
		{true, false, false, false, false},
		{true, false, false, false, false},
		{false, true, true, true, false},
		
		// L
		{true, false, false, false, false},
		{true, false, false, false, false},
		{true, false, false, false, false},
		{true, false, false, false, false},
		{true, true, true, true, true},
		
		// E
		{true, true, true, true, true},
		{true, false, false, false, false},
		{true, true, true, true, false},
		{true, false, false, false, false},
		{true, true, true, true, true},
		
		// A
		{false, true, true, true, false},
		{true, false, false, false, true},
		{true, true, true, true, true},
		{true, false, false, false, true},
		{true, false, false, false, true},
		
		// R
		{true, true, true, true, false},
		{true, false, false, false, true},
		{true, true, true, true, false},
		{true, false, false, true, false},
		{true, false, false, false, true},
	}
	
	pixelSize := float32(3)
	letterSpacing := float32(18)
	startX := centerX - float32(11*int(letterSpacing))/2
	
	for letterIndex := 0; letterIndex < 11; letterIndex++ {
		letterStartX := startX + float32(letterIndex)*letterSpacing
		
		for row := 0; row < 5; row++ {
			for col := 0; col < 5; col++ {
				letterArrayIndex := letterIndex * 5 + row
				if letterArrayIndex < len(letters) && col < len(letters[letterArrayIndex]) && letters[letterArrayIndex][col] {
					x := letterStartX + float32(col)*pixelSize
					y := centerY - 10 + float32(row)*pixelSize
					vector.DrawFilledRect(screen, x, y, pixelSize, pixelSize, color.RGBA{R: 255, G: 255, B: 255, A: 255}, false)
				}
			}
		}
	}
}

type Game struct {
	currentScene Scene
}

func (g *Game) Update() error {
	g.currentScene = g.currentScene.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.currentScene.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 32 * TileSize, 17 * TileSize
}

func main() {
	gameScene := &GameScene{
		maze: [][]int{
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			{1, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 1},
			{1, 2, 1, 1, 1, 1, 3, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 3, 1, 1, 1, 1, 2, 1},
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
			X:               TileSize*16 + TileSize/2,
			Y:               TileSize*9 + TileSize/2,
			Speed:           1.5,
			DirX:            1.0,
			DirY:            0.0,
			State:           Normal,
			FrightenedTimer: 0,
			InitialX:        TileSize*16 + TileSize/2,
			InitialY:        TileSize*9 + TileSize/2,
		},
	}
	
	game := &Game{
		currentScene: gameScene,
	}
	
	ebiten.SetWindowTitle("PackMan Game")
	ebiten.SetWindowSize(32*TileSize, 17*TileSize)
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}