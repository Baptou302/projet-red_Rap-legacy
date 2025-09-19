package game

import "github.com/hajimehoshi/ebiten/v2"

// Définition de la structure Player
type Player struct {
	X, Y                  float64       // Position du joueur sur l'axe X et Y
	Ego                   int           // Valeur d'ego du joueur (sa "vie" ou énergie)
	Flow                  int           // Niveau de flow
	Charisma              int           // Charisme du joueur
	BonusEgo              int           // Bonus temporaire d'ego pour le prochain combat
	PendingEnemyEgoDebuff int           // Malus d'ego appliqué à l'ennemi lors du prochain combat
	sprite                *ebiten.Image // Image représentant le joueur
	class                 string        // Classe ou type de joueur
}

// Coordonnées fixes du joueur au spawn
const (
	SpawnX = 100
	SpawnY = 150
)

// Constructeur pour créer un nouveau joueur
func NewPlayer(x, y float64, class string) *Player {
	if x == 0 && y == 0 { // valeurs par défaut si 0
		x = SpawnX
		y = SpawnY
	}
	return &Player{
		X:        x,
		Y:        y,
		Ego:      100,
		Flow:     10,
		Charisma: 5,
		class:    class,
		sprite:   LoadImage("assets/player_idle.png"),
	}
}

// Fonction pour réinitialiser la position du joueur
func (p *Player) ResetPosition() {
	p.X = SpawnX
	p.Y = SpawnY
}

// Fonction Update pour gérer les déplacements du joueur
func (p *Player) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyW) { // Monter
		p.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) { // Descendre
		p.Y += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) { // Aller à gauche
		p.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) { // Aller à droite
		p.X += 2
	}
}

// Fonction Draw pour dessiner le joueur à l'écran
func (p *Player) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(p.X, p.Y)
	screen.DrawImage(p.sprite, opts)
}
