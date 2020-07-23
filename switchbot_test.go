package switchbot

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/JuulLabs-OSS/ble"
)

type MockDevice struct {
	ble.Device

	scan func(ctx context.Context, allowDup bool, h ble.AdvHandler) error
}

func (d *MockDevice) Scan(ctx context.Context, allowDup bool, h ble.AdvHandler) error {
	return d.scan(ctx, allowDup, h)
}

func (d *MockDevice) Dial(ctx context.Context, a ble.Addr) (ble.Client, error) {
	return struct{ ble.Client }{}, nil
}

type MockAdvertisement struct {
	ble.Advertisement

	addr     string
	services []ble.UUID
}

func (a MockAdvertisement) Addr() ble.Addr {
	return ble.NewAddr(a.addr)
}

func (a MockAdvertisement) Services() []ble.UUID {
	return a.services
}

func Test_Scan_Error(t *testing.T) {
	d := &MockDevice{}
	d.scan = func(ctx context.Context, allowDup bool, h ble.AdvHandler) error {
		return errors.New("Scan Error")
	}
	bleDevice = d

	if _, err := Scan(context.Background(), 0); err == nil {
		t.Fatal("Must return error")
	}
}

func Test_Scan_Not_Found(t *testing.T) {
	d := &MockDevice{}
	d.scan = func(ctx context.Context, allowDup bool, h ble.AdvHandler) error {
		return context.DeadlineExceeded
	}
	bleDevice = d

	if _, err := Scan(context.Background(), 0); err != nil {
		t.Fatal("Must not return DeadlineExceeded error")
	}
}

func Test_Scan_Found(t *testing.T) {
	advs := []MockAdvertisement{
		*&MockAdvertisement{
			addr:     "7F:8E:6B:F5:CA:91",
			services: []ble.UUID{ble.MustParse("dc67b962-ccbe-11ea-87d0-0242ac130003")},
		},
		*&MockAdvertisement{
			addr:     "9D:76:72:29:40:83",
			services: []ble.UUID{ble.MustParse("cba20d00-224d-11e6-9fb8-0002a5d5c51b")},
		},
		*&MockAdvertisement{
			addr:     "4D:24:A8:D9:43:6C",
			services: []ble.UUID{ble.MustParse("cba20d00-224d-11e6-9fb8-0002a5d5c51b")},
		},
	}

	d := &MockDevice{}
	d.scan = func(ctx context.Context, allowDup bool, h ble.AdvHandler) error {
		for _, a := range advs {
			h(a)
		}
		return nil
	}
	bleDevice = d

	addrs, err := Scan(context.Background(), 0)
	if err != nil {
		t.Fatal("Must not return error")
	}

	if !reflect.DeepEqual(addrs, []string{"9d:76:72:29:40:83", "4d:24:a8:d9:43:6c"}) {
		t.Fatal("Must return found mac addresses")
	}
}

func Test_Connect_NG(t *testing.T) {
	addr := "7f:8e:6b:f5:ca:91"
	d := &MockDevice{}
	d.scan = func(ctx context.Context, allowDup bool, h ble.AdvHandler) error {
		return context.DeadlineExceeded
	}
	bleDevice = d

	bot, err := Connect(context.Background(), addr, 0)
	if err == nil {
		t.Fatal("Must return error")
	}
	if bot != nil {
		t.Fatal("Incorrect bot")
	}
}

func Test_Connect_OK(t *testing.T) {
	addr := "9d:76:72:29:40:83"
	a := &MockAdvertisement{
		addr: addr,
	}
	d := &MockDevice{}
	d.scan = func(ctx context.Context, allowDup bool, h ble.AdvHandler) error {
		go h(a)
		return nil
	}
	bleDevice = d

	bot, err := Connect(context.Background(), addr, 0)
	if err != nil {
		t.Fatal("Must not return error")
	}
	if bot.Addr != addr {
		t.Fatal("Incorrect bot")
	}
}
