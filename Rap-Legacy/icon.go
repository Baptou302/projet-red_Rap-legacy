// icon.go
package main

import (
	"image"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func SetGameIcon(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Println("Impossible d'ouvrir l'icône :", err)
		return
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.Println("Erreur décodage icône :", err)
		return
	}

	ebiten.SetWindowIcon([]image.Image{img})
}
