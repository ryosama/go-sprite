package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/ryosama/go-sprite"
	"log"
	"fmt"
)

const (
	WINDOW_WIDTH  		= 320			// Width of the window
	WINDOW_HEIGHT 		= 240			// Height of the window
	SCALE         		= 2 			// Scale of the window
	CARACTERE_SPEED		= WINDOW_WIDTH/160
)

var (
	girl *sprite.Sprite
	explosion1,explosion2,explosion3,explosion4, zoom1, rotate1, skew1 *sprite.Sprite
)

// update at every frame
func update(surface *ebiten.Image) error {

	// manage controle
	binding()

	// reset position if outside of the screen
	if girl.X 						> WINDOW_WIDTH 	{ girl.X = 0 - girl.GetWidth()  }
	if girl.X +   girl.GetWidth() 	< 0  			{ girl.X = WINDOW_WIDTH }
	if girl.Y +   girl.GetHeight() 	< 0 			{ girl.Y = WINDOW_HEIGHT + 2*girl.GetHeight() }
	if girl.Y - 2*girl.GetHeight() 	> WINDOW_HEIGHT { girl.Y = 0 - girl.GetHeight() }

	// frame skip
	if ebiten.IsDrawingSkipped() { return nil }

	// draw sprite
	girl.Draw(surface)
	explosion1.Draw(surface)
	explosion2.Draw(surface)
	explosion3.Draw(surface)
	explosion4.Draw(surface)

	zoom1.Draw(surface)

	rotate1.Draw(surface)

	skew1.Draw(surface)

	// display some informations
	drawFPS(surface)

	return nil
}


func main() {

	// create new sprite and load animations
	girl = sprite.NewSprite()
	girl.AddAnimation("stand-right","gfx/som_girl_stand_right.png",	  0,1, ebiten.FilterDefault)
	girl.AddAnimation("walk-right",	"gfx/som_girl_walk_right.png",	700,6, ebiten.FilterDefault)
	girl.AddAnimation("stand-left",	"gfx/som_girl_stand_left.png",	  0,1, ebiten.FilterDefault)
	girl.AddAnimation("walk-left",	"gfx/som_girl_walk_left.png",	700,6, ebiten.FilterDefault)
	girl.AddAnimation("stand-up",	"gfx/som_girl_stand_up.png",	  0,1, ebiten.FilterDefault)
	girl.AddAnimation("walk-up",	"gfx/som_girl_walk_up.png",	    500,4, ebiten.FilterDefault)
	girl.AddAnimation("stand-down",	"gfx/som_girl_stand_down.png",	  0,1, ebiten.FilterDefault)
	girl.AddAnimation("walk-down",	"gfx/som_girl_walk_down.png",   500,4, ebiten.FilterDefault)

	// set position and first animation
	girl.Position(WINDOW_WIDTH/2, WINDOW_HEIGHT/2)
	girl.CurrentAnimation = "stand-right"
	girl.Start()

	explosionDuration := 500
	// create some explosions
	explosion1 = sprite.NewSprite()
	explosion1.AddAnimation("default","gfx/explosion1.png",	explosionDuration, 5, ebiten.FilterDefault)
	explosion1.Position(10, WINDOW_HEIGHT/3*2)
	explosion1.Start()

	explosion2 = sprite.NewSprite()
	explosion2.AddAnimation("default","gfx/explosion2.png",	 explosionDuration, 7, ebiten.FilterDefault)
	explosion2.Position(WINDOW_WIDTH/2 -24, WINDOW_HEIGHT/3*2)
	explosion2.Start()

	explosion3 = sprite.NewSprite()
	explosion3.AddAnimation("default","gfx/explosion3.png",	 explosionDuration, 9, ebiten.FilterDefault)
	explosion3.Position(WINDOW_WIDTH -10 -48, WINDOW_HEIGHT/3*2)
	explosion3.Start()

	explosion4 = sprite.NewSprite()
	explosion4.AddAnimation("default","gfx/explosion3.png",	 explosionDuration, 9, ebiten.FilterDefault)
	explosion4.Position(WINDOW_WIDTH -10 -48, 50)
	explosion4.RunOnce(afterRunOnce)

	zoom1 = sprite.NewSprite()
	zoom1.AddAnimation("default","gfx/som_girl_stand_down.png",	 1, 1, ebiten.FilterDefault)
	zoom1.Position(10, 20)
	zoom1.Zoom(2)
	zoom1.Start()

	rotate1 = sprite.NewSprite()
	rotate1.AddAnimation("default","gfx/som_girl_stand_down.png",	 1, 1, ebiten.FilterDefault)
	rotate1.Position(200, 70)
	rotate1.Rotate(45) // in degres
	rotate1.Start()

	skew1 = sprite.NewSprite()
	skew1.AddAnimation("default","gfx/som_girl_stand_down.png",	 1, 1, ebiten.FilterDefault)
	skew1.Position(60, 70)
	skew1.Skew(45,23) // in degres
	skew1.Start()

	// infinite loop
	if err := ebiten.Run(update, WINDOW_WIDTH, WINDOW_WIDTH, SCALE, "Sprite demo"); err != nil { log.Fatal(err) }
}

func afterRunOnce(s *sprite.Sprite) {
	fmt.Printf("Execute after run once %v\n",s)
}

// display some stuff on the screen
func drawFPS(surface *ebiten.Image) {
	ebitenutil.DebugPrint(surface,
		fmt.Sprintf("FPS:%0.1f  X:%d Y:%d %s\nLeft:%v Right:%v Up:%v Down:%v",
			ebiten.CurrentFPS(),
			int(girl.X), int(girl.Y),
			girl.CurrentAnimation,
			ebiten.IsKeyPressed(ebiten.KeyLeft),
			ebiten.IsKeyPressed(ebiten.KeyRight),
			ebiten.IsKeyPressed(ebiten.KeyUp),
			ebiten.IsKeyPressed(ebiten.KeyDown),
	))
}

func binding() {

//////////////////////////// GO THE RIGHT
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) { 

		if 		  ebiten.IsKeyPressed(ebiten.KeyUp) {		// Right+Up
			girl.Direction = 45
			girl.Speed 		= CARACTERE_SPEED+1
		} else if ebiten.IsKeyPressed(ebiten.KeyDown)	{	// Right+Down		
			girl.Direction = -45
			girl.Speed 		= CARACTERE_SPEED+1
		} else {											// Right
			girl.Direction = 0
			girl.Speed 		= CARACTERE_SPEED
		}
		girl.CurrentAnimation = "walk-right"
		girl.Start() // Show, Reset, Resume
	}

	if inpututil.IsKeyJustReleased(ebiten.KeyRight) {
		girl.Speed = 0
		girl.CurrentAnimation = "stand-right"
	}


//////////////////////////// GO THE LEFT
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {

		if 		  ebiten.IsKeyPressed(ebiten.KeyUp) {		// Left+Up
			girl.Direction  = 135
			girl.Speed 		= CARACTERE_SPEED+1
		} else if ebiten.IsKeyPressed(ebiten.KeyDown)	{	// Left+Down		
			girl.Direction  = 225
			girl.Speed 		= CARACTERE_SPEED+1
		} else {											// Left
			girl.Speed 		= CARACTERE_SPEED
			girl.Direction 	= 180
		}

		girl.CurrentAnimation = "walk-left"
		girl.Start() // Show, Reset, Resume
	}

	if inpututil.IsKeyJustReleased(ebiten.KeyLeft) {
		girl.Speed = 0
		girl.CurrentAnimation = "stand-left"
	}


//////////////////////////// GO THE TOP
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {

		if 		  ebiten.IsKeyPressed(ebiten.KeyRight) {	// Up+Right
			girl.Direction = 45
			girl.Speed 		= CARACTERE_SPEED+1
		} else if ebiten.IsKeyPressed(ebiten.KeyLeft)	{	// Up+Left
			girl.Direction = 135
			girl.Speed 		= CARACTERE_SPEED+1
		} else {											// Up
			girl.Direction = 90
			girl.Speed 		= CARACTERE_SPEED
		}

		girl.CurrentAnimation = "walk-up"
		girl.Start() // Show, Reset, Resume
	}

	if inpututil.IsKeyJustReleased(ebiten.KeyUp) {
		girl.Speed = 0
		girl.CurrentAnimation = "stand-up"
	}


//////////////////////////// GO THE BOTTOM
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {

		if 		  ebiten.IsKeyPressed(ebiten.KeyRight) {	// Down+Right
			girl.Direction  = -45
			girl.Speed 		= CARACTERE_SPEED+1
		} else if ebiten.IsKeyPressed(ebiten.KeyLeft)	{	// Down+Left
			girl.Direction  = 225
			girl.Speed 		= CARACTERE_SPEED+1
		} else {											// Down
			girl.Speed 		= CARACTERE_SPEED
			girl.Direction 	= 270
		}

		girl.CurrentAnimation = "walk-down"
		girl.Start() // Show, Reset, Resume
	}

	if inpututil.IsKeyJustReleased(ebiten.KeyDown) {
		girl.Speed = 0
		girl.CurrentAnimation = "stand-down"
	}
}