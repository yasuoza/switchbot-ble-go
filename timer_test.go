package switchbot

import (
	"testing"
)

func testTimer(t *testing.T, timer *Timer, enabled bool, weekdays [7]bool, hour int, minutes, action int) {
	if timer.Enabled != enabled {
		t.Errorf("timer.Enabled expected %v, got %v\n", timer.Enabled, enabled)
	}
	if timer.Weekdays != weekdays {
		t.Errorf("timer.Weekday expected %v, got %v\n", timer.Weekdays, weekdays)
	}
	if timer.Hour != hour {
		t.Errorf("timer.Hour expected %v, got %v\n", timer.Hour, hour)
	}
	if timer.Minutes != minutes {
		t.Errorf("timer.Minutes expected %v, got %v\n", timer.Minutes, minutes)
	}
	if timer.Action != action {
		t.Errorf("timer.Action expected %v, got %v\n", timer.Action, action)
	}
}

func TestParseTimerReturnsNil(t *testing.T) {
	d := []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	timer := ParseTimerBytes(d)
	if timer != nil {
		t.Fatal("expected timer == nil, but got", timer)
	}
}

func TestParseTimerPressRepeatEnabled(t *testing.T) {
	// - Enabled
	// - Repeat Sunday, Monday, Thursday, Friday, Saturday
	// - Hour 10
	// - Minutes 11
	d := []byte{1, 2, 0, 121, 10, 11, 0, 0, 0, 0, 0, 0}

	timer := ParseTimerBytes(d)
	testTimer(
		t,
		timer,
		true,
		[7]bool{true, true, false, false, true, true, true},
		10,
		11,
		0,
	)
}

func TestParseTimerBytesPressRepeatDisabled(t *testing.T) {
	// - Disabled
	// - Repeat Sunday
	// - Hour 12
	// - Minutes 00
	d := []byte{1, 2, 0, 0, 12, 0, 64, 0, 0, 0, 0, 0}

	timer := ParseTimerBytes(d)
	testTimer(
		t,
		timer,
		false,
		[7]bool{true, false, false, false, false, false, false},
		12,
		0,
		0,
	)
}

func TestParseTimerBytesOnRepeatEnabled(t *testing.T) {
	// - Enabled
	// - Repeat Monday, Tuesday, Wednesday, Thursday, Friday
	// - Hour 10
	// - Minutes 11
	// - On
	d := []byte{1, 3, 0, 31, 10, 11, 0, 1, 0, 0, 0, 0}

	timer := ParseTimerBytes(d)
	testTimer(
		t,
		timer,
		true,
		[7]bool{false, true, true, true, true, true, false},
		10,
		11,
		1,
	)
}

func TestParseTimerBytesOffRepeastEnabled(t *testing.T) {
	// - Enabled
	// - Repeat Wednesday
	// - Hour 15
	// - Minutes 33
	// - Off
	d := []byte{1, 3, 0, 4, 15, 33, 0, 2, 0, 0, 0, 0}

	timer := ParseTimerBytes(d)
	testTimer(
		t,
		timer,
		true,
		[7]bool{false, false, false, true, false, false, false},
		15,
		33,
		2,
	)
}

func TestParseTimerBytesOffOnceEnabled(t *testing.T) {
	// - Enabled
	// - Once
	// - Hour 18
	// - Minutes 45
	// - Off
	d := []byte{1, 2, 0, 128, 18, 45, 0, 2, 0, 0, 0, 0}

	timer := ParseTimerBytes(d)
	testTimer(
		t,
		timer,
		true,
		[7]bool{false, false, false, false, false, false, false},
		18,
		45,
		2,
	)
}

func TestParseTimerBytesOffOnceDisabled(t *testing.T) {
	// - Disabled
	// - Once
	// - Hour 18
	// - Minutes 45
	// - Off
	d := []byte{1, 2, 0, 0, 18, 45, 128, 2, 0, 0, 0, 0}

	timer := ParseTimerBytes(d)
	testTimer(
		t,
		timer,
		false,
		[7]bool{false, false, false, false, false, false, false},
		18,
		45,
		2,
	)
}
