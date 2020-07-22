package switchbot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/JuulLabs-OSS/ble"
	"github.com/JuulLabs-OSS/ble/linux"
	"github.com/pkg/errors"
)

var (
	serviceUUID = ble.MustParse("cba20d00224d11e69fb80002a5d5c51b")
)

func Scan(ctx context.Context, timeout time.Duration) error {
	err := setupBleDevice()
	if err != nil {
		return errors.Wrap(err, "Could not initialize a device")
	}

	ctx = ble.WithSigHandler(context.WithTimeout(ctx, timeout))
	err = ble.Scan(ctx, false, func(a ble.Advertisement) {
		if contains(a.Services(), serviceUUID) {
			fmt.Printf("Found Switch Bot. MAC Address: %s\n", a.Addr().String())
		}
	}, nil)
	return handleErr(err)
}

func Connect(ctx context.Context, addr string, timeout time.Duration) (*Bot, error) {
	err := setupBleDevice()
	if err != nil {
		return nil, errors.Wrap(err, "Could not initialize a device")
	}

	ctx = ble.WithSigHandler(context.WithTimeout(ctx, timeout))
	addr = strings.ToLower(addr)
	cl, err := ble.Connect(ctx, func(a ble.Advertisement) bool {
		return a.Addr().String() == addr
	})
	if err != nil {
		return nil, err
	}
	bot := &Bot{Addr: addr, cl: cl}
	return bot, nil
}

func setupBleDevice() error {
	d, err := linux.NewDevice()
	if err != nil {
		return errors.Wrap(err, "Could not initialize a device")
	}
	ble.SetDefaultDevice(d)
	return nil
}

func contains(arr []ble.UUID, uuid ble.UUID) bool {
	for _, u := range arr {
		if u.Equal(uuid) {
			return true
		}
	}
	return false
}

func handleErr(err error) error {
	switch errors.Cause(err) {
	case nil, context.DeadlineExceeded, context.Canceled:
		return nil
	default:
		return err
	}
}
