package switchbot

import (
	"context"
	"strings"
	"time"

	"tinygo.org/x/bluetooth"
)

var (
	serviceUUID, _   = bluetooth.ParseUUID("cba20d00-224d-11e6-9fb8-0002a5d5c51b")
	subscribeUUID, _ = bluetooth.ParseUUID("cba20003-224d-11e6-9fb8-0002a5d5c51b")
	commandUUID, _   = bluetooth.ParseUUID("cba20002-224d-11e6-9fb8-0002a5d5c51b")

	adapter *bluetooth.Adapter
)

func init() {
	adapter = bluetooth.DefaultAdapter
}

// Scan scans nearby SwitchBots.
// Callback function will be executed with MAC address once a SwitchBot is found.
// If any SwitchBots are not found, it returns nothing(no timeout error).
func Scan(ctx context.Context, timeout time.Duration, callback func(addr string)) error {
	if err := adapter.Enable(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var err error
	go func() {
		err = adapter.Scan(func(a *bluetooth.Adapter, res bluetooth.ScanResult) {
			if res.LocalName() == "WoHand" {
				callback(res.Address.String())
			}
		})
	}()

	for {
		select {
		case <-ctx.Done():
			adapter.StopScan()
			return scanError(ctx.Err())
		default:
			if err != nil {
				adapter.StopScan()
				return scanError(err)
			}
		}
	}
}

// Connect connects to SwitchBot filter by addr argument.
// If connection failed within timeout, Connect returns error.
func Connect(ctx context.Context, addr string, timeout time.Duration) (*Bot, error) {
	if err := adapter.Enable(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var res *bluetooth.ScanResult
	var err error
	go func() {
		err = adapter.Scan(func(a *bluetooth.Adapter, sres bluetooth.ScanResult) {
			if strings.ToUpper(sres.Address.String()) != addr {
				return
			}
			a.StopScan()
			res = &sres
		})
		if err != nil {
			cancel()
		}
	}()

wait:
	for {
		select {
		case <-ctx.Done():
			adapter.StopScan()
			return nil, ctx.Err()
		default:
			if err != nil {
				return nil, err
			}
			if res != nil {
				break wait
			}
		}
	}

	device, err := adapter.Connect(res.Address, bluetooth.ConnectionParams{})
	if err != nil {
		return nil, err
	}

	srvcs, err := device.DiscoverServices([]bluetooth.UUID{serviceUUID})
	if err != nil {
		return nil, err
	}

	var srvc *bluetooth.DeviceService
	var bot *Bot
	var cmdchar *bluetooth.DeviceCharacteristic
	var subschar *bluetooth.DeviceCharacteristic

	for _, dsvc := range srvcs {
		if dsvc.UUID().String() == serviceUUID.String() {
			srvc = &dsvc

			chars, err := srvc.DiscoverCharacteristics([]bluetooth.UUID{commandUUID})
			if err != nil {
				return nil, err
			}
			cmdchar = &chars[0]

			chars, err = srvc.DiscoverCharacteristics([]bluetooth.UUID{subscribeUUID})
			if err != nil {
				return nil, err
			}
			subschar = &chars[0]

			break
		}
	}

	bot = &Bot{
		Addr: res.Address.String(),

		dev: device,

		cmdchar:  cmdchar,
		subschar: subschar,

		subsque:    make(chan []byte),
		subscribed: false,
	}
	return bot, nil
}

func scanError(err error) error {
	switch err {
	case nil, context.DeadlineExceeded, context.Canceled:
		return nil
	default:
		return err
	}
}
