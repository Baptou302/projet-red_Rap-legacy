package main

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Battle struct {
	player         *Player
	enemy          *Enemy
	turn           int // 0 = joueur, 1 = ennemi
	over           bool
	selectedAttack int // 0=Punchline,1=Flow,2=Diss
}

func NewBattle(p *Player, e *Enemy) *Battle {
	return &Battle{
		player:         p,
		enemy:          e,
		turn:           0,
		over:           false,
		selectedAttack: 0,
	}
}

func (b *Battle) Update() {
	if b.over {
		return
	}

	if b.turn == 0 { // Tour du joueur
		// Changement de sélection avec flèches
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			b.selectedAttack--
			if b.selectedAttack < 0 {
				b.selectedAttack = 1
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			b.selectedAttack++
			if b.selectedAttack > 1 {
				b.selectedAttack = 0
			}
		}

		// Validation avec Enter
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			b.PlayerAttack()
			b.turn = 1
		}
	} else { // Tour de l'ennemi
		b.EnemyAttack()
		b.turn = 0
	}

	// Vérification de fin
	if b.player.Ego <= 0 || b.enemy.Ego <= 0 {
		b.over = true
	}
}

func (b *Battle) PlayerAttack() {
	switch b.selectedAttack {
	case 0: // Punchline
		b.enemy.Ego -= b.player.Flow
	case 1: // Flow
		b.enemy.Ego -= b.player.Flow / 2
		b.player.Flow += 1 // boost temporaire
	case 2: // Diss Track
		b.enemy.Ego -= b.player.Flow * 2
		b.player.Charisma -= 1
	}
}

func (b *Battle) EnemyAttack() {
	b.player.Ego -= 5
}

func (b *Battle) Draw(screen *ebiten.Image) {
	// Fond combat
	ebitenutil.DrawRect(screen, 0, 0, 640, 480, color.RGBA{50, 50, 50, 255})

	// Stats
	ebitenutil.DebugPrintAt(screen, "Votre égo: "+strconv.Itoa(b.player.Ego), 10, 10)
	ebitenutil.DebugPrintAt(screen, "égo rapeur adverse"+strconv.Itoa(b.enemy.Ego), 500, 10)

	// Menu d'attaques
	attacks := []string{"Punchline", "Flow", "Diss Track"}
	for i, a := range attacks {
		text := a
		if i == b.selectedAttack {
			text = "> " + a
		}
		ebitenutil.DebugPrintAt(screen, text, 50, 350+i*20)
	}

	if b.over {
		if b.player.Ego <= 0 {
			ebitenutil.DebugPrintAt(screen, "Tu vas te prendre une sauce sur les réseaux !", 250, 200)
		} else {
			ebitenutil.DebugPrintAt(screen, "Tu vas avoir un gros buzz !", 250, 200)
		}
	}
}

func (b *Battle) IsOver() bool {
	return b.over
}
