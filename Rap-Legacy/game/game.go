package game

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// -----------------
// Button structure
// -----------------
type Button struct {
	Rect   image.Rectangle
	Label  string
	Action func()
}

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
	menuButtons  []*Button
	menuBg       *ebiten.Image
	menuSelected int
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
		state:  StateMenu,
		volume: 5,
	}
	g.menuBg = LoadImage("assets/menu_bg.png")

	// Ici tu peux configurer les boutons
	g.menuButtons = []*Button{
		{
			Rect:   image.Rect(689, 489, 1152, 453), // Coin supérieur gauche (800,400), coin inférieur droit (1120,460)
			Label:  "Play",
			Action: g.StartGame,
		},
		{
			Rect:   image.Rect(800, 500, 1120, 560),
			Label:  "Settings",
			Action: func() { g.state = StateSettings },
		},
		{
			Rect:   image.Rect(800, 600, 1120, 660),
			Label:  "Quit",
			Action: func() { os.Exit(0) },
		},
	}

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
	mx, my := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		for _, btn := range g.menuButtons {
			if mx >= btn.Rect.Min.X && mx <= btn.Rect.Max.X &&
				my >= btn.Rect.Min.Y && my <= btn.Rect.Max.Y {
				btn.Action()
			}
		}
	}

	// Optionnel : navigation clavier
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.menuSelected--
		if g.menuSelected < 0 {
			g.menuSelected = len(g.menuButtons) - 1
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.menuSelected++
		if g.menuSelected >= len(g.menuButtons) {
			g.menuSelected = 0
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		g.menuButtons[g.menuSelected].Action()
	}
}

func (g *Game) DrawMenu(screen *ebiten.Image) {
	// Dessine le fond
	if g.menuBg != nil {
		opts := &ebiten.DrawImageOptions{}
		screen.DrawImage(g.menuBg, opts)
	} else {
		screen.Fill(color.RGBA{30, 30, 30, 255})
	}

	// Dessine les zones des boutons pour config (semi-transparent rouge)
	for _, btn := range g.menuButtons {
		ebitenutil.DrawRect(screen,
			float64(btn.Rect.Min.X),
			float64(btn.Rect.Min.Y),
			float64(btn.Rect.Dx()),
			float64(btn.Rect.Dy()),
			color.RGBA{255, 0, 0, 100}, // Rouge semi-transparent
		)
		// Affiche le label
		ebitenutil.DebugPrintAt(screen, btn.Label, btn.Rect.Min.X+10, btn.Rect.Min.Y+5)
	}

	// Affiche la position de la souris pour aider à configurer
	mx, my := ebiten.CursorPosition()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Cursor X: %d Y: %d", mx, my))
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

	for _, e := range g.enemies {
		if int(g.player.X) == int(e.X) && int(g.player.Y) == int(e.Y) {
			g.inBattle = true
			g.currentEnemy = e
			g.battle = NewBattle(g.player, e)
		}
	}
}
