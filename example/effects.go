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

	//sprites[0].Draw(surface)

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


	// Zoom in and out
	sprites[i].AddEffect(&sprite.EffectOptions{
							Animation:"default", // optional
							Effect: sprite.ZOOM,
							Zoom:1.3,
							Duration:1000,
							Repeat:true,
							GoBack:true,
							Callback:func(){ print("ZOOM Callback\n") },
						})
	i++

	// Flip X and go back
	sprites[i].AddEffect(&sprite.EffectOptions{ Effect: sprite.FLIPX, Duration:1000, Repeat:true, GoBack:true})
	i++

	// Flip Y and Go back
	sprites[i].AddEffect(&sprite.EffectOptions{ Effect: sprite.FLIPY, Duration:1000, Repeat:true, GoBack:true })
	i++

	// Fade in and out
	sprites[i].AddEffect(&sprite.EffectOptions{ Effect: sprite.FADE, FadeFrom:1 , FadeTo:0.1, Duration:2000, Repeat:true, GoBack:true })
	i++

	// Turn and go back clockwise
	sprites[i].AddEffect(&sprite.EffectOptions{ Effect: sprite.TURN, Angle:360, Duration:2000, Repeat:true, GoBack:true, Clockwise:true })
	i++

	// Hue Yellow
	sprites[i].AddEffect(&sprite.EffectOptions{ Effect: sprite.HUE, Red:5,  Duration:1000, Repeat:true, GoBack:true })
	i++

	// Move Relative X and Relative Y
	sprites[i].AddEffect(&sprite.EffectOptions{ Effect: sprite.MOVE, X: sprites[i].X +10, Y:sprites[i].Y +10, Duration:1000, Repeat:true, GoBack:true })
	i++


	// infinite loop
	if err := ebiten.Run(update, WINDOW_WIDTH, WINDOW_WIDTH, SCALE, "Sprite demo"); err != nil { log.Fatal(err) }
}