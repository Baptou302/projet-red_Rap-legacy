package game

import (
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/projet-red_rap-legacy/menu"
)

// -----------------
// GameState
// -----------------
type GameState int

const (
	StateMenu GameState = iota
	StateSettings
	StatePlaying
)

// -----------------
// Game structure
// -----------------
type Game struct {
	state        GameState
	menuSelected int
	menuOptions  []string
	menuBg       *ebiten.Image
	menu         *menu.Menu
	volume       int

	player       *Player
	mapData      *Map
	enemies      []*Enemy
	inBattle     bool
	currentEnemy *Enemy
	battle       *Battle
}

// -----------------
// LoadImage helper
// -----------------
func LoadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

// -----------------
// NewGame
// -----------------
func NewGame() *Game {
	g := &Game{
		state:       StateMenu,
		menu:        menu.NewMenu(),
		volume:      5,
		menuOptions: []string{"Play", "Settings", "Quit"},
	}
	g.menuBg = LoadImage("assets/menu_bg.png")
	return g
}

// -----------------
// Ebiten methods
// -----------------
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

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case StateMenu:
		g.DrawMenu(screen)
	case StateSettings:
		g.DrawSettings(screen)
	case StatePlaying:
		if g.mapData != nil {
			g.mapData.Draw(screen)
		}
		if g.player != nil {
			g.player.Draw(screen)
		}
		for _, e := range g.enemies {
			e.Draw(screen)
		}
		if g.inBattle && g.battle != nil {
			g.battle.Draw(screen)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1920, 1080
}

// -----------------
// Menu methods
// -----------------
func (g *Game) UpdateMenu() {
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
	// Dessine le fond du menu
	if g.menuBg != nil {
		opts := &ebiten.DrawImageOptions{}
		screen.DrawImage(g.menuBg, opts)
	} else {
		screen.Fill(color.RGBA{30, 30, 30, 255})
	}

	// Dessine les options du menu
	for i, option := range g.menuOptions {
		text := option
		if i == g.menuSelected {
			text = "> " + text
		}
		ebitenutil.DebugPrintAt(screen, text, 400, 200+60*i)
	}
}

// -----------------
// Settings methods
// -----------------
func (g *Game) UpdateSettings() {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.volume++
		if g.volume > 10 {
			g.volume = 10
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.volume--
		if g.volume < 0 {
			g.volume = 0
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.state = StateMenu
	}
}

func (g *Game) DrawSettings(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "SETTINGS", 250, 100)
	ebitenutil.DebugPrintAt(screen, "Volume: "+string(rune(g.volume+'0')), 250, 150)
	ebitenutil.DebugPrintAt(screen, "Press ESC to return", 200, 200)
}

// -----------------
// StartGame
// -----------------
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

// -----------------
// Playing methods
// -----------------
func (g *Game) UpdatePlaying() {
	if g.inBattle && g.battle != nil {
		g.battle.Update()
		if g.battle.IsOver() {
			g.inBattle = false
		}
		return
	}

	if g.player != nil {
		g.player.Update()
	}

	// Détection collision avec les ennemis
	for _, e := range g.enemies {
		if int(g.player.X) == int(e.X) && int(g.player.Y) == int(e.Y) {
			g.inBattle = true
			g.currentEnemy = e
			g.battle = NewBattle(g.player, e)
		}
	}
}
