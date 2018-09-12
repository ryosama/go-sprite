Go-sprite
=======

A simple library for playing with sprites and animations

It use the [Ebiten](https://github.com/hajimehoshi/ebiten) library for the 2D graphics engine

Install
=======

```bash
$ go get -u github.com/hajimehoshi/ebiten
$ go get -u github.com/ryosama/go-sprite
```

Screenshot
===========

![Screenshot](https://github.com/ryosama/go-sprite/raw/master/screenshot1.png "Screenshot")

Quick Start
===========

```Go
import "github.com/ryosama/go-sprite"

mySprite = sprite.NewSprite()
mySprite.AddAnimation("walk-right",	"walk_right.png", 700, 6, ebiten.FilterDefault)
mySprite.Position(WINDOW_WIDTH/2, WINDOW_HEIGHT/2)
mySprite.CurrentAnimation = "walk-right"
mySprite.Speed = 2
mySprite.Start()
```

Documentation
=============

The documentation can be found here : https://godoc.org/github.com/ryosama/go-sprite

Or export with this command

```bash
$ godoc github.com/ryosama/go-sprite
```

TODO
====

- Add a video

- Add pre-defined animation (inflate, deflate, ...), like in font-awesome-animation