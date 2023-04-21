package main

import (
	"image/color"
	"machine"

	"tinygo.org/x/drivers/ssd1306"
)

var BLACK = color.RGBA{0, 0, 0, 255}
var WHITE = color.RGBA{255, 255, 255, 255}

type Display struct {
	device ssd1306.Device
}

func newDisplay() *Display {
	return &Display{}
}

func (d *Display) Configure() {
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: machine.TWI_FREQ_400KHZ,
		SDA:       machine.I2C0_SDA_PIN,
		SCL:       machine.I2C0_SCL_PIN,
	})
	d.device = ssd1306.NewI2C(machine.I2C0)
	d.device.Configure(ssd1306.Config{
		Address: ssd1306.Address_128_32,
		Width:   128,
		Height:  32,
	})
	d.device.ClearDisplay()
}
