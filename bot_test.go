package switchbot

import (
	"reflect"
	"testing"

	"github.com/JuulLabs-OSS/ble"
)

type MockBleClient struct {
	ble.Client

	writeCharacteristics func(c *ble.Characteristic, v []byte, noRsp bool) error
}

func (p *MockBleClient) WriteCharacteristic(c *ble.Characteristic, v []byte, noRsp bool) error {
	return p.writeCharacteristics(c, v, noRsp)
}

func Test_Press(t *testing.T) {
	cl := &MockBleClient{}
	cl.writeCharacteristics = func(c *ble.Characteristic, v []byte, noRsp bool) error {
		if c.ValueHandle != 0x16 {
			t.Fatal("Incorrect VHandle")
		}
		if !reflect.DeepEqual(v, []byte{0x57, 0x01, 0x00}) {
			t.Fatal("Incorrect value")
		}
		if noRsp != false {
			t.Fatal("Incorrect noRsp")
		}
		return nil
	}

	bot := &Bot{"ADDR", cl}
	bot.cl = cl
	if err := bot.Press(); err != nil {
		t.Fatal("test failed")
	}
}

func Test_On(t *testing.T) {
	cl := &MockBleClient{}
	cl.writeCharacteristics = func(c *ble.Characteristic, v []byte, noRsp bool) error {
		if c.ValueHandle != 0x16 {
			t.Fatal("Incorrect VHandle")
		}
		if !reflect.DeepEqual(v, []byte{0x57, 0x01, 0x01}) {
			t.Fatal("Incorrect value")
		}
		if noRsp != false {
			t.Fatal("Incorrect noRsp")
		}
		return nil
	}

	bot := &Bot{"ADDR", cl}
	bot.cl = cl
	if err := bot.On(); err != nil {
		t.Fatal("test failed")
	}
}

func Test_Off(t *testing.T) {
	cl := &MockBleClient{}
	cl.writeCharacteristics = func(c *ble.Characteristic, v []byte, noRsp bool) error {
		if c.ValueHandle != 0x16 {
			t.Fatal("Incorrect VHandle")
		}
		if !reflect.DeepEqual(v, []byte{0x57, 0x01, 0x02}) {
			t.Fatal("Incorrect value")
		}
		if noRsp != false {
			t.Fatal("Incorrect noRsp")
		}
		return nil
	}

	bot := &Bot{"ADDR", cl}
	if err := bot.Off(); err != nil {
		t.Fatal("test failed")
	}
}
