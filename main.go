package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	proto "github.com/huin/mqtt"
	"github.com/jeffallen/mqtt"
	"github.com/stianeikeland/go-rpio"
)

const (
	button1Pin = 4
	button2Pin = 5
	button3Pin = 6

	vibrationPin = 18

	learningIndicatorPin = 21
	buttonIndicatorPin   = 26
)

var (
	mqttClient *mqtt.ClientConn
)

func init() {

	var host = flag.String("host", "mqtt.makeandbuild.info:1883", "hostname of broker")
	var id = flag.String("id", "", "client id")
	var user = flag.String("user", "raspberry", "username")
	var pass = flag.String("pass", "pi", "password")

	conn, err := net.Dial("tcp", *host)
	if err != nil {
		fmt.Fprint(os.Stderr, "dial: ", err)
		return
	}
	mqttClient = mqtt.NewClientConn(conn)

	mqttClient.ClientId = *id

	tq := []proto.TopicQos{
		{Topic: "tick", Qos: proto.QosAtMostOnce},
	}

	if err := mqttClient.Connect(*user, *pass); err != nil {
		fmt.Fprintf(os.Stderr, "connect: %v\n", err)
		os.Exit(1)
	}

	mqttClient.Subscribe(tq)

}

func main() {

	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer rpio.Close()

	fmt.Println("BOOT SEQUENCE COMPLETED")
	fmt.Println("WAITING FOR VIBRATIONS")

	broadcast("impact/initialize", "hi!")

	button1 := button(button1Pin)
	button2 := button(button2Pin)
	button3 := button(button3Pin)

	vibrations := vibration(vibrationPin)

	var lastVibration seismic
	learningMode := false

	for {

		select {
		case <-button1:

			if learningMode {
				fmt.Printf("Vibration Pattern: (%d ms) (%d data points) is a MISS\n", lastVibration.Duration, len(lastVibration.Pattern))
				broadcast("learn/miss", lastVibration.Pattern)

				lightOff(learningIndicatorPin)
				blink(buttonIndicatorPin)
				learningMode = false
			}

		case <-button2:
			if learningMode {
				fmt.Printf("Vibration Pattern: (%d ms) (%d data points) is a HIT\n", lastVibration.Duration, len(lastVibration.Pattern))
				broadcast("learn/hit", lastVibration.Pattern)

				lightOff(learningIndicatorPin)
				blink(buttonIndicatorPin)
				learningMode = false
			}

		case <-button3:
			if learningMode {
				fmt.Printf("Vibration Pattern: (%d ms) (%d data points) is a POSSIBLE FLASH PINT\n", lastVibration.Duration, len(lastVibration.Pattern))
				broadcast("learn/flash", lastVibration.Pattern)

				lightOff(learningIndicatorPin)
				blink(buttonIndicatorPin)
				learningMode = false
			}

		case vib := <-vibrations:
			fmt.Printf("Vibration Detected: (%d ms) (%d data points)\n", vib.Duration, len(lastVibration.Pattern))
			lightOn(learningIndicatorPin)
			lastVibration = vib
			learningMode = true

		default:
		}
	}

}
