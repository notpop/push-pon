package display

import (
	"fmt"
	"image/color"
	"time"

	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/proggy"
)

type DisplayController struct {
	device *ssd1306.Device
	width  int16
	height int16
}

func NewDisplayController(device *ssd1306.Device) *DisplayController {
	width, height := device.Size()
	return &DisplayController{
		device: device,
		width:  width,
		height: height,
	}
}

func (d *DisplayController) Initialize() {
	d.device.ClearDisplay()
	d.SetHorizontalFlip(false)
	d.SetVerticalFlip(false)
	d.device.Display()
}

func (d *DisplayController) SetHorizontalFlip(enable bool) {
	if enable {
		d.device.Command(ssd1306.SEGREMAP | 0x1)
	} else {
		d.device.Command(ssd1306.SEGREMAP)
	}
}

func (d *DisplayController) SetVerticalFlip(enable bool) {
	if enable {
		d.device.Command(ssd1306.COMSCANDEC)
	} else {
		d.device.Command(ssd1306.COMSCANINC)
	}
}

func (d *DisplayController) ScrollMessage(text string, delay time.Duration, checkRotary func() bool) {
	width := int(d.width)
	messageWidth, _ := tinyfont.LineWidth(&proggy.TinySZ8pt7b, text)
	startX := width

	for startX > -int(messageWidth) {
		if checkRotary() {
			return
		}
		d.device.ClearBuffer()
		tinyfont.WriteLine(d.device, &proggy.TinySZ8pt7b, int16(startX), d.height/2, text, color.RGBA{255, 255, 255, 255})
		d.device.Display()
		startX -= 2
		time.Sleep(delay)
	}
	d.device.ClearBuffer()
	d.device.Display()
}

func (d *DisplayController) ShowMessage(text string) {
	d.device.ClearBuffer()
	tinyfont.WriteLine(d.device, &proggy.TinySZ8pt7b, 0, 16, text, color.RGBA{255, 255, 255, 255})
	d.device.Display()
}

func (d *DisplayController) ShowResult(totalScore int, success bool) {
	d.device.ClearBuffer()
	var message string
	if success {
		message = fmt.Sprintf("Congratulations!\nTotal Score: %d", totalScore)
	} else {
		message = fmt.Sprintf("Game Over\nTotal Score: %d", totalScore)
	}
	tinyfont.WriteLine(d.device, &proggy.TinySZ8pt7b, 0, 16, message, color.RGBA{255, 255, 255, 255})
	d.device.Display()
}
