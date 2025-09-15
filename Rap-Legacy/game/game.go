package game

import (
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

	// Zone de combat
	combatZone image.Rectangle
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

	// Configuration des boutons
	g.menuButtons = []*Button{
		{
			Rect:   image.Rect(685, 490, 1160, 550), // New Game
			Label:  "Play",
			Action: g.StartGame,
		},
		{
			Rect:   image.Rect(730, 600, 1110, 660), // Settings
			Label:  "Settings",
			Action: func() { g.state = StateSettings },
		},
		{
			Rect:   image.Rect(820, 700, 1010, 765), // Quit
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
		if g.inBattle && g.battle != nil {
			// Affiche uniquement le combat
			g.battle.Draw(screen)
		} else {
			// Affiche map, joueur et ennemis hors combat
			if g.mapData != nil {
				g.mapData.Draw(screen)
			}
			if g.player != nil {
				g.player.Draw(screen)
			}
			for _, e := range g.enemies {
				e.Draw(screen)
			}

			// Dessine la zone de combat en rouge
			red := color.RGBA{255, 0, 0, 100} // semi-transparent
			ebitenutil.DrawRect(screen,
				float64(g.combatZone.Min.X),
				float64(g.combatZone.Min.Y),
				float64(g.combatZone.Dx()),
				float64(g.combatZone.Dy()),
				red,
			)

			// Notification si le joueur est dans la zone
			playerRect := image.Rect(int(g.player.X), int(g.player.Y), int(g.player.X)+32, int(g.player.Y)+32)
			if playerRect.Overlaps(g.combatZone) && !g.inBattle {
				ebitenutil.DebugPrintAt(screen, "Appuie sur E pour lancer un combat !", 200, 180)
			}
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
	if g.menuBg != nil {
		opts := &ebiten.DrawImageOptions{}
		screen.DrawImage(g.menuBg, opts)
	} else {
		screen.Fill(color.RGBA{30, 30, 30, 255})
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

	// Zone de combat
	g.combatZone = image.Rect(200, 200, 300, 300)
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

	// Vérifie si le joueur est dans la zone de combat
	playerRect := image.Rect(int(g.player.X), int(g.player.Y), int(g.player.X)+32, int(g.player.Y)+32)
	if playerRect.Overlaps(g.combatZone) {
		if ebiten.IsKeyPressed(ebiten.KeyE) && !g.inBattle {
			// Commence le combat
			g.inBattle = true
			if len(g.enemies) > 0 {
				g.currentEnemy = g.enemies[0]
				g.battle = NewBattle(g.player, g.currentEnemy)
				g.battle.active = true // ⚠️ Important pour que le combat s'affiche
			}
		}
	}
}
