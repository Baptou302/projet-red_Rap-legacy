package game

import "github.com/hajimehoshi/ebiten/v2"

type Player struct {
	X, Y                  float64
	Ego                   int
	Flow                  int
	Charisma              int
	BonusEgo              int // Bonus d'ego pour le prochain combat
	PendingEnemyEgoDebuff int // Malus d'ego appliqué à l'ennemi au prochain combat
	sprite                *ebiten.Image
	class                 string
}

func NewPlayer(x, y float64, class string) *Player {
	return &Player{
		X: x, Y: y,
		Ego: 100, Flow: 10, Charisma: 5,
		class:  class, // 👈 Sauvegarde la classe
		sprite: LoadImage("assets/player_idle.png"),
	}
}

func (p *Player) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.Y += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.X += 2
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(p.X, p.Y)
	screen.DrawImage(p.sprite, opts)
}
