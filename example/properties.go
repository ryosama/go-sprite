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
	sprites [7]*sprite.Sprite
)

// update at every frame
func update(surface *ebiten.Image) error {

	// frame skip
	if ebiten.IsDrawingSkipped() { return nil }

	// draw sprites
	for i:=0 ; i<len(sprites) ; i++ {
		sprites[i].Draw(surface)
	}

	return nil
}


func main() {

	x := 1.0
	y := 1.0
	for i:=0 ; i<len(sprites) ; i++ {
		sprites[i] = sprite.NewSprite()
		sprites[i].CenterCoordonnates = true
		sprites[i].AddAnimation("default","gfx/som_girl_stand_down.png",	 1, 1, ebiten.FilterDefault)
		sprites[i].Position(WINDOW_WIDTH/4 * x, WINDOW_HEIGHT/4 * y)
		sprites[i].Start()

		x++
		if x > 3 {
			x=1
			y++
		}
	}

	i:=0

	sprites[i].Zoom(2, 1.5) // set the multiplier ZoomX and ZoomY
	i++

	sprites[i].Skew(30, 10) // set SkewX and SkewY in degres
	i++

	sprites[i].Rotate(45) // set Angle in degres
	i++

	sprites[i].Alpha = 0.5 // Between 0 and 1
	i++

	sprites[i].Red = 5 // Multiplier
	i++

	sprites[i].Red = 3 // become Yellow
	sprites[i].Green = 3
	i++

	sprites[i].Borders = true // Debug Borders
	i++

	// infinite loop
	if err := ebiten.Run(update, WINDOW_WIDTH, WINDOW_WIDTH, SCALE, "Sprite demo"); err != nil { log.Fatal(err) }
}