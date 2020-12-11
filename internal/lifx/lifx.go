package lifx

import (
	"sort"

	"github.com/jnovack/dmx2lifx/pkg/golifx"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/rs/zerolog/log"
)

var bulbs []*golifx.Bulb

func label(bulb *golifx.Bulb) string {
	label, _ := bulb.GetLabel()
	return label
}

// Count the number of bulbs
func Count() int {
	return len(bulbs)
}

func init() {
	//Lookup all bulbs
	log.Warn().Msg("scanning for lifx bulbs (this could take some time)")
	bulbs, _ = golifx.LookupBulbs()

	// Sort by labels for consistent ordering
	sort.SliceStable(bulbs, func(i, j int) bool {
		return label(bulbs[i]) < label(bulbs[j])
	})

	for b := range bulbs {
		log.Info().Str("bulb", label(bulbs[b])).Int("index", b).Msg("found bulb")
		//Get power state
		powerState, _ := bulbs[b].GetPowerState()

		//Turn if off
		if !powerState {
			bulbs[b].SetPowerState(true)
		}
	}
}

// Set a bulb to a specific RGB color
func Set(bulb int, red int, green int, blue int, white int) {
	r := float64(red) / 255
	g := float64(green) / 255
	b := float64(blue) / 255
	w := white

	rgb := colorful.Color{R: r, G: g, B: b}

	h, s, v := rgb.Hsv()

	hue := uint16(h / 360 * 65535)
	saturation := uint16(s * 65535)
	brightness := uint16(v * 65535)
	kelvin := (uint16(w * 6500 / 255)) + 2500

	log.Debug().
		Int("bulb", bulb).
		Int("red", red).
		Int("green", green).
		Int("blue", blue).
		Int("white", white).
		Uint16("hue", hue).
		Uint16("saturation", saturation).
		Uint16("brightness", brightness).
		Uint16("kelvin", kelvin).
		Msg("setting bulb")

	hsbk := &golifx.HSBK{
		Hue:        hue,
		Saturation: saturation,
		Brightness: brightness,
		Kelvin:     kelvin,
	}

	go func(bulb int, hsbk *golifx.HSBK) {
		bulbs[bulb].SetColorState(hsbk, 0)
	}(bulb, hsbk)
}
