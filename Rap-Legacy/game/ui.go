package game // Déclare le package "game", utilisé pour organiser le code

import (
	"strconv" // Pour convertir des entiers en chaînes de caractères

	"github.com/hajimehoshi/ebiten/v2"            // Bibliothèque principale pour le jeu
	"github.com/hajimehoshi/ebiten/v2/ebitenutil" // Utilitaires de Ebiten (ex : DebugPrintAt)
)

// DrawHUD affiche les stats du joueur à l'écran
func DrawHUD(screen *ebiten.Image, player *Player) {
	// Affiche l'Ego du joueur en haut à gauche (x=10, y=10)
	ebitenutil.DebugPrintAt(screen, "Ego: "+strconv.Itoa(player.Ego), 10, 10)

	// Affiche le Flow du joueur juste en dessous (x=10, y=30)
	ebitenutil.DebugPrintAt(screen, "Flow: "+strconv.Itoa(player.Flow), 10, 30)

	// Affiche le Charisma du joueur juste en dessous (x=10, y=50)
	ebitenutil.DebugPrintAt(screen, "Charisma: "+strconv.Itoa(player.Charisma), 10, 50)
}
