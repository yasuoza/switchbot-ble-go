package switchbot

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/JuulLabs-OSS/ble"
)

type MockBleClient struct {
	ble.Client

	writeCharacteristics func(c *ble.Characteristic, v []byte, noRsp bool) error
	subscribe            func(c *ble.Characteristic, ind bool, h ble.NotificationHandler) error
}

func (cl *MockBleClient) DiscoverProfile(force bool) (*ble.Profile, error) {
	s := ble.NewService(ble.MustParse("7b74bec2-ce6f-11ea-87d0-0242ac130003"))
	s.Characteristics = []*ble.Characteristic{ble.NewCharacteristic(ble.MustParse("cba20003-224d-11e6-9fb8-0002a5d5c51b"))}
	pr := &ble.Profile{}
	pr.Services = []*ble.Service{s}
	return pr, nil
}

func (cl *MockBleClient) WriteCharacteristic(c *ble.Characteristic, v []byte, noRsp bool) error {
	return cl.writeCharacteristics(c, v, noRsp)
}

func (cl *MockBleClient) Subscribe(c *ble.Characteristic, ind bool, h ble.NotificationHandler) error {
	return cl.subscribe(c, ind, h)
}

type MockProfile struct {
	ble.Profile

	mockCharactercteristic *ble.Characteristic
}

func (p *MockProfile) FindCharacteristic(c *ble.Characteristic) *ble.Characteristic {
	return p.mockCharactercteristic
}

func newBot(addr string, cl *MockBleClient) *Bot {
	b := NewBot("ADDR")
	b.cl = cl
	return b
}

func TestNewBot(t *testing.T) {
	b := NewBot("ADDR")
	if b.Addr != "addr" {
		t.Fatal("Addr must be conterted to lowercase")
	}
}

func TestSetPassword(t *testing.T) {
	cl := &MockBleClient{}
	bot := newBot("ADDR", cl)
	bot.SetPassword("password")
	psw := fmt.Sprintf("% x", bot.pw)
	if psw != "35 c2 46 d5" {
		t.Fatal("Incorrect password")
	}
}

func TestSubscribeError(t *testing.T) {
	cl := &MockBleClient{}
	cl.subscribe = func(c *ble.Characteristic, ind bool, h ble.NotificationHandler) error {
		return errors.New("Subscribe error")
	}
	bot := newBot("ADDR", cl)
	err := bot.Subscribe()
	if err == nil {
		t.Fatal("Must not return error")
	}
}

func TestSubscribeOk(t *testing.T) {
	cl := &MockBleClient{}
	cl.subscribe = func(c *ble.Characteristic, ind bool, h ble.NotificationHandler) error {
		go h([]byte{12, 12, 12})
		return nil
	}
	bot := newBot("ADDR", cl)
	err := bot.Subscribe()
	if err != nil {
		t.Fatal("Must not return error")
	}
}

func TestGetInfoOk(t *testing.T) {
	cl := &MockBleClient{}
	raw := []byte{1, 79, 45, 100, 0, 0, 0, 152, 0, 0, 3, 72, 0}
	cl.subscribe = func(c *ble.Characteristic, ind bool, h ble.NotificationHandler) error {
		go h(raw)
		return nil
	}
	cl.writeCharacteristics = func(c *ble.Characteristic, v []byte, noRsp bool) error {
		if !bytes.Equal(v, []byte{0x57, 0x02}) {
			t.Fatal("Invalid cmd")
		}
		return nil
	}
	bot := newBot("ADDR", cl)
	info, err := bot.GetInfo()
	if err != nil {
		t.Fatal("Must not return error")
	}
	if info == nil {
		t.Errorf("info must not to be nil, but %v", info)
	}
}

func TestGetTimersOk(t *testing.T) {
	cl := &MockBleClient{}
	raw := [][]byte{
		{1, 2, 0, 121, 10, 11, 0, 0, 0, 0, 0, 0},
		{1, 2, 0, 0, 12, 0, 64, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // nil
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // nil
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // nil
	}
	scnt := 0
	cl.subscribe = func(c *ble.Characteristic, ind bool, h ble.NotificationHandler) error {
		go func() {
			h(raw[scnt])
			scnt++
		}()
		return nil
	}
	wcnt := 0
	cl.writeCharacteristics = func(c *ble.Characteristic, v []byte, noRsp bool) error {
		if !bytes.Equal(v, append([]byte{0x57, 0x08}, []byte{byte(wcnt*16 + 3)}...)) {
			t.Fatal("Invalid cmd")
		}
		wcnt++
		return nil
	}
	bot := newBot("ADDR", cl)
	timers, err := bot.GetTimers(len(raw))
	if err != nil {
		t.Fatal("Must not return error")
	}
	for i := 0; i < len(raw); i++ {
		timer := timers[i]
		if i < 2 && timer == nil {
			t.Errorf("timer at index %d is expected not nil, but got nil\n", i)
		}
		if i > 2 && timer != nil {
			t.Errorf("timer at index %d is expected nil, but got %v\n", i, timer)
		}
	}
}

func TestPress(t *testing.T) {
	cl := &MockBleClient{}
	cl.subscribe = func(c *ble.Characteristic, ind bool, h ble.NotificationHandler) error {
		go h([]byte{1, 255, 0})
		return nil
	}
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

	bot := newBot("ADDR", cl)
	bot.cl = cl
	if err := bot.Press(false); err != nil {
		t.Fatal("test failed")
	}
}

func TestPressWithPassword(t *testing.T) {
	cl := &MockBleClient{}
	cl.subscribe = func(c *ble.Characteristic, ind bool, h ble.NotificationHandler) error {
		go h([]byte{1, 255, 0})
		return nil
	}
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

	bot := newBot("ADDR", cl)
	bot.SetPassword("password")
	bot.cl = cl
	if err := bot.Press(false); err != nil {
		t.Fatal("test failed")
	}
}

func TestPressWaitFail(t *testing.T) {
	cl := &MockBleClient{}
	cl.subscribe = func(c *ble.Characteristic, ind bool, h ble.NotificationHandler) error {
		go h([]byte{0, 255, 0})
		return nil
	}
	cl.writeCharacteristics = func(c *ble.Characteristic, v []byte, noRsp bool) error {
		return nil
	}

	bot := newBot("ADDR", cl)
	bot.cl = cl
	if err := bot.Press(true); err == nil {
		t.Fatal("Must return error")
	}
}

func TestPressWaitSuccess(t *testing.T) {
	cl := &MockBleClient{}
	cl.subscribe = func(c *ble.Characteristic, ind bool, h ble.NotificationHandler) error {
		go h([]byte{1, 255, 0})
		return nil
	}
	cl.writeCharacteristics = func(c *ble.Characteristic, v []byte, noRsp bool) error {
		return nil
	}

	bot := newBot("ADDR", cl)
	bot.cl = cl
	if err := bot.Press(true); err != nil {
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

	bot := newBot("ADDR", cl)
	bot.cl = cl
	if err := bot.On(false); err != nil {
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

	bot := newBot("ADDR", cl)
	bot.SetPassword("password")
	bot.cl = cl
	if err := bot.On(false); err != nil {
		t.Fatal("test failed")
	}
}

func TestOnWaitFail(t *testing.T) {
	cl := &MockBleClient{}
	cl.subscribe = func(c *ble.Characteristic, ind bool, h ble.NotificationHandler) error {
		go h([]byte{0, 255, 0})
		return nil
	}
	cl.writeCharacteristics = func(c *ble.Characteristic, v []byte, noRsp bool) error {
		return nil
	}

	bot := newBot("ADDR", cl)
	bot.cl = cl
	if err := bot.On(true); err == nil {
		t.Fatal("Must return error")
	}
}

func TestOnWaitSuccess(t *testing.T) {
	cl := &MockBleClient{}
	cl.subscribe = func(c *ble.Characteristic, ind bool, h ble.NotificationHandler) error {
		go h([]byte{1, 255, 0})
		return nil
	}
	cl.writeCharacteristics = func(c *ble.Characteristic, v []byte, noRsp bool) error {
		return nil
	}

	bot := newBot("ADDR", cl)
	bot.cl = cl
	if err := bot.On(true); err != nil {
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

	bot := newBot("ADDR", cl)
	if err := bot.Off(false); err != nil {
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

	bot := newBot("ADDR", cl)
	bot.SetPassword("password")
	bot.cl = cl
	if err := bot.Off(false); err != nil {
		t.Fatal("test failed")
	}
}

func TestOffWaitFail(t *testing.T) {
	cl := &MockBleClient{}
	cl.subscribe = func(c *ble.Characteristic, ind bool, h ble.NotificationHandler) error {
		go h([]byte{0, 255, 0})
		return nil
	}
	cl.writeCharacteristics = func(c *ble.Characteristic, v []byte, noRsp bool) error {
		return nil
	}

	bot := newBot("ADDR", cl)
	bot.cl = cl
	if err := bot.Off(true); err == nil {
		t.Fatal("Must return error")
	}
}

func TestOffWaitSuccess(t *testing.T) {
	cl := &MockBleClient{}
	cl.subscribe = func(c *ble.Characteristic, ind bool, h ble.NotificationHandler) error {
		go h([]byte{1, 255, 0})
		return nil
	}
	cl.writeCharacteristics = func(c *ble.Characteristic, v []byte, noRsp bool) error {
		return nil
	}

	bot := newBot("ADDR", cl)
	bot.cl = cl
	if err := bot.Off(true); err != nil {
		t.Fatal("test failed")
	}
}
