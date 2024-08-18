package keyboard

import (
	"machine"
	"time"

	pio "github.com/tinygo-org/pio/rp2-pio"
	"github.com/tinygo-org/pio/rp2-pio/piolib"
)

type WS2812B struct {
	Pin machine.Pin
	ws  *piolib.WS2812B
}

func NewWS2812B(pin machine.Pin) *WS2812B {
	s, _ := pio.PIO0.ClaimStateMachine()
	ws, _ := piolib.NewWS2812B(s, pin)
	ws.EnableDMA(true)
	return &WS2812B{
		ws: ws,
	}
}

func (ws *WS2812B) WriteRaw(rawGRB []uint32) error {
	return ws.ws.WriteRaw(rawGRB)
}

type KeyboardController struct {
	rowPins []machine.Pin
	colPins []machine.Pin
	ws      *WS2812B
}

func NewKeyboardController(rowPins []machine.Pin, colPins []machine.Pin, ws *WS2812B) *KeyboardController {
	initializePins(rowPins)
	initializePins(colPins)

	return &KeyboardController{
		rowPins: rowPins,
		colPins: colPins,
		ws:      ws,
	}
}

func initializePins(pins []machine.Pin) {
	time.Sleep(10 * time.Millisecond)
	for _, pin := range pins {
		pin.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}
	time.Sleep(10 * time.Millisecond)
}

func (kb *KeyboardController) RowPins() []machine.Pin {
	return kb.rowPins
}

func (kb *KeyboardController) ColPins() []machine.Pin {
	return kb.colPins
}

func (kb *KeyboardController) LightKeys(ledArray []bool, color uint32) {
	rawGRB := make([]uint32, len(ledArray))

	for i, isLit := range ledArray {
		if isLit {
			rawGRB[i] = color
		}
	}

	kb.ws.WriteRaw(rawGRB)
}

func (kb *KeyboardController) TurnOffKey(keyIndex int) {
	rawGRB := make([]uint32, len(kb.rowPins)*len(kb.colPins))
	kb.ws.WriteRaw(rawGRB)
}

func (kb *KeyboardController) DisableAllKeys() {
	rawGRB := make([]uint32, len(kb.rowPins)*len(kb.colPins))
	kb.ws.WriteRaw(rawGRB)
}

// ScanKeys scans the matrix keyboard for any pressed key and returns the key index and whether a key was pressed.
func (kb *KeyboardController) ScanKeys() (int, bool) {
	for colIndex, colPin := range kb.colPins {
		// Set the current column to LOW (active)
		colPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
		colPin.Low()

		// Check each row to see if it is LOW (pressed)
		for rowIndex, rowPin := range kb.rowPins {
			if rowPin.Get() {
				// Calculate the key index based on row and column
				keyIndex := rowIndex*len(kb.colPins) + colIndex
				// Reset column pin
				colPin.High()
				colPin.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
				return keyIndex, true
			}
		}

		// Reset column pin
		colPin.High()
		colPin.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}

	return -1, false
}
