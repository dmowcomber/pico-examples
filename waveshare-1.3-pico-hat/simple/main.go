package main

import (
	"image/color"
	"machine"
	"time"

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

	time.Sleep(time.Second * 1)

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

	handleButtons(&display)
}

/*
	buttons:
*/

var (
	// button GPIO pins are defined in the waveshare wiki: https://www.waveshare.com/wiki/Pico-LCD-1.3
	buttonUp     = newButton(machine.GP2)
	buttonDown   = newButton(machine.GP18)
	buttonLeft   = newButton(machine.GP16)
	buttonRight  = newButton(machine.GP20)
	buttonCenter = newButton(machine.GP3)
	buttonA      = newButton(machine.GP15)
	buttonB      = newButton(machine.GP17)
	buttonX      = newButton(machine.GP19)
	buttonY      = newButton(machine.GP21)
)

func handleButtons(display *st7789.Device) {
	// led buttons
	buttonA.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		machine.LED.Set(p.Get())
	})
	buttonB.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		machine.LED.Set(!p.Get())
	})

	// screen buttons
	buttonX.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		display.FillScreen(color.RGBA{255, 255, 255, 255})
		display.EnableBacklight(true)
	})
	buttonY.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		display.EnableBacklight(false)
	})
	buttonUp.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		display.FillScreen(color.RGBA{0, 255, 0, 255})
	})
	buttonLeft.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		display.FillScreen(color.RGBA{255, 0, 0, 255})
	})
	buttonRight.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		display.FillScreen(color.RGBA{0, 0, 255, 255})
	})
	buttonDown.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		display.FillScreen(color.RGBA{0, 0, 0, 255})
	})
}

func newButton(pin machine.Pin) machine.Pin {
	button := pin
	// all of the waveshare 1.3 inch pico hat buttons are pull up buttons
	button.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	return button
}
