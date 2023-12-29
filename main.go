package main

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

var (
	// for a list of codes please check https://www.bluetooth.com/specifications/assigned-numbers/
	airPodsManufacturer uint16 = 76
	airPodsDataLength          = 27
	airpodsMessage      byte   = 7
)

func main() {
	err := adapter.Enable()
	if err != nil {
		panic(err)
	}

	fmt.Println("Scanning AirPods")

	for {
		err := scanAirpods()
		if err != nil {
			panic(err)
		}

		time.Sleep(time.Second * 1)
	}
}

func scanAirpods() error {
	err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		// signal, higher number means the device is closer
		if device.RSSI < -50 {
			return
		}

		data := device.ManufacturerData()
		for key, value := range data {
			if key != airPodsManufacturer {
				return
			}

			if len(value) < airPodsDataLength {
				return
			}

			if value[0] != airpodsMessage {
				return
			}

			hexString := hex.EncodeToString(value)

			fmt.Println(getBatteryFromData(hexString))

			_ = adapter.StopScan()
		}
	})
	return err
}

func getBatteryFromData(data string) string {
	var output string

	status, err := strconv.ParseInt(string(rune(data[15])), 16, 0)
	if err != nil {
		panic(err)
	}

	if status >= 10 {
		output = "Case battery: 100"
	} else {
		output = fmt.Sprintf("Case battery: %v", status*10)
	}

	return output
}
