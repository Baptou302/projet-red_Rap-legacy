package game // Déclare le package "game", qui contient tous les éléments du jeu

import (
	"fmt"           // Pour formater du texte (ex: fmt.Sprintf)
	"math/rand"     // Pour générer des nombres aléatoires
	"path/filepath" // Pour créer des chemins de fichiers portables
	"strconv"       // Pour convertir des int en string
	"time"          // Pour gérer durées et timestamps

	"github.com/hajimehoshi/ebiten/v2"            // Ebiten, moteur 2D
	"github.com/hajimehoshi/ebiten/v2/ebitenutil" // Pour afficher texte et debug facilement
)

// LoadAnimation charge une série d’images pour une animation
func LoadAnimation(prefix string, count int) []*ebiten.Image {
	var frames []*ebiten.Image    // Slice pour stocker toutes les frames
	for i := 1; i <= count; i++ { // Boucle pour charger chaque frame
		path := filepath.Join("assets", prefix+strconv.Itoa(i)+".png") // Construire chemin
		frames = append(frames, LoadImage(path))                       // Charger et ajouter à frames
	}
	return frames // Retourne les frames chargées
}

// Battle contient toutes les infos d’un combat
type Battle struct {
	bg *ebiten.Image // Image de fond

	playerEgo int // Ego joueur
	enemyEgo  int // Ego ennemi

	menuOptions    []string // Options du menu de combat
	selectedOption int      // Option sélectionnée
	attackDamages  []int    // Dégâts attaques joueur
	enemyDamages   []int    // Dégâts attaques ennemies

	// Animations
	playerIdle *ebiten.Image // Sprite idle joueur
	enemyIdle  *ebiten.Image // Sprite idle ennemi

	playerAtk  []*ebiten.Image // Animation attaque joueur
	enemyAtk   []*ebiten.Image // Animation attaque ennemi
	playerHit  []*ebiten.Image // Animation hit joueur
	enemyHit   []*ebiten.Image // Animation hit ennemi
	playerDead []*ebiten.Image // Animation mort joueur
	enemyDead  []*ebiten.Image // Animation mort ennemi

	// Animation contrôle
	currentFrames []*ebiten.Image // Frames actuellement jouées
	currentIndex  int             // Index frame courante
	animStart     time.Time       // Timestamp début animation
	animSpeed     time.Duration   // Délai entre frames
	animPlaying   bool            // Animation en cours ?
	attacker      string          // Qui attaque ("player", "enemy", "dead_player"...)

	// IA
	lastPlayerAttack int // Dernière attaque sélectionnée par le joueur

	// Dialogues
	playerLines    [][]string    // Lignes de dialogues joueur
	enemyLines     [][]string    // Lignes de dialogues IA
	currentLine    string        // Ligne affichée
	lineStart      time.Time     // Début affichage dialogue
	lineDuration   time.Duration // Durée affichage
	dialogCooldown time.Duration // Délai entre dialogues
	lastDialogTime time.Time     // Dernier dialogue affiché

	// Gestion mort + sortie
	deadFinished  bool          // Animation mort terminée ?
	deadFrame     *ebiten.Image // Dernière frame de mort
	endMsg        *ebiten.Image // Image fin combat
	exitRequested bool          // Sortie demandée ?

	Winner string // "player" ou "enemy"
}

// NewBattle initialise un combat avec un joueur et un ennemi
func NewBattle(player *Player, enemy *Enemy) *Battle {
	egoFinal := player.Ego // Ego de base joueur
	if player != nil {
		egoFinal += player.BonusEgo // Ajouter bonus temporaire
		player.BonusEgo = 0         // Réinitialiser bonus
	}

	enemyEgo := enemy.Ego // Ego de base ennemi
	if player != nil && player.PendingEnemyEgoDebuff > 0 {
		enemyEgo -= player.PendingEnemyEgoDebuff // Appliquer malus
		if enemyEgo < 0 {
			enemyEgo = 0 // Minimum 0
		}
		player.PendingEnemyEgoDebuff = 0
	}

	// Crée la structure Battle
	b := &Battle{
		bg:               LoadImage("assets/battle_bg.png"),           // Fond combat
		playerEgo:        egoFinal,                                    // Ego joueur
		enemyEgo:         enemyEgo,                                    // Ego ennemi
		menuOptions:      []string{"Punchline", "Flow", "Diss Track"}, // Menu
		attackDamages:    []int{10, 5, 30},                            // Dégâts joueur
		enemyDamages:     []int{10, 5, 30},                            // Dégâts ennemis
		animSpeed:        150 * time.Millisecond,                      // Vitesse animations
		animPlaying:      false,
		currentIndex:     0,
		deadFinished:     false,
		selectedOption:   0,
		lastPlayerAttack: 0,
		attacker:         "",
		lineDuration:     2000 * time.Millisecond,          // 2s affichage dialogues
		dialogCooldown:   2500 * time.Millisecond,          // 2,5s entre dialogues
		lastDialogTime:   time.Now().Add(-2 * time.Second), // Permet dialogue immédiat
	}

	rand.Seed(time.Now().UnixNano()) // Initialiser RNG

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

	return b // Retourne la structure initialisée
}

// ChooseEnemyAttack sélectionne aléatoirement l’attaque ennemie
func (b *Battle) ChooseEnemyAttack() int {
	r := rand.Intn(100) // 0-99
	if r < 50 {         // 50%
		return 0 // Punchline
	} else if r < 80 { // 30%
		return 1 // Flow
	}
	return 2 // 20% Diss Track
}

// LaunchAttack démarre l’animation d’une attaque
func (b *Battle) LaunchAttack(attacker string) {
	b.animPlaying = true     // Lancer animation
	b.animStart = time.Now() // Timestamp début
	b.currentIndex = 0       // Première frame
	b.attacker = attacker    // Stocke attaquant

	now := time.Now()
	if now.Sub(b.lastDialogTime) >= b.dialogCooldown { // Si cooldown ok
		if attacker == "player" { // Joueur attaque
			b.currentFrames = b.playerAtk
			b.lastPlayerAttack = b.selectedOption
			lines := b.playerLines[b.selectedOption]
			b.currentLine = lines[rand.Intn(len(lines))] // Ligne aléatoire
		} else { // Ennemi attaque
			idx := b.ChooseEnemyAttack()
			b.currentFrames = b.enemyAtk
			lines := b.enemyLines[idx]
			b.currentLine = lines[rand.Intn(len(lines))]
		}
		b.lineStart = now      // Début affichage dialogue
		b.lastDialogTime = now // Reset cooldown dialogue
	} else if attacker == "player" {
		b.currentFrames = b.playerAtk
		b.lastPlayerAttack = b.selectedOption
	} else {
		b.currentFrames = b.enemyAtk
	}
}

// LaunchDeath démarre animation mort et définit le gagnant
func (b *Battle) LaunchDeath(who string) {
	b.animPlaying = true
	b.animStart = time.Now()
	b.currentIndex = 0

	if who == "player" {
		b.attacker = "dead_player"
		b.currentFrames = b.playerDead
		b.Winner = "enemy" // ✅ l'ennemi gagne
	} else {
		b.attacker = "dead_enemy"
		b.currentFrames = b.enemyDead
		b.Winner = "player" // ✅ joueur gagne
	}
}

func (b *Battle) Update() {
	// Vérifie si le combat est terminé (mort + fin animation)
	if b.deadFinished {
		// Si le joueur appuie sur Enter, demande de sortie
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			b.exitRequested = true
		}
		// On ne fait rien d'autre si le combat est terminé
		return
	}

	// Gestion des animations en cours (attaque ou mort)
	if b.animPlaying {
		// Vérifie si le temps écoulé depuis le début de l'animation est supérieur à la vitesse définie
		if time.Since(b.animStart) > b.animSpeed {
			// Réinitialise le temps de départ pour la frame suivante
			b.animStart = time.Now()
			// Passe à la frame suivante
			b.currentIndex++
			// Vérifie si on est arrivé à la fin de l'animation
			if b.currentIndex >= len(b.currentFrames) {
				// Selon qui effectue l'attaque ou la mort
				switch b.attacker {
				case "player":
					// Dégâts infligés à l'ennemi
					dmg := b.attackDamages[b.lastPlayerAttack]
					b.enemyEgo -= dmg
					// Vérifie si l'ennemi est mort
					if b.enemyEgo <= 0 {
						b.LaunchDeath("enemy") // Déclenche l'animation de mort de l'ennemi
						return
					}
					// Sinon, fin de l'animation et préparation de l'attaque ennemie
					b.animPlaying = false
					b.currentIndex = 0
					b.LaunchAttack("enemy")
					return

				case "enemy":
					// Choix aléatoire de l'attaque ennemie
					idx := b.ChooseEnemyAttack()
					dmg := b.enemyDamages[idx]
					b.playerEgo -= dmg
					// Vérifie si le joueur est mort
					if b.playerEgo <= 0 {
						b.LaunchDeath("player") // Déclenche l'animation de mort du joueur
						return
					}
					// Sinon, fin de l'animation et reset de l'attaquant
					b.animPlaying = false
					b.currentIndex = 0
					b.attacker = ""

				case "dead_player", "dead_enemy":
					// Si animation de mort terminée
					b.deadFinished = true
					// Stocke la dernière frame de l'animation de mort
					if len(b.currentFrames) > 0 {
						b.deadFrame = b.currentFrames[len(b.currentFrames)-1]
					}
					// Fin de l'animation
					b.animPlaying = false
					b.currentIndex = 0
				}
			}
		}
		// Ne rien faire d'autre pendant l'animation
		return
	}

	// Gestion de la sélection du menu joueur si personne n'attaque
	if b.attacker == "" {
		// Déplacement vers le bas dans le menu
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
			b.selectedOption++
			// Bouclage sur la première option si on dépasse la fin
			if b.selectedOption >= len(b.menuOptions) {
				b.selectedOption = 0
			}
			// Petit délai pour éviter plusieurs mouvements rapides
			time.Sleep(150 * time.Millisecond)
		}
		// Déplacement vers le haut dans le menu
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
			b.selectedOption--
			// Bouclage sur la dernière option si on dépasse le début
			if b.selectedOption < 0 {
				b.selectedOption = len(b.menuOptions) - 1
			}
			time.Sleep(150 * time.Millisecond)
		}
		// Validation de l'attaque avec Enter
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			b.LaunchAttack("player")
		}
	}
}

func (b *Battle) Draw(screen *ebiten.Image) {
	// Dessine le fond si disponible
	if b.bg != nil {
		screen.DrawImage(b.bg, &ebiten.DrawImageOptions{})
	}

	// Récupération des dimensions de l'écran
	screenW, screenH := screen.Size()
	scale := 3.0
	// Positions X des sprites joueur et ennemi
	playerX := float64(screenW/2) - 400
	enemyX := float64(screenW/2) + 150
	// Position Y du sol
	groundY := float64(screenH - 400)

	// Dessin des animations en cours
	if b.animPlaying && len(b.currentFrames) > 0 {
		idx := b.currentIndex
		// S'assure qu'on ne dépasse pas le nombre de frames
		if idx >= len(b.currentFrames) {
			idx = len(b.currentFrames) - 1
		}
		frame := b.currentFrames[idx]

		// Dessin selon l'attaquant
		switch b.attacker {
		case "player":
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(playerX, groundY)
			screen.DrawImage(frame, op)
			// Dessine l'ennemi qui prend un coup si encore vivant
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
			// Dessine le joueur qui prend un coup si encore vivant
			if b.playerEgo > 0 && len(b.playerHit) > 0 {
				hitFrame := b.playerHit[idx%len(b.playerHit)]
				op2 := &ebiten.DrawImageOptions{}
				op2.GeoM.Scale(scale, scale)
				op2.GeoM.Translate(playerX, groundY)
				screen.DrawImage(hitFrame, op2)
			}
		case "dead_player":
			// Animation de mort du joueur
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(playerX, groundY)
			screen.DrawImage(frame, op)
			// Ennemi reste idle
			op2 := &ebiten.DrawImageOptions{}
			op2.GeoM.Scale(scale, scale)
			op2.GeoM.Translate(enemyX, groundY)
			screen.DrawImage(b.enemyIdle, op2)
		case "dead_enemy":
			// Animation de mort de l'ennemi
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(enemyX, groundY)
			screen.DrawImage(frame, op)
			// Joueur reste idle
			op2 := &ebiten.DrawImageOptions{}
			op2.GeoM.Scale(scale, scale)
			op2.GeoM.Translate(playerX, groundY)
			screen.DrawImage(b.playerIdle, op2)
		}
	} else if b.deadFinished && b.deadFrame != nil {
		// Dessin de la frame finale si le combat est terminé
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

		// Dessin de l'image de fin du combat si disponible
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
		// Dessin idle des deux personnages si pas d'animation
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(playerX, groundY)
		screen.DrawImage(b.playerIdle, op)
		op2 := &ebiten.DrawImageOptions{}
		op2.GeoM.Scale(scale, scale)
		op2.GeoM.Translate(enemyX, groundY)
		screen.DrawImage(b.enemyIdle, op2)
	}

	// Affiche les valeurs d'ego
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Votre égo: %d", b.playerEgo), 10, 10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("égo de Lil’ Patafix: %d", b.enemyEgo), screenW-155, 10)

	// Affiche le dialogue en cours si encore actif
	if b.currentLine != "" && time.Since(b.lineStart) < b.lineDuration {
		x := float64((screenW - len(b.currentLine)*7) / 2)
		y := float64(screenH/2 - 10)
		ebitenutil.DebugPrintAt(screen, b.currentLine, int(x), int(y))
	}

	// Dessin du menu joueur si pas de combat en cours
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

// Fonction qui retourne si le combat est terminé
func (b *Battle) IsOver() bool {
	return b.exitRequested
}
