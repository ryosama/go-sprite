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
	sprites [12]*sprite.Sprite
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

	// Inflate
	sprites[i].AddEffect(&sprite.EffectOptions{ Animation:"default" , Effect: sprite.INFLATE, Zoom:2, Duration:2000, Repeat:true })
	i++
	
	// Defalte
	sprites[i].AddEffect(&sprite.EffectOptions{ Effect: sprite.INFLATE, Zoom:0.5, Duration:2000, Repeat:true})
	i++

	// Breathe
	sprites[i].AddEffect(&sprite.EffectOptions{ Effect: sprite.INFLATE, Zoom:1.3, Duration:1000, Repeat:true, GoBack:true })
	i++

	// Flip X
	sprites[i].AddEffect(&sprite.EffectOptions{ Effect: sprite.FLIPX, Duration:1000, Repeat:true })
	i++

	// Flip X and go back
	sprites[i].AddEffect(&sprite.EffectOptions{ Effect: sprite.FLIPX, Duration:1000, Repeat:true, GoBack:true })
	i++

	// Flip Y
	sprites[i].AddEffect(&sprite.EffectOptions{ Effect: sprite.FLIPY, Duration:1000, Repeat:true })
	i++

	// Flip Y and Go back
	sprites[i].AddEffect(&sprite.EffectOptions{ Effect: sprite.FLIPY, Duration:1000, Repeat:true, GoBack:true })
	i++

	// Fade
	sprites[i].AddEffect(&sprite.EffectOptions{ Effect: sprite.FADE, FadeFrom:1 , FadeTo:0.5, Duration:1000, Repeat:true })
	i++

	// Fade in and out
	sprites[i].AddEffect(&sprite.EffectOptions{ Effect: sprite.FADE, FadeFrom:1 , FadeTo:0.1, Duration:2000, Repeat:true, GoBack:true })
	i++

	// Turn
	sprites[i].AddEffect(&sprite.EffectOptions{ Effect: sprite.TURN, Angle:90, Duration:2000, Repeat:true })
	i++

	// Turn and go back
	sprites[i].AddEffect(&sprite.EffectOptions{ Effect: sprite.TURN, Angle:90, Duration:2000, Repeat:true, GoBack:true })
	i++

	// Turn and go back clockwise
	sprites[i].AddEffect(&sprite.EffectOptions{ Effect: sprite.TURN, Angle:90, Duration:2000, Repeat:true, GoBack:true, Clockwise:true })
	i++

	// infinite loop
	if err := ebiten.Run(update, WINDOW_WIDTH, WINDOW_WIDTH, SCALE, "Sprite demo"); err != nil { log.Fatal(err) }
}