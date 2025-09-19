package game

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func LoadAnimation(prefix string, count int) []*ebiten.Image {
	var frames []*ebiten.Image
	for i := 1; i <= count; i++ {
		path := filepath.Join("assets", prefix+strconv.Itoa(i)+".png")
		frames = append(frames, LoadImage(path))
	}
	return frames
}

type Battle struct {
	bg *ebiten.Image

	playerEgo int
	enemyEgo  int

	menuOptions    []string
	selectedOption int
	attackDamages  []int
	enemyDamages   []int

	// Animations
	playerIdle *ebiten.Image
	enemyIdle  *ebiten.Image

	playerAtk  []*ebiten.Image
	enemyAtk   []*ebiten.Image
	playerHit  []*ebiten.Image
	enemyHit   []*ebiten.Image
	playerDead []*ebiten.Image
	enemyDead  []*ebiten.Image

	// Animation contrôle
	currentFrames []*ebiten.Image
	currentIndex  int
	animStart     time.Time
	animSpeed     time.Duration
	animPlaying   bool
	attacker      string

	// IA
	lastPlayerAttack int

	// Dialogues
	playerLines    [][]string
	enemyLines     [][]string
	currentLine    string
	lineStart      time.Time
	lineDuration   time.Duration
	dialogCooldown time.Duration
	lastDialogTime time.Time

	// Gestion mort + sortie
	deadFinished  bool
	deadFrame     *ebiten.Image
	endMsg        *ebiten.Image
	exitRequested bool

	// ✅ Nouveau champ pour savoir qui a gagné
	Winner string // "player" ou "enemy"
}

func NewBattle(player *Player, enemy *Enemy) *Battle {
	// Ego du joueur avec bonus temporaire
	egoFinal := player.Ego
	if player != nil {
		egoFinal += player.BonusEgo
		player.BonusEgo = 0
	}

	// Ego de l’ennemi avec éventuel malus
	enemyEgo := enemy.Ego
	if player != nil && player.PendingEnemyEgoDebuff > 0 {
		enemyEgo -= player.PendingEnemyEgoDebuff
		if enemyEgo < 0 {
			enemyEgo = 0
		}
		player.PendingEnemyEgoDebuff = 0
	}

	b := &Battle{
		bg:               LoadImage("assets/battle_bg.png"),
		playerEgo:        egoFinal,
		enemyEgo:         enemyEgo,
		menuOptions:      []string{"Punchline", "Flow", "Diss Track"},
		attackDamages:    []int{10, 5, 30},
		enemyDamages:     []int{10, 5, 30},
		animSpeed:        150 * time.Millisecond,
		animPlaying:      false,
		currentIndex:     0,
		deadFinished:     false,
		selectedOption:   0,
		lastPlayerAttack: 0,
		attacker:         "",
		lineDuration:     2000 * time.Millisecond, // 2s affichage
		dialogCooldown:   2500 * time.Millisecond, // 2,5s entre dialogues
		lastDialogTime:   time.Now().Add(-2 * time.Second),
	}

	rand.Seed(time.Now().UnixNano())

	// Idle
	b.playerIdle = LoadImage("assets/player_idle.png")
	b.enemyIdle = LoadImage("assets/enemy_idle.png")

	// Attaques
	b.playerAtk = LoadAnimation("player_attack", 5)
	b.enemyAtk = LoadAnimation("enemy_attack", 5)

	// Hit
	b.playerHit = LoadAnimation("player_hited", 4)
	b.enemyHit = LoadAnimation("enemy_hited", 4)

	// Dead
	b.playerDead = LoadAnimation("player_dead", 5)
	b.enemyDead = LoadAnimation("enemy_dead", 5)

	// Dialogues joueur
	b.playerLines = [][]string{
		{"Yo je te pète la rime !", "C’est chaud comme le freestyle !"},    // Punchline
		{"Mon flow te fait trembler !", "Tu peux pas suivre mon rythme !"}, // Flow
		{"Diss track incoming ! je vais ruiner ta carrière !"},             // Diss Track
	}

	// Dialogues IA
	b.enemyLines = [][]string{
		{"Tu crois pouvoir me punchliner ?", "J'te mets KO avec mes rimes !"}, // Punchline
		{"Mon flow est supérieur !", "Trop lent pour moi !"},                  // Flow
		{"Diss Track ! je vais te faire regretter !"},                         // Diss Track
	}

	// Image de fin
	b.endMsg = LoadImage("assets/combat_end.png")
	if b.endMsg == nil {
		println("⚠️ Impossible de charger assets/combat_end.png")
	}

	return b
}

func (b *Battle) ChooseEnemyAttack() int {
	r := rand.Intn(100)
	if r < 50 {
		return 0 // Punchline 50%
	} else if r < 80 {
		return 1 // Flow 30%
	}
	return 2 // Diss Track 20%
}

func (b *Battle) LaunchAttack(attacker string) {
	b.animPlaying = true
	b.animStart = time.Now()
	b.currentIndex = 0
	b.attacker = attacker

	now := time.Now()
	if now.Sub(b.lastDialogTime) >= b.dialogCooldown {
		if attacker == "player" {
			b.currentFrames = b.playerAtk
			b.lastPlayerAttack = b.selectedOption
			lines := b.playerLines[b.selectedOption]
			b.currentLine = lines[rand.Intn(len(lines))]
		} else {
			idx := b.ChooseEnemyAttack()
			b.currentFrames = b.enemyAtk
			lines := b.enemyLines[idx]
			b.currentLine = lines[rand.Intn(len(lines))]
		}
		b.lineStart = now
		b.lastDialogTime = now
	} else if attacker == "player" {
		b.currentFrames = b.playerAtk
		b.lastPlayerAttack = b.selectedOption
	} else {
		b.currentFrames = b.enemyAtk
	}
}

func (b *Battle) LaunchDeath(who string) {
	b.animPlaying = true
	b.animStart = time.Now()
	b.currentIndex = 0

	if who == "player" {
		b.attacker = "dead_player"
		b.currentFrames = b.playerDead
		b.Winner = "enemy" // ✅ l'ennemi a gagné
	} else {
		b.attacker = "dead_enemy"
		b.currentFrames = b.enemyDead
		b.Winner = "player" // ✅ le joueur a gagné
	}
}

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
				switch b.attacker {
				case "player":
					dmg := b.attackDamages[b.lastPlayerAttack]
					b.enemyEgo -= dmg
					if b.enemyEgo <= 0 {
						b.LaunchDeath("enemy") // 👈 Winner défini ici
						return
					}
					b.animPlaying = false
					b.currentIndex = 0
					b.LaunchAttack("enemy")
					return

				case "enemy":
					idx := b.ChooseEnemyAttack()
					dmg := b.enemyDamages[idx]
					b.playerEgo -= dmg
					if b.playerEgo <= 0 {
						b.LaunchDeath("player") // 👈 Winner défini ici
						return
					}
					b.animPlaying = false
					b.currentIndex = 0
					b.attacker = ""

				case "dead_player", "dead_enemy":
					b.deadFinished = true
					if len(b.currentFrames) > 0 {
						b.deadFrame = b.currentFrames[len(b.currentFrames)-1]
					}
					b.animPlaying = false
					b.currentIndex = 0
				}
			}
		}
		return
	}

	if b.attacker == "" {
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
}

func (b *Battle) Draw(screen *ebiten.Image) {
	if b.bg != nil {
		screen.DrawImage(b.bg, &ebiten.DrawImageOptions{})
	}

	screenW, screenH := screen.Size()
	scale := 3.0
	playerX := float64(screenW/2) - 400
	enemyX := float64(screenW/2) + 150
	groundY := float64(screenH - 400)

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
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		if b.attacker == "dead_enemy" {
			op.GeoM.Translate(enemyX, groundY)
			screen.DrawImage(b.deadFrame, op)
			op2 := &ebiten.DrawImageOptions{}
			op2.GeoM.Scale(scale, scale)
			op2.GeoM.Translate(playerX, groundY)
			screen.DrawImage(b.playerIdle, op2)
		} else {
			op.GeoM.Translate(playerX, groundY)
			screen.DrawImage(b.deadFrame, op)
			op2 := &ebiten.DrawImageOptions{}
			op2.GeoM.Scale(scale, scale)
			op2.GeoM.Translate(enemyX, groundY)
			screen.DrawImage(b.enemyIdle, op2)
		}

		if b.endMsg != nil {
			opMsg := &ebiten.DrawImageOptions{}
			w, h := b.endMsg.Size()
			endScale := 0.6
			opMsg.GeoM.Scale(endScale, endScale)
			opMsg.GeoM.Translate(
				float64(screenW/2)-(float64(w)*endScale)/2,
				float64(screenH/2)-(float64(h)*endScale)/2,
			)
			screen.DrawImage(b.endMsg, opMsg)
		}
	} else {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(playerX, groundY)
		screen.DrawImage(b.playerIdle, op)
		op2 := &ebiten.DrawImageOptions{}
		op2.GeoM.Scale(scale, scale)
		op2.GeoM.Translate(enemyX, groundY)
		screen.DrawImage(b.enemyIdle, op2)
	}

	// Ego
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Votre égo: %d", b.playerEgo), 10, 10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Égo adverse: %d", b.enemyEgo), screenW-125, 10)

	// Dialogue centré
	if b.currentLine != "" && time.Since(b.lineStart) < b.lineDuration {
		x := float64((screenW - len(b.currentLine)*7) / 2)
		y := float64(screenH/2 - 10)
		ebitenutil.DebugPrintAt(screen, b.currentLine, int(x), int(y))
	}

	// Menu joueur
	if !b.deadFinished && b.attacker == "" {
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
