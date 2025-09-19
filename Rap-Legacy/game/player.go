package game

import "github.com/hajimehoshi/ebiten/v2"

// Définition de la structure Player
type Player struct {
	X, Y                  float64       // Position du joueur sur l'axe X et Y
	Ego                   int           // Valeur d'ego du joueur (sa "vie" ou énergie)
	Flow                  int           // Niveau de flow (peut influencer certaines actions)
	Charisma              int           // Charisme du joueur
	BonusEgo              int           // Bonus temporaire d'ego pour le prochain combat
	PendingEnemyEgoDebuff int           // Malus d'ego appliqué à l'ennemi lors du prochain combat
	sprite                *ebiten.Image // Image représentant le joueur
	class                 string        // Classe ou type de joueur (ex: "rappeur")
}

// Constructeur pour créer un nouveau joueur
func NewPlayer(x, y float64, class string) *Player {
	return &Player{
		X: x, Y: y, // Position initiale du joueur
		Ego: 100, Flow: 10, Charisma: 5, // Stats initiales
		class:  class,                               // Sauvegarde la classe du joueur
		sprite: LoadImage("assets/player_idle.png"), // Charge le sprite du joueur
	}
}

// Fonction Update pour gérer les déplacements du joueur
func (p *Player) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyW) { // Si la touche W est pressée → monter
		p.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) { // Si la touche S est pressée → descendre
		p.Y += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) { // Si la touche A est pressée → aller à gauche
		p.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) { // Si la touche D est pressée → aller à droite
		p.X += 2
	}
}

// Fonction Draw pour dessiner le joueur à l'écran
func (p *Player) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{} // Crée des options de dessin
	opts.GeoM.Translate(p.X, p.Y)      // Positionne le sprite aux coordonnées X, Y du joueur
	screen.DrawImage(p.sprite, opts)   // Dessine le sprite sur l'écran
}
