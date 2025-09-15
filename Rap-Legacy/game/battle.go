package game

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// -----------------
// Battle structure
// -----------------
type Battle struct {
	player         *Player
	enemy          *Enemy
	turn           int
	over           bool
	selectedAttack int
	active         bool

	// Animations
	playerIdle []*ebiten.Image
	enemyIdle  []*ebiten.Image

	playerAttack []*ebiten.Image
	enemyAttack  []*ebiten.Image

	playerHited []*ebiten.Image
	enemyHited  []*ebiten.Image

	playerDead []*ebiten.Image
	enemyDead  []*ebiten.Image

	battleBg *ebiten.Image

	// Animation state
	animFrames   []*ebiten.Image
	animIndex    int
	animSpeed    time.Duration
	lastAnimTime time.Time
	animPlaying  bool
	animForEnemy bool
	animType     string
}

// -----------------
// Load frames util
// -----------------
func loadFrames(prefix string, count int) []*ebiten.Image {
	var frames []*ebiten.Image
	if count == 1 {
		frames = append(frames, LoadImage(fmt.Sprintf("assets/%s.png", prefix)))
	} else {
		for i := 1; i <= count; i++ {
			frames = append(frames, LoadImage(fmt.Sprintf("assets/%s%d.png", prefix, i)))
		}
	}
	return frames
}

// -----------------
// NewBattle
// -----------------
func NewBattle(p *Player, e *Enemy) *Battle {
	return &Battle{
		player:         p,
		enemy:          e,
		turn:           0,
		over:           false,
		selectedAttack: 0,
		active:         false,

		playerIdle:   loadFrames("player_idle", 1),
		enemyIdle:    loadFrames("enemy_idle", 1),
		playerAttack: loadFrames("player_attack", 5),
		enemyAttack:  loadFrames("enemy_attack", 5),
		playerHited:  loadFrames("player_hited", 3),
		enemyHited:   loadFrames("enemy_hited", 4),
		playerDead:   loadFrames("player_dead", 5),
		enemyDead:    loadFrames("enemy_dead", 5),

		battleBg: LoadImage("assets/battle_bg.png"),

		animSpeed: 200 * time.Millisecond,
	}
}

// -----------------
// Update
// -----------------
func (b *Battle) Update() {
	if !b.active || b.over {
		return
	}

	// Si une animation est en cours
	if b.animPlaying {
		if time.Since(b.lastAnimTime) >= b.animSpeed {
			b.animIndex++
			b.lastAnimTime = time.Now()
			if b.animIndex >= len(b.animFrames) {
				b.animPlaying = false
				b.animIndex = 0

				// Action après animation
				if b.animType == "player_attack" {
					b.PlayerAttack()
					b.turn++
				} else if b.animType == "enemy_attack" {
					b.EnemyAttack()
					b.turn = 0
				}
			}
		}
		return
	}

	// Fin du combat
	if b.player.Ego <= 0 {
		b.animFrames = b.playerDead
		b.animPlaying = true
		b.animIndex = 0
		b.animType = "player_dead"
		b.over = true
		return
	}
	if b.enemy.Ego <= 0 {
		b.animFrames = b.enemyDead
		b.animPlaying = true
		b.animIndex = 0
		b.animType = "enemy_dead"
		b.over = true
		return
	}

	// Tour joueur
	if b.turn == 0 {
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
			// Lance anim attaque joueur + ennemi hit
			b.animFrames = b.playerAttack
			b.animPlaying = true
			b.animIndex = 0
			b.lastAnimTime = time.Now()
			b.animType = "player_attack"
		}
	} else {
		// Tour ennemi (attaque auto)
		b.animFrames = b.enemyAttack
		b.animPlaying = true
		b.animIndex = 0
		b.lastAnimTime = time.Now()
		b.animType = "enemy_attack"
	}
}

// -----------------
// Attaques
// -----------------
func (b *Battle) PlayerAttack() {
	switch b.selectedAttack {
	case 0:
		b.enemy.Ego -= b.player.Flow
	case 1:
		b.enemy.Ego -= b.player.Flow / 2
		b.player.Flow++
	case 2:
		b.enemy.Ego -= b.player.Flow * 2
		b.player.Charisma--
	}
}

func (b *Battle) EnemyAttack() {
	b.player.Ego -= 5
}

// -----------------
// Draw
// -----------------
func (b *Battle) Draw(screen *ebiten.Image) {
	if !b.active {
		ebitenutil.DebugPrintAt(screen, "Appuie sur E pour lancer un combat !", 200, 200)
		return
	}

	// Fond de combat
	if b.battleBg != nil {
		opts := &ebiten.DrawImageOptions{}
		screen.DrawImage(b.battleBg, opts)
	}

	// Stats
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Votre égo: %d", b.player.Ego), 50, 50)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Égo adverse: %d", b.enemy.Ego), 1600, 50)

	// Affichage animations
	if b.animPlaying {
		frame := b.animFrames[b.animIndex]
		op := &ebiten.DrawImageOptions{}
		if b.animType == "player_attack" || b.animType == "player_dead" {
			op.GeoM.Translate(200, 700)
		} else if b.animType == "enemy_attack" || b.animType == "enemy_dead" {
			op.GeoM.Translate(1400, 700)
		}
		screen.DrawImage(frame, op)
	} else {
		// Idle
		opPlayer := &ebiten.DrawImageOptions{}
		opPlayer.GeoM.Translate(200, 700)
		screen.DrawImage(b.playerIdle[0], opPlayer)

		opEnemy := &ebiten.DrawImageOptions{}
		opEnemy.GeoM.Translate(1400, 700)
		screen.DrawImage(b.enemyIdle[0], opEnemy)
	}

	// Menu attaques
	attacks := []string{"Punchline", "Flow", "Diss Track"}
	for i, a := range attacks {
		text := a
		if i == b.selectedAttack {
			text = "> " + a
		}
		ebitenutil.DebugPrintAt(screen, text, 50, 900+i*30)
	}

	// Message fin
	if b.over {
		if b.player.Ego <= 0 {
			ebitenutil.DebugPrintAt(screen, "Tu vas te prendre une sauce sur les réseaux !", 700, 500)
		} else {
			ebitenutil.DebugPrintAt(screen, "Tu vas avoir un gros buzz !", 700, 500)
		}
	}
}

// -----------------
// IsOver
// -----------------
func (b *Battle) IsOver() bool {
	return b.over
}
