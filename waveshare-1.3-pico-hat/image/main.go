package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/color"
	"machine"

	"github.com/sergeymakinen/go-bmp"
	"tinygo.org/x/drivers/st7789"
)

//go:embed mario.bmp
var marioData []byte

func main() {
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
	// display.FillScreen(color.RGBA{255, 255, 255, 1})
	display.FillScreen(color.RGBA{10, 200, 10, 1})
	display.EnableBacklight(false)

	fmt.Println(len(marioData))
	display.SetRotation(st7789.ROTATION_270)

	// fmt.Println("writing")
	// display.Tx(marioData, false)

	display.EnableBacklight(true)
	reader := bytes.NewReader(marioData)

	fmt.Println("decoding")
	// png decoding was way to slow. like never finishes after 30 minutes slow. this may be a limitation of the Pico
	// so I'm using bmp instead
	// marioImg, err := png.Decode(reader)
	marioImg, err := bmp.Decode(reader)
	if err != nil {
		fmt.Printf("failed to decode image: %s", err.Error())
		return
	}
	fmt.Println("done")

	maxX, maxY := int16(marioImg.Bounds().Max.X), int16(marioImg.Bounds().Max.Y)
	var x, y int16
	for y = 0; y < maxY; y++ {
		for x = 0; x < maxX; x++ {
			r, g, b, a := marioImg.At(int(x), int(y)).RGBA()
			if a == 0 {
				continue // keep transparent pixels transparent
			}
			clr := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}

			// // this enlarges the image out from 24x24 to 240x240.
			// // note that if a 240x240 image is loaded from a file directly it might be too much memory for the Pico to handle.
			// // this makes some assumptions about the image size. it may break if a different image is used.

			colors := make([]color.RGBA, 0, 100)
			for i := 0; i < 100; i++ {
				colors = append(colors, clr)
			}
			err = display.FillRectangleWithBuffer(x*10, y*10, 10, 10, colors)
			if err != nil {
				fmt.Printf("failed to FillRectangleWithBuffer: %s", err.Error())
			}

			// var i int16
			// for ; i < 10; i++ {
			// 	var j int16
			// 	for ; j < 10; j++ {
			// 		// set j to i for a fun look :)
			// 		display.SetPixel((x*10)+i, (y*10)+j, clr)
			// 	}
			// }
		}
	}
}
