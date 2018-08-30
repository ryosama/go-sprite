/*
This package permits animations of sprites via the Ebiten library (http://www.github.com/hajimehoshi/ebiten)

Basic Usage :

import "github.com/ryosama/go-sprite"

mySprite = sprite.NewSprite()
mySprite.AddAnimation("walk-right",	"walk_right.png", 700, 6, ebiten.FilterDefault)
mySprite.Position(WINDOW_WIDTH/2, WINDOW_HEIGHT/2)
mySprite.CurrentAnimation = "walk-right"
mySprite.Speed = 2
mySprite.Start()
*/
package sprite

import (
	"image"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"time"
	"math"
	"log"
)

//////////////////////////////////////////// TYPES ////////////////////////////////////////////

type Sprite struct {
	// Animation label currently displayed
	CurrentAnimation	string

	// Array of animations
	Animations 			map[string]*SpriteAnimation

	// X coordinates of the sprite (in pixel)
	X 					int

	// Y coordinate of the sprite (in pixel)
	Y 					int

	// Speed is in pixel/frame
	Speed				float64

	// Direction is an Angle in degres
	Direction			float64

	// Visibility of the sprite
	Visible				bool

	// Animated or not
	Animated			bool
}

type SpriteAnimation struct {
	// File path of the animation
	// Step of the animation must be the same width on one line
	Path 							string

	// ebiten.Image generated
	Image 							*ebiten.Image

	// Number of steps for the total animation
	Steps 							int

	// Current step displayed
	CurrentStep 	 				int

	// Where to start the animation
	FirstStep 						int

	// Width of the animation steps (in pixel)
	StepWidth 						int

	// Height of the animation steps (in pixel)
	StepHeight 						int

	// Total duration of the animation in millisecond
	Duration						time.Duration

	// Total time for one step in millisecond
	OneStepDuration					time.Duration

	// Start time of the current step
	currentStepTimeStart 			time.Time
}

//////////////////////////////////////////// CONSTRUCTORS ////////////////////////////////////////////

/*
Create a new sprite
*/
func NewSprite() *Sprite {
	this := new(Sprite)
	this.Animations = make(map[string]*SpriteAnimation)
	this.Visible = true
	this.Animated   = true
	return this
}

func newSpriteAnimation(path string,duration int, steps int, filter ebiten.Filter) *SpriteAnimation {
	var err error
	this 			:= new(SpriteAnimation)
	this.Path 		= path
	this.Image, _,err = ebitenutil.NewImageFromFile(path, filter)
	if err != nil { log.Fatal(err) }
	this.Steps 		= steps
	this.Duration 	= time.Millisecond * time.Duration(duration)

	width, height := this.Image.Size()
	this.StepWidth 	= width/this.Steps
	this.StepHeight = height

	this.currentStepTimeStart = time.Now()
	this.OneStepDuration = time.Duration(int(this.Duration) / this.Steps)

	return this
}

//////////////////////////////////////////// METHODS ////////////////////////////////////////////

/*
Add an animation to the sprite

"label" is the tag for the animation

"path" is the path for the image file

"duration" is in millisecond

"steps" is the number of step for the animation

"filter" is ebiten.FilterDefault or ebiten.FilterNearest  or ebiten.FilterLinear

Example : 

mySprite.AddAnimation("walk-right",	"walk_right.png", 700, 6, ebiten.FilterDefault)
*/
func (this *Sprite) AddAnimation(label string, path string, duration int, steps int, filter ebiten.Filter) {
	this.Animations[label] = newSpriteAnimation(path,duration,steps,filter)
}

/*
Return width of the current animation displayed
*/
func (this *Sprite) GetWidth() int {
	currentAnimation := this.Animations[this.CurrentAnimation]
	return currentAnimation.StepWidth
}

/*
Return height of the current animation displayed
*/
func (this *Sprite) GetHeight() int {
	currentAnimation := this.Animations[this.CurrentAnimation]
	return currentAnimation.StepHeight
}

/*
Hide the sprite
*/
func (this *Sprite) Hide() {
	this.Visible = false
}

/*
Show the sprite
*/
func (this *Sprite) Show() {
	this.Visible = true
}

/*
Toogle visibility of the sprite
*/
func (this *Sprite) ToogleVisibility() {
	if this.Visible {
		this.Hide()
	} else {
		this.Show()
	}
}

/*
Set X and Y coordonnates of the sprite

Return X and Y coordonnates

Exemple :
mySprite.Position(WINDOW_WIDTH/2, WINDOW_HEIGHT/2)

or 

x,y := mySprite.Position()
*/
func (this *Sprite) Position(arg... int) (int,int) {
	if len(arg)==2 {
		this.X = arg[0]
		this.Y = arg[1]
	}
	return this.X, this.Y
}

/*
Calculate new coordonnates and draw the sprite on the screen, after drawing, go to the next step of animation
*/
func (this *Sprite) Draw(surface *ebiten.Image) {
	if this.Visible {
		currentAnimation := this.Animations[this.CurrentAnimation] // SpriteAnimation object

		options := &ebiten.DrawImageOptions{}
		
		// move sprite x,y
		angleRad := (this.Direction * math.Pi / 180) + math.Pi/2
		this.X += int(this.Speed * math.Sin(angleRad))
		this.Y += int(this.Speed * math.Cos(angleRad))

		options.GeoM.Translate(float64(this.X), float64(this.Y)) 

		// Choose current image inside animation
		x0 := currentAnimation.CurrentStep * currentAnimation.StepWidth
		x1 := currentAnimation.CurrentStep * currentAnimation.StepWidth + currentAnimation.StepWidth
		r := image.Rect( x0 , 0, x1 , currentAnimation.StepHeight)
		options.SourceRect = &r

		surface.DrawImage(currentAnimation.Image, options)

		this.NextStep()
	}
}

/*
Start the animation (Reset+Show+Resume)
*/
func (this *Sprite) Start() {
	this.Reset()
	this.Show()
	this.Resume()
}

/*
Stop the animation (Reset+Pause)
*/
func (this *Sprite) Stop() {
	this.Reset()
	this.Pause()
}

/*
Reset current step to the first step of the animation
*/
func (this *Sprite) Reset() {
	currentAnimation := this.Animations[this.CurrentAnimation]
	currentAnimation.CurrentStep = currentAnimation.FirstStep
}

/*
Pause the animation
*/
func (this *Sprite) Pause() {
	this.Animated = false
}

/*
Resume the animation
*/
func (this *Sprite) Resume() {
	this.Animated = true
}

/*
Toogle animation status
*/
func (this *Sprite) ToogleAnimation() {
	if this.Animated {
		this.Pause()
	} else {
		this.Resume()
	}
}

/*
Go to the next step of animation

Return true if animation go to the next step or false if step duration is not finish
*/
func (this *Sprite) NextStep() bool {
	currentAnimation := this.Animations[this.CurrentAnimation]
	if this.Animated {
		now 		:= time.Now()
		nextStepAt 	:= currentAnimation.currentStepTimeStart.Add(currentAnimation.OneStepDuration)

		if now.Sub(nextStepAt) > 0 { // time to change the current step
			currentAnimation.CurrentStep++ // next step
			if currentAnimation.CurrentStep+1 > currentAnimation.Steps {
				this.Reset() // restart at the end of the animation
			}
			currentAnimation.currentStepTimeStart = now
			return true
		}
	}
	return false
}