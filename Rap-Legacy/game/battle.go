package game

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Charge une série d’images (ex: "player_attack1.png", "player_attack2.png", …)
func LoadAnimation(prefix string, count int) []*ebiten.Image {
	var frames []*ebiten.Image
	for i := 1; i <= count; i++ {
		path := filepath.Join("assets", prefix+strconv.Itoa(i)+".png")
		frames = append(frames, LoadImage(path)) // utilise LoadImage définie dans game.go
	}
	return frames
}

// Battle (frame-by-frame images)
type Battle struct {
	bg *ebiten.Image

	playerEgo int
	enemyEgo  int

	menuOptions    []string
	selectedOption int

	// Animations
	playerIdle *ebiten.Image
	enemyIdle  *ebiten.Image

	playerAtk  [][]*ebiten.Image
	enemyAtk   [][]*ebiten.Image
	playerHit  []*ebiten.Image
	enemyHit   []*ebiten.Image
	playerDead []*ebiten.Image
	enemyDead  []*ebiten.Image

	// Contrôle animation
	currentFrames []*ebiten.Image
	currentIndex  int
	animStart     time.Time
	animSpeed     time.Duration
	animPlaying   bool
	attacker      string // "player" / "enemy" / "dead_player" / "dead_enemy"

	// Gestion mort + sortie
	deadFinished  bool
	deadFrame     *ebiten.Image
	endMsg        *ebiten.Image
	exitRequested bool
}

func NewBattle(player *Player, enemy *Enemy) *Battle {
	baseEgo := 100
	egoFinal := baseEgo
	if player != nil {
		egoFinal += player.BonusEgo
		// Important : reset le bonus pour que ça ne reste pas après
		player.BonusEgo = 0
	}

	b := &Battle{
		bg:             LoadImage("assets/battle_bg.png"), // ✅ fond du combat
		playerEgo:      egoFinal,
		enemyEgo:       baseEgo,
		menuOptions:    []string{"Punchline", "Flow", "Diss Track"},
		animSpeed:      150 * time.Millisecond,
		animPlaying:    false,
		currentIndex:   0,
		deadFinished:   false,
		selectedOption: 0,
	}

	// Idle
	b.playerIdle = LoadImage("assets/player_idle.png")
	b.enemyIdle = LoadImage("assets/enemy_idle.png")

	// Attaques
	b.playerAtk = [][]*ebiten.Image{
		LoadAnimation("player_attack", 5),
	}
	b.enemyAtk = [][]*ebiten.Image{
		LoadAnimation("enemy_attack", 5),
	}

	// Hit
	b.playerHit = LoadAnimation("player_hited", 4)
	b.enemyHit = LoadAnimation("enemy_hited", 4)

	// Dead
	b.playerDead = LoadAnimation("player_dead", 5)
	b.enemyDead = LoadAnimation("enemy_dead", 5)

	// Image de fin (mets ton image dans assets sous ce nom exact)
	b.endMsg = LoadImage("assets/combat_end.png")
	if b.endMsg == nil {
		println("⚠️ Impossible de charger assets/combat_end.png (vérifie le nom et le dossier)")
	}

	return b
}

// Update logique
func (b *Battle) Update() {
	if b.deadFinished {
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			b.exitRequested = true
		}
		return
	}

	if b.animPlaying {
		if time.Since(b.animStart) > b.animSpeed {
			b.animStart = time.Now()
			b.currentIndex++
			if b.currentIndex >= len(b.currentFrames) {
				if b.attacker == "dead_enemy" || b.attacker == "dead_player" {
					b.deadFinished = true
					if len(b.currentFrames) > 0 {
						b.deadFrame = b.currentFrames[len(b.currentFrames)-1]
					}
					b.animPlaying = false
					b.currentIndex = 0
					return
				}

				if b.attacker == "player" {
					b.enemyEgo -= 10
					if b.enemyEgo <= 0 {
						b.LaunchDeath("enemy")
						return
					}
				} else if b.attacker == "enemy" {
					b.playerEgo -= 5
					if b.playerEgo <= 0 {
						b.LaunchDeath("player")
						return
					}
				}

				b.animPlaying = false
				b.currentIndex = 0
			}
		}
		return
	}

	// Navigation menu
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		b.selectedOption++
		if b.selectedOption >= len(b.menuOptions) {
			b.selectedOption = 0
		}
		time.Sleep(150 * time.Millisecond)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		b.selectedOption--
		if b.selectedOption < 0 {
			b.selectedOption = len(b.menuOptions) - 1
		}
		time.Sleep(150 * time.Millisecond)
	}
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		b.LaunchAttack("player")
	}
}

func (b *Battle) LaunchAttack(attacker string) {
	b.animPlaying = true
	b.animStart = time.Now()
	b.currentIndex = 0
	b.attacker = attacker

	if attacker == "player" {
		b.currentFrames = b.playerAtk[0]
	} else {
		b.currentFrames = b.enemyAtk[0]
	}
}

func (b *Battle) LaunchDeath(who string) {
	b.animPlaying = true
	b.animStart = time.Now()
	b.currentIndex = 0

	if who == "player" {
		b.attacker = "dead_player"
		b.currentFrames = b.playerDead
	} else {
		b.attacker = "dead_enemy"
		b.currentFrames = b.enemyDead
	}
}

// Draw
func (b *Battle) Draw(screen *ebiten.Image) {
	if b.bg != nil {
		screen.DrawImage(b.bg, &ebiten.DrawImageOptions{})
	}

	screenW, screenH := screen.Size()
	scale := 3.0
	playerX := float64(screenW/2) - 400
	enemyX := float64(screenW/2) + 150
	groundY := float64(screenH - 400)

	// Affichage pendant animation
	if b.animPlaying && len(b.currentFrames) > 0 {
		idx := b.currentIndex
		if idx >= len(b.currentFrames) {
			idx = len(b.currentFrames) - 1
		}
		frame := b.currentFrames[idx]

		switch b.attacker {
		case "player":
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(playerX, groundY)
			screen.DrawImage(frame, op)

			// Ennemi encaisse (frame hit si dispo)
			if b.enemyEgo > 0 && len(b.enemyHit) > 0 {
				hitFrame := b.enemyHit[idx%len(b.enemyHit)]
				op2 := &ebiten.DrawImageOptions{}
				op2.GeoM.Scale(scale, scale)
				op2.GeoM.Translate(enemyX, groundY)
				screen.DrawImage(hitFrame, op2)
			}

		case "enemy":
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(enemyX, groundY)
			screen.DrawImage(frame, op)

			// Player encaisse
			if b.playerEgo > 0 && len(b.playerHit) > 0 {
				hitFrame := b.playerHit[idx%len(b.playerHit)]
				op2 := &ebiten.DrawImageOptions{}
				op2.GeoM.Scale(scale, scale)
				op2.GeoM.Translate(playerX, groundY)
				screen.DrawImage(hitFrame, op2)
			}

		case "dead_player":
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(playerX, groundY)
			screen.DrawImage(frame, op)

			op2 := &ebiten.DrawImageOptions{}
			op2.GeoM.Scale(scale, scale)
			op2.GeoM.Translate(enemyX, groundY)
			screen.DrawImage(b.enemyIdle, op2)

		case "dead_enemy":
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(enemyX, groundY)
			screen.DrawImage(frame, op)

			op2 := &ebiten.DrawImageOptions{}
			op2.GeoM.Scale(scale, scale)
			op2.GeoM.Translate(playerX, groundY)
			screen.DrawImage(b.playerIdle, op2)
		}

	} else if b.deadFinished && b.deadFrame != nil {
		// Affiche dernier sprite de mort (fixe)
		if b.attacker == "dead_enemy" {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(enemyX, groundY)
			screen.DrawImage(b.deadFrame, op)

			op2 := &ebiten.DrawImageOptions{}
			op2.GeoM.Scale(scale, scale)
			op2.GeoM.Translate(playerX, groundY)
			screen.DrawImage(b.playerIdle, op2)

		} else if b.attacker == "dead_player" {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(playerX, groundY)
			screen.DrawImage(b.deadFrame, op)

			op2 := &ebiten.DrawImageOptions{}
			op2.GeoM.Scale(scale, scale)
			op2.GeoM.Translate(enemyX, groundY)
			screen.DrawImage(b.enemyIdle, op2)
		}

		// Affiche image de fin (centrée + réduite si besoin)
		if b.endMsg != nil {
			opMsg := &ebiten.DrawImageOptions{}
			w, h := b.endMsg.Size()

			endScale := 0.6 // ajuste si trop grand/petit
			opMsg.GeoM.Scale(endScale, endScale)
			opMsg.GeoM.Translate(
				float64(screenW/2)-(float64(w)*endScale)/2,
				float64(screenH/2)-(float64(h)*endScale)/2,
			)
			screen.DrawImage(b.endMsg, opMsg)
		}

	} else {
		// Idle : afficher idles
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(playerX, groundY)
		screen.DrawImage(b.playerIdle, op)

		op2 := &ebiten.DrawImageOptions{}
		op2.GeoM.Scale(scale, scale)
		op2.GeoM.Translate(enemyX, groundY)
		screen.DrawImage(b.enemyIdle, op2)
	}

	// UI : ego en haut
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Votre égo: %d", b.playerEgo), 10, 10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Égo adverse: %d", b.enemyEgo), screenW-125, 10)

	// Menu attaques
	if !b.deadFinished {
		for i, option := range b.menuOptions {
			y := screenH - 60 + i*20
			prefix := "  "
			if i == b.selectedOption {
				prefix = "> "
			}
			ebitenutil.DebugPrintAt(screen, prefix+option, 10, y)
		}
	}
}

func (b *Battle) IsOver() bool {
	return b.exitRequested
}
