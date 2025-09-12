package game

import "github.com/hajimehoshi/ebiten/v2"

type Map struct {
	bgImage *ebiten.Image
}

func NewMap() *Map {
	return &Map{
		bgImage: LoadImage("assets/image3.png"),
	}
}

func (m *Map) Draw(screen *ebiten.Image) {
	if m.bgImage != nil {
		opts := &ebiten.DrawImageOptions{}
		screen.DrawImage(m.bgImage, opts)
	}
}
