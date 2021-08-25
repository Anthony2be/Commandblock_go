package main

import "github.com/anthony2be/commandblock_go/datapack"

func main() {
	e := datapack.New("hi", "lllllllllllll", 7, "pp", "poopoo", "idk")
	var h [1]string
	h[0] = "h"
	e.RegisterFunction("pp", h[0])
	e.Abort(true)
}