package game

import (
	"image"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Animation struct {
	spriteSheet *ebiten.Image
	frameWidth  int
	frameHeight int
	frameCount  int
	current     int
	lastUpdate  time.Time
	frameDelay  time.Duration
	nbLignes    int
	nbColonnes  int
}

// NewAnimation crée une animation avec sprites répartis sur plusieurs lignes et colonnes
func NewAnimation(img *ebiten.Image, nbLignes, nbColonnes int, frameDelay time.Duration) *Animation {
	w, h := img.Size()
	frameWidth := w / nbColonnes
	frameHeight := h / nbLignes
	frameCount := nbLignes * nbColonnes

	return &Animation{
		spriteSheet: img,
		frameWidth:  frameWidth,
		frameHeight: frameHeight,
		frameCount:  frameCount,
		current:     0,
		lastUpdate:  time.Now(),
		frameDelay:  frameDelay,
		nbLignes:    nbLignes,
		nbColonnes:  nbColonnes,
	}
}

// Update avance l'animation
func (a *Animation) Update() {
	if time.Since(a.lastUpdate) > a.frameDelay {
		a.current = (a.current + 1) % a.frameCount
		a.lastUpdate = time.Now()
	}
}

// Draw affiche l'animation sur l'écran à la position x,y
func (a *Animation) Draw(screen *ebiten.Image, x, y float64) {
	col := a.current % a.nbColonnes
	row := a.current / a.nbColonnes

	sx := col * a.frameWidth
	sy := row * a.frameHeight

	rect := image.Rect(sx, sy, sx+a.frameWidth, sy+a.frameHeight)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(x, y)
	screen.DrawImage(a.spriteSheet.SubImage(rect).(*ebiten.Image), opts)
}
