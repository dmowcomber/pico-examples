package main

import (
	"image/color"
	"machine"

	"tinygo.org/x/drivers/st7789"
)

func main() {
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	machine.SPI1.Configure(machine.SPIConfig{
		Frequency: 62500000,
		SCK:       machine.SPI1_SCK_PIN,
		SDO:       machine.SPI1_SDO_PIN,
		SDI:       machine.SPI1_SDI_PIN,
		Mode:      2,
	})

	display := st7789.New(machine.SPI1,
		machine.GP12,
		machine.GP8,
		machine.GP9,
		machine.GP13,
	)

	display.Configure(st7789.Config{
		Width:     240,
		Height:    240,
		RowOffset: 80, // without the RowOffset there's static on the lft side of the screen
	})
	display.EnableBacklight(true)
	var i, j int16
	for ; i < 255; i++ {
		for j = 0; j < 255; j++ {
			display.SetPixel(i, j, color.RGBA{uint8(i), uint8(j), 0, 1})
		}
	}
}
