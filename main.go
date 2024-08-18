package main

import (
	"machine"
	"push-pon/display"
	"push-pon/game"
	"push-pon/keyboard"
	"push-pon/rotary"
	"time"

	"tinygo.org/x/drivers/ssd1306"
)

func initializeDisplay(i2c *machine.I2C) *display.DisplayController {
	i2c.Configure(machine.I2CConfig{
		Frequency: 400 * machine.KHz,
		SDA:       machine.GPIO12,
		SCL:       machine.GPIO13,
	})

	displayDevice := ssd1306.NewI2C(i2c)
	displayDevice.Configure(ssd1306.Config{
		Address: 0x3C,
		Width:   128,
		Height:  64,
	})

	displayCtrl := display.NewDisplayController(&displayDevice)
	displayCtrl.Initialize()
	return displayCtrl
}

func initializeKeyboard() *keyboard.KeyboardController {
	ws := keyboard.NewWS2812B(machine.GPIO1)

	// https://kicanvas.org/?github=https%3A%2F%2Fgithub.com%2Fsago35%2Fkeyboards%2Ftree%2Fmain%2Fzero-kb02%2Fzero-kb02
	rowPins := []machine.Pin{
		machine.GPIO9,  // ROW1
		machine.GPIO10, // ROW2
		machine.GPIO11, // ROW3
	}

	colPins := []machine.Pin{
		machine.GPIO5, // COL1
		machine.GPIO6, // COL2
		machine.GPIO7, // COL3
		machine.GPIO8, // COL4
	}

	return keyboard.NewKeyboardController(rowPins, colPins, ws)
}

func initializeRotaryButton() *rotary.RotaryButton {
	return rotary.NewRotaryButton(machine.GPIO2)
}

func mainLoop(gameCtrl *game.GameController, rotaryBtn *rotary.RotaryButton) {
	for {
		if !gameCtrl.IsGameStarted() {
			if rotaryBtn.IsPressed() {
				rotaryBtn.WaitForRelease()
				gameCtrl.StartGame()
			}
		} else {
			gameCtrl.CheckKeyPress()
		}

		if rotaryBtn.IsPressed() {
			rotaryBtn.WaitForRelease()
			gameCtrl.ResetGame()
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	i2c := machine.I2C0
	displayCtrl := initializeDisplay(i2c)
	keyboardCtrl := initializeKeyboard()
	rotaryBtn := initializeRotaryButton()
	gameCtrl := game.NewGameController(keyboardCtrl, displayCtrl)

	// 安定しない（途中で止まる）のでコメントアウト
	// displayCtrl.ScrollMessage("Welcome to the game!!\nPlease push rotary button!!", 100*time.Millisecond, func() bool {
	// 	return rotaryBtn.IsPressed()
	// })

	displayCtrl.ShowMessage("Welcome to the game!!\nPush rotary!!")

	mainLoop(gameCtrl, rotaryBtn)
}
