package game

import "github.com/hajimehoshi/ebiten/v2"

// Définition de la structure Enemy
type Enemy struct {
	X, Y   float64       // Position de l'ennemi sur l'écran (coordonnées X et Y)
	Name   string        // Nom de l'ennemi
	Ego    int           // Niveau d'égo (points de vie) de l'ennemi
	sprite *ebiten.Image // Image représentant l'ennemi
}

// Constructeur pour créer un nouvel ennemi
func NewEnemy(x, y float64, name string) *Enemy {
	return &Enemy{
		X:      x,                                  // Initialise la position X
		Y:      y,                                  // Initialise la position Y
		Name:   name,                               // Initialise le nom
		Ego:    100,                                // Initialise l'égo par défaut à 100
		sprite: LoadImage("assets/enemy_idle.png"), // Charge l'image de l'ennemi
	}
}

// Fonction pour dessiner l'ennemi à l'écran
func (e *Enemy) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{} // Crée une structure d'options de dessin
	opts.GeoM.Translate(e.X, e.Y)      // Positionne l'image aux coordonnées X, Y
	screen.DrawImage(e.sprite, opts)   // Dessine l'image sur l'écran
}
