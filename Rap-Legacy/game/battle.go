package game

import (
	"image/png"
	"os"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Battle struct {
	player         *Player
	enemy          *Enemy
	turn           int
	over           bool
	selectedAttack int
	active         bool

	playerAnim   *Animation
	enemyAnim    *Animation
	animPlaying  bool
	animForEnemy bool

	playerIdle *ebiten.Image
	enemyIdle  *ebiten.Image
	background *ebiten.Image
}

// Charge une image et panic si problème
func LoadSprite(path string) *ebiten.Image {
	f, err := os.Open(path)
	if err != nil {
		panic("❌ Impossible d’ouvrir le fichier : " + path + " | " + err.Error())
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		panic("❌ Impossible de décoder l’image : " + path + " | " + err.Error())
	}

	return ebiten.NewImageFromImage(img)
}

func NewBattle(p *Player, e *Enemy) *Battle {
	nbLignes := 3
	nbColonnes := 5
	frameDelay := 300 * time.Millisecond

	playerSheet := LoadSprite("assets/player_spritesheet.png")
	playerAnim := NewAnimation(playerSheet, nbLignes, nbColonnes, frameDelay)

	enemySheet := LoadSprite("assets/enemy_spritesheet.png")
	enemyAnim := NewAnimation(enemySheet, nbLignes, nbColonnes, frameDelay)

	// Idle séparés
	playerIdle := LoadSprite("assets/player_idle.png")
	enemyIdle := LoadSprite("assets/enemy_idle.png")

	// Fond combat
	background := LoadSprite("assets/battle_bg.png")

	return &Battle{
		player:         p,
		enemy:          e,
		turn:           0,
		over:           false,
		selectedAttack: 0,
		active:         false,
		playerAnim:     playerAnim,
		enemyAnim:      enemyAnim,
		playerIdle:     playerIdle,
		enemyIdle:      enemyIdle,
		background:     background,
	}
}

func (b *Battle) Update() {
	if !b.active || b.over {
		return
	}

	if b.animPlaying {
		if b.animForEnemy {
			b.enemyAnim.Update()
			if b.enemyAnim.current == b.enemyAnim.frameCount-1 {
				b.EnemyAttack()
				b.turn = 0
				b.animPlaying = false
			}
		} else {
			b.playerAnim.Update()
			if b.playerAnim.current == b.playerAnim.frameCount-1 {
				b.PlayerAttack()
				b.turn = 1
				b.animPlaying = false
			}
		}
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
			b.animPlaying = true
			b.animForEnemy = false
			b.playerAnim.current = 0
			b.playerAnim.lastUpdate = time.Now()
		}
	} else {
		// Tour ennemi
		b.animPlaying = true
		b.animForEnemy = true
		b.enemyAnim.current = 0
		b.enemyAnim.lastUpdate = time.Now()
	}

	if b.player.Ego <= 0 || b.enemy.Ego <= 0 {
		b.over = true
	}
}

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

func (b *Battle) Draw(screen *ebiten.Image) {
	if !b.active {
		ebitenutil.DebugPrintAt(screen, "Appuie sur E pour lancer un combat !", 200, 200)
		return
	}

	screenW, screenH := screen.Size()

	// --- Fond combat ---
	opBg := &ebiten.DrawImageOptions{}
	bgW, bgH := b.background.Size()
	scaleX := float64(screenW) / float64(bgW)
	scaleY := float64(screenH) / float64(bgH)
	opBg.GeoM.Scale(scaleX, scaleY)
	screen.DrawImage(b.background, opBg)

	// --- Infos Ego ---
	ebitenutil.DebugPrintAt(screen, "Votre égo: "+strconv.Itoa(b.player.Ego), 10, 10)
	ebitenutil.DebugPrintAt(screen, "Égo adverse: "+strconv.Itoa(b.enemy.Ego), screenW-200, 10)

	// --- Placement persos ---
	scale := 3.0
	pw, ph := b.playerIdle.Size()
	ew, eh := b.enemyIdle.Size()

	// positions collées au sol
	groundYPlayer := float64(screenH)
	groundYEnemy := float64(screenH)

	if b.animPlaying {
		if b.animForEnemy {
			b.enemyAnim.Draw(screen, float64(3*screenW/4), groundYEnemy, scale, true)
		} else {
			b.playerAnim.Draw(screen, float64(screenW/4), groundYPlayer, scale, true)
		}
	} else {
		// Joueur idle
		opPlayer := &ebiten.DrawImageOptions{}
		opPlayer.GeoM.Scale(scale, scale)
		opPlayer.GeoM.Translate(
			float64(screenW/4)-float64(pw)*scale/2,
			groundYPlayer-float64(ph)*scale,
		)
		screen.DrawImage(b.playerIdle, opPlayer)

		// Ennemi idle
		opEnemy := &ebiten.DrawImageOptions{}
		opEnemy.GeoM.Scale(scale, scale)
		opEnemy.GeoM.Translate(
			float64(3*screenW/4)-float64(ew)*scale/2,
			groundYEnemy-float64(eh)*scale,
		)
		screen.DrawImage(b.enemyIdle, opEnemy)
	}

	// --- Menu attaques en bas à gauche ---
	attacks := []string{"Punchline", "Flow", "Diss Track"}
	for i, a := range attacks {
		text := a
		if i == b.selectedAttack {
			text = "> " + a
		}
		ebitenutil.DebugPrintAt(screen, text, 10, screenH-80+i*20)
	}

	// --- Message fin ---
	if b.over {
		if b.player.Ego <= 0 {
			ebitenutil.DebugPrintAt(screen, "Tu vas te prendre une sauce sur les réseaux !", screenW/2-150, screenH/2)
		} else {
			ebitenutil.DebugPrintAt(screen, "Tu vas avoir un gros buzz !", screenW/2-150, screenH/2)
		}
	}
}

func (b *Battle) IsOver() bool {
	return b.over
}
