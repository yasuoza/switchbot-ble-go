package switchbot

import (
	"encoding/binary"
	"hash/crc32"

	"github.com/JuulLabs-OSS/ble"
)

const (
	characteristics = "cba20002-224d-11e6-9fb8-0002a5d5c51b"
	handle          = 0x16
)

// Bot represents SwitchBot device.
type Bot struct {
	Addr string

	cl ble.Client
	pw []byte
}

// Press triggers press function for the SwitchBot.
// SwitchBot must be set to press mode.
func (b *Bot) Press() error {
	var cmd []byte
	if b.encrypted() {
		cmd = append([]byte{0x57, 0x11}, b.pw...)
	} else {
		cmd = []byte{0x57, 0x01}
	}
	return b.trigger(cmd)
}

// On triggers on function for the SwitchBot.
// SwitchBot must be set to On/Off mode.
func (b *Bot) On() error {
	var cmd []byte
	if b.encrypted() {
		cmd = append(append([]byte{0x57, 0x11}, b.pw...), []byte{0x01}...)
	} else {
		cmd = []byte{0x57, 0x01, 0x01}
	}
	return b.trigger(cmd)
}

// Off triggers off function for the SwitchBot.
// SwitchBot must be set to On/Off mode.
func (b *Bot) Off() error {
	var cmd []byte
	if b.encrypted() {
		cmd = append(append([]byte{0x57, 0x11}, b.pw...), []byte{0x02}...)
	} else {
		cmd = []byte{0x57, 0x01, 0x02}
	}
	return b.trigger(cmd)
}

// SetPSetPSetPassword sets SwitchBot's password.
// If SwitchBot is configured to use password authentication,
// you need to call SetPassword before calling Press/On/Off function.
func (b *Bot) SetPassword(pw string) {
	crc := crc32.ChecksumIEEE([]byte(pw))
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs[0:], crc)
	b.pw = bs
}

func (b *Bot) encrypted() bool {
	return len(b.pw) != 0
}

func (b *Bot) trigger(cmd []byte) error {
	c := ble.NewCharacteristic(ble.MustParse(characteristics))
	c.ValueHandle = handle
	if err := b.cl.WriteCharacteristic(c, cmd, false); err != nil {
		return err
	}
	return nil
}
