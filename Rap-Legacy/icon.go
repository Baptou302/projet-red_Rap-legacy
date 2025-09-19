// icon.go
package main // Déclare le package principal pour le jeu

import (
	"image"       // Gestion des images en Go
	_ "image/png" // Import pour le décodage PNG (underscore = import pour l'effet secondaire)
	"log"         // Pour afficher les erreurs
	"os"          // Gestion des fichiers

	"github.com/hajimehoshi/ebiten/v2" // Bibliothèque Ebiten pour le jeu
)

// SetGameIcon charge et définit l'icône de la fenêtre du jeu
func SetGameIcon(path string) {
	// Ouvre le fichier image
	f, err := os.Open(path)
	if err != nil {
		log.Println("Impossible d'ouvrir l'icône :", err) // Affiche une erreur si le fichier est introuvable
		return
	}
	defer f.Close() // Ferme le fichier à la fin de la fonction

	// Décode l'image (PNG, JPEG, etc.)
	img, _, err := image.Decode(f)
	if err != nil {
		log.Println("Erreur décodage icône :", err) // Affiche une erreur si le fichier n'est pas une image valide
		return
	}

	// Définit l'icône de la fenêtre du jeu
	ebiten.SetWindowIcon([]image.Image{img})
}
