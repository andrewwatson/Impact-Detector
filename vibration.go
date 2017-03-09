package main

import (
	"strconv"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

const (
	stateStill           = 0
	stateMoving          = 1
	stateChangeThreshold = 2000
	cycleDelay           = 20 * time.Microsecond
)

type seismic struct {
	Duration int    // in ns
	Pattern  string // pattern of highs and lows recorded
}

func vibration(pin int) chan seismic {

	var lastReading, currentReading rpio.State
	var currentState int
	event := make(chan seismic)

	go func() {

		sensor := rpio.Pin(pin)
		sensor.Input()

		lastReading = sensor.Read()

		var lastStateChange time.Time
		unchangedState := 0

		pulses := ""

		for {

			rightNow := time.Now()
			currentReading = sensor.Read()

			if currentState == stateMoving {
				pulses += strconv.Itoa(int(currentReading))
			}

			if lastReading != currentReading {

				// the pin state has changed!
				unchangedState = 0
				if currentState == stateStill {
					currentState = stateMoving
					lastStateChange = rightNow

				}

			} else {

				unchangedState++

				if unchangedState > stateChangeThreshold {
					// we've been in the same same condition long enough to assume we're not moving

					if currentState == stateMoving {
						// we were moving but we've stopped now
						currentState = stateStill

						duration := rightNow.Sub(lastStateChange).Nanoseconds() / 1000000

						event <- seismic{
							Duration: int(duration),
							Pattern:  pulses,
						}

						lastStateChange = rightNow
						pulses = ""

						// "de-bounce" the switch because sometimes you get little blips
						time.Sleep(100 * time.Millisecond)

					}

					unchangedState = 0
				}
			}

			lastReading = currentReading
			time.Sleep(cycleDelay)
		}
	}()

	return event
}
