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
// NewInventaire (ancienne méthode — garde un contenu par défaut)
// -----------------
func NewInventaire() *Inventaire {
	// Liste des objets par défaut (ton ancien mapping)
	objets := map[string]string{
		"Téléphone":                 "assets/téléphone.png",
		"RandM - 9000K":             "assets/puff.png",
		"Cristalline - mystérieuse": "assets/cristalline.png",
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
// NewInventaireFromItems (construit un inventaire à partir d'une slice d'items)
// -----------------
func NewInventaireFromItems(items []string) *Inventaire {
	icons := map[string]*ebiten.Image{}

	// mapping item -> asset path (si tu as des noms différents adapte ici)
	paths := map[string]string{
		"Micro":                     "assets/micro.png",
		"Cigarette électronique":    "assets/puff.png",
		"Cigarette Electronique":    "assets/puff.png", // variantes
		"Cristalline - mystérieuse": "assets/cristalline.png",
		"Cristalline - tonic":       "assets/cristalline_tonic.png",
		"Cristalline - suspicieuse": "assets/cristalline_suspicieuse.png",
		"Téléphone":                 "assets/téléphone.png",
		"RandM - 9000K":             "assets/puff.png",
	}

	for _, item := range items {
		if path, ok := paths[item]; ok {
			icons[item] = LoadImage(path)
		} else {
			// si pas de path défini, on peut laisser nil (pas d'icône)
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
	// tente de charger icône si existe un mapping (reuse NewInventaireFromItems mapping)
	paths := map[string]string{
		"Micro":                     "assets/micro.png",
		"Cigarette électronique":    "assets/puff.png",
		"Cigarette Electronique":    "assets/puff.png", // variantes
		"Cristalline - mystérieuse": "assets/cristalline.png",
		"Cristalline - tonic":       "assets/cristalline_tonic.png",
		"Cristalline - suspicieuse": "assets/cristalline_suspicieuse.png",
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
func (inv *Inventaire) Update() {
	// Rien pour l'instant
}

// -----------------
// Draw
// -----------------
func (inv *Inventaire) Draw(screen *ebiten.Image) {
	if !inv.Open {
		return
	}

	if inv.Bg != nil {
		opts := &ebiten.DrawImageOptions{}
		screen.DrawImage(inv.Bg, opts)
	} else {
		screen.Fill(color.RGBA{0, 0, 0, 200})
	}

	// Paramètres
	lineHeight := 100 // espacement vertical
	iconScale := 0.08 // taille des icônes
	iconSpacing := 10 // espace entre texte et icône

	// Calculer largeur max du texte pour positionner les icônes
	maxTextWidth := 0
	for _, item := range inv.Items {
		w := text.BoundString(basicfont.Face7x13, item).Dx()
		if w > maxTextWidth {
			maxTextWidth = w
		}
	}

	// Calculer largeur max des icônes
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

	// Centrer horizontalement et verticalement
	screenWidth, screenHeight := screen.Size()
	startX := (screenWidth - totalWidth) / 2
	startY := (screenHeight - lineHeight*len(inv.Items)) / 2

	// Colonne fixe pour les icônes
	iconColumnX := startX + maxTextWidth + iconSpacing

	for i, item := range inv.Items {
		textX := startX
		textY := startY + i*lineHeight

		// Afficher le texte
		ebitenutil.DebugPrintAt(screen, item, textX, textY)

		// Afficher l'icône à la colonne fixe
		if icon, ok := inv.Icons[item]; ok && icon != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(iconScale, iconScale)

			_, ih := icon.Size()
			iconY := float64(textY) + (float64(lineHeight)-float64(ih)*iconScale)/2

			op.GeoM.Translate(float64(iconColumnX), iconY)
			screen.DrawImage(icon, op)
		}
	}
}

// -----------------
// DrawNote (toujours affichée en haut à gauche)
// -----------------
func (inv *Inventaire) DrawNote(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "Appuie sur [TAB] pour ouvrir la FAUSSE sacoche Gucci", 20, 20)
}
