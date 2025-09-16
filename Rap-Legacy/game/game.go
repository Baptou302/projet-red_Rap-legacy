package game

import (
	"fmt"
	"image"
	"image/color"
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
	StateSaveSelect
	StateCreateSave
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

	// Sauvegardes
	saves        []Save
	saveSelected int
	newSaveName  string
	newSaveClass string
	cursorTimer  int

	// Gameplay
	player       *Player
	mapData      *Map
	enemies      []*Enemy
	inBattle     bool
	currentEnemy *Enemy
	battle       *Battle

	// Zone de combat
	combatZone image.Rectangle

	// Inventaire
	Inventaire *Inventaire
}

// -----------------
// LoadImage helper
// -----------------
func LoadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		panic(err)
	}
	return img
}

// -----------------
// Input helper (pour éviter le "repeat")
// -----------------
var prevInput = map[ebiten.Key]bool{}

func IsKeyJustPressed(key ebiten.Key) bool {
	pressed := ebiten.IsKeyPressed(key)
	wasPressed := prevInput[key]
	prevInput[key] = pressed
	return pressed && !wasPressed
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

	// Boutons du menu principal
	g.menuButtons = []*Button{
		{
			Rect:   image.Rect(685, 490, 1160, 550),
			Label:  "New Game",
			Action: func() { g.openSaveSelect() },
		},
		{
			Rect:   image.Rect(730, 600, 1110, 660),
			Label:  "Settings",
			Action: func() { g.state = StateSettings },
		},
		{
			Rect:   image.Rect(820, 700, 1010, 765),
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
		g.updateMenu()
	case StateSettings:
		g.updateSettings()
	case StateSaveSelect:
		g.updateSaveSelect()
	case StateCreateSave:
		g.updateCreateSave()
	case StatePlaying:
		g.updatePlaying()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case StateMenu:
		g.drawMenu(screen)
	case StateSettings:
		g.drawSettings(screen)
	case StateSaveSelect:
		g.drawSaveSelect(screen)
	case StateCreateSave:
		g.drawCreateSave(screen)
	case StatePlaying:
		if g.inBattle && g.battle != nil {
			g.battle.Draw(screen)
		} else {
			if g.mapData != nil {
				g.mapData.Draw(screen)
			}
			if g.player != nil {
				g.player.Draw(screen)
			}
			for _, e := range g.enemies {
				e.Draw(screen)
			}

			// Zone combat visuelle
			red := color.RGBA{255, 0, 0, 100}
			ebitenutil.DrawRect(screen,
				float64(g.combatZone.Min.X),
				float64(g.combatZone.Min.Y),
				float64(g.combatZone.Dx()),
				float64(g.combatZone.Dy()),
				red,
			)

			// Message si joueur dans la zone
			if g.player != nil {
				playerRect := image.Rect(int(g.player.X), int(g.player.Y), int(g.player.X)+32, int(g.player.Y)+32)
				if playerRect.Overlaps(g.combatZone) && !g.inBattle {
					ebitenutil.DebugPrintAt(screen, "Appuie sur E pour lancer un combat !", 200, 180)
				}
			}

			// Inventaire
			if g.Inventaire != nil {
				g.Inventaire.DrawNote(screen)
				g.Inventaire.Draw(screen)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1920, 1080
}

// -----------------
// Menu principal
// -----------------
func (g *Game) updateMenu() {
	mx, my := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		for _, btn := range g.menuButtons {
			if mx >= btn.Rect.Min.X && mx <= btn.Rect.Max.X &&
				my >= btn.Rect.Min.Y && my <= btn.Rect.Max.Y {
				if btn.Action != nil {
					btn.Action()
				}
			}
		}
	}

	// navigation clavier
	if IsKeyJustPressed(ebiten.KeyUp) {
		g.menuSelected--
		if g.menuSelected < 0 {
			g.menuSelected = len(g.menuButtons) - 1
		}
	}
	if IsKeyJustPressed(ebiten.KeyDown) {
		g.menuSelected++
		if g.menuSelected >= len(g.menuButtons) {
			g.menuSelected = 0
		}
	}
	if IsKeyJustPressed(ebiten.KeyEnter) {
		if g.menuButtons[g.menuSelected].Action != nil {
			g.menuButtons[g.menuSelected].Action()
		}
	}
}

func (g *Game) drawMenu(screen *ebiten.Image) {
	if g.menuBg != nil {
		opts := &ebiten.DrawImageOptions{}
		screen.DrawImage(g.menuBg, opts)
	} else {
		screen.Fill(color.RGBA{30, 30, 30, 255})
	}
	for i, btn := range g.menuButtons {
		label := btn.Label
		if i == g.menuSelected {
			label = "> " + label
		}
		ebitenutil.DebugPrintAt(screen, label, btn.Rect.Min.X, btn.Rect.Min.Y-20)
	}
}

// -----------------
// Settings
// -----------------
func (g *Game) updateSettings() {
	if IsKeyJustPressed(ebiten.KeyUp) {
		g.volume++
		if g.volume > 10 {
			g.volume = 10
		}
	}
	if IsKeyJustPressed(ebiten.KeyDown) {
		g.volume--
		if g.volume < 0 {
			g.volume = 0
		}
	}
	if IsKeyJustPressed(ebiten.KeyEscape) {
		g.state = StateMenu
	}
}

func (g *Game) drawSettings(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "SETTINGS", 250, 100)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Volume: %d", g.volume), 250, 150)
	ebitenutil.DebugPrintAt(screen, "Press ESC to return", 200, 200)
}

// -----------------
// Save selection
// -----------------
func (g *Game) openSaveSelect() {
	saves, err := LoadAllSaves()
	if err != nil {
		saves = []Save{}
	}
	g.saves = saves
	g.saveSelected = 0
	g.state = StateSaveSelect
}

func (g *Game) updateSaveSelect() {
	if g.saves == nil {
		all, _ := LoadAllSaves()
		g.saves = all
	}

	if IsKeyJustPressed(ebiten.KeyUp) {
		g.saveSelected--
		if g.saveSelected < 0 {
			g.saveSelected = len(g.saves)
		}
	}
	if IsKeyJustPressed(ebiten.KeyDown) {
		g.saveSelected++
		if g.saveSelected > len(g.saves) {
			g.saveSelected = 0
		}
	}

	if IsKeyJustPressed(ebiten.KeyEnter) {
		if g.saveSelected == len(g.saves) {
			g.newSaveName = ""
			g.newSaveClass = ""
			g.cursorTimer = 0
			g.state = StateCreateSave
		} else {
			s := g.saves[g.saveSelected]
			g.startGameFromSave(s)
		}
	}

	if IsKeyJustPressed(ebiten.KeyEscape) {
		g.state = StateMenu
	}
}

func (g *Game) drawSaveSelect(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 60, 255})
	ebitenutil.DebugPrintAt(screen, "Sélectionne une sauvegarde :", 600, 200)

	for i, s := range g.saves {
		text := fmt.Sprintf("%s (%s)", s.Name, s.Class)
		if i == g.saveSelected {
			text = "> " + text
		}
		ebitenutil.DebugPrintAt(screen, text, 600, 260+i*36)
	}

	newSaveText := "Créer une nouvelle sauvegarde"
	if g.saveSelected == len(g.saves) {
		newSaveText = "> " + newSaveText
	}
	ebitenutil.DebugPrintAt(screen, newSaveText, 600, 260+len(g.saves)*36)

	ebitenutil.DebugPrintAt(screen, "Appuie sur ESC pour revenir", 600, 700)
}

// -----------------
// Create Save
// -----------------
func (g *Game) updateCreateSave() {
	for _, r := range ebiten.InputChars() {
		if r == '\n' || r == '\r' {
			continue
		}
		g.newSaveName += string(r)
	}
	if IsKeyJustPressed(ebiten.KeyBackspace) && len(g.newSaveName) > 0 {
		g.newSaveName = g.newSaveName[:len(g.newSaveName)-1]
	}

	if IsKeyJustPressed(ebiten.Key1) {
		g.newSaveClass = "Lyricistes"
	}
	if IsKeyJustPressed(ebiten.Key2) {
		g.newSaveClass = "Performeurs"
	}
	if IsKeyJustPressed(ebiten.Key3) {
		g.newSaveClass = "Hitmakers"
	}

	if IsKeyJustPressed(ebiten.KeyEnter) && g.newSaveClass != "" {
		name := g.newSaveName
		if name == "" {
			name = fmt.Sprintf("Player-%d", len(g.saves)+1)
		}
		s, err := CreateSave(name, g.newSaveClass)
		if err != nil {
			fmt.Println("Erreur création save:", err)
		} else {
			all, _ := LoadAllSaves()
			g.saves = all
			g.startGameFromSave(s)
		}
	}

	if IsKeyJustPressed(ebiten.KeyEscape) {
		g.state = StateSaveSelect
	}

	g.cursorTimer++
}

func (g *Game) drawCreateSave(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 10, 40, 255})
	ebitenutil.DebugPrintAt(screen, "Création d'une nouvelle sauvegarde", 600, 200)
	ebitenutil.DebugPrintAt(screen, "Nom du personnage :", 600, 260)

	cursor := "_"
	if (g.cursorTimer/30)%2 == 0 {
		cursor = " "
	}
	ebitenutil.DebugPrintAt(screen, g.newSaveName+cursor, 600, 300)

	ebitenutil.DebugPrintAt(screen, "Choisis une classe :", 600, 360)
	ebitenutil.DebugPrintAt(screen, "[1] Lyricistes", 600, 400)
	ebitenutil.DebugPrintAt(screen, "[2] Performeurs", 600, 440)
	ebitenutil.DebugPrintAt(screen, "[3] Hitmakers", 600, 480)

	if g.newSaveClass != "" {
		ebitenutil.DebugPrintAt(screen, "Classe choisie: "+g.newSaveClass, 600, 540)
		ebitenutil.DebugPrintAt(screen, "Appuie sur Entrée pour valider", 600, 580)
	}

	ebitenutil.DebugPrintAt(screen, "Appuie sur ECHAP pour annuler", 600, 640)
}

// -----------------
// Start game from save
// -----------------
func (g *Game) startGameFromSave(s Save) {
	g.player = NewPlayer(s.PlayerX, s.PlayerY)
	g.player.Ego = s.Ego
	g.player.Flow = s.Flow
	g.player.Charisma = s.Charisma

	g.Inventaire = NewInventaireFromItems(s.Inventory)

	g.mapData = NewMap()
	g.enemies = []*Enemy{
		NewEnemy(200, 200, "Rival Rapper"),
		NewEnemy(400, 300, "Boss Rapper"),
	}
	g.inBattle = false
	g.combatZone = image.Rect(200, 200, 300, 300)

	g.state = StatePlaying
}

// -----------------
// Playing update
// -----------------
func (g *Game) updatePlaying() {
	if g.Inventaire != nil {
		if IsKeyJustPressed(ebiten.KeyTab) {
			g.Inventaire.Open = true
		}
		if IsKeyJustPressed(ebiten.KeyEscape) {
			g.Inventaire.Open = false
		}
	}

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

	if g.player != nil {
		playerRect := image.Rect(int(g.player.X), int(g.player.Y), int(g.player.X)+32, int(g.player.Y)+32)
		if playerRect.Overlaps(g.combatZone) {
			if IsKeyJustPressed(ebiten.KeyE) && !g.inBattle {
				g.inBattle = true
				g.battle = NewBattle()
			}
		}
	}
}
