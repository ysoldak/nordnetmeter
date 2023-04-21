package main

import (
	"device/arm"
	"fmt"
	"machine"
	"time"

	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/proggy"
)

var Version string

var led = machine.LED
var servo = machine.D12
var button = machine.D11

var display *Display

const nordnetId = 17385289

var idx = 0
var periods = []string{"DAY_1", "WEEK_1", "MONTH_1", "YEAR_1", "ALL"}
var returns = []float64{0, 0, 0, 0, 0}

var last = float64(0)

var servoValue = 1500 * time.Microsecond

func main() {

	// Indicate wake up
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.High()

	servo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	servo.Low()

	button.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	time.Sleep(time.Second)
	prev := int64(0)
	button.SetInterrupt(machine.PinFalling, func(pin machine.Pin) {
		now := time.Now().UnixMilli()
		if prev != 0 && now-prev < 1000 {
			return
		}
		prev = now
		idx++
		if idx > len(periods)-1 {
			idx = 0
		}
		println(idx)
		servoValue = scale(returns[idx])
		show()
	})

	display = newDisplay()
	display.Configure()
	show()
	go func() {
		for {
			show()
			time.Sleep(5 * time.Second)
		}
	}()

	time.Sleep(3 * time.Second)

	// Connect to Wifi
	err := setupWifi(wifiSsid, wifiPass)
	if err != nil {
		println(err.Error())
		time.Sleep(time.Second)
		arm.SystemReset()
	}

	nordnet := newNordnet()

	go func() {
		for {
			servo.High()
			time.Sleep(servoValue)
			servo.Low()
			time.Sleep(20_000*time.Microsecond - servoValue)
		}
	}()

	for {
		led.Set(!led.Get())
		returns, _ = nordnet.getReturns(periods, nordnetId)
		last, _ = nordnet.getLast(nordnetId)
		servoValue = scale(returns[idx])
		time.Sleep(5 * time.Second)
	}

}

func scale(percent float64) time.Duration {
	// percent = -7 // for calibration
	percent *= -1 // inversion
	if percent < -5 {
		return 700 * time.Microsecond
	}
	if percent > 5 {
		return 2250 * time.Microsecond
	}
	calibration := float64(-25)
	x := (percent+5)*140 + calibration
	return time.Duration(800+x) * time.Microsecond
}

func show() {
	display.device.ClearDisplay()
	tinyfont.WriteLineRotated(&display.device, &proggy.TinySZ8pt7b, 64, 12, fmt.Sprintf("%.02f", returns[idx]), WHITE, tinyfont.NO_ROTATION)
	tinyfont.WriteLineRotated(&display.device, &proggy.TinySZ8pt7b, 64, 28, periods[idx], WHITE, tinyfont.NO_ROTATION)
	tinyfont.WriteLineRotated(&display.device, &proggy.TinySZ8pt7b, 14, 28, "SAVE", WHITE, tinyfont.NO_ROTATION)
	tinyfont.WriteLineRotated(&display.device, &proggy.TinySZ8pt7b, 14, 12, fmt.Sprintf("%.02f", last), WHITE, tinyfont.NO_ROTATION)
	display.device.Display()
}
