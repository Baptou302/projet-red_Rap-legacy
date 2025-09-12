package game

import (
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func DrawHUD(screen *ebiten.Image, player *Player) {
	ebitenutil.DebugPrintAt(screen, "Ego: "+strconv.Itoa(player.Ego), 10, 10)
	ebitenutil.DebugPrintAt(screen, "Flow: "+strconv.Itoa(player.Flow), 10, 30)
	ebitenutil.DebugPrintAt(screen, "Charisma: "+strconv.Itoa(player.Charisma), 10, 50)
}
