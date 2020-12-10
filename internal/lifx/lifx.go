package lifx

import (
	"fmt"
	"sort"

	"github.com/2tvenom/golifx"
	"github.com/lucasb-eyer/go-colorful"
)

var bulbs []*golifx.Bulb

func label(bulb *golifx.Bulb) string {
	label, _ := bulb.GetLabel()
	return label
}

func init() {
	//Lookup all bulbs
	bulbs, _ = golifx.LookupBulbs()

	// Sort by labels for consistent ordering
	sort.SliceStable(bulbs, func(i, j int) bool {
		return label(bulbs[i]) < label(bulbs[j])
	})

	for b := range bulbs {
		//Get power state
		powerState, _ := bulbs[b].GetPowerState()

		//Turn if off
		if !powerState {
			bulbs[b].SetPowerState(true)
		}

		fmt.Printf("Found Bulb: %s", label(bulbs[b]))
	}
}

// Set a bulb to a specific RGB color
func Set(bulb int, red int, green int, blue int) {
	r := float64(red) / 255
	g := float64(green) / 255
	b := float64(blue) / 255

	rgb := colorful.Color{R: r, G: g, B: b}

	h, s, v := rgb.Hsv()

	fmt.Printf("%d, %d, %d\t\t%d\t%d\t%d\n", red, green, blue, int(h/360*65535), int(s*65535), int(v*65535))

	hsbk := &golifx.HSBK{
		Hue:        uint16(h / 360 * 65535),
		Saturation: uint16(s * 65535),
		Brightness: uint16(v * 65535),
		Kelvin:     3200,
	}

	bulbs[bulb].SetColorState(hsbk, 500)

}
