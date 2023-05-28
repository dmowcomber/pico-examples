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

	// // Default Serial Clock Bus 1 for SPI communications
	// SPI1_SCK_PIN	= GPIO10
	// // Default Serial Out Bus 1 for SPI communications
	// SPI1_SDO_PIN	= GPIO11	// Tx
	// // Default Serial In Bus 1 for SPI communications
	// SPI1_SDI_PIN	= GPIO12	// Rx
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

	// without the RowOffset there's static on the lft side of the screen
	display.Configure(st7789.Config{RowOffset: 80})

	handleButtons(&display)

	buttonX.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		display.FillScreen(color.RGBA{255, 255, 255, 255})
		display.EnableBacklight(true)
		machine.LED.Set(p.Get())
	})
	buttonY.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		// display.FillScreen(color.RGBA{0, 0, 255, 255})
		display.EnableBacklight(false)
		machine.LED.Set(!p.Get())
	})

	// display.FillScreen(color.RGBA{255, 255, 255, 255})
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

	display.EnableBacklight(true)
	var i int16 = 0
	var j int16 = 0
	for ; i < 240; i++ {
		for ; j < 240; j++ {
			display.SetPixel(i, j, color.RGBA{0, 255, 0, 255})
		}
	}

	// for {
	// 	// 	// led.Low()
	// 	// 	display.EnableBacklight(false)
	// 	// 	time.Sleep(time.Millisecond * 500)

	// 	// 	// led.High()
	// 	display.FillScreen(color.RGBA{0, 255, 0, 255})
	// 	// 	display.EnableBacklight(true)
	// 	time.Sleep(time.Millisecond * 500)

	// 	display.FillScreen(color.RGBA{0, 0, 255, 255})
	// 	time.Sleep(time.Millisecond * 500)
	// }
}

// buttons:

var (
	// MicroPython Example: https://github.com/russhughes/st7789_mpy/blob/518562c7f44a40f5ab8f18b255eb6564284be64c/examples/configs/ws_pico_13/tft_buttons.py#L2
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
	buttonA.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		machine.LED.Set(p.Get())
	})
	buttonB.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		machine.LED.Set(!p.Get())
	})
}

// all of the waveshare 1.3 inch pico hat buttons are pull up buttons
func newButton(pin machine.Pin) machine.Pin {
	button := pin
	button.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	return button
}
