package switchbot

// Timer represents Timer configuration.
type Timer struct {
	Enabled  bool
	Weekdays [7]bool
	Hour     int
	Minutes  int
	Action   int
}

// ParseTimerBytes parses bytes to timer object.
func ParseTimerBytes(val []byte) *Timer {
	enabled := val[3] != 0
	h := int(val[4])
	m := int(val[5])
	ac := int(val[7] & 15)

	if !enabled && h == 0 && m == 0 {
		return nil
	}

	var weekdays [7]bool
	if !enabled {
		rep := (val[6] & 240) | ((val[7] & 240) >> 4)
		weekdays = parseWeekDays(rep)
	} else {
		weekdays = parseWeekDays(val[3])
	}

	return &Timer{
		Enabled:  enabled,
		Weekdays: weekdays,
		Hour:     h,
		Minutes:  m,
		Action:   ac,
	}
}

func parseWeekDays(b byte) [7]bool {
	return [7]bool{
		(b & 64) != 0, // Sun
		(b & 1) != 0,  // Mon
		(b & 2) != 0,  // Tue
		(b & 4) != 0,  // Wed
		(b & 8) != 0,  // Thu
		(b & 16) != 0, // Fri
		(b & 32) != 0, // Sat
	}
}
