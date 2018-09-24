/*Package sprite permits animations of sprites via the Ebiten library (http://www.github.com/hajimehoshi/ebiten)

Basic Usage :

import "github.com/ryosama/go-sprite"

mySprite = sprite.NewSprite()
mySprite.AddAnimation("walk-right",	"walk_right.png", 700, 6, ebiten.FilterDefault)
mySprite.Position(WINDOW_WIDTH/2, WINDOW_HEIGHT/2)
mySprite.CurrentAnimation = "walk-right"
mySprite.Speed = 2
mySprite.Start()
*/package sprite

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image"
	"image/color"
	"log"
	"math"
	"time"
	//"fmt"
)

// Constant effects
const (
	NoEffect = iota

	// Zoom multiplier on the sprite
	Zoom

	// Flip the sprite horizontaly or verticaly
	Flip

	// Fade off or fade in the sprite
	Fade

	// Rotate the sprite
	Turn

	// Play with Hue top change color of the sprite
	Hue

	// Move absoluty the sprite
	Move

	// For Flip effect
	Horizontaly = false
	Verticaly   = true
)

var violet = color.RGBA{R: 255, G: 0, B: 255, A: 255}

/*Sprite contains the sprite, animations and effects */
type Sprite struct {
	// Animation label currently displayed
	CurrentAnimation string

	// Array of animations
	Animations map[string]*Animation

	// X coordinates of the sprite (in pixel)
	X float64

	// Y coordinate of the sprite (in pixel)
	Y float64

	// Speed is in pixel/frame
	Speed float64

	// Direction is an Angle in degres
	Direction float64

	// Zoom in or out on X axis
	ZoomX float64

	// Zoom in or out on Y axis
	ZoomY float64

	// Colors multipliers
	Red, Green, Blue float64

	// Transparency
	Alpha float64

	// Angle of rotation in degres
	Angle float64

	// Skew on X axis in degres
	SkewX float64

	// Skew on Y axis in degres
	SkewY float64

	// Visibility of the sprite
	Visible bool

	// Animated or not
	Animated bool

	// Displace X and Y coordonnate to the center of the sprite
	CenterCoordonnates bool

	// Draw debug borders around sprite
	Borders bool
}

/*Animation contains animations and effects */
type Animation struct {
	// File path of the animation
	// Step of the animation must be the same width on one line
	Path string

	// ebiten.Image generated
	Image *ebiten.Image

	// Number of steps for the total animation
	Steps int

	// Current step displayed
	CurrentStep int

	// Where to start the animation
	FirstStep int

	// Width of the animation steps (in pixel)
	StepWidth int

	// Height of the animation steps (in pixel)
	StepHeight int

	// Total duration of the animation in millisecond
	Duration time.Duration

	// Total time for one step in millisecond
	OneStepDuration time.Duration

	// Effects object
	Effects []*animationEffect

	// Animation once and disapared
	RunOnce bool

	// Callback after run once
	callbackAfterRunOnce func(*Sprite)

	// Start time of the current step
	currentStepTimeStart time.Time
}

type animationEffect struct {
	options                                     *EffectOptions
	zoomStart                                   float64
	redStart, greenStart, blueStart, alphaStart float64
	angleStart                                  float64
	xStart, yStart                              float64
	duration                                    time.Duration
	timeStart, timeEnd                          time.Time
	repeatCallback                              func()
}

//EffectOptions contains options for the effect
type EffectOptions struct {
	// Name of animation (default is omitted)
	Animation string

	// Effect= Zoom, FlipX, FlipY, Fade, Turn, Move
	Effect int

	// For Fade and FadeINOUT effects
	FadeFrom, FadeTo float64

	// For Zoom effects
	Zoom float64

	// For Turn effect
	Clockwise bool

	// For Turn effect (in degres)
	Angle float64

	// Horizontaly or Verticaly
	Axis bool

	// For Hue effect
	Red, Green, Blue float64

	// For Move effect
	X, Y float64

	// Duration of the effect
	Duration int

	// Duration of the effect (in time.Duration)
	durationTime time.Duration

	// Redo the animation on the counter way
	GoBack bool

	// Repeat or not at the end of effect
	Repeat bool

	// function to launch afert one complete effect
	Callback func()

	// index of the effect in the stack
	index int

	// Count the number of loop in animation
	loopCounter int64
}

//////////////////////////////////////////// CONSTRUCTORS ////////////////////////////////////////////

//NewSprite creates a new sprite
func NewSprite() *Sprite {
	sprite := new(Sprite)
	sprite.Animations = make(map[string]*Animation)
	sprite.Visible = true
	sprite.Animated = true
	sprite.CurrentAnimation = "default"
	sprite.ZoomX = 1
	sprite.ZoomY = 1
	sprite.Red = 1
	sprite.Green = 1
	sprite.Blue = 1
	sprite.Alpha = 1
	return sprite
}

func newAnimation(path string, duration int, steps int, filter ebiten.Filter) *Animation {
	var err error
	animation := new(Animation)
	animation.Path = path
	animation.Image, _, err = ebitenutil.NewImageFromFile(path, filter)
	if err != nil {
		log.Fatal(err)
	}
	animation.Steps = steps
	animation.Duration = time.Millisecond * time.Duration(duration)

	width, height := animation.Image.Size()
	animation.StepWidth = width / animation.Steps
	animation.StepHeight = height

	animation.currentStepTimeStart = time.Now()
	animation.OneStepDuration = time.Duration(int(animation.Duration) / animation.Steps)

	animation.Effects = make([]*animationEffect, 0)

	return animation
}

//////////////////////////////////////////// METHODS ////////////////////////////////////////////

/*
AddAnimation adds an animation to the sprite

"label" is the tag for the animation

"path" is the path for the image file

"duration" is in millisecond

"steps" is the number of step for the animation

"filter" is ebiten.FilterDefault or ebiten.FilterNearest  or ebiten.FilterLinear

Example :

mySprite.AddAnimation("walk-right",	"walk_right.png", 700, 6, ebiten.FilterDefault)
*/
func (sprite *Sprite) AddAnimation(label string, path string, duration int, steps int, filter ebiten.Filter) {
	sprite.Animations[label] = newAnimation(path, duration, steps, filter)
}

/*
AddEffect adds an effect to the sprite. You can cumulate effects at the same time

Example :

sprites[i].AddEffect(&sprite.EffectOptions{ Effect: sprite.Zoom, Zoom:3, Duration:2000, Repeat:true, GoBack:true })

*/
func (sprite *Sprite) AddEffect(options *EffectOptions) {
	if options.Animation == "" {
		options.Animation = "default"
	}

	options.durationTime = time.Millisecond * time.Duration(options.Duration)

	switch options.Effect {
	case Zoom:
		sprite.zoom(options)
	case Flip:
		sprite.flip(options)
	case Fade:
		sprite.fade(options)
	case Turn:
		sprite.turn(options)
	case Hue:
		sprite.hue(options)
	case Move:
		sprite.move(options)
	}
}

func (sprite *Sprite) zoom(options *EffectOptions) {
	e := new(animationEffect)
	e.options = options
	e.zoomStart = sprite.ZoomX

	if options.loopCounter == 0 { // first loop
		sprite.Animations[options.Animation].Effects = append(sprite.Animations[options.Animation].Effects, e)
		options.index = len(sprite.Animations[options.Animation].Effects) - 1 // store index
	} else {
		sprite.Animations[options.Animation].Effects[options.index] = e
	}

	if options.Repeat == true {
		e.repeatCallback = func() {
			sprite.ZoomX = e.zoomStart // reset zoom
			e = nil                    // erase previous effect
			sprite.zoom(options)
		}
	}

	options.loopCounter++
}

func (sprite *Sprite) flip(options *EffectOptions) {
	e := new(animationEffect)
	e.options = options
	if e.options.Axis == Horizontaly {
		e.zoomStart = sprite.ZoomX
	} else {
		e.zoomStart = sprite.ZoomY
	}

	if options.loopCounter == 0 { // first loop
		sprite.Animations[options.Animation].Effects = append(sprite.Animations[options.Animation].Effects, e)
		options.index = len(sprite.Animations[options.Animation].Effects) - 1 // store index
	} else {
		sprite.Animations[options.Animation].Effects[options.index] = e
	}

	if options.Repeat == true {
		e.repeatCallback = func() {
			if e.options.Axis == Horizontaly {
				sprite.ZoomX = e.zoomStart // reset zoom
			} else {
				sprite.ZoomY = e.zoomStart // reset zoom
			}
			e = nil // erase previous effect
			sprite.flip(options)
		}
	}

	options.loopCounter++
}

func (sprite *Sprite) fade(options *EffectOptions) {
	e := new(animationEffect)
	e.options = options
	e.alphaStart = sprite.Alpha

	if options.loopCounter == 0 { // first loop
		sprite.Animations[options.Animation].Effects = append(sprite.Animations[options.Animation].Effects, e)
		options.index = len(sprite.Animations[options.Animation].Effects) - 1 // store index
	} else {
		sprite.Animations[options.Animation].Effects[options.index] = e
	}

	if options.Repeat == true {
		e.repeatCallback = func() {
			sprite.Alpha = e.alphaStart // reset alpha
			e = nil                     // erase previous effect
			sprite.fade(options)
		}
	}

	options.loopCounter++
}

func (sprite *Sprite) turn(options *EffectOptions) {
	e := new(animationEffect)
	e.options = options
	e.angleStart = sprite.Angle

	if options.loopCounter == 0 { // first loop
		sprite.Animations[options.Animation].Effects = append(sprite.Animations[options.Animation].Effects, e)
		options.index = len(sprite.Animations[options.Animation].Effects) - 1 // store index
	} else {
		sprite.Animations[options.Animation].Effects[options.index] = e
	}

	if options.Repeat == true {
		e.repeatCallback = func() {
			sprite.Angle = e.angleStart // reset alpha
			e = nil                     // erase previous effect
			sprite.turn(options)
		}
	}

	options.loopCounter++
}

func (sprite *Sprite) hue(options *EffectOptions) {
	e := new(animationEffect)
	e.options = options

	// init value
	if e.options.Red == 0 {
		e.options.Red = 1
	}
	if e.options.Green == 0 {
		e.options.Green = 1
	}
	if e.options.Blue == 0 {
		e.options.Blue = 1
	}

	e.redStart = sprite.Red
	e.greenStart = sprite.Green
	e.blueStart = sprite.Blue

	if options.loopCounter == 0 { // first loop
		sprite.Animations[options.Animation].Effects = append(sprite.Animations[options.Animation].Effects, e)
		options.index = len(sprite.Animations[options.Animation].Effects) - 1 // store index
	} else {
		sprite.Animations[options.Animation].Effects[options.index] = e
	}

	if options.Repeat == true {
		e.repeatCallback = func() {
			//fmt.Printf("DEBUG z.redStart:%f  z.greenStart:%f  z.blueStart:%f\n", e.redStart, e.greenStart, e.blueStart)
			sprite.Red = e.redStart
			sprite.Green = e.greenStart
			sprite.Blue = e.blueStart
			e = nil // erase previous effect
			sprite.hue(options)
		}
	}

	options.loopCounter++
}

func (sprite *Sprite) move(options *EffectOptions) {
	e := new(animationEffect)
	e.options = options

	if e.options.X == 0 {
		e.options.X = sprite.X
	}
	if e.options.Y == 0 {
		e.options.Y = sprite.Y
	}

	e.xStart = sprite.X
	e.yStart = sprite.Y

	if options.loopCounter == 0 { // first loop
		sprite.Animations[options.Animation].Effects = append(sprite.Animations[options.Animation].Effects, e)
		options.index = len(sprite.Animations[options.Animation].Effects) - 1 // store index
	} else {
		sprite.Animations[options.Animation].Effects[options.index] = e
	}

	if options.Repeat == true {
		e.repeatCallback = func() {
			sprite.X = e.xStart // reset x position
			sprite.Y = e.yStart // reset y position
			e = nil             // erase previous effect
			sprite.move(options)
		}
	}

	options.loopCounter++
}

//GetWidth returns width of the current animation displayed
func (sprite *Sprite) GetWidth() float64 {
	currentAnimation := sprite.Animations[sprite.CurrentAnimation]
	return float64(currentAnimation.StepWidth)
}

//GetHeight returns height of the current animation displayed
func (sprite *Sprite) GetHeight() float64 {
	currentAnimation := sprite.Animations[sprite.CurrentAnimation]
	return float64(currentAnimation.StepHeight)
}

//Hide the sprite
func (sprite *Sprite) Hide() {
	sprite.Visible = false
}

//Show the sprite
func (sprite *Sprite) Show() {
	sprite.Visible = true
}

//ToogleVisibility toogle visibility of the sprite
func (sprite *Sprite) ToogleVisibility() {
	if sprite.Visible {
		sprite.Hide()
	} else {
		sprite.Show()
	}
}

/*
Position sets or retrieve X and Y coordonnates of the sprite

Return X and Y coordonnates

Exemple :
mySprite.Position(WINDOW_WIDTH/2, WINDOW_HEIGHT/2)

or

x,y := mySprite.Position()
*/
func (sprite *Sprite) Position(arg ...float64) (float64, float64) {
	if len(arg) == 2 {
		sprite.X = arg[0]
		sprite.Y = arg[1]
	}
	return sprite.X, sprite.Y
}

/*
Zoom sets or retrieve Zoom factor

Exemple :
mySprite.Zoom(1.5)    // set both ZoomX and ZoomY to 1.5

mySprite.Zoom(1.5, 2) // set ZoomX to 1.5 and ZoomY to 2

zoomX, zoomY := mySprite.Zoom()
*/
func (sprite *Sprite) Zoom(arg ...float64) (float64, float64) {
	if len(arg) == 1 {
		sprite.ZoomX = arg[0]
		sprite.ZoomY = arg[0]
	} else if len(arg) == 2 {
		sprite.ZoomX = arg[0]
		sprite.ZoomY = arg[1]
	}
	return sprite.ZoomX, sprite.ZoomY
}

/*
Rotate sets rotation angle (in degres)

Exemple :
mySprite.Rotate(45)    // the same as mySprite.Angle = 45
*/
func (sprite *Sprite) Rotate(angle float64) {
	sprite.Angle = angle
}

/*
Skew sets or retrieve Skew factor (in degres)

Exemple :
mySprite.Skew(20)    // set both SkewX and SkewY to 20

mySprite.Skew(20, 40) // set SkewX to 20 and SkewY to 40

skewX, skewY := mySprite.Skew()
*/
func (sprite *Sprite) Skew(arg ...float64) (float64, float64) {
	if len(arg) == 1 {
		sprite.SkewX = arg[0]
		sprite.SkewY = arg[0]
	} else if len(arg) == 2 {
		sprite.SkewX = arg[0]
		sprite.SkewY = arg[1]
	}
	return sprite.SkewX, sprite.SkewY
}

//Draw calculates new coordonnates and draw the sprite on the screen, after drawing, go to the next step of animation
func (sprite *Sprite) Draw(surface *ebiten.Image) {
	if sprite.Visible {
		currentAnimation := sprite.Animations[sprite.CurrentAnimation] // Animation object

		options := &ebiten.DrawImageOptions{}

		// move sprite x,y
		angleRad := sprite.Direction * math.Pi / 180 // convert degres into radians
		sprite.Y -= sprite.Speed * math.Sin(angleRad)
		sprite.X += sprite.Speed * math.Cos(angleRad)

		// apply diffrents effects
		sprite.applyEffects(surface)

		// apply modification
		if sprite.CenterCoordonnates {
			options.GeoM.Translate(-float64(sprite.GetWidth())/2, -float64(sprite.GetHeight())/2)
		}
		options.GeoM.Scale(sprite.ZoomX, sprite.ZoomY)
		options.GeoM.Rotate(deg2rad(sprite.Angle))
		options.GeoM.Translate(sprite.X, sprite.Y)

		options.GeoM.Skew(deg2rad(sprite.SkewX), deg2rad(sprite.SkewY))

		// change Hue and Alpha
		options.ColorM.Scale(sprite.Red, sprite.Green, sprite.Blue, sprite.Alpha)

		// Choose current image inside animation
		x0 := currentAnimation.CurrentStep * currentAnimation.StepWidth
		x1 := x0 + currentAnimation.StepWidth
		r := image.Rect(x0, 0, x1, currentAnimation.StepHeight)
		options.SourceRect = &r

		if sprite.Borders {
			sprite.DrawBorders(surface, violet)
		}

		surface.DrawImage(currentAnimation.Image, options)

		sprite.NextStep()
	}
}

//DrawBorders draw debug borders around the sprite
func (sprite *Sprite) DrawBorders(surface *ebiten.Image, c color.Color) {
	var x, y, x1, y1 float64
	if sprite.CenterCoordonnates {
		x = math.Round(sprite.X - sprite.GetWidth()/2*sprite.ZoomX)
		y = math.Round(sprite.Y - sprite.GetHeight()/2*sprite.ZoomY)

	} else {
		x = sprite.X
		y = sprite.Y
	}
	x1 = math.Round(x + sprite.GetWidth()*sprite.ZoomX)
	y1 = math.Round(y + sprite.GetHeight()*sprite.ZoomY)

	ebitenutil.DrawLine(surface, x, y, x1, y, c)   // top
	ebitenutil.DrawLine(surface, x, y1, x1, y1, c) // bottom
	ebitenutil.DrawLine(surface, x, y, x, y1, c)   // left
	ebitenutil.DrawLine(surface, x1, y, x1, y1, c) // right
}

//Start the animation (Reset+Show+Resume)
func (sprite *Sprite) Start() {
	sprite.Reset()
	sprite.Show()
	sprite.Resume()
}

/*
RunOnce start the animation only one time (Reset+Show+Resume)

After running animation, call the callback and pass the sprite pointer as argument
*/
func (sprite *Sprite) RunOnce(c func(*Sprite)) {
	currentAnimation := sprite.Animations[sprite.CurrentAnimation]
	currentAnimation.RunOnce = true
	currentAnimation.callbackAfterRunOnce = c
	sprite.Reset()
	sprite.Show()
	sprite.Resume()
}

//Stop the animation (Reset+Pause)
func (sprite *Sprite) Stop() {
	sprite.Reset()
	sprite.Pause()
}

//Reset current step to the first step of the animation
func (sprite *Sprite) Reset() {
	currentAnimation := sprite.Animations[sprite.CurrentAnimation]
	currentAnimation.CurrentStep = currentAnimation.FirstStep
}

//Pause the animation
func (sprite *Sprite) Pause() {
	sprite.Animated = false
}

//Resume the animation
func (sprite *Sprite) Resume() {
	sprite.Animated = true
}

//ToogleAnimation toogle animation status
func (sprite *Sprite) ToogleAnimation() {
	if sprite.Animated {
		sprite.Pause()
	} else {
		sprite.Resume()
	}
}

/*
NextStep go to the next step of animation

Return true if animation go to the next step or false if step duration is not finish
*/
func (sprite *Sprite) NextStep() bool {
	currentAnimation := sprite.Animations[sprite.CurrentAnimation]
	if sprite.Animated {
		now := time.Now()
		nextStepAt := currentAnimation.currentStepTimeStart.Add(currentAnimation.OneStepDuration)

		if now.Sub(nextStepAt) > 0 { // time to change the current step
			currentAnimation.CurrentStep++ // next step
			if currentAnimation.CurrentStep+1 > currentAnimation.Steps {
				if currentAnimation.RunOnce { // run only one time
					sprite.Stop()
					sprite.Hide()
					currentAnimation.callbackAfterRunOnce(sprite)

				} else {
					sprite.Reset() // restart at the end of the animation
				}
			}
			currentAnimation.currentStepTimeStart = now
			return true
		}
	}
	return false
}

func (sprite *Sprite) applyEffects(surface *ebiten.Image) {
	currentAnimation := sprite.Animations[sprite.CurrentAnimation]

	for _, e := range currentAnimation.Effects { // foreach Effects in the stack

		// if an animation is defined
		//e := currentAnimation.Effect
		if e != nil {
			if e.options.Effect > 0 {
				// first drawing ? defined the time for first step
				if e.timeStart.IsZero() || e.timeStart.Unix() == 0 {
					e.timeStart = time.Now()
					e.timeEnd = e.timeStart.Add(e.options.durationTime)
					//fmt.Printf("Demarre une animation %v\n            et la fin %v\n", e.timeStart, e.timeEnd)
				}

				now := time.Now()
				durationFromStart := now.Sub(e.timeStart)

				// animation not finished
				if e.timeEnd.Sub(now) > 0 {

					where := float64(durationFromStart.Nanoseconds()) / float64(e.options.durationTime.Nanoseconds())
					zoomFactor := 1.0
					//fmt.Printf("Effect:%d\n",e.options.Effect )

					switch e.options.Effect {
					case Zoom:
						if e.options.GoBack { // go and return
							step := 0.5
							if where < step {
								zoomFactor = convertScale(where, &scale{min: 0, max: step}, &scale{min: e.zoomStart, max: e.options.Zoom})
							} else {
								zoomFactor = convertScale(where, &scale{min: step, max: 1}, &scale{min: e.options.Zoom, max: e.zoomStart})
							}

						} else { // only one way
							zoomFactor = convertScale(where, &scale{min: 0, max: 1}, &scale{min: e.zoomStart, max: e.options.Zoom})
						}
						sprite.ZoomX = zoomFactor
						sprite.ZoomY = zoomFactor
						///////////////////////////////////////////////

					case Flip:
						if e.options.GoBack { // go and return
							step := 0.25
							if where < step*1 {
								zoomFactor = convertScale(where, &scale{min: step * 0, max: step * 1}, &scale{min: 1, max: 0})
							} else if where < step*2 {
								zoomFactor = convertScale(where, &scale{min: step * 1, max: step * 2}, &scale{min: 0, max: -1})
							} else if where < step*3 {
								zoomFactor = convertScale(where, &scale{min: step * 2, max: step * 3}, &scale{min: -1, max: 0})
							} else {
								zoomFactor = convertScale(where, &scale{min: step * 3, max: step * 4}, &scale{min: 0, max: 1})
							}

						} else { // only one way
							step := 0.5
							if where < step*1 {
								zoomFactor = convertScale(where, &scale{min: step * 0, max: step * 1}, &scale{min: 1, max: 0})
							} else {
								zoomFactor = convertScale(where, &scale{min: step * 1, max: step * 2}, &scale{min: 0, max: -1})
							}
						}
						if e.options.Axis == Horizontaly {
							sprite.ZoomX = zoomFactor
						} else {
							sprite.ZoomY = zoomFactor
						}
						///////////////////////////////////////////////

					case Fade:
						if e.options.GoBack { // go and return
							step := 0.5
							if where < step {
								sprite.Alpha = convertScale(where, &scale{min: 0, max: step}, &scale{min: e.options.FadeFrom, max: e.options.FadeTo})
							} else {
								sprite.Alpha = convertScale(where, &scale{min: step, max: 1}, &scale{min: e.options.FadeTo, max: e.options.FadeFrom})
							}

						} else { // only one way
							sprite.Alpha = convertScale(where, &scale{min: 0, max: 1}, &scale{min: e.options.FadeFrom, max: e.options.FadeTo})
						}
						///////////////////////////////////////////////

					case Turn:
						clockwise := 1.0
						if e.options.Clockwise {
							clockwise = -1.0
						}

						if e.options.GoBack { // go and return
							step := 0.5
							if where < step {
								sprite.Angle = convertScale(where, &scale{min: 0, max: step}, &scale{min: 0, max: e.options.Angle * clockwise})
							} else {
								sprite.Angle = convertScale(where, &scale{min: step * 1, max: step * 2}, &scale{min: e.options.Angle * clockwise, max: 0})
							}

						} else {
							sprite.Angle = convertScale(where, &scale{min: 0, max: 1}, &scale{min: 0, max: e.options.Angle * clockwise})
						}
					///////////////////////////////////////////////

					case Hue:
						if e.options.GoBack { // go and return
							step := 0.5
							if where < step {
								sprite.Red = convertScale(where, &scale{min: 0, max: step}, &scale{min: e.redStart, max: e.options.Red})
								sprite.Green = convertScale(where, &scale{min: 0, max: step}, &scale{min: e.greenStart, max: e.options.Green})
								sprite.Blue = convertScale(where, &scale{min: 0, max: step}, &scale{min: e.blueStart, max: e.options.Blue})
							} else {
								sprite.Red = convertScale(where, &scale{min: step, max: 1}, &scale{min: e.options.Red, max: e.redStart})
								sprite.Green = convertScale(where, &scale{min: step, max: 1}, &scale{min: e.options.Green, max: e.greenStart})
								sprite.Blue = convertScale(where, &scale{min: step, max: 1}, &scale{min: e.options.Blue, max: e.blueStart})
							}

						} else { // only one way
							sprite.Red = convertScale(where, &scale{min: 0, max: 1}, &scale{min: e.redStart, max: e.options.Red})
							sprite.Green = convertScale(where, &scale{min: 0, max: 1}, &scale{min: e.greenStart, max: e.options.Green})
							sprite.Blue = convertScale(where, &scale{min: 0, max: 1}, &scale{min: e.blueStart, max: e.options.Blue})
						}
						///////////////////////////////////////////////

					case Move:
						if e.options.GoBack { // go and return
							step := 0.5
							if where < step {
								sprite.X = convertScale(where, &scale{min: 0, max: step}, &scale{min: e.xStart, max: e.options.X})
								sprite.Y = convertScale(where, &scale{min: 0, max: step}, &scale{min: e.yStart, max: e.options.Y})
							} else {
								sprite.X = convertScale(where, &scale{min: step, max: 1}, &scale{min: e.options.X, max: e.xStart})
								sprite.Y = convertScale(where, &scale{min: step, max: 1}, &scale{min: e.options.Y, max: e.yStart})
							}

						} else { // only one way
							sprite.X = convertScale(where, &scale{min: 0, max: 1}, &scale{min: e.xStart, max: e.options.X})
							sprite.Y = convertScale(where, &scale{min: 0, max: 1}, &scale{min: e.yStart, max: e.options.Y})
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
	} // foreach Effect
}

//////////////////////////////////////////// TOOLS ////////////////////////////////////////////////:

func deg2rad(angle float64) float64 {
	return angle * math.Pi / -180
}

type scale struct {
	min, max float64
}

func convertScale(oldValue float64, oldRange, newRange *scale) float64 {
	oldDelta := (oldRange.max - oldRange.min)
	newDelta := (newRange.max - newRange.min)
	return (((oldValue - oldRange.min) * newDelta) / oldDelta) + newRange.min
}
