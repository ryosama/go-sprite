package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/ryosama/go-sprite"
	"log"
)

const (
	WINDOW_WIDTH  		= 320			// Width of the window
	WINDOW_HEIGHT 		= 240			// Height of the window
	SCALE         		= 2 			// Scale of the window
)

var (
	skew1 *sprite.Sprite
)

// update at every frame
func update(surface *ebiten.Image) error {

	// frame skip
	if ebiten.IsDrawingSkipped() { return nil }

	skew1.Draw(surface)

	return nil
}


func main() {

	skew1 = sprite.NewSprite()
	skew1.AddAnimation("default","gfx/som_girl_stand_down.png",	 1, 1, ebiten.FilterDefault)
	skew1.Position(WINDOW_WIDTH/2, WINDOW_HEIGHT/2)
	skew1.Skew(20, 23) // in degres
	skew1.Start()

	// infinite loop
	if err := ebiten.Run(update, WINDOW_WIDTH, WINDOW_WIDTH, SCALE, "Sprite demo"); err != nil { log.Fatal(err) }
}