package rotary

import (
	"machine"
	"time"
)

type RotaryButton struct {
	btn machine.Pin
}

func NewRotaryButton(pin machine.Pin) *RotaryButton {
	pin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	return &RotaryButton{
		btn: pin,
	}
}

func (r *RotaryButton) IsPressed() bool {
	// デバウンス処理: ボタンが確実に押されたことを確認する
	if !r.btn.Get() {
		time.Sleep(20 * time.Millisecond) // 20ms待機してノイズを除去
		if !r.btn.Get() {                 // もう一度確認
			return true
		}
	}
	return false
}

func (r *RotaryButton) WaitForRelease() {
	// ボタンが放されるのを確実に待機
	for !r.btn.Get() { // ボタンが押されたままなら待機
		time.Sleep(10 * time.Millisecond)
	}
}
