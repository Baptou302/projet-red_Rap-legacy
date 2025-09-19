package game // Déclare le package "game" pour regrouper les fichiers liés au jeu

import (
	"image" // Import pour manipuler des rectangles et sous-images
	"time"  // Import pour gérer le temps (durées, horodatage)

	"github.com/hajimehoshi/ebiten/v2" // Ebiten pour le rendu et les images du jeu
)

// Animation représente une animation découpée en plusieurs frames
type Animation struct {
	frames     []*ebiten.Image // Liste des images (frames) de l'animation
	frameCount int             // Nombre total de frames
	current    int             // Index de la frame courante
	frameDelay time.Duration   // Délai entre chaque frame
	lastUpdate time.Time       // Dernier moment où la frame a été changée
}

// NewAnimation crée une nouvelle animation à partir d'une spritesheet
func NewAnimation(sheet *ebiten.Image, rows, cols int, frameDelay time.Duration) *Animation {
	w, h := sheet.Size() // Récupère la largeur et hauteur totale de la spritesheet
	frameW := w / cols   // Largeur d'une seule frame
	frameH := h / rows   // Hauteur d'une seule frame

	var frames []*ebiten.Image  // Initialise le slice pour stocker toutes les frames
	for y := 0; y < rows; y++ { // Parcourt les lignes
		for x := 0; x < cols; x++ { // Parcourt les colonnes
			sx := x * frameW                                                                  // Coordonnée X du coin supérieur gauche de la frame
			sy := y * frameH                                                                  // Coordonnée Y du coin supérieur gauche de la frame
			frame := sheet.SubImage(image.Rect(sx, sy, sx+frameW, sy+frameH)).(*ebiten.Image) // Découpe la frame
			frames = append(frames, frame)                                                    // Ajoute la frame à la liste
		}
	}

	return &Animation{ // Retourne l'objet Animation initialisé
		frames:     frames,      // Toutes les frames découpées
		frameCount: len(frames), // Nombre total de frames
		current:    0,           // Commence à la première frame
		frameDelay: frameDelay,  // Temps entre chaque frame
		lastUpdate: time.Now(),  // Date/heure de la création
	}
}

// Update passe à la frame suivante si le délai entre frames est écoulé
func (a *Animation) Update() {
	if time.Since(a.lastUpdate) >= a.frameDelay { // Vérifie si le délai est dépassé
		a.current++                    // Passe à la frame suivante
		if a.current >= a.frameCount { // Si on dépasse le nombre de frames
			a.current = 0 // Repart de la première frame
		}
		a.lastUpdate = time.Now() // Met à jour le timestamp du dernier changement
	}
}

// Draw affiche la frame courante sur l'écran
// x : position horizontale (centre X du sprite)
// y : position verticale (référence)
// scale : facteur d'agrandissement
// anchorBottom : si true, y correspond au bas du sprite, sinon y est le centre vertical
func (a *Animation) Draw(screen *ebiten.Image, x, y, scale float64, anchorBottom bool) {
	frame := a.frames[a.current] // Récupère la frame actuelle
	w, h := frame.Size()         // Largeur et hauteur de la frame

	op := &ebiten.DrawImageOptions{} // Options de dessin pour Ebiten
	op.GeoM.Scale(scale, scale)      // Applique le facteur d'échelle

	if anchorBottom { // Si le y correspond au bas du sprite
		op.GeoM.Translate(x-float64(w)*scale/2, y-float64(h)*scale) // Centre horizontal et bas aligné
	} else { // Si y correspond au centre
		op.GeoM.Translate(x-float64(w)*scale/2, y-float64(h)*scale/2) // Centre horizontal et vertical
	}

	screen.DrawImage(frame, op) // Dessine la frame sur l'écran
}
