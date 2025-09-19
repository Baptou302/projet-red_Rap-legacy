package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// -----------------
// Inventaire struct
// -----------------
type Inventaire struct {
	Open     bool
	Bg       *ebiten.Image
	Items    []string
	Icons    map[string]*ebiten.Image
	selected int
}

// -----------------
// NewInventaire
// -----------------
func NewInventaire() *Inventaire {
	// Liste des objets par défaut
	objets := map[string]string{
		"Téléphone":                 "assets/téléphone.png",
		"RandM - 9000K":             "assets/puff.png",
		"Cristalline - mystérieuse": "assets/cristalline.png",
		"Cristalline - big":         "assets/cristalline_big.png",
	}

	items := []string{}
	icons := map[string]*ebiten.Image{}

	for name, path := range objets {
		img := LoadImage(path)
		items = append(items, name)
		icons[name] = img
	}

	return &Inventaire{
		Open:  false,
		Bg:    LoadImage("assets/inventaire.png"),
		Items: items,
		Icons: icons,
	}
}

// -----------------
// NewInventaireFromItems
// -----------------
func NewInventaireFromItems(items []string) *Inventaire {
	icons := map[string]*ebiten.Image{}

	paths := map[string]string{
		"Micro":                     "assets/micro.png",
		"Cigarette électronique":    "assets/puff.png",
		"Cigarette Electronique":    "assets/puff.png",
		"Cristalline - mystérieuse": "assets/cristalline.png",
		"Cristalline - tonic":       "assets/cristalline_tonic.png",
		"Cristalline - suspicieuse": "assets/cristalline_suspicieuse.png",
		"Cristalline - big":         "assets/cristalline_big.png",
		"Téléphone":                 "assets/téléphone.png",
		"RandM - 9000K":             "assets/puff.png",
	}

	for _, item := range items {
		if path, ok := paths[item]; ok {
			icons[item] = LoadImage(path)
		} else {
			icons[item] = nil
		}
	}

	return &Inventaire{
		Open:  false,
		Bg:    LoadImage("assets/inventaire.png"),
		Items: items,
		Icons: icons,
	}
}

// -----------------
// AddItem
// -----------------
func (inv *Inventaire) AddItem(item string) {
	inv.Items = append(inv.Items, item)
	paths := map[string]string{
		"Micro":                     "assets/micro.png",
		"Cigarette électronique":    "assets/puff.png",
		"Cigarette Electronique":    "assets/puff.png",
		"Cristalline - mystérieuse": "assets/cristalline.png",
		"Cristalline - tonic":       "assets/cristalline_tonic.png",
		"Cristalline - suspicieuse": "assets/cristalline_suspicieuse.png",
		"Cristalline - big":         "assets/cristalline_big.png",
		"Téléphone":                 "assets/téléphone.png",
		"RandM - 9000K":             "assets/puff.png",
	}

	if inv.Icons == nil {
		inv.Icons = map[string]*ebiten.Image{}
	}
	if path, ok := paths[item]; ok {
		inv.Icons[item] = LoadImage(path)
	}
}

// -----------------
// Update
// -----------------
func (inv *Inventaire) Update(player *Player) {
	if !inv.Open {
		return
	}

	// Navigation
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && inv.selected > 0 {
		inv.selected--
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && inv.selected < len(inv.Items)-1 {
		inv.selected++
	}

	// Consommer l'objet sélectionné
	if ebiten.IsKeyPressed(ebiten.KeyEnter) && len(inv.Items) > 0 {
		selectedItem := inv.Items[inv.selected]

		if selectedItem == "Cristalline - mystérieuse" ||
			selectedItem == "Cristalline - tonic" ||
			selectedItem == "Cristalline - suspicieuse" {

			// ✅ Cristallines normales → +50 ego
			player.BonusEgo = 50
			AddNotification("Tu as bu une Cristalline, ton ego est boosté !")

			// Retirer après usage
			inv.Items = append(inv.Items[:inv.selected], inv.Items[inv.selected+1:]...)

		} else if selectedItem == "Cristalline - big" {
			// ✅ BIG Cristalline → +100 ego
			player.BonusEgo = 100
			AddNotification("Tu as bu une BIG Cristalline ! Ton ego sera boosté de +100 au prochain combat.")

			// Retirer après usage
			inv.Items = append(inv.Items[:inv.selected], inv.Items[inv.selected+1:]...)

		} else if selectedItem == "Micro" {
			// ✅ Micro → +10 ego
			player.BonusEgo = 10
			AddNotification("Tu as utilisé le Micro, ton ego sera boosté de +10 au prochain combat !")

			inv.Items = append(inv.Items[:inv.selected], inv.Items[inv.selected+1:]...)

		} else if selectedItem == "Cigarette électronique" {
			// ✅ Cigarette → -15 ego ennemi
			player.PendingEnemyEgoDebuff = 15
			AddNotification("Tu as utilisé la Cigarette électronique, l'ennemi commencera avec -15 ego !")

			inv.Items = append(inv.Items[:inv.selected], inv.Items[inv.selected+1:]...)

		}

		// ✅ Ajuster la sélection après suppression
		if inv.selected >= len(inv.Items) {
			inv.selected = len(inv.Items) - 1
		}
	}
}

func (inv *Inventaire) Draw(screen *ebiten.Image, g *Game) {
	if !inv.Open {
		return
	}

	// Fond
	if inv.Bg != nil {
		opts := &ebiten.DrawImageOptions{}
		screen.DrawImage(inv.Bg, opts)
	} else {
		screen.Fill(color.RGBA{0, 0, 0, 200})
	}

	// Paramètres
	lineHeight := 100 // espacement vertical entre chaque item
	iconScale := 0.08 // taille des icônes
	iconSpacing := 20 // espace entre texte et icône

	// Police : on utilise PressStart2P
	face := g.fontSmall
	if face == nil {
		face = basicfont.Face7x13 // fallback si jamais la police ne charge pas
	}

	// Calculer largeur max du texte
	maxTextWidth := 0
	for _, item := range inv.Items {
		w := text.BoundString(face, item).Dx()
		if w > maxTextWidth {
			maxTextWidth = w
		}
	}

	// Largeur max des icônes
	maxIconWidth := 0
	for _, icon := range inv.Icons {
		if icon == nil {
			continue
		}
		w, _ := icon.Size()
		if w > maxIconWidth {
			maxIconWidth = w
		}
	}

	totalWidth := maxTextWidth + iconSpacing + int(float64(maxIconWidth)*iconScale)

	// Centrage
	screenWidth, screenHeight := screen.Size()
	startX := (screenWidth - totalWidth) / 2
	startY := (screenHeight - lineHeight*len(inv.Items)) / 2

	iconColumnX := startX + maxTextWidth + iconSpacing

	for i, item := range inv.Items {
		textX := startX
		textY := startY + i*lineHeight

		// --- Texte ---
		text.Draw(screen, item, face, textX, textY, color.White)

		// --- Flèche de sélection jaune ▶ ---
		if i == inv.selected {
			arrow := "▶"
			arrowX := textX - 40 // un peu plus proche
			yellow := color.RGBA{255, 255, 0, 255}
			text.Draw(screen, arrow, face, arrowX, textY, yellow)
		}

		// --- Icône alignée avec le texte ---
		if icon, ok := inv.Icons[item]; ok && icon != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(iconScale, iconScale)

			// On centre l’icône sur la ligne de texte
			_, ih := icon.Size()
			iconY := float64(textY) - float64(ih)*iconScale/2 + float64(face.Metrics().Ascent.Round())/2

			op.GeoM.Translate(float64(iconColumnX), iconY)
			screen.DrawImage(icon, op)
		}
	}
}

// -----------------
// DrawNote
// -----------------
func (inv *Inventaire) DrawNote(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "Appuie sur [TAB] pour ouvrir la FAUSSE sacoche Gucci", 20, 20)
	ebitenutil.DebugPrintAt(screen, "Appuie sur [F] pour ouvrir craft", 20, 40)
}

// -----------------
// Fonctions utilitaires pour l'inventaire
// -----------------

// Vérifie si l'inventaire contient au moins une Cristalline (peu importe laquelle)
func (inv *Inventaire) HasCristalline() bool {
	for _, item := range inv.Items {
		if item == "Cristalline - mystérieuse" ||
			item == "Cristalline - tonic" ||
			item == "Cristalline - suspicieuse" ||
			item == "Cristalline - big" {
			return true
		}
	}
	return false
}

// Vérifie si l'inventaire contient un objet précis
func (inv *Inventaire) HasItem(name string) bool {
	for _, item := range inv.Items {
		if item == name {
			return true
		}
	}
	return false
}

// Retire la première Cristalline simple trouvée (mystérieuse / tonic / suspicieuse)
// Retourne true si une cristalline a été supprimée, false sinon.
func (inv *Inventaire) RemoveCristalline() bool {
	for i, item := range inv.Items {
		if item == "Cristalline - mystérieuse" ||
			item == "Cristalline - tonic" ||
			item == "Cristalline - suspicieuse" {
			// supprimer l'icône si présente
			if inv.Icons != nil {
				delete(inv.Icons, item)
			}
			// retirer de la slice
			inv.Items = append(inv.Items[:i], inv.Items[i+1:]...)
			// ajuster sélection
			if inv.selected >= len(inv.Items) {
				inv.selected = len(inv.Items) - 1
				if inv.selected < 0 {
					inv.selected = 0
				}
			}
			return true
		}
	}
	return false
}

// Retire un objet précis (par nom). Retourne true si supprimé, false sinon.
func (inv *Inventaire) RemoveItem(name string) bool {
	for i, item := range inv.Items {
		if item == name {
			if inv.Icons != nil {
				delete(inv.Icons, item)
			}
			inv.Items = append(inv.Items[:i], inv.Items[i+1:]...)
			if inv.selected >= len(inv.Items) {
				inv.selected = len(inv.Items) - 1
				if inv.selected < 0 {
					inv.selected = 0
				}
			}
			return true
		}
	}
	return false
}
