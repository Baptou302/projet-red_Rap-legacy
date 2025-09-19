package game

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	StateIntro GameState = iota
	StateMenu
	StateSettings
	StateSaveSelect
	StateCreateSave
	StatePlaying
	StateMerchantMenu
)

// -----------------
// Merchant Items
// -----------------
type MerchantItem struct {
	Name   string
	Price  int
	Action func(g *Game)
}

// Liste des articles vendus par le marchand
var MerchantItems = []MerchantItem{
	{
		Name:  "Cristalline",
		Price: 5,
		Action: func(g *Game) {
			if g.player != nil {
				switch g.player.class {
				case "Lyricistes":
					g.Inventaire.AddItem("Cristalline - mystérieuse")
				case "Performeurs":
					g.Inventaire.AddItem("Cristalline - tonic")
				case "Hitmakers":
					g.Inventaire.AddItem("Cristalline - suspicieuse")
				default:
					g.Inventaire.AddItem("Cristalline")
				}
			}
		},
	},

	{
		Name:  "Followers",
		Price: 20,
		Action: func(g *Game) {
			g.Followers++
		},
	},
}

// -----------------
// Notification system
// -----------------
type Notification struct {
	Text      string
	ExpiresAt time.Time
}

var notifications []Notification

func AddNotification(msg string) {
	notifications = append(notifications, Notification{
		Text:      msg,
		ExpiresAt: time.Now().Add(3 * time.Second), // 3 secondes
	})
}

func UpdateNotifications() {
	now := time.Now()
	active := []Notification{}
	for _, n := range notifications {
		if n.ExpiresAt.After(now) {
			active = append(active, n)
		}
	}
	notifications = active
}

func DrawNotifications(screen *ebiten.Image, fontFace font.Face) {
	for i, n := range notifications {
		text.Draw(screen, n.Text, fontFace, 20, 40+i*30, color.RGBA{255, 255, 0, 255})
	}
}

// -----------------
// Game structure
// -----------------
type Game struct {
	state        GameState
	menuButtons  []*Button
	menuBg       *ebiten.Image
	settingsBg   *ebiten.Image
	audioContext *audio.Context
	bgmPlayer    *audio.Player
	menuSelected int
	volume       int
	moneyIcon    *ebiten.Image   // ✅ icône argent
	followerIcon *ebiten.Image   // ✅ icône followers
	MerchantZone image.Rectangle // Zone interaction marchand
	// Intro
	introTimer int
	introText  string

	// Sauvegardes
	saves         []Save
	saveSelected  int
	newSaveName   string
	newSaveClass  string
	cursorTimer   int
	pendingDelete string // confirmation suppression

	// Gameplay
	player                *Player
	mapData               *Map
	enemies               []*Enemy
	inBattle              bool
	currentEnemy          *Enemy
	battle                *Battle
	Merchant              *Merchant
	Money                 int
	Followers             int
	SelectedMerchantIndex int

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
		state:      StateIntro, // commence par l'intro
		volume:     50,
		introText:  "Bienvenue dans Rap Legacy !",
		introTimer: 0,
	}
	g.moneyIcon = LoadImage("assets/money.png")
	g.followerIcon = LoadImage("assets/followers.png")

	// Background menu
	g.menuBg = LoadImage("assets/menu_bg.png")

	// Background paramètres
	g.settingsBg = LoadImage("assets/image3.png")

	// Boutons menu principal
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

	// Charger la police externe
	ttfBytes, err := os.ReadFile("assets/PressStart2P.ttf")
	if err == nil {
		tt, err := opentype.Parse(ttfBytes)
		if err == nil {
			g.fontSmall, err = opentype.NewFace(tt, &opentype.FaceOptions{
				Size:    24,
				DPI:     72,
				Hinting: font.HintingFull,
			})
			if err != nil {
				g.fontSmall = nil
			}
			g.fontBig, err = opentype.NewFace(tt, &opentype.FaceOptions{
				Size:    36,
				DPI:     72,
				Hinting: font.HintingFull,
			})
			if err != nil {
				g.fontBig = nil
			}
		}
	}

	// Initialiser l'audio
	g.audioContext = audio.NewContext(44100) // fréquence 44.1 kHz

	// Charger le MP3
	mp3Data, err := os.ReadFile("menu/menu.mp3")
	if err != nil {
		log.Println("Impossible de charger la musique :", err)
	} else {
		d, err := mp3.Decode(g.audioContext, bytes.NewReader(mp3Data))
		if err != nil {
			log.Println("Erreur decode mp3 :", err)
		} else {
			g.bgmPlayer, err = audio.NewPlayer(g.audioContext, d)
			if err != nil {
				log.Println("Erreur création player :", err)
			} else {
				g.bgmPlayer.SetVolume(0.5) // volume 0.0 à 1.0
				g.bgmPlayer.Play()         // lancer la musique
			}
		}
	} // ✅ Ajout du marchand
	g.Merchant = &Merchant{
		X:      500,
		Y:      750,
		Items:  []string{"Cristaline", "Followers"},
		Active: true,
		Sprite: LoadImage("assets/marchand.png"),
	}
	g.MerchantZone = image.Rect(500, 700, 500+128, 750+128)
	g.Followers = 0
	g.Money = 100 // ton joueur commence avec 100 pièces
	g.Inventaire = &Inventaire{
		Items: []string{},
	}
	return g
}

// Marchand représente un vendeur avec quelques items
type Merchant struct {
	X, Y   float64
	Items  []string
	Active bool
	Sprite *ebiten.Image
}

// Affiche le marchand (sprite PNG)
func (m *Merchant) Draw(screen *ebiten.Image) {
	if !m.Active || m.Sprite == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(m.X, m.Y)
	screen.DrawImage(m.Sprite, op)
}

// -----------------
// Ebiten methods
// -----------------
func (g *Game) Update() error {
	UpdateNotifications()
	// === Mise à jour de la zone du marchand pour coller au sprite ===
	if g.Merchant != nil && g.Merchant.Sprite != nil {
		w, h := g.Merchant.Sprite.Size() // largeur et hauteur du PNG du marchand
		g.MerchantZone = image.Rect(
			int(g.Merchant.X),
			int(g.Merchant.Y),
			int(g.Merchant.X)+w,
			int(g.Merchant.Y)+h,
		)
	}

	switch g.state {
	case StateIntro:
		g.updateIntro()
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
	case StateIntro:
		g.drawIntro(screen)

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
			// ✅ Ajout : dessiner le marchand
			if g.Merchant != nil {
				g.Merchant.Draw(screen)
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
					if g.fontSmall != nil {
						text.Draw(screen, "Appuie sur E pour lancer un combat !", g.fontSmall, 200, 180, color.White)
					} else {
						ebitenutil.DebugPrintAt(screen, "Appuie sur E pour lancer un combat !", 200, 180)
					}
				}
			}
			// Message marchand
			if g.player != nil && g.Merchant != nil {
				playerRect := image.Rect(
					int(g.player.X),
					int(g.player.Y),
					int(g.player.X)+32,
					int(g.player.Y)+32,
				)
				if g.MerchantZone.Overlaps(playerRect) {
					if g.fontSmall != nil {
						text.Draw(
							screen,
							"Appuie sur E pour parler au marchand",
							g.fontSmall,
							int(g.Merchant.X)-40, // position horizontale
							int(g.Merchant.Y)-20, // position verticale (au-dessus du marchand)
							color.White,
						)
					} else {
						ebitenutil.DebugPrintAt(
							screen,
							"Appuie sur E pour parler au marchand",
							int(g.Merchant.X)-40,
							int(g.Merchant.Y)-20,
						)
					}
				}
			}

			// Inventaire
			if g.Inventaire != nil {
				g.Inventaire.DrawNote(screen)
				g.Inventaire.Draw(screen)

				if g.Inventaire.Open && len(g.Inventaire.Items) > 0 {
					sel := g.Inventaire.selected
					lineHeight := 120
					screenW, screenH := screen.Size()
					startY := (screenH - lineHeight*len(g.Inventaire.Items)) / 2
					y := startY + sel*lineHeight + 8
					textX := screenW/2 - 160
					fnt := g.fontSmall
					if fnt == nil && g.fontBig != nil {
						fnt = g.fontBig
					}
					if fnt != nil {
						text.Draw(screen, "▶", fnt, textX-24, y, color.RGBA{255, 255, 0, 255})
					} else {
						ebitenutil.DebugPrintAt(screen, ">", textX-24, y-6)
					}
				}
			}
		}
		// ✅ HUD Argent & Followers (uniquement dans la map)
		screenW, _ := screen.Size()

		opMoney := &ebiten.DrawImageOptions{}
		opFollower := &ebiten.DrawImageOptions{}

		// On place les icônes en haut à droite
		opMoney.GeoM.Translate(float64(screenW-200), 20)    // position pour money.png
		opFollower.GeoM.Translate(float64(screenW-200), 60) // position pour followers.png

		// Dessiner les icônes si elles existent
		if g.moneyIcon != nil {
			screen.DrawImage(g.moneyIcon, opMoney)
		}
		if g.followerIcon != nil {
			screen.DrawImage(g.followerIcon, opFollower)
		}

		// Dessiner les valeurs à côté
		if g.fontSmall != nil {
			text.Draw(screen, fmt.Sprintf("%d", g.Money), g.fontSmall, screenW-160, 45, color.White)
			text.Draw(screen, fmt.Sprintf("%d", g.Followers), g.fontSmall, screenW-160, 85, color.White)
		} else {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", g.Money), screenW-160, 45)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", g.Followers), screenW-160, 85)
		}

	case StateMerchantMenu: // ✅ Menu du marchand
		if g.fontSmall != nil {
			text.Draw(screen, "=== Marchand ===", g.fontSmall, 200, 150, color.White)
			text.Draw(screen, "1. Cristalline - 5$", g.fontSmall, 200, 200, color.White)
			text.Draw(screen, "2. Followers  - 20$", g.fontSmall, 200, 250, color.White)
			text.Draw(screen, "Appuie sur ESC pour quitter", g.fontSmall, 200, 320, color.White)
		} else {
			ebitenutil.DebugPrintAt(screen, "=== Marchand ===", 200, 150)
			ebitenutil.DebugPrintAt(screen, "1. Cristalline - 5$", 200, 200)
			ebitenutil.DebugPrintAt(screen, "2. Followers  - 20$", 200, 250)
			ebitenutil.DebugPrintAt(screen, "Appuie sur ESC pour quitter", 200, 320)
		}

		// Gestion des achats
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			if g.Money >= MerchantItems[0].Price {
				g.Money -= MerchantItems[0].Price
				MerchantItems[0].Action(g)
				AddNotification("Tu as acheté : " + MerchantItems[0].Name)
			} else {
				AddNotification("Pas assez d'argent !")
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.Key2) {
			if g.Money >= MerchantItems[1].Price {
				g.Money -= MerchantItems[1].Price
				MerchantItems[1].Action(g)
				AddNotification("Tu as acheté : " + MerchantItems[1].Name)
			} else {
				AddNotification("Pas assez d'argent !")
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.state = StatePlaying // retourne au jeu
		}

		// HUD argent / followers
		screenW, _ := screen.Size()
		opMoney := &ebiten.DrawImageOptions{}
		opFollower := &ebiten.DrawImageOptions{}
		opMoney.GeoM.Translate(float64(screenW-200), 20)
		opFollower.GeoM.Translate(float64(screenW-200), 60)

		if g.moneyIcon != nil {
			screen.DrawImage(g.moneyIcon, opMoney)
		}
		if g.followerIcon != nil {
			screen.DrawImage(g.followerIcon, opFollower)
		}

		if g.fontSmall != nil {
			text.Draw(screen, fmt.Sprintf("%d", g.Money), g.fontSmall, screenW-160, 45, color.White)
			text.Draw(screen, fmt.Sprintf("%d", g.Followers), g.fontSmall, screenW-160, 85, color.White)
		} else {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", g.Money), screenW-160, 45)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", g.Followers), screenW-160, 85)
		}
	}

	// Notifications (affichées par-dessus tout)
	if g.fontSmall != nil {
		DrawNotifications(screen, g.fontSmall)
	} else {
		for i, n := range notifications {
			ebitenutil.DebugPrintAt(screen, n.Text, 20, 40+i*30)
		}
	}
}

// -----------------
// Intro update/draw
// -----------------
func (g *Game) updateIntro() {
	g.introTimer++
	if g.introTimer > 180 || IsKeyJustPressed(ebiten.KeyEnter) { // 3 secondes
		g.state = StateMenu
	}
}

func (g *Game) drawIntro(screen *ebiten.Image) {
	screen.Fill(color.Black)
	if g.fontBig != nil {
		text.Draw(screen, g.introText, g.fontBig, 600, 400, color.White)
	} else {
		ebitenutil.DebugPrintAt(screen, g.introText, 600, 400)
	}
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
			x := btn.Rect.Min.X - 40
			y := (btn.Rect.Min.Y+btn.Rect.Max.Y)/2 + 10
			if g.fontSmall != nil {
				text.Draw(screen, "▶", g.fontSmall, x, y, color.RGBA{255, 255, 0, 255})
			} else {
				ebitenutil.DebugPrintAt(screen, ">", x, y)
			}
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

	if g.fontBig != nil {
		text.Draw(screen, title, g.fontBig, w/2-(len(title)*18), h/2-100, color.White)
		text.Draw(screen, vol, g.fontBig, w/2-(len(vol)*18), h/2, color.White)
	} else {
		ebitenutil.DebugPrintAt(screen, title, w/2-60, h/2-100)
		ebitenutil.DebugPrintAt(screen, vol, w/2-60, h/2)
	}
	if g.fontSmall != nil {
		text.Draw(screen, info, g.fontSmall, w/2-(len(info)*9), h/2+100, color.RGBA{200, 200, 200, 255})
	} else {
		ebitenutil.DebugPrintAt(screen, info, w/2-80, h/2+100)
	}
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

// ... le reste de ton game.go reste inchangé (updateSaveSelect, drawSaveSelect, updateCreateSave, drawCreateSave, startGameFromSave, updatePlaying)

// Suppression avec confirmation et navigation dans la sélection de sauvegarde
func (g *Game) updateSaveSelect() {
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

	title := "Sélectionne une sauvegarde :"
	if g.fontBig != nil {
		text.Draw(screen, title, g.fontBig, 600, 200, color.White)
	} else {
		ebitenutil.DebugPrintAt(screen, title, 600, 200)
	}

	for i, s := range g.saves {
		line := fmt.Sprintf("%s (%s)", s.Name, s.Class)
		if i == g.saveSelected {
			if g.pendingDelete == s.Name {
				line = "> " + line + "   (Appuie encore sur Suppr pour CONFIRMER)"
			} else {
				line = "> " + line + "   (Suppr = supprimer)"
			}
		}
		if g.fontSmall != nil {
			text.Draw(screen, line, g.fontSmall, 600, 260+i*40, color.White)
		} else {
			ebitenutil.DebugPrintAt(screen, line, 600, 260+i*20)
		}
	}

	newSaveText := "Créer une nouvelle sauvegarde"
	if g.saveSelected == len(g.saves) {
		newSaveText = "> " + newSaveText
	}
	if g.fontSmall != nil {
		text.Draw(screen, newSaveText, g.fontSmall, 600, 260+len(g.saves)*40, color.White)
		text.Draw(screen, "Appuie sur ESC pour revenir", g.fontSmall, 600, 700, color.RGBA{200, 200, 200, 255})
	} else {
		ebitenutil.DebugPrintAt(screen, newSaveText, 600, 260+len(g.saves)*20)
		ebitenutil.DebugPrintAt(screen, "Appuie sur ESC pour revenir", 600, 700)
	}
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
	if g.fontBig != nil {
		text.Draw(screen, title, g.fontBig, 600, 200, color.White)
	} else {
		ebitenutil.DebugPrintAt(screen, title, 600, 200)
	}

	prompt := "Nom du personnage :"
	if g.fontSmall != nil {
		text.Draw(screen, prompt, g.fontSmall, 600, 260, color.White)
	} else {
		ebitenutil.DebugPrintAt(screen, prompt, 600, 260)
	}

	cursor := "_"
	if (g.cursorTimer/30)%2 == 0 {
		cursor = " "
	}
	if g.fontSmall != nil {
		text.Draw(screen, g.newSaveName+cursor, g.fontSmall, 600, 300, color.White)
	} else {
		ebitenutil.DebugPrintAt(screen, g.newSaveName+cursor, 600, 300)
	}

	if g.fontSmall != nil {
		text.Draw(screen, "Choisis une classe :", g.fontSmall, 600, 360, color.White)
		text.Draw(screen, "[1] Lyricistes", g.fontSmall, 600, 400, color.White)
		text.Draw(screen, "[2] Performeurs", g.fontSmall, 600, 440, color.White)
		text.Draw(screen, "[3] Hitmakers", g.fontSmall, 600, 480, color.White)
	} else {
		ebitenutil.DebugPrintAt(screen, "Choisis une classe :", 600, 360)
		ebitenutil.DebugPrintAt(screen, "[1] Lyricistes", 600, 400)
		ebitenutil.DebugPrintAt(screen, "[2] Performeurs", 600, 440)
		ebitenutil.DebugPrintAt(screen, "[3] Hitmakers", 600, 480)
	}

	if g.newSaveClass != "" {
		if g.fontSmall != nil {
			text.Draw(screen, "Classe choisie: "+g.newSaveClass, g.fontSmall, 600, 540, color.White)
			text.Draw(screen, "Appuie sur Entrée pour valider", g.fontSmall, 600, 580, color.White)
		} else {
			ebitenutil.DebugPrintAt(screen, "Classe choisie: "+g.newSaveClass, 600, 540)
		}
	}

	if g.fontSmall != nil {
		text.Draw(screen, "Appuie sur ECHAP pour annuler", g.fontSmall, 600, 640, color.White)
	} else {
		ebitenutil.DebugPrintAt(screen, "Appuie sur ECHAP pour annuler", 600, 640)
	}
}

// -----------------
// Start game from save
// -----------------
func (g *Game) startGameFromSave(s Save) {
	g.player = NewPlayer(s.PlayerX, s.PlayerY, s.Class) // 👈 ajoute s.Class
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
	// Player update (sans param — ta Player.Update n'attend pas mapData)
	if g.player != nil {
		g.player.Update()
	}

	// Ouvrir/fermer inventaire avec TAB
	if IsKeyJustPressed(ebiten.KeyTab) && g.Inventaire != nil {
		g.Inventaire.Open = !g.Inventaire.Open
	}

	// Si inventaire ouvert → navigation + actions
	if g.Inventaire != nil && g.Inventaire.Open {
		// navigation Up/Down
		if IsKeyJustPressed(ebiten.KeyUp) {
			g.Inventaire.selected--
			if g.Inventaire.selected < 0 && len(g.Inventaire.Items) > 0 {
				g.Inventaire.selected = len(g.Inventaire.Items) - 1
			}
		}
		if IsKeyJustPressed(ebiten.KeyDown) {
			g.Inventaire.selected++
			if g.Inventaire.selected >= len(g.Inventaire.Items) && len(g.Inventaire.Items) > 0 {
				g.Inventaire.selected = 0
			}
		}

		// Consommation : Enter
		if IsKeyJustPressed(ebiten.KeyEnter) && len(g.Inventaire.Items) > 0 {
			idx := g.Inventaire.selected
			if idx >= 0 && idx < len(g.Inventaire.Items) {
				item := g.Inventaire.Items[idx]
				switch item {
				case "Cristalline - mystérieuse", "Cristalline - tonic", "Cristalline - suspicieuse":
					if g.player != nil {
						// On stocke un bonus pour le prochain combat
						g.player.BonusEgo += 50
						AddNotification(fmt.Sprintf("%s consommée : +50 Ego au prochain combat", item))
					}
					// supprimer l'item de l'inventaire
					g.Inventaire.Items = append(g.Inventaire.Items[:idx], g.Inventaire.Items[idx+1:]...)
					if g.Inventaire.Icons != nil {
						delete(g.Inventaire.Icons, item)
					}
					// ajuster selected
					if idx >= len(g.Inventaire.Items) {
						g.Inventaire.selected = len(g.Inventaire.Items) - 1
					}
					if g.Inventaire.selected < 0 {
						g.Inventaire.selected = 0
					}
				default:
					// Pour les autres objets, juste message
					AddNotification("Tu as choisi : " + item)
				}
			}
		}
	}

	// Détection entrée zone marchand + E pour ouvrir le menu
	if g.player != nil {
		playerRect := image.Rect(
			int(g.player.X), int(g.player.Y),
			int(g.player.X)+32, int(g.player.Y)+32, // adapte à la taille du sprite joueur
		)

		if g.MerchantZone.Overlaps(playerRect) {

			// Si touche E pressée → ouvre le menu du marchand
			if IsKeyJustPressed(ebiten.KeyE) {
				g.state = StateMerchantMenu
			}
		}
	}

	// Détection entrée zone combat + E pour lancer combat
	if g.player != nil && g.combatZone.Overlaps(image.Rect(int(g.player.X), int(g.player.Y), int(g.player.X)+32, int(g.player.Y)+32)) && !g.inBattle {
		if IsKeyJustPressed(ebiten.KeyE) {
			g.inBattle = true
			if len(g.enemies) > 0 {
				// On passe BonusEgo à NewBattle via g.player
				g.battle = NewBattle(g.player, g.enemies[0])
			}
		}
	}

	// Si en bataille → update combat
	if g.inBattle && g.battle != nil {
		g.battle.Update()
		if g.battle.IsOver() {
			// Reset état après combat
			g.inBattle = false
			g.battle = nil
			g.player.BonusEgo = 0 // le bonus ne s’applique qu’une fois
		}
	}
}

// -----------------
// Layout (obligatoire pour Ebiten)
// -----------------
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 1920, 1080 // taille interne du rendu
}
