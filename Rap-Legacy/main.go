package main // Déclare le package principal du jeu

import (
	"log" // Pour afficher les erreurs critiques

	"github.com/hajimehoshi/ebiten/v2"      // Bibliothèque Ebiten pour le jeu
	"github.com/projet-red_rap-legacy/game" // Import du package local "game" contenant la logique du jeu
)

func main() {
	// Définit l'icône de la fenêtre avec notre fonction SetGameIcon
	SetGameIcon("assets/icon.png")

	// Crée une nouvelle instance du jeu
	g := game.NewGame()

	// Supprime la barre de fenêtre (bordure et boutons)
	ebiten.SetWindowDecorated(false)

	// Définit la taille de la fenêtre à 1920x1080
	ebiten.SetWindowSize(1920, 1080)

	// Définit le titre de la fenêtre
	ebiten.SetWindowTitle("Rap Legacy")

	// Démarre la boucle principale du jeu
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err) // Affiche une erreur et termine si la boucle du jeu échoue
	}
}
