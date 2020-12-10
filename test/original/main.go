package main

import (
	"sort"
	"time"

	"github.com/2tvenom/golifx"
)

func label(bulb *golifx.Bulb) string {
	label, _ := bulb.GetLabel()
	return label
}

func main() {
	//Lookup all bulbs
	bulbs, _ := golifx.LookupBulbs()

	// Sort by labels for consistent ordering
	sort.SliceStable(bulbs, func(i, j int) bool {
		return label(bulbs[i]) < label(bulbs[j])
	})

	//Get power state
	powerState, _ := bulbs[0].GetPowerState()

	//Turn if off
	if !powerState {
		bulbs[0].SetPowerState(true)
	}

	ticker := time.NewTicker(time.Second)
	counter := 0

	hsbk := &golifx.HSBK{
		Hue:        2000,
		Saturation: 13106,
		Brightness: 65535,
		Kelvin:     3200,
	}
	//Change color every second
	for range ticker.C {
		bulbs[0].SetColorState(hsbk, 500)
		counter++
		hsbk.Hue += 5000
		if counter > 10 {
			ticker.Stop()
			break
		}
	}

	//Turn off
	bulbs[0].SetPowerState(false)
}
