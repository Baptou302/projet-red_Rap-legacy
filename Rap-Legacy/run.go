package main

import (
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// --- GameState ---
type GameState int

const (
	StateMenu GameState = iota
	StateSettings
	StatePlaying
)

// --- Game ---
type Game struct {
	state        GameState
	menuSelected int
	menuOptions  []*ebiten.Image
	menuBg       *ebiten.Image

	player       *Player
	mapData      *Map
	enemies      []*Enemy
	inBattle     bool
	currentEnemy *Enemy
	battle       *Battle
}

// --- LoadImage helper ---
func LoadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

// --- NewGame ---
func NewGame() *Game {
	g := &Game{
		state: StateMenu,
	}

	// Charger les images du menu
	g.menuBg = LoadImage("assets/menu_bg.png")
	return g
}

// --- Update ---
func (g *Game) Update() error {
	switch g.state {
	case StateMenu:
		g.UpdateMenu()
	case StateSettings:
		g.UpdateSettings()
	case StatePlaying:
		g.UpdatePlaying()
	}
	return nil
}

// --- Draw ---
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case StateMenu:
		g.DrawMenu(screen)
	case StateSettings:
		g.DrawSettings(screen)
	case StatePlaying:
		g.mapData.Draw(screen)
		g.player.Draw(screen)
		for _, e := range g.enemies {
			e.Draw(screen)
		}
		if g.inBattle {
			g.battle.Draw(screen)
		}
	}
}

// --- Menu ---
func (g *Game) UpdateMenu() {
	// Navigation
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.menuSelected--
		if g.menuSelected < 0 {
			g.menuSelected = len(g.menuOptions) - 1
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.menuSelected++
		if g.menuSelected >= len(g.menuOptions) {
			g.menuSelected = 0
		}
	}

	// Validation
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		switch g.menuSelected {
		case 0:
			g.StartGame()
		case 1:
			g.state = StateSettings
		case 2:
			os.Exit(0)
		}
	}
}

func (g *Game) DrawMenu(screen *ebiten.Image) {
	// Fond
	if g.menuBg != nil {
		opts := &ebiten.DrawImageOptions{}
		screen.DrawImage(g.menuBg, opts)
	} else {
		screen.Fill(color.RGBA{30, 30, 30, 255})
	}

	// Options
	for i, img := range g.menuOptions {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(220, 200+float64(i*60))
		if i == g.menuSelected {
			opts.ColorM.Scale(1.2, 1.2, 1.2, 1) // surlignage léger
		}
		screen.DrawImage(img, opts)
	}
}

// --- Settings simple ---
var volume int = 5

func (g *Game) UpdateSettings() {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		volume++
		if volume > 10 {
			volume = 10
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		volume--
		if volume < 0 {
			volume = 0
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.state = StateMenu
	}
}

func (g *Game) DrawSettings(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "SETTINGS", 250, 100)
	ebitenutil.DebugPrintAt(screen, "Volume: "+string(rune(volume+'0')), 250, 150)
	ebitenutil.DebugPrintAt(screen, "Press ESC to return", 200, 200)
}

// --- StartGame ---
func (g *Game) StartGame() {
	g.state = StatePlaying
	g.player = NewPlayer(100, 100)
	g.mapData = NewMap()
	g.enemies = []*Enemy{
		NewEnemy(200, 200, "Rival Rapper"),
		NewEnemy(400, 300, "Boss Rapper"),
	}
	g.inBattle = false
}

// --- UpdatePlaying ---
func (g *Game) UpdatePlaying() {
	if g.inBattle {
		g.battle.Update()
		if g.battle.IsOver() {
			g.inBattle = false
		}
	} else {
		g.player.Update()
		for _, e := range g.enemies {
			if g.player.X == e.X && g.player.Y == e.Y {
				g.inBattle = true
				g.currentEnemy = e
				g.battle = NewBattle(g.player, e)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1920, 1080
}

func main() {
	game := NewGame()
	ebiten.SetWindowSize(1920, 1080)
	ebiten.SetWindowTitle("Bars & Bosses RPG")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
