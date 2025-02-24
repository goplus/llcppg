package main

import (
	"raylib"

	"github.com/goplus/llgo/c"
)

func main() {
	const screenWidth = 800
	const screenHeight = 450
	blue := raylib.GetColor(0xFF00F179)
	white := raylib.GetColor(0xFFFFFFFF)
	raylib.InitWindow(screenWidth, screenHeight, c.Str("Raylib DEMO"))
	startTime := raylib.GetTime()
	for raylib.WindowShouldClose() == 0 {
		currentTime := raylib.GetTime()
		if currentTime-startTime >= 1.0 {
			break
		}
		raylib.BeginDrawing()
		raylib.ClearBackground(white)
		raylib.DrawRectangle(screenWidth/2-50, screenHeight/2-50, 100, 100, blue)
		raylib.EndDrawing()
	}
}
