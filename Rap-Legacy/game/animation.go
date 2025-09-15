package game

import (
	"image"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Animation struct {
	frames     []*ebiten.Image
	frameCount int
	current    int
	frameDelay time.Duration
	lastUpdate time.Time
}

// Découpe une spritesheet en frames
func NewAnimation(sheet *ebiten.Image, rows, cols int, frameDelay time.Duration) *Animation {
	w, h := sheet.Size()
	frameW := w / cols
	frameH := h / rows

	var frames []*ebiten.Image
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			sx := x * frameW
			sy := y * frameH
			frame := sheet.SubImage(image.Rect(sx, sy, sx+frameW, sy+frameH)).(*ebiten.Image)
			frames = append(frames, frame)
		}
	}

	return &Animation{
		frames:     frames,
		frameCount: len(frames),
		current:    0,
		frameDelay: frameDelay,
		lastUpdate: time.Now(),
	}
}

// Passe à la frame suivante si le délai est écoulé
func (a *Animation) Update() {
	if time.Since(a.lastUpdate) >= a.frameDelay {
		a.current++
		if a.current >= a.frameCount {
			a.current = 0
		}
		a.lastUpdate = time.Now()
	}
}

// Affiche la frame courante
// x, y = position
// scale = facteur d'agrandissement
// anchorBottom = true → colle au sol
func (a *Animation) Draw(screen *ebiten.Image, x, y float64, scale float64, anchorBottom bool) {
	frame := a.frames[a.current]
	w, h := frame.Size()

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)

	if anchorBottom {
		// On centre en X mais colle en bas
		op.GeoM.Translate(x-float64(w)*scale/2, y-float64(h)*scale)
	} else {
		// Ancien comportement (centré milieu)
		op.GeoM.Translate(x-float64(w)*scale/2, y-float64(h)*scale/2)
	}

	screen.DrawImage(frame, op)
}
