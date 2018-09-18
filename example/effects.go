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

	// Inflate
	sprites[0].AddEffect(&sprite.EffectOptions{ Animation:"default" , Effect: sprite.INFLATE, Zoom:2, Duration:200, Repeat:true })
	
	// Defalte
	sprites[1].AddEffect(&sprite.EffectOptions{ Effect: sprite.DEFLATE, Zoom:2, Duration:2000, Repeat:true })

	// Breathe
	sprites[2].AddEffect(&sprite.EffectOptions{ Effect: sprite.BREATHE, Zoom:1.3, Duration:1000, Repeat:true })

	// Flip X
	sprites[3].AddEffect(&sprite.EffectOptions{ Effect: sprite.FLIPX, Duration:1000, Repeat:true })

	// Flip Y
	sprites[4].AddEffect(&sprite.EffectOptions{ Effect: sprite.FLIPY, Duration:1000, Repeat:true })

	// Fade
	sprites[5].AddEffect(&sprite.EffectOptions{ Effect: sprite.FADE, FadeFrom:1 , FadeTo:0.5, Duration:1000, Repeat:true })

	// Fade in and out
	sprites[6].AddEffect(&sprite.EffectOptions{ Effect: sprite.FADEINOUT, FadeFrom:1 , FadeTo:0.1, Duration:2000, Repeat:true })

	// infinite loop
	if err := ebiten.Run(update, WINDOW_WIDTH, WINDOW_WIDTH, SCALE, "Sprite demo"); err != nil { log.Fatal(err) }
}