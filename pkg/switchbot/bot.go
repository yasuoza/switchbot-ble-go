package switchbot

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"strings"

	"tinygo.org/x/bluetooth"
)

const (
	characteristics = "cba20002-224d-11e6-9fb8-0002a5d5c51b"
	handle          = 0x16
)

// Bot represents SwitchBot device.
type Bot struct {
	Addr string

	dev *bluetooth.Device

	subschar *bluetooth.DeviceCharacteristic
	cmdchar  *bluetooth.DeviceCharacteristic

	pw []byte

	subsque    chan []byte
	subscribed bool
}

// NewBot initializes bot object.
func NewBot(addr string) *Bot {
	b := &Bot{Addr: strings.ToLower(addr)}
	b.subsque = make(chan []byte)
	b.subscribed = false
	return b
}

// SetPassword sets SwitchBot's password.
// If SwitchBot is configured to use password authentication,
// you need to call SetPassword before calling Press/On/Off function.
func (b *Bot) SetPassword(pw string) {
	crc := crc32.ChecksumIEEE([]byte(pw))
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs[0:], crc)
	b.pw = bs
}

// Subscribe subscribes to bot and waiting notification from SwitchBot.
func (b *Bot) Subscribe() error {
	err := b.subschar.EnableNotifications(func(info []byte) {
		b.subsque <- info
	})
	if err != nil {
		return err
	}
	b.subscribed = true
	return nil
}

// Disconnect  disconnects current SwitchBot connection.
func (b *Bot) Disconnect() error {
	return b.dev.Disconnect()
}

// Press triggers press function for the SwitchBot.
// SwitchBot must be set to press mode.
func (b *Bot) Press(wait bool) error {
	var cmd []byte
	if b.encrypted() {
		cmd = append([]byte{0x57, 0x11}, b.pw...)
	} else {
		cmd = []byte{0x57, 0x01}
	}
	_, err := b.trigger(cmd, wait)
	return err
}

// On triggers on function for the SwitchBot.
// SwitchBot must be set to On/Off mode.
func (b *Bot) On(wait bool) error {
	var cmd []byte
	if b.encrypted() {
		cmd = append(append([]byte{0x57, 0x11}, b.pw...), []byte{0x01}...)
	} else {
		cmd = []byte{0x57, 0x01, 0x01}
	}
	_, err := b.trigger(cmd, wait)
	return err
}

// Off triggers off function for the SwitchBot.
// SwitchBot must be set to On/Off mode.
func (b *Bot) Off(wait bool) error {
	var cmd []byte
	if b.encrypted() {
		cmd = append(append([]byte{0x57, 0x11}, b.pw...), []byte{0x02}...)
	} else {
		cmd = []byte{0x57, 0x01, 0x02}
	}
	_, err := b.trigger(cmd, wait)
	return err
}

// Down triggers down function for the SwitchBot.
func (b *Bot) Down(wait bool) error {
	var cmd []byte
	if b.encrypted() {
		cmd = append(append([]byte{0x57, 0x11}, b.pw...), []byte{0x03}...)
	} else {
		cmd = []byte{0x57, 0x01, 0x03}
	}
	_, err := b.trigger(cmd, wait)
	return err
}

// Up triggers down function for the SwitchBot.
func (b *Bot) Up(wait bool) error {
	var cmd []byte
	if b.encrypted() {
		cmd = append(append([]byte{0x57, 0x11}, b.pw...), []byte{0x04}...)
	} else {
		cmd = []byte{0x57, 0x01, 0x04}
	}
	_, err := b.trigger(cmd, wait)
	return err
}

// GetInfo retrieves bot's settings.
func (b *Bot) GetInfo() (*BotInfo, error) {
	var cmd []byte
	if len(b.pw) != 0 {
		cmd = append([]byte{0x57, 0x12}, b.pw...)
	} else {
		cmd = []byte{0x57, 0x02}
	}

	res, err := b.trigger(cmd, true)
	if err != nil {
		return nil, err
	}
	return NewBotInfoWithRawInfo(res), nil
}

// GetTimers retrieves bot's timer settings.
func (b *Bot) GetTimers(cnt int) ([]*Timer, error) {
	ret := []*Timer{}

	for i := 0; i < cnt; i++ {
		var cmd []byte
		if len(b.pw) != 0 {
			cmd = append([]byte{0x57, 0x18}, b.pw...)
		} else {
			cmd = []byte{0x57, 0x08}
		}
		cmd = append(cmd, []byte{byte(i*16 + 3)}...)
		r, err := b.trigger(cmd, true)
		if err != nil {
			return ret, err
		}
		t := ParseTimerBytes(r)
		ret = append(ret, t)
	}

	return ret, nil
}

func (b *Bot) encrypted() bool {
	return len(b.pw) != 0
}

// trigger executes write characteristics againt SwitchBot.
// response []byte represents following status.
// []byte{1}: trigger success.
// []byte{0}: trigger failure.
func (b *Bot) trigger(cmd []byte, wait bool) ([]byte, error) {
	if wait && !b.subscribed {
		if err := b.Subscribe(); err != nil {
			return []byte{0}, err
		}
	}

	_, err := b.cmdchar.WriteWithoutResponse(cmd)
	if err != nil {
		return []byte{0}, err
	}

	if !wait {
		return []byte{1}, nil
	}

	res := <-b.subsque
	if res[0] != byte(1) {
		return res, errors.New("Failed to send command to SwitchBot")
	}

	return res, nil
}
