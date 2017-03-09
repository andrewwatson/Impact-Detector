package main

import (
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

func lightOn(pin int) {
	theLED := rpio.Pin(pin)
	theLED.Output()
	theLED.High()
}

func lightOff(pin int) {
	theLED := rpio.Pin(pin)
	theLED.Output()
	theLED.Low()
}

func blink(pin int) {

	theLED := rpio.Pin(pin)
	theLED.Output()

	theLED.High()
	time.Sleep(500 * time.Millisecond)
	theLED.Low()
}

func button(pin int) chan bool {

	event := make(chan bool)

	theButton := rpio.Pin(pin)
	theButton.Input()

	go func() {

		depressed := false
		for {

			switch theButton.Read() {
			case rpio.Low:
				if !depressed {
					event <- true
					depressed = true
				}

			case rpio.High:
				depressed = false
			}

			time.Sleep(1 * time.Millisecond)
		}
	}()

	return event
}
