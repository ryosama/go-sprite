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
	"image/color"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"time"
	"math"
	"log"
	"fmt"
)

// create constant for effects
const (
	NO_EFFECT = iota
	ZOOM
	FLIPX
	FLIPY
	FADE
	TURN
	HUE
)

//////////////////////////////////////////// TYPES ////////////////////////////////////////////

type Sprite struct {
	// Animation label currently displayed
	CurrentAnimation	string

	// Array of animations
	Animations 			map[string]*SpriteAnimation

	// X coordinates of the sprite (in pixel)
	X 					float64

	// Y coordinate of the sprite (in pixel)
	Y 					float64

	// Speed is in pixel/frame
	Speed				float64

	// Direction is an Angle in degres
	Direction			float64

	// Zoom in or out on X axis
	ZoomX				float64

	// Zoom in or out on Y axis
	ZoomY				float64
	
	// Red multiplier
	Red					float64

	// Green multiplier
	Green				float64

	// Blue multiplier
	Blue				float64

	// Transparency
	Alpha				float64

	// Angle of rotation in degres
	Angle				float64

	// Skew on X axis in degres
	SkewX				float64

	// Skew on Y axis in degres
	SkewY				float64

	// Visibility of the sprite
	Visible				bool

	// Animated or not
	Animated			bool

	// Displace X and Y coordonnate to the center of the sprite
	CenterCoordonnates	bool

	// Draw debug borders around sprite
	Borders				bool
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

	// Effect object
	Effect 							*AnimationEffect

	// Animation once and disapared
	RunOnce 						bool

	// Callback after run once
	callbackAfterRunOnce 			func(*Sprite)

	// Start time of the current step
	currentStepTimeStart 			time.Time
}

type AnimationEffect struct {
	options 				*EffectOptions
	zoomStart										float64
	redStart, greenStart, blueStart, alphaStart	 	float64
	angleStart										float64
	duration 				time.Duration
	timeStart,timeEnd		time.Time
	repeatCallback			func()
}

/*
Create effect on sprite
*/
type EffectOptions struct {
	// Name of animation (default is omitted)
	Animation 				string

	// Effect= ZOOM, FLIPX, FLIPY, FADE, TURN
	Effect 					int

	// For FADE and FADEINOUT effects
	FadeFrom,FadeTo			float64

	// For ZOOM effects
	Zoom 					float64

	// For TURN effect
	Clockwise 				bool

	// For TURN effect (in degres)
	Angle 					float64

	// For HUE effect
	Red, Green, Blue		float64	

	// Duration of the effect
	Duration 				int

	// Duration of the effect (in time.Duration)
	durationTime 			time.Duration

	// Redo the animation on the counter way
	GoBack 					bool

	// Repeat or not at the end of effect
	Repeat 					bool

	// function to launch afert one complete effect
	Callback 				func()
}

//////////////////////////////////////////// CONSTRUCTORS ////////////////////////////////////////////

/*
Create a new sprite
*/
func NewSprite() *Sprite {
	this := new(Sprite)
	this.Animations 		= make(map[string]*SpriteAnimation)
	this.Visible 			= true
	this.Animated  			= true
	this.CurrentAnimation 	= "default"
	this.ZoomX 				= 1
	this.ZoomY 				= 1
	this.Red 				= 1
	this.Green 				= 1
	this.Blue 				= 1
	this.Alpha 				= 1
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


func (this *Sprite) AddEffect(options *EffectOptions) {
	if options.Animation == "" {
		options.Animation = "default"
	}

	options.durationTime = time.Millisecond * time.Duration(options.Duration)

	switch options.Effect {
		case ZOOM :		this.zoom(options)
		case FLIPX : 	this.flipX(options)
		case FLIPY :	this.flipY(options)
		case FADE : 	this.fade(options)
		case TURN:		this.turn(options)
		case HUE:		this.hue(options)
	}
}


func (this *Sprite) zoom(options *EffectOptions) {
	e := new(AnimationEffect)
	e.options 				= options
	e.zoomStart				= this.ZoomX
	this.Animations[options.Animation].Effect = e

	if options.Repeat == true {
		e.repeatCallback = func() {
			this.ZoomX = e.zoomStart 	// reset zoom
			e = nil 					// erase previous effect
			this.zoom(options)
		}
	}
}


func (this *Sprite) flipX(options *EffectOptions) {
	e := new(AnimationEffect)
	e.options 				= options
	e.zoomStart				= this.ZoomX
	this.Animations[options.Animation].Effect = e

	if options.Repeat == true {
		e.repeatCallback = func() {
			this.ZoomX = e.zoomStart 	// reset zoom
			e = nil 					// erase previous effect
			this.flipX(options)
		}
	}
}

func (this *Sprite) flipY(options *EffectOptions) {
	e := new(AnimationEffect)
	e.options 				= options
	e.zoomStart				= this.ZoomY
	this.Animations[options.Animation].Effect = e

	if options.Repeat == true {
		e.repeatCallback = func() {
			this.ZoomY = e.zoomStart 	// reset zoom
			e = nil 					// erase previous effect
			this.flipY(options)
		}
	}
}


func (this *Sprite) fade(options *EffectOptions) {
	e := new(AnimationEffect)
	e.options 				= options
	e.alphaStart			= this.Alpha
	this.Animations[options.Animation].Effect = e

	if options.Repeat == true {
		e.repeatCallback = func() {
			this.Alpha = e.alphaStart 	// reset alpha
			e = nil 					// erase previous effect
			this.fade(options)
		}
	}
}


func (this *Sprite) turn(options *EffectOptions) {
	e := new(AnimationEffect)
	e.options 				= options
	e.angleStart			= this.Angle
	this.Animations[options.Animation].Effect = e

	if options.Repeat == true {
		e.repeatCallback = func() {
			this.Angle = e.angleStart 	// reset alpha
			e = nil 				// erase previous effect
			this.turn(options)
		}
	}
}

func (this *Sprite) hue(options *EffectOptions) {
	e := new(AnimationEffect)
	e.options 			= options

	if e.options.Red == 0 {
		e.options.Red = 1
	}
	if e.options.Green == 0 {
		e.options.Green = 1
	}
	if e.options.Blue == 0 {
		e.options.Blue = 1
	}
	e.redStart			= this.Red
	e.greenStart		= this.Green
	e.blueStart			= this.Blue
	this.Animations[options.Animation].Effect = e

	if options.Repeat == true {
		e.repeatCallback = func() {
			this.Red 	= e.redStart
			this.Green 	= e.greenStart
			this.Blue 	= e.blueStart
			e = nil 					// erase previous effect
			this.hue(options)
		}
	}
}


/*
Return width of the current animation displayed
*/
func (this *Sprite) GetWidth() float64 {
	currentAnimation := this.Animations[this.CurrentAnimation]
	return float64(currentAnimation.StepWidth)
}

/*
Return height of the current animation displayed
*/
func (this *Sprite) GetHeight() float64 {
	currentAnimation := this.Animations[this.CurrentAnimation]
	return float64(currentAnimation.StepHeight)
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
func (this *Sprite) Position(arg... float64) (float64,float64) {
	if len(arg)==2 {
		this.X = arg[0]
		this.Y = arg[1]
	}
	return this.X, this.Y
}


/*
Set or retrieve Zoom factor

Exemple :
mySprite.Zoom(1.5)    // set both ZoomX and ZoomY to 1.5

mySprite.Zoom(1.5, 2) // set ZoomX to 1.5 and ZoomY to 2

zoomX, zoomY := mySprite.Zoom()
*/
func (this *Sprite) Zoom(arg... float64) (float64,float64) {
	if len(arg)==1 {
		this.ZoomX = arg[0]
		this.ZoomY = arg[0]
	} else if len(arg)==2 {
		this.ZoomX = arg[0]
		this.ZoomY = arg[1]
	}
	return this.ZoomX, this.ZoomY
}


/*
Set rotation angle (in degres)

Exemple :
mySprite.Rotate(45)    // the same as mySprite.Angle = 45
*/
func (this *Sprite) Rotate(angle float64) {
	this.Angle = angle
}


/*
Set or retrieve Skew factor (in degres)

Exemple :
mySprite.Skew(20)    // set both SkewX and SkewY to 20

mySprite.Skew(20, 40) // set SkewX to 20 and SkewY to 40

skewX, skewY := mySprite.Skew()
*/
func (this *Sprite) Skew(arg... float64) (float64,float64) {
	if len(arg)==1 {
		this.SkewX = arg[0]
		this.SkewY = arg[0]
	} else if len(arg)==2 {
		this.SkewX = arg[0]
		this.SkewY = arg[1]
	}
	return this.SkewX, this.SkewY
}


/*
Calculate new coordonnates and draw the sprite on the screen, after drawing, go to the next step of animation
*/
func (this *Sprite) Draw(surface *ebiten.Image) {
	if this.Visible {
		currentAnimation := this.Animations[this.CurrentAnimation] // SpriteAnimation object

		options := &ebiten.DrawImageOptions{}
		
		// move sprite x,y
		angleRad := this.Direction * math.Pi / 180 // convert degres into radians
		this.Y -= this.Speed * math.Sin(angleRad)
		this.X += this.Speed * math.Cos(angleRad)
		
		// if an animation is defined
		e := currentAnimation.Effect
		if e != nil {
			if e.options.Effect > 0 {
				// first drawing ? defined the time for first step
				if e.timeStart.IsZero() || e.timeStart.Unix()==0 {
					e.timeStart = time.Now()
					e.timeEnd 	= e.timeStart.Add(e.options.durationTime)
					//fmt.Printf("Demarre une animation %v\n            et la fin %v\n", e.timeStart, e.timeEnd)
				}

				now := time.Now()
				durationFromStart := now.Sub(e.timeStart)

				// animation not finished
				if e.timeEnd.Sub(now) > 0 {

					where := float64(durationFromStart.Nanoseconds()) / float64(e.options.durationTime.Nanoseconds())
					zoomFactor  := 1.0
					//fmt.Printf("Effect:%d\n",e.options.Effect )

					switch e.options.Effect {
						case ZOOM :
							if e.options.GoBack { // go and return
								var step float64 = 0.5
								if where < step {
									zoomFactor = convertRange(where,
													&Range{min:0,max:step},
													&Range{min:e.zoomStart,max:e.options.Zoom})
								} else {
									zoomFactor = convertRange(where,
													&Range{min:step,max:1},
													&Range{min:e.options.Zoom,max:e.zoomStart})
								}

							} else { // only one way
								zoomFactor = convertRange(where,
													&Range{min:0,max:1},
													&Range{min:e.zoomStart,max:e.options.Zoom})
							}
							this.ZoomX = zoomFactor
							this.ZoomY = zoomFactor
							///////////////////////////////////////////////

						case FLIPX :
							if e.options.GoBack { // go and return
								var step float64 = 0.25
								if 			where < step*1 {
									zoomFactor = convertRange(where, &Range{min:step*0,max:step*1},	&Range{min:1,max:0} )
								} else if 	where < step*2 {
								 	zoomFactor = convertRange(where, &Range{min:step*1,max:step*2}, &Range{min:0,max:-1} )
								} else if 	where < step*3 {
								 	zoomFactor = convertRange(where, &Range{min:step*2,max:step*3}, &Range{min:-1,max:0} )
								} else {
								 	zoomFactor = convertRange(where, &Range{min:step*3,max:step*4}, &Range{min:0,max:1} )
								}

							} else { // only one way
								var step float64 = 0.5
								if 			where < step*1 {
									zoomFactor = convertRange(where, &Range{min:step*0,max:step*1},	&Range{min:1,max:0} )
								} else {
								 	zoomFactor = convertRange(where, &Range{min:step*1,max:step*2}, &Range{min:0,max:-1} )
								}
							}
							this.ZoomX = zoomFactor
							///////////////////////////////////////////////

						case FLIPY :
							if e.options.GoBack { // go and return
								var step float64 = 0.25
								if 			where < step*1 {
									zoomFactor = convertRange(where, &Range{min:step*0,max:step*1},	&Range{min:1,max:0} )
								} else if 	where < step*2 {
								 	zoomFactor = convertRange(where, &Range{min:step*1,max:step*2}, &Range{min:0,max:-1} )
								} else if 	where < step*3 {
								 	zoomFactor = convertRange(where, &Range{min:step*2,max:step*3}, &Range{min:-1,max:0} )
								} else {
								 	zoomFactor = convertRange(where, &Range{min:step*3,max:step*4}, &Range{min:0,max:1} )
								}

							} else { // only one way
								var step float64 = 0.5
								if 			where < step*1 {
									zoomFactor = convertRange(where, &Range{min:step*0,max:step*1},	&Range{min:1,max:0} )
								} else {
								 	zoomFactor = convertRange(where, &Range{min:step*1,max:step*2}, &Range{min:0,max:-1} )
								}
							}
							this.ZoomY = zoomFactor
							///////////////////////////////////////////////

						case FADE :
							if e.options.GoBack { // go and return
								var step float64 = 0.5
								if 			where < step {
									this.Alpha = convertRange(where,
																&Range{min:0,max:step},
																&Range{min:e.options.FadeFrom,max:e.options.FadeTo})
								} else {
									this.Alpha = convertRange(where,
																&Range{min:step,max:1},
																&Range{min:e.options.FadeTo,max:e.options.FadeFrom})
								}

							} else { // only one way
								this.Alpha = convertRange(where,
															&Range{min:0,max:1},
															&Range{min:e.options.FadeFrom, max:e.options.FadeTo})
							}
							///////////////////////////////////////////////

						case TURN :
							clockwise := 1.0
							if e.options.Clockwise {
								clockwise = -1.0
							}

							if e.options.GoBack { // go and return
								var step float64 = 0.5
								if 			where < step {
									this.Angle = convertRange(where,
														&Range{min:0,max:step},
														&Range{min:0,max:e.options.Angle * clockwise})
								} else {
									this.Angle = convertRange(where,
														&Range{min:step*1,max:step*2},
														&Range{min:e.options.Angle * clockwise,max:0})
								}

							} else {
								this.Angle = convertRange(where,
														&Range{min:0,max:1},
														&Range{min:0,max:e.options.Angle * clockwise})
							}
							///////////////////////////////////////////////

							case HUE :
							if e.options.GoBack { // go and return
								var step float64 = 0.5
								if 			where < step {
									if e.options.Red != 1 {
										this.Red 	= convertRange(where, &Range{min:0,max:step}, &Range{min:1,max:e.options.Red})
									}
									if e.options.Green != 1 {
										this.Green 	= convertRange(where, &Range{min:0,max:step}, &Range{min:1,max:e.options.Green})
									}
									if e.options.Blue != 1 {
										this.Blue 	= convertRange(where, &Range{min:0,max:step}, &Range{min:1,max:e.options.Blue})
									}
								} else {
									if e.options.Red != 1 {
										this.Red 	= convertRange(where, &Range{min:step,max:1}, &Range{min:e.options.Red,max:1})
									}
									if e.options.Green != 1 {
										this.Green 	= convertRange(where, &Range{min:step,max:1}, &Range{min:e.options.Green,max:1})
									}
									if e.options.Blue != 1 {
										this.Blue 	= convertRange(where, &Range{min:step,max:1}, &Range{min:e.options.Blue,max:1})
									}
								}

							} else { // only one way
								if e.options.Red != 1 {
									this.Red 	= convertRange(where, &Range{min:0,max:1}, &Range{min:1, max:e.options.Red})
								}
								if e.options.Green != 1 {
									this.Green 	= convertRange(where, &Range{min:0,max:1}, &Range{min:1, max:e.options.Green})
								}
								if e.options.Blue != 1 {
									this.Blue 	= convertRange(where, &Range{min:0,max:1}, &Range{min:1, max:e.options.Blue})
								}
							}
							
							///////////////////////////////////////////////

					} // switch case

				// animation finished
				} else {
					// repeat animation
					if e.repeatCallback != nil {
						e.repeatCallback()
					}

					// laucnh user Callback
					if e.options.Callback != nil {
						e.options.Callback()
					}
				}
			} // if e.effect
		} // if e != nil

		
		// apply modification
		if this.CenterCoordonnates {
			options.GeoM.Translate(-float64(this.GetWidth())/2, -float64(this.GetHeight())/2)
		}
		options.GeoM.Scale(this.ZoomX, this.ZoomY)
		options.GeoM.Rotate( deg2rad(this.Angle) )
		options.GeoM.Translate(this.X , this.Y)

		options.GeoM.Skew( deg2rad(this.SkewX), deg2rad(this.SkewY) )
		
		// change HUE and Alpha
		options.ColorM.Scale(this.Red, this.Green, this.Blue, this.Alpha)

		// Choose current image inside animation
		x0 := currentAnimation.CurrentStep * currentAnimation.StepWidth
		x1 := x0 + currentAnimation.StepWidth
		r := image.Rect( x0 , 0, x1 , currentAnimation.StepHeight)
		options.SourceRect = &r

		surface.DrawImage(currentAnimation.Image, options)

		this.NextStep()
	}
}

func (this *Sprite) DrawBorders(surface *ebiten.Image, c color.Color) {
	var x,y,x1,y1 float64
	if this.CenterCoordonnates {
		x 	= math.Round(this.X - this.GetWidth()/2 * this.ZoomX )
		y 	= math.Round(this.Y - this.GetHeight()/2* this.ZoomY )
		
	} else {
		x 	= this.X
		y 	= this.Y
	}
	x1 	= math.Round(x + this.GetWidth() * this.ZoomX)
	y1 	= math.Round(y + this.GetHeight() * this.ZoomY)

	ebitenutil.DrawLine(surface, x,  y, x1,  y, c) 		// top
	ebitenutil.DrawLine(surface, x, y1, x1, y1, c)		// bottom
	ebitenutil.DrawLine(surface, x,  y,  x, y1, c)		// left
	ebitenutil.DrawLine(surface,x1,  y, x1, y1, c)		// right
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
Start the animation only one time (Reset+Show+Resume)

After running this, call the callback and pass the sprite pointer as argument
*/
func (this *Sprite) RunOnce( c func(*Sprite) ) {
	currentAnimation := this.Animations[this.CurrentAnimation]
	currentAnimation.RunOnce = true
	currentAnimation.callbackAfterRunOnce = c
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
				if currentAnimation.RunOnce {  // run only one time
					this.Stop()
					this.Hide()
					currentAnimation.callbackAfterRunOnce(this)

				} else {
					this.Reset() // restart at the end of the animation
				}
			}
			currentAnimation.currentStepTimeStart = now
			return true
		}
	}
	return false
}



func deg2rad(angle float64) float64 {
	return angle * math.Pi / -180
}


type Range struct {
	min, max float64
}

func convertRange(oldValue float64, oldRange, newRange *Range) float64 {
	oldDelta := (oldRange.max - oldRange.min)
	newDelta := (newRange.max - newRange.min) 
	return (((oldValue - oldRange.min) * newDelta) / oldDelta) + newRange.min
}