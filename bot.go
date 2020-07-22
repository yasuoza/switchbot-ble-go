package switchbot

import (
	"github.com/JuulLabs-OSS/ble"
)

const (
	characteristics = "cba20002-224d-11e6-9fb8-0002a5d5c51b"
	handle          = 0x16
)

var (
	press = []byte{0x57, 0x01, 0x00}
	on    = []byte{0x57, 0x01, 0x01}
	off   = []byte{0x57, 0x01, 0x02}
)

type Bot struct {
	Addr string

	cl ble.Client
}

func NewBot(Addr string) *Bot {
	return &Bot{Addr: Addr}
}

func (b *Bot) Press() error {
	return b.trigger(press)
}

func (b *Bot) On() error {
	return b.trigger(on)
}

func (b *Bot) Off() error {
	return b.trigger(off)
}

func (b *Bot) trigger(op []byte) error {
	c := ble.NewCharacteristic(ble.MustParse(characteristics))
	c.ValueHandle = handle
	if err := b.cl.WriteCharacteristic(c, op, false); err != nil {
		return err
	}
	return nil
}
