package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/ryosama/go-sprite"
	"log"
)

const (
	windowWidth  = 320 // Width of the window
	windowHeight = 240 // Height of the window
	scale        = 2   // Scale of the window
)

var (
	explosion4 *sprite.Sprite
)

// update at every frame
func update(surface *ebiten.Image) error {

	// manage controle
	binding()

	// frame skip
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// draw sprite
	explosion4.Draw(surface)

	ebitenutil.DebugPrint(surface, "Press 'Space' to restart animation")

	return nil
}

func main() {

	explosion4 = sprite.NewSprite()
	explosion4.AddAnimation("default", "gfx/explosion3.png", 500, 9, ebiten.FilterDefault)
	explosion4.Position(windowWidth/2, windowHeight/2)
	explosion4.RunOnce(afterRunOnce)

	// infinite loop
	if err := ebiten.Run(update, windowWidth, windowWidth, scale, "Sprite demo"); err != nil {
		log.Fatal(err)
	}
}

func afterRunOnce(s *sprite.Sprite) {
	print("Execute after run once\n")
}

func binding() {

	//////////////////////////// GO THE RIGHT
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		explosion4.RunOnce(afterRunOnce)
	}
}
