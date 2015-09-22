// A trivial implementation of Snake, just to give Gopher a try.
// "Do what you like with it"-license. (c) Richard van Roy

package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gopherjs/gopherjs/js"
)

type Point struct {
	x, y int
}

const (
	RIGHT = 0
	LEFT  = 1
	UP    = 2
	DOWN  = 3
)

type Direction int

type State struct {
	speed          int
	worldW, worldH int
	direction      Direction
	snake          []Point
	apples         []Point
}

func clear(canvas *js.Object) {
	w := canvas.Get("width")
	h := canvas.Get("height")
	context := canvas.Call("getContext", "2d")
	context.Set("fillStyle", "black")
	context.Call("fillRect", 0, 0, w, h)
}

func drawSnake(canvas *js.Object, state State) {
	context := canvas.Call("getContext", "2d")
	for _, point := range state.snake {
		context.Set("fillStyle", "red")
		context.Call("fillRect", point.x*20, point.y*20, 20, 20)
	}
}

func drawApple(canvas *js.Object, state State) {
	context := canvas.Call("getContext", "2d")
	for _, point := range state.apples {
		context.Set("fillStyle", "green")
		context.Call("fillRect", point.x*20, point.y*20, 20, 20)
	}
}

func randomPointNotInSnake(state State) Point {
	randomPoint := Point{rand.Intn(state.worldW), rand.Intn(state.worldH)}
	for _, point := range state.snake {
		if point == randomPoint {
			return randomPointNotInSnake(state)
		} else {
			return randomPoint
		}
	}
	return randomPoint
}

func nextPoint(state State) Point {
	x := state.snake[0].x
	y := state.snake[0].y
	w := state.worldW
	h := state.worldH

	plusX := 0
	plusY := 0
	switch state.direction {
	case LEFT:
		plusX = -1
	case RIGHT:
		plusX = 1
	case UP:
		plusY = -1
	case DOWN:
		plusY = 1
	}

	newX := (x + plusX) % w
	newY := (y + plusY) % h
	if newX < 0 {
		newX += w
	}
	if newY < 0 {
		newY += h
	}
	return Point{newX, newY}
}

func moveSnake(state *State) {
	for i := len(state.snake) - 1; i > 0; i-- {
		state.snake[i] = state.snake[i-1]
	}

	tail := state.snake[len((*state).snake)-1]
	onApple := false
	var newApples []Point
	for _, apple := range state.apples {
		if state.snake[0] == apple {
			onApple = true
		} else {
			newApples = append(newApples, apple)
		}
	}
	(*state).apples = newApples
	if onApple {
		state.snake = append(state.snake, tail)
		fmt.Println(state.snake)
	}
	(*state).snake[0] = nextPoint(*state)

}

func isGameOver(state State) bool {

	head := state.snake[0]
	for _, point := range state.snake[1:] {
		if head == point {
			return true
		}
	}
	return false
}

func newState() State {
	snake := []Point{Point{0, 0}}
	state := State{100, 10, 10, RIGHT, snake, []Point{}}
	return state
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	document := js.Global.Get("document")
	canvas := document.Call("createElement", "canvas")
	document.Get("body").Call("appendChild", canvas)

	state := newState()

	update := func() {
		if len(state.apples) == 0 {
			state.apples = append(state.apples, randomPointNotInSnake(state))
		}
		if isGameOver(state) {
			state = newState()
		}
		moveSnake(&state)
		clear(canvas)
		drawApple(canvas, state)
		drawSnake(canvas, state)
	}

	canvas.Set("width", state.worldW*20)
	canvas.Set("height", state.worldH*20)

	update()

	window := js.Global.Get("window")
	window.Call("setInterval", update, state.speed)

	window.Set("onkeydown", func(ev *js.Object) {
		keyCode := ev.Get("keyCode").Int()
		fmt.Println(keyCode)
		switch {
		case keyCode == 37:
			state.direction = LEFT

		case keyCode == 38:
			state.direction = UP
		case keyCode == 39:
			state.direction = RIGHT
		case keyCode == 40:
			state.direction = DOWN
		}
	})
}
