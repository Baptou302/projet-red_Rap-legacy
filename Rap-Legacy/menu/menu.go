package menu // Déclare le package "menu", utilisé pour organiser le code du menu principal

import (
	"image/color" // Pour gérer les couleurs (ex : fond du menu)
	"log"         // Pour afficher les erreurs de manière lisible
	"os"          // Pour gérer les fichiers et quitter le programme

	"github.com/hajimehoshi/ebiten/v2"            // Bibliothèque principale du jeu
	"github.com/hajimehoshi/ebiten/v2/audio"      // Gestion audio
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"  // Décodage MP3
	"github.com/hajimehoshi/ebiten/v2/ebitenutil" // Utilitaires (ex : DebugPrintAt)
)

// -----------------
// Input helper partagé
// -----------------

var prevInput = map[ebiten.Key]bool{} // Map pour mémoriser l'état précédent de chaque touche

// IsKeyJustPressed renvoie true si la touche vient juste d'être pressée
func IsKeyJustPressed(key ebiten.Key) bool {
	pressed := ebiten.IsKeyPressed(key) // Etat actuel de la touche
	wasPressed := prevInput[key]        // Etat précédent de la touche
	prevInput[key] = pressed            // Sauvegarde l'état actuel pour la prochaine frame
	return pressed && !wasPressed       // True si la touche est pressée maintenant mais pas avant
}

const sampleRate = 44100 // Fréquence audio standard

var audioContext *audio.Context // Contexte audio partagé pour le menu

type Menu struct {
	menuSelected                       int           // Option sélectionnée dans le menu
	menuOptions                        []string      // Liste des options
	buttonX, buttonY, buttonW, buttonH int           // Position et dimensions des boutons
	audioPlayer                        *audio.Player // Player pour la musique de fond
}

// NewMenu crée un nouveau menu
func NewMenu() *Menu {
	// Initialise le contexte audio si ce n'est pas déjà fait
	if audioContext == nil {
		audioContext = audio.NewContext(sampleRate)
	}

	m := &Menu{
		menuOptions: []string{"New game", "Settings", "Quit"}, // Options du menu
		buttonX:     700,                                      // Position X des boutons
		buttonY:     400,                                      // Position Y des boutons
		buttonW:     500,                                      // Largeur des boutons
		buttonH:     80,                                       // Hauteur des boutons
	}

	// Ouverture du fichier MP3 pour la musique du menu
	f, err := os.Open("../menu/menu.mp3")
	if err != nil {
		log.Println("Impossible d'ouvrir le fichier MP3 :", err)
		return m
	}

	// Décodage du MP3 avec le sample rate défini
	stream, err := mp3.DecodeWithSampleRate(sampleRate, f)
	if err != nil {
		log.Println("Impossible de décoder le MP3 :", err)
		return m
	}

	// Boucle infinie de la musique
	loop := audio.NewInfiniteLoop(stream, stream.Length())
	m.audioPlayer, err = audioContext.NewPlayer(loop)
	if err != nil {
		log.Println("Impossible de créer le player audio :", err)
		return m
	}
	m.audioPlayer.Play() // Démarre la lecture

	return m
}

// Update gère l'input clavier et souris et retourne la prochaine scène à lancer
func (m *Menu) Update() string {
	// Flèche haut : sélection précédente
	if IsKeyJustPressed(ebiten.KeyUp) {
		m.menuSelected--
		if m.menuSelected < 0 {
			m.menuSelected = len(m.menuOptions) - 1 // Reboucle à la fin
		}
	}

	// Flèche bas : sélection suivante
	if IsKeyJustPressed(ebiten.KeyDown) {
		m.menuSelected++
		if m.menuSelected >= len(m.menuOptions) {
			m.menuSelected = 0 // Reboucle au début
		}
	}

	// Touche Entrée : valider l'option sélectionnée
	if IsKeyJustPressed(ebiten.KeyEnter) {
		switch m.menuSelected {
		case 0:
			return "intro" // Lancer l'intro au lieu de démarrer directement le jeu
		case 1:
			return "settings"
		case 2:
			os.Exit(0) // Quitter le jeu
		}
	}

	// Gestion clic souris
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition() // Coordonnées du curseur
		for i := range m.menuOptions {
			bx := m.buttonX
			by := m.buttonY + i*(m.buttonH+20)
			bw := m.buttonW
			bh := m.buttonH
			// Vérifie si le curseur est sur un bouton
			if x >= bx && x <= bx+bw && y >= by && y <= by+bh {
				switch i {
				case 0:
					return "intro"
				case 1:
					return "settings"
				case 2:
					os.Exit(0)
				}
			}
		}
	}

	return "" // Aucune action effectuée
}

// Draw affiche le menu à l'écran
func (m *Menu) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 30, 255}) // Fond gris foncé
	for i, option := range m.menuOptions {
		x := m.buttonX
		y := m.buttonY + i*(m.buttonH+20)
		text := option
		if i == m.menuSelected {
			text = "> " + option // Marque l'option sélectionnée
		}
		// Affiche le texte du bouton
		ebitenutil.DebugPrintAt(screen, text, x+20, y+20)
	}
}
