package game

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
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
	settingsBg   *ebiten.Image
	menuSelected int
	volume       int

	// Sauvegardes
	saves         []Save
	saveSelected  int
	newSaveName   string
	newSaveClass  string
	cursorTimer   int
	pendingDelete string // confirmation suppression

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

	// Polices
	fontSmall font.Face
	fontBig   font.Face
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
// Helper : touche "just pressed"
// -----------------
var prevKeyState = make(map[ebiten.Key]bool)

func IsKeyJustPressed(key ebiten.Key) bool {
	pressed := ebiten.IsKeyPressed(key)
	was := prevKeyState[key]
	prevKeyState[key] = pressed
	return pressed && !was
}

// -----------------
// NewGame
// -----------------
func NewGame() *Game {
	g := &Game{
		state:  StateMenu,
		volume: 50,
	}
	// background menu
	g.menuBg = LoadImage("assets/menu_bg.png")

	// background paramètres
	g.settingsBg = LoadImage("assets/image3.png")

	// boutons menu principal
	g.menuButtons = []*Button{
		{
			Rect:   image.Rect(685, 490, 1160, 550),
			Label:  "New Game",
			Action: func() { g.openSaveSelect() },
		},
		{
			Rect:   image.Rect(730, 600, 1110, 660),
			Label:  "Options",
			Action: func() { g.state = StateSettings },
		},
		{
			Rect:   image.Rect(820, 700, 1010, 765),
			Label:  "Quit",
			Action: func() { os.Exit(0) },
		},
	}

	// Charger la police
	ttfBytes, err := os.ReadFile("assets/PressStart2P.ttf")
	if err != nil {
		panic("Impossible de charger la police: " + err.Error())
	}
	tt, err := opentype.Parse(ttfBytes)
	if err != nil {
		panic("Impossible de parser la police: " + err.Error())
	}

	g.fontSmall, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}

	g.fontBig, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    36,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
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

			// Zone combat
			red := color.RGBA{255, 0, 0, 100}
			ebitenutil.DrawRect(screen,
				float64(g.combatZone.Min.X),
				float64(g.combatZone.Min.Y),
				float64(g.combatZone.Dx()),
				float64(g.combatZone.Dy()),
				red,
			)

			// Message combat
			if g.player != nil {
				playerRect := image.Rect(int(g.player.X), int(g.player.Y), int(g.player.X)+32, int(g.player.Y)+32)
				if playerRect.Overlaps(g.combatZone) && !g.inBattle {
					text.Draw(screen, "Appuie sur E pour lancer un combat !", g.fontSmall, 200, 180, color.White)
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
		if i == g.menuSelected {
			// flèche à gauche du bouton sélectionné
			x := btn.Rect.Min.X - 60
			y := (btn.Rect.Min.Y+btn.Rect.Max.Y)/2 + 10
			text.Draw(screen, "▶", g.fontBig, x, y, color.RGBA{255, 255, 0, 255})
		}
	}
}

// -----------------
// Settings
// -----------------
func (g *Game) updateSettings() {
	if IsKeyJustPressed(ebiten.KeyRight) || IsKeyJustPressed(ebiten.KeyUp) {
		g.volume += 5
		if g.volume > 100 {
			g.volume = 100
		}
	}
	if IsKeyJustPressed(ebiten.KeyLeft) || IsKeyJustPressed(ebiten.KeyDown) {
		g.volume -= 5
		if g.volume < 0 {
			g.volume = 0
		}
	}
	if IsKeyJustPressed(ebiten.KeyEscape) {
		g.state = StateMenu
	}
}

func (g *Game) drawSettings(screen *ebiten.Image) {
	if g.settingsBg != nil {
		opts := &ebiten.DrawImageOptions{}
		screen.DrawImage(g.settingsBg, opts)
	} else {
		screen.Fill(color.RGBA{0, 0, 0, 255})
	}

	w, h := screen.Size()

	title := "SETTINGS"
	vol := fmt.Sprintf("Volume: %d", g.volume)
	info := "Press ESC to return"

	// centrer les textes
	text.Draw(screen, title, g.fontBig, w/2-(len(title)*18), h/2-100, color.White)
	text.Draw(screen, vol, g.fontBig, w/2-(len(vol)*18), h/2, color.White)
	text.Draw(screen, info, g.fontSmall, w/2-(len(info)*9), h/2+100, color.RGBA{200, 200, 200, 255})
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
	g.pendingDelete = ""
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

	// Suppression avec confirmation
	if g.saveSelected < len(g.saves) {
		if IsKeyJustPressed(ebiten.KeyDelete) {
			s := g.saves[g.saveSelected]
			if g.pendingDelete == s.Name {
				if err := DeleteSave(s.Name); err != nil {
					fmt.Println("Erreur suppression save:", err)
				} else {
					all, _ := LoadAllSaves()
					g.saves = all
					if g.saveSelected >= len(g.saves) {
						g.saveSelected = len(g.saves) - 1
						if g.saveSelected < 0 {
							g.saveSelected = 0
						}
					}
				}
				g.pendingDelete = ""
			} else {
				g.pendingDelete = s.Name
			}
		}
	} else {
		g.pendingDelete = ""
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

	title := "Sélectionne une sauvegarde :"
	text.Draw(screen, title, g.fontBig, 600, 200, color.White)

	for i, s := range g.saves {
		line := fmt.Sprintf("%s (%s)", s.Name, s.Class)
		if i == g.saveSelected {
			if g.pendingDelete == s.Name {
				line = "> " + line + "   (Appuie encore sur Suppr pour CONFIRMER)"
			} else {
				line = "> " + line + "   (Suppr = supprimer)"
			}
		}
		text.Draw(screen, line, g.fontSmall, 600, 260+i*40, color.White)
	}

	newSaveText := "Créer une nouvelle sauvegarde"
	if g.saveSelected == len(g.saves) {
		newSaveText = "> " + newSaveText
	}
	text.Draw(screen, newSaveText, g.fontSmall, 600, 260+len(g.saves)*40, color.White)

	info := "Appuie sur ESC pour revenir"
	text.Draw(screen, info, g.fontSmall, 600, 700, color.RGBA{200, 200, 200, 255})
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

	title := "Création d'une nouvelle sauvegarde"
	text.Draw(screen, title, g.fontBig, 600, 200, color.White)

	prompt := "Nom du personnage :"
	text.Draw(screen, prompt, g.fontSmall, 600, 260, color.White)

	cursor := "_"
	if (g.cursorTimer/30)%2 == 0 {
		cursor = " "
	}
	text.Draw(screen, g.newSaveName+cursor, g.fontSmall, 600, 300, color.White)

	text.Draw(screen, "Choisis une classe :", g.fontSmall, 600, 360, color.White)
	text.Draw(screen, "[1] Lyricistes", g.fontSmall, 600, 400, color.White)
	text.Draw(screen, "[2] Performeurs", g.fontSmall, 600, 440, color.White)
	text.Draw(screen, "[3] Hitmakers", g.fontSmall, 600, 480, color.White)

	if g.newSaveClass != "" {
		text.Draw(screen, "Classe choisie: "+g.newSaveClass, g.fontSmall, 600, 540, color.White)
		text.Draw(screen, "Appuie sur Entrée pour valider", g.fontSmall, 600, 580, color.White)
	}

	text.Draw(screen, "Appuie sur ECHAP pour annuler", g.fontSmall, 600, 640, color.White)
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
			g.Inventaire.Open = !g.Inventaire.Open
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
