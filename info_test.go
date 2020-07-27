package switchbot

import "testing"

func TestNewInfoWithRawInfo(t *testing.T) {
	r := []byte{1, 79, 45, 100, 0, 0, 0, 152, 3, 0, 3, 72, 0}

	got := NewBotInfoWithRawInfo(r)
	if got.Battery != 79 {
		t.Errorf("Battery is not correct, got %v", got.Battery)
	}
	if got.Firmware != float64(45)/10 {
		t.Errorf("Firmware is not correct, got %v", got.Firmware)
	}
	if got.Firmware != float64(45)/10 {
		t.Errorf("Firmware is not correct, got %v", got.Firmware)
	}
	if got.TimerCount != 3 {
		t.Errorf("TimerCount is not correct, got %v", got.TimerCount)
	}
	if got.StateMode != false {
		t.Errorf("StateMode is not correct, got %v", got.TimerCount)
	}
	if got.Inverse != false {
		t.Errorf("Inverse is not correct, got %v", got.TimerCount)
	}
	if got.HoldSec != 3 {
		t.Errorf("HoldSec is not correct, got %v", got.HoldSec)
	}
}
