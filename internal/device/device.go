package device

import (
	"fmt"
	"os"
	"strings"
)

// DeviceType represents the type of HID device
type DeviceType int

const (
	Mouse DeviceType = iota
	Keyboard
)

// Device represents a HID device interface
type Device interface {
	Open() error
	Close() error
	Write([]byte) error
	SendRelease() error
}

// DeviceConfig holds device configuration
type DeviceConfig struct {
	InputPath  string
	OutputPath string
	Type       DeviceType
}

var FindInputDeviceFunc = FindInputDevice
var readFile = os.ReadFile

// handlerKeywords maps device types to their Linux input handler prefixes.
// The kernel assigns these handlers regardless of the device's product name.
var handlerKeywords = map[string]string{
	"mouse":    "mouse",
	"keyboard": "kbd",
}

func FindInputDevice(deviceType string) (string, error) {
	data, err := readFile("/proc/bus/input/devices")
	if err != nil {
		return "", fmt.Errorf("failed to read devices: %v", err)
	}

	var nameMatched bool

	for _, line := range strings.Split(string(data), "\n") {
		// Reset on new device block (blocks separated by blank lines)
		if strings.TrimSpace(line) == "" {
			nameMatched = false
			continue
		}

		// Check device name for match
		if strings.HasPrefix(line, "N: Name=") {
			if strings.Contains(strings.ToLower(line), deviceType) {
				nameMatched = true
			}
			continue
		}

		// Check handlers line for matching device
		if strings.HasPrefix(line, "H: Handlers=") {
			handlers := strings.TrimPrefix(line, "H: Handlers=")
			matched := nameMatched

			// Also check handler keywords (e.g. "mouse0" for mice, "kbd" for keyboards)
			if !matched {
				if kw, ok := handlerKeywords[deviceType]; ok {
					for _, handler := range strings.Fields(handlers) {
						if strings.HasPrefix(handler, kw) {
							matched = true
							break
						}
					}
				}
			}

			if matched {
				for _, word := range strings.Fields(handlers) {
					if strings.HasPrefix(word, "event") {
						return fmt.Sprintf("/dev/input/%s", word), nil
					}
				}
			}
			nameMatched = false
		}
	}

	return "", fmt.Errorf("%s not found", deviceType)
}
