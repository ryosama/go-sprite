package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/ryosama/go-sprite"
	"log"
)

const (
	windowWidth  = 320 // Width of the window
	windowHeight = 240 // Height of the window
	scale        = 2   // Scale of the window
)

var (
	sprites [8]*sprite.Sprite
)

// update at every frame
func update(surface *ebiten.Image) error {

	// frame skip
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// draw sprites
	for i := 0; i < len(sprites); i++ {
		sprites[i].Draw(surface)
	}

	return nil
}

func main() {

	x := 1.0
	y := 1.0
	for i := 0; i < len(sprites); i++ {
		sprites[i] = sprite.NewSprite()
		sprites[i].CenterCoordonnates = true
		sprites[i].AddAnimation("default", "gfx/som_girl_stand_down.png", 1, 1, ebiten.FilterDefault)
		sprites[i].Position(windowWidth/4*x, windowHeight/4*y)
		sprites[i].Start()

		x++
		if x > 3 {
			x = 1
			y++
		}
	}

	i := 0

	// Zoom in and out
	sprites[i].AddEffect(&sprite.EffectOptions{
		Animation: "default", // optional
		Effect:    sprite.Zoom,
		Zoom:      1.3,
		Duration:  1000,
		Repeat:    true,
		GoBack:    true,
		Callback:  func() { print("ZOOM Callback\n") },
	})
	i++

	// Flip X and go back
	sprites[i].AddEffect(&sprite.EffectOptions{Effect: sprite.Flip, Axis: sprite.Horizontaly, Duration: 1000, Repeat: true, GoBack: true})
	i++

	// Flip Y and Go back
	sprites[i].AddEffect(&sprite.EffectOptions{Effect: sprite.Flip, Axis: sprite.Verticaly, Duration: 1000, Repeat: true, GoBack: true})
	i++

	// Fade in and out
	sprites[i].AddEffect(&sprite.EffectOptions{Effect: sprite.Fade, FadeFrom: 1, FadeTo: 0.1, Duration: 2000, Repeat: true, GoBack: true})
	i++

	// Turn and go back clockwise
	sprites[i].AddEffect(&sprite.EffectOptions{Effect: sprite.Turn, Angle: 360, Duration: 2000, Repeat: true, GoBack: true, Clockwise: true})
	i++

	// Hue Red
	sprites[i].AddEffect(&sprite.EffectOptions{Effect: sprite.Hue, Red: 5, Duration: 1000, Repeat: true, GoBack: true})
	i++

	// Move Relative X and Relative Y
	sprites[i].AddEffect(&sprite.EffectOptions{Effect: sprite.Move, X: sprites[i].X + 10, Y: sprites[i].Y + 10, Duration: 1000, Repeat: true, GoBack: true})
	i++

	// multiple effects : Zoom->x3, HUE->Yellow, TURN->360°, MOVE-> x+100
	sprites[i].AddEffect(&sprite.EffectOptions{Effect: sprite.Zoom, Zoom: 3, Duration: 2000, Repeat: true, GoBack: true})
	sprites[i].AddEffect(&sprite.EffectOptions{Effect: sprite.Hue, Red: 5, Green: 5, Duration: 2000, Repeat: true, GoBack: true})
	sprites[i].AddEffect(&sprite.EffectOptions{Effect: sprite.Turn, Angle: 360, Duration: 2000, Repeat: true})
	sprites[i].AddEffect(&sprite.EffectOptions{Effect: sprite.Move, X: sprites[i].X + 100, Duration: 2000, Repeat: true, GoBack: true})
	i++

	// infinite loop
	if err := ebiten.Run(update, windowWidth, windowWidth, scale, "Sprite demo"); err != nil {
		log.Fatal(err)
	}
}
