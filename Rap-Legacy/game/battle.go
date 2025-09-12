package game

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Battle représente un combat entre le joueur et un ennemi
type Battle struct {
	player         *Player
	enemy          *Enemy
	turn           int // 0 = joueur, 1 = ennemi
	over           bool
	selectedAttack int
}

// NewBattle crée une nouvelle instance de Battle
func NewBattle(p *Player, e *Enemy) *Battle {
	return &Battle{
		player:         p,
		enemy:          e,
		turn:           0,
		over:           false,
		selectedAttack: 0,
	}
}

// Update gère les tours et la logique du combat
func (b *Battle) Update() {
	if b.over {
		return
	}

	if b.turn == 0 {
		// Tour du joueur
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			b.selectedAttack--
			if b.selectedAttack < 0 {
				b.selectedAttack = 2
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			b.selectedAttack++
			if b.selectedAttack > 2 {
				b.selectedAttack = 0
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			b.PlayerAttack()
			b.turn = 1
		}
	} else {
		// Tour de l'ennemi
		b.EnemyAttack()
		b.turn = 0
	}

	// Vérification de fin du combat
	if b.player.Ego <= 0 || b.enemy.Ego <= 0 {
		b.over = true
	}
}

// PlayerAttack applique l'attaque du joueur
func (b *Battle) PlayerAttack() {
	switch b.selectedAttack {
	case 0: // Punchline
		b.enemy.Ego -= b.player.Flow
	case 1: // Flow
		b.enemy.Ego -= b.player.Flow / 2
		b.player.Flow++
	case 2: // Diss Track
		b.enemy.Ego -= b.player.Flow * 2
		b.player.Charisma--
	}
}

// EnemyAttack applique l'attaque de l'ennemi
func (b *Battle) EnemyAttack() {
	b.player.Ego -= 5
}

// Draw affiche le combat à l'écran
func (b *Battle) Draw(screen *ebiten.Image) {
	// Fond combat
	ebitenutil.DrawRect(screen, 0, 0, 640, 480, color.RGBA{50, 50, 50, 255})

	// Stats
	ebitenutil.DebugPrintAt(screen, "Votre égo: "+strconv.Itoa(b.player.Ego), 10, 10)
	ebitenutil.DebugPrintAt(screen, "Égo adverse: "+strconv.Itoa(b.enemy.Ego), 500, 10)

	// Menu d'attaques
	attacks := []string{"Punchline", "Flow", "Diss Track"}
	for i, a := range attacks {
		text := a
		if i == b.selectedAttack {
			text = "> " + a
		}
		ebitenutil.DebugPrintAt(screen, text, 50, 350+i*20)
	}

	// Message de fin
	if b.over {
		if b.player.Ego <= 0 {
			ebitenutil.DebugPrintAt(screen, "Tu vas te prendre une sauce sur les réseaux !", 250, 200)
		} else {
			ebitenutil.DebugPrintAt(screen, "Tu vas avoir un gros buzz !", 250, 200)
		}
	}
}

// IsOver retourne true si le combat est terminé
func (b *Battle) IsOver() bool {
	return b.over
}
