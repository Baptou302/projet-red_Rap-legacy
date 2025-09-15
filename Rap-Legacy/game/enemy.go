package game

import "github.com/hajimehoshi/ebiten/v2"

type Enemy struct {
	X, Y   float64
	Name   string
	Ego    int
	sprite *ebiten.Image
}

func NewEnemy(x, y float64, name string) *Enemy {
	return &Enemy{
		X: x, Y: y, Name: name, Ego: 50,
		sprite: LoadImage("assets/enemy_idle.png"),
	}
}

func (e *Enemy) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(e.X, e.Y)
	screen.DrawImage(e.sprite, opts)
}
