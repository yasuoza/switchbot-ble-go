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

func TestParseTimer(t *testing.T) {
	type Want struct {
		enabled  bool
		weekdays [7]bool
		hour     int
		minutes  int
		action   int
	}
	tests := []struct {
		got  []byte
		want Want
	}{
		{
			got: []byte{1, 2, 0, 121, 10, 11, 0, 0, 0, 0, 0, 0},
			want: Want{
				enabled:  true,
				weekdays: [7]bool{true, true, false, false, true, true, true},
				hour:     10,
				minutes:  11,
				action:   0,
			},
		},
		{
			got: []byte{1, 2, 0, 0, 12, 0, 64, 0, 0, 0, 0, 0},
			want: Want{
				enabled:  false,
				weekdays: [7]bool{true, false, false, false, false, false, false},
				hour:     12,
				minutes:  0,
				action:   0,
			},
		},
		{
			got: []byte{1, 3, 0, 31, 10, 11, 0, 1, 0, 0, 0, 0},
			want: Want{
				enabled:  true,
				weekdays: [7]bool{false, true, true, true, true, true, false},
				hour:     10,
				minutes:  11,
				action:   1,
			},
		},
		{
			got: []byte{1, 3, 0, 4, 15, 33, 0, 2, 0, 0, 0, 0},
			want: Want{
				enabled:  true,
				weekdays: [7]bool{false, false, false, true, false, false, false},
				hour:     15,
				minutes:  33,
				action:   2,
			},
		},
		{
			got: []byte{1, 2, 0, 128, 18, 45, 0, 2, 0, 0, 0, 0},
			want: Want{
				enabled:  true,
				weekdays: [7]bool{false, false, false, false, false, false, false},
				hour:     18,
				minutes:  45,
				action:   2,
			},
		},
		{
			got: []byte{1, 2, 0, 0, 18, 45, 128, 2, 0, 0, 0, 0},
			want: Want{
				enabled:  false,
				weekdays: [7]bool{false, false, false, false, false, false, false},
				hour:     18,
				minutes:  45,
				action:   2,
			},
		},
	}

	for _, tt := range tests {
		got, want := tt.got, tt.want
		timer := ParseTimerBytes(got)
		testTimer(t, timer, want.enabled, want.weekdays, want.hour, want.minutes, want.action)
	}
}
