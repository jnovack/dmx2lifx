package main

import (
	"fmt"

	"github.com/lucasb-eyer/go-colorful"
)

func main() {
	red := 0
	green := 0
	blue := 0
	for i := 32; i < 257; i = i + 32 {
		blue = i - 1
		r := float64(red) / 255
		g := float64(green) / 255
		b := float64(blue) / 255

		rgb := colorful.Color{R: r, G: g, B: b}

		h, s, v := rgb.Hsv()

		fmt.Printf("%d, %d, %d\t\t%d\t%d\t%d\n", red, green, blue, int(h/360*65535), int(s*65535), int(v*65535))
	}
}
