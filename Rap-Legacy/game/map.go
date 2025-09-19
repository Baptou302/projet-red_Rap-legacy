package game

import "github.com/hajimehoshi/ebiten/v2"

// Définition de la structure Map
type Map struct {
	bgImage *ebiten.Image // Image de fond de la carte
}

// Constructeur pour créer une nouvelle Map
func NewMap() *Map {
	return &Map{
		bgImage: LoadImage("assets/image3.png"), // Charge l'image de fond depuis le dossier assets
	}
}

// Fonction pour dessiner la map à l'écran
func (m *Map) Draw(screen *ebiten.Image) {
	if m.bgImage != nil { // Vérifie si l'image de fond a été chargée
		opts := &ebiten.DrawImageOptions{} // Crée une structure d'options de dessin
		screen.DrawImage(m.bgImage, opts)  // Dessine l'image de fond sur l'écran
	}
}
