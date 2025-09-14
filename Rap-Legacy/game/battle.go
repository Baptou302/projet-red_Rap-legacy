package game

import (
	"image/png"
	"os"
	"strconv"
	"time"

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
	active         bool // savoir si un combat est actif

	// --- Animations ---
	playerAnim   *Animation
	enemyAnim    *Animation
	animPlaying  bool
	animForEnemy bool // true = anim de l’ennemi, false = anim du joueur
}

// --------- UTILS POUR CHARGER SPRITESHEET ----------
func LoadSpriteSheet(path string) *ebiten.Image {
	f, _ := os.Open(path)
	defer f.Close()
	img, _ := png.Decode(f)
	return ebiten.NewImageFromImage(img)
}

// NewBattle crée une nouvelle instance de Battle
func NewBattle(p *Player, e *Enemy) *Battle {
	// ⚠️ adapter selon ton spritesheet
	nbLignes := 3
	nbColonnes := 5
	frameDelay := 100 * time.Millisecond

	playerSheet := LoadSpriteSheet("assets/player_spritesheet.png")
	playerAnim := NewAnimation(playerSheet, nbLignes, nbColonnes, frameDelay)

	enemySheet := LoadSpriteSheet("assets/enemy_spritesheet.png")
	enemyAnim := NewAnimation(enemySheet, nbLignes, nbColonnes, frameDelay)

	return &Battle{
		player:         p,
		enemy:          e,
		turn:           0,
		over:           false,
		selectedAttack: 0,
		active:         false,
		playerAnim:     playerAnim,
		enemyAnim:      enemyAnim,
		animPlaying:    false,
		animForEnemy:   false,
	}
}

// Update gère les tours et la logique du combat
func (b *Battle) Update() {
	if !b.active {
		if ebiten.IsKeyPressed(ebiten.KeyE) {
			b.active = true
		}
		return
	}

	if b.over {
		return
	}

	// --- Si une animation est en cours ---
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

	// --- Tour du joueur ---
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
		// --- Tour de l’ennemi ---
		b.animPlaying = true
		b.animForEnemy = true
		b.enemyAnim.current = 0
		b.enemyAnim.lastUpdate = time.Now()
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
	if !b.active {
		ebitenutil.DebugPrintAt(screen, "Appuie sur E pour lancer un combat !", 200, 200)
		return
	}

	// Stats
	ebitenutil.DebugPrintAt(screen, "Votre égo: "+strconv.Itoa(b.player.Ego), 10, 10)
	ebitenutil.DebugPrintAt(screen, "Égo adverse: "+strconv.Itoa(b.enemy.Ego), 500, 10)

	// --- Afficher animations ---
	if b.animPlaying {
		if b.animForEnemy {
			b.enemyAnim.Draw(screen, 400, 200) // position de l’ennemi
		} else {
			b.playerAnim.Draw(screen, 100, 200) // position du joueur
		}
	}

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
