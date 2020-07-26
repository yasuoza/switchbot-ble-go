package switchbot

import (
	"fmt"
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

func TestSetPassword(t *testing.T) {
	cl := &MockBleClient{}
	bot := &Bot{"ADDR", cl, []byte{}}
	bot.SetPassword("password")
	psw := fmt.Sprintf("% x", bot.pw)
	if psw != "35 c2 46 d5" {
		t.Fatal("Incorrect password")
	}
}

func TestPress(t *testing.T) {
	cl := &MockBleClient{}
	cl.writeCharacteristics = func(c *ble.Characteristic, v []byte, noRsp bool) error {
		if c.ValueHandle != 0x16 {
			t.Fatal("Incorrect VHandle")
		}
		if !reflect.DeepEqual(v, []byte{0x57, 0x01}) {
			t.Fatal("Incorrect value")
		}
		if noRsp != false {
			t.Fatal("Incorrect noRsp")
		}
		return nil
	}

	bot := &Bot{"ADDR", cl, []byte{}}
	bot.cl = cl
	if err := bot.Press(); err != nil {
		t.Fatal("test failed")
	}
}

func TestPressWithPassword(t *testing.T) {
	cl := &MockBleClient{}
	cl.writeCharacteristics = func(c *ble.Characteristic, v []byte, noRsp bool) error {
		if c.ValueHandle != 0x16 {
			t.Fatal("Incorrect VHandle")
		}
		if !reflect.DeepEqual(v, append([]byte{0x57, 0x11}, []byte{0x35, 0xc2, 0x46, 0xd5}...)) {
			t.Fatal("Incorrect value")
		}
		if noRsp != false {
			t.Fatal("Incorrect noRsp")
		}
		return nil
	}

	bot := &Bot{"ADDR", cl, []byte{}}
	bot.SetPassword("password")
	bot.cl = cl
	if err := bot.Press(); err != nil {
		t.Fatal("test failed")
	}
}

func TestOn(t *testing.T) {
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

	bot := &Bot{"ADDR", cl, []byte{}}
	bot.cl = cl
	if err := bot.On(); err != nil {
		t.Fatal("test failed")
	}
}

func TestOnWithPassword(t *testing.T) {
	cl := &MockBleClient{}
	cl.writeCharacteristics = func(c *ble.Characteristic, v []byte, noRsp bool) error {
		if c.ValueHandle != 0x16 {
			t.Fatal("Incorrect VHandle")
		}
		cmd := append(append([]byte{0x57, 0x11}, []byte{0x35, 0xc2, 0x46, 0xd5}...), []byte{0x01}...)
		if !reflect.DeepEqual(v, cmd) {
			t.Fatal("Incorrect value")
		}
		if noRsp != false {
			t.Fatal("Incorrect noRsp")
		}
		return nil
	}

	bot := &Bot{"ADDR", cl, []byte{}}
	bot.SetPassword("password")
	bot.cl = cl
	if err := bot.On(); err != nil {
		t.Fatal("test failed")
	}
}

func TestOff(t *testing.T) {
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

	bot := &Bot{"ADDR", cl, []byte{}}
	if err := bot.Off(); err != nil {
		t.Fatal("test failed")
	}
}

func TestOffWithPassword(t *testing.T) {
	cl := &MockBleClient{}
	cl.writeCharacteristics = func(c *ble.Characteristic, v []byte, noRsp bool) error {
		if c.ValueHandle != 0x16 {
			t.Fatal("Incorrect VHandle")
		}
		cmd := append(append([]byte{0x57, 0x11}, []byte{0x35, 0xc2, 0x46, 0xd5}...), []byte{0x02}...)
		if !reflect.DeepEqual(v, cmd) {
			t.Fatal("Incorrect value")
		}
		if noRsp != false {
			t.Fatal("Incorrect noRsp")
		}
		return nil
	}

	bot := &Bot{"ADDR", cl, []byte{}}
	bot.SetPassword("password")
	bot.cl = cl
	if err := bot.Off(); err != nil {
		t.Fatal("test failed")
	}
}
