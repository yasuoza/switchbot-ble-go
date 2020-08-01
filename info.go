package switchbot

import (
	"bytes"
	"fmt"
)

// BotInfo represents current SwitchBot's information.
type BotInfo struct {
	Battery    int
	Firmware   float64
	TimerCount int
	StateMode  bool
	Inverse    bool
	HoldSec    int
}

// NewBotInfoWithRawInfo initialize BotInfo with raw byte data.
// This works with switchbot.GetInfo.
func NewBotInfoWithRawInfo(info []byte) *BotInfo {
	batt := int(info[1])
	firm := float64(info[2]) / 10
	tc := int(info[8])
	st := (info[9] & 16) != 0
	iv := (info[9] & 1) != 0
	hs := int(info[10])

	return &BotInfo{
		Battery:    batt,
		Firmware:   firm,
		TimerCount: tc,
		StateMode:  st,
		Inverse:    iv,
		HoldSec:    hs,
	}
}

// String returns formatted information
func (i *BotInfo) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("Battery: %d", i.Battery))
	buf.WriteString(fmt.Sprintf(", Firmware: %0.1f", i.Firmware))
	buf.WriteString(fmt.Sprintf(", TimerCount: %d", i.TimerCount))
	buf.WriteString(fmt.Sprintf(", StateMode: %t", i.StateMode))
	buf.WriteString(fmt.Sprintf(", Inverse: %t", i.Inverse))
	buf.WriteString(fmt.Sprintf(", HoldSec: %d", i.HoldSec))
	return buf.String()
}
