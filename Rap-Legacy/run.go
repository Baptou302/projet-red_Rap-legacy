package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var game *Game

type Game struct {
	player       *Player
	mapData      *Map
	enemies      []*Enemy
	inBattle     bool
	currentEnemy *Enemy
	battle       *Battle
}

func (g *Game) Update() error {
	if g.inBattle {
		g.battle.Update()
		if g.battle.IsOver() {
			g.inBattle = false
		}
	} else {
		g.player.Update()
		for _, e := range g.enemies {
			if g.player.X == e.X && g.player.Y == e.Y {
				g.inBattle = true
				g.currentEnemy = e
				g.battle = NewBattle(g.player, e)
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.mapData.Draw(screen)
	g.player.Draw(screen)
	for _, e := range g.enemies {
		e.Draw(screen)
	}
	if g.inBattle {
		g.battle.Draw(screen)
	}
	ebitenutil.DebugPrint(screen, "Rap Legacy")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	player := NewPlayer(100, 100)
	mapData := NewMap()
	enemies := []*Enemy{
		NewEnemy(200, 200, "Rapeur Rivale"),
		NewEnemy(400, 300, "Star main stream"),
	}

	game = &Game{
		player:  player,
		mapData: mapData,
		enemies: enemies,
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Rap Legacy")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
