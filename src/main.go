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

const nordnetId = 17385289 // instrument id

var outpwm = machine.PWM2

var idx = 0
var periods = []string{"DAY_1", "WEEK_1", "MONTH_1", "YEAR_1", "ALL"}
var returns = []float64{0, 0, 0, 0, 0}

var last = float64(0)

var servoValue = 1500 * time.Microsecond

func main() {

	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.High()

	setupButton()
	setupServo()
	setupDisplay()

	time.Sleep(3 * time.Second)

	setupDataFetch()

	for {
		led.Set(!led.Get()) // heartbeat
		time.Sleep(1 * time.Second)
	}

}

func setupButton() {
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
}

func setupServo() {
	servo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	servo.Low()

	err := outpwm.Configure(machine.PWMConfig{
		Period: uint64(20_000 * time.Microsecond),
	})
	if err != nil {
		println("failed to configure PWM")
		return
	}
	outch, err := outpwm.Channel(servo)
	if err != nil {
		println("failed to configure PWM channel")
		return
	}
	outpwm.Set(outch, 0)
	go func() {
		for {
			servoValue = scale(returns[idx])
			value := float64(outpwm.Top()) / 20_000 * float64(servoValue.Microseconds())
			outpwm.Set(outch, uint32(value))
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func setupDisplay() {
	display = newDisplay()
	display.Configure()
	go func() {
		for {
			show()
			time.Sleep(5 * time.Second)
		}
	}()
}

func setupDataFetch() {
	// Connect to Wifi
	err := setupWifi(wifiSsid, wifiPass)
	if err != nil {
		println(err.Error())
		time.Sleep(time.Second)
		arm.SystemReset()
	}
	nordnet := newNordnet()
	for {
		returns, _ = nordnet.getReturns(periods, nordnetId)
		last, _ = nordnet.getLast(nordnetId)
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
