package menu

import (
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// -----------------
// Input helper partagé
// -----------------
var prevInput = map[ebiten.Key]bool{}

func IsKeyJustPressed(key ebiten.Key) bool {
	pressed := ebiten.IsKeyPressed(key)
	wasPressed := prevInput[key]
	prevInput[key] = pressed
	return pressed && !wasPressed
}

const sampleRate = 44100

var audioContext *audio.Context

type Menu struct {
	menuSelected                       int
	menuOptions                        []string
	buttonX, buttonY, buttonW, buttonH int
	audioPlayer                        *audio.Player
}

func NewMenu() *Menu {
	if audioContext == nil {
		audioContext = audio.NewContext(sampleRate)
	}

	m := &Menu{
		menuOptions: []string{"New game", "Settings", "Quit"},
		buttonX:     700,
		buttonY:     400,
		buttonW:     500,
		buttonH:     80,
	}

	f, err := os.Open("../menu/menu.mp3")
	if err != nil {
		log.Println("Impossible d'ouvrir le fichier MP3 :", err)
		return m
	}
	stream, err := mp3.DecodeWithSampleRate(sampleRate, f)
	if err != nil {
		log.Println("Impossible de décoder le MP3 :", err)
		return m
	}

	loop := audio.NewInfiniteLoop(stream, stream.Length())
	m.audioPlayer, err = audioContext.NewPlayer(loop)
	if err != nil {
		log.Println("Impossible de créer le player audio :", err)
		return m
	}
	m.audioPlayer.Play()

	return m
}

func (m *Menu) Update() string {
	if IsKeyJustPressed(ebiten.KeyUp) {
		m.menuSelected--
		if m.menuSelected < 0 {
			m.menuSelected = len(m.menuOptions) - 1
		}
	}
	if IsKeyJustPressed(ebiten.KeyDown) {
		m.menuSelected++
		if m.menuSelected >= len(m.menuOptions) {
			m.menuSelected = 0
		}
	}
	if IsKeyJustPressed(ebiten.KeyEnter) {
		switch m.menuSelected {
		case 0:
			// 🚀 Lancer la cinématique au lieu du "play" direct
			return "intro"
		case 1:
			return "settings"
		case 2:
			os.Exit(0)
		}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		for i := range m.menuOptions {
			bx := m.buttonX
			by := m.buttonY + i*(m.buttonH+20)
			bw := m.buttonW
			bh := m.buttonH
			if x >= bx && x <= bx+bw && y >= by && y <= by+bh {
				switch i {
				case 0:
					// 🚀 Idem avec la souris
					return "intro"
				case 1:
					return "settings"
				case 2:
					os.Exit(0)
				}
			}
		}
	}

	return ""
}

func (m *Menu) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 30, 255})
	for i, option := range m.menuOptions {
		x := m.buttonX
		y := m.buttonY + i*(m.buttonH+20)
		text := option
		if i == m.menuSelected {
			text = "> " + option
		}
		ebitenutil.DebugPrintAt(screen, text, x+20, y+20)
	}
}
