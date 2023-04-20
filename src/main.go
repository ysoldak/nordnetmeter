package main

import (
	"device/arm"
	"machine"
	"time"
)

var Version string

var led = machine.LED
var servo = machine.D12
var button = machine.D11

const nordnetId = 17385289

var periods = []string{"DAY_1", "WEEK_1", "MONTH_1", "YEAR_1", "ALL"}
var periodIdx = 0

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
		periodIdx++
		if periodIdx > len(periods)-1 {
			periodIdx = 0
		}
		println(periodIdx)
	})

	time.Sleep(3 * time.Second)

	// Connect to Wifi
	err := setupWifi(wifiSsid, wifiPass)
	if err != nil {
		println(err.Error())
		time.Sleep(time.Second)
		arm.SystemReset()
	}

	nordnet := newNordnet()

	servoValue := 1500 * time.Microsecond

	go func() {
		for {
			servo.High()
			time.Sleep(servoValue)
			servo.Low()
			time.Sleep(20_000*time.Microsecond - servoValue)
		}
	}()

	// delta := 100 * time.Microsecond
	for {
		led.Set(!led.Get())
		ret, _ := nordnet.getReturn(periods[periodIdx], nordnetId)
		println(ret)

		// servoValue += delta
		// if servoValue > 2000*time.Microsecond {
		// 	servoValue = 1000 * time.Microsecond
		// }
		servoValue = scale(ret)
		println(servoValue)

		time.Sleep(5 * time.Second)
	}

}

func scale(percent float64) time.Duration {
	percent *= -1 // inversion
	if percent < -5 {
		return 1000 * time.Microsecond
	}
	if percent > 5 {
		return 2000 * time.Microsecond
	}
	x := (percent + 5) * 100
	return time.Duration(1000+x) * time.Microsecond
}
