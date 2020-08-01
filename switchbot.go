package switchbot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/JuulLabs-OSS/ble"
	"github.com/JuulLabs-OSS/ble/linux"
)

var (
	serviceUUID = ble.MustParse("cba20d00224d11e69fb80002a5d5c51b")
	bleDevice   ble.Device
)

// Scan scans nearby SwitchBots.
// Callback function will be executed with MAC address once a SwitchBot is found.
// If any SwitchBots are not found, it returns nothing(no timeout error).
func Scan(ctx context.Context, timeout time.Duration, f func(addr string)) error {
	err := setupDefaultDevice()
	if err != nil {
		return fmt.Errorf("Cound not initialize a device: %w", err)
	}

	ctx = ble.WithSigHandler(context.WithTimeout(ctx, timeout))
	err = ble.Scan(ctx, false, func(a ble.Advertisement) {
		if contains(a.Services(), serviceUUID) {
			f(a.Addr().String())
		}
	}, nil)
	return scanError(err)
}

// Connect connects to SwitchBot filter by addr argument.
// If connection failed within timeout, Connect returns error.
func Connect(ctx context.Context, addr string, timeout time.Duration) (*Bot, error) {
	if err := setupDefaultDevice(); err != nil {
		return nil, fmt.Errorf("Cound not initialize a device: %w", err)
	}

	ctx = ble.WithSigHandler(context.WithTimeout(ctx, timeout))
	addr = strings.ToLower(addr)
	cl, err := ble.Connect(ctx, func(a ble.Advertisement) bool {
		return a.Addr().String() == addr
	})
	if err != nil {
		return nil, err
	}
	bot := NewBot(addr)
	bot.cl = cl
	return bot, nil
}

func setupDefaultDevice() error {
	if bleDevice == nil {
		d, err := linux.NewDevice()
		if err != nil {
			return err
		}
		bleDevice = d
	}
	ble.SetDefaultDevice(bleDevice)
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

func scanError(err error) error {
	switch err {
	case nil, context.DeadlineExceeded, context.Canceled:
		return nil
	default:
		return err
	}
}
