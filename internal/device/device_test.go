package device

import (
	"os"
	"strings"
	"testing"
)

func TestFindInputDevice(t *testing.T) {
	// Create a temporary test file with devices that have "mouse"/"keyboard" in name
	tempContent := `I: Bus=0003 Vendor=046d Product=c534 Version=0111
N: Name="Logitech Gaming Mouse G502"
P: Phys=usb-0000:00:14.0-1/input0
S: Sysfs=/devices/pci0000:00/0000:00:14.0/usb1/1-1/1-1:1.0/0003:046D:C534.0001/input/input20
U: Uniq=
H: Handlers=mouse0 event4
B: PROP=0
B: EV=17
B: KEY=ffff0000 0 0 0 0
B: REL=903
B: MSC=10

I: Bus=0003 Vendor=04d9 Product=0024 Version=0110
N: Name="USB Keyboard"
P: Phys=usb-0000:00:14.0-2/input0
S: Sysfs=/devices/pci0000:00/0000:00:14.0/usb1/1-2/1-2:1.0/0003:04D9:0024.0002/input/input21
U: Uniq=
H: Handlers=sysrq kbd event5 leds
B: PROP=0
B: EV=120013`

	tmpfile, err := os.CreateTemp("", "devices")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if err := os.WriteFile(tmpfile.Name(), []byte(tempContent), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	// Mock the original file path
	originalReadFile := readFile
	readFile = func(name string) ([]byte, error) {
		if name == "/proc/bus/input/devices" {
			return os.ReadFile(tmpfile.Name())
		}
		return originalReadFile(name)
	}
	defer func() {
		readFile = originalReadFile
	}()

	tests := []struct {
		name        string
		deviceType  string
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name:       "Find mouse device",
			deviceType: "mouse",
			want:       "/dev/input/event4",
			wantErr:    false,
		},
		{
			name:       "Find keyboard device",
			deviceType: "keyboard",
			want:       "/dev/input/event5",
			wantErr:    false,
		},
		{
			name:        "Device not found",
			deviceType:  "nonexistent",
			want:        "",
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindInputDevice(tt.deviceType)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindInputDevice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("FindInputDevice() error = %v, want error containing %v", err, tt.errContains)
				return
			}
			if got != tt.want {
				t.Errorf("FindInputDevice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindInputDevice_HandlerFallback(t *testing.T) {
	// Devices without "mouse"/"keyboard" in name but with correct handler types.
	// This matches BLE mice and Bluetooth keyboards with brand-only names.
	tempContent := `I: Bus=0005 Vendor=046d Product=b034 Version=0006
N: Name="Logitech MX Master 3S"
P: Phys=b8:27:eb:34:20:7b
S: Sysfs=/devices/virtual/misc/uhid/0005:046D:B034.0001/input/input2
U: Uniq=d2:72:87:25:a5:8b
H: Handlers=mouse0 event0
B: PROP=0
B: EV=17
B: KEY=ffff0000 0 0 0 0
B: REL=1943
B: MSC=10

I: Bus=0005 Vendor=0a5c Product=8502 Version=011b
N: Name="Majestouch Convertible 2"
P: Phys=b8:27:eb:34:20:7b
S: Sysfs=/devices/virtual/misc/uhid/0005:0A5C:8502.0002/input/input3
U: Uniq=00:18:00:3e:92:a0
H: Handlers=sysrq kbd leds event1
B: PROP=0
B: EV=12001b`

	tmpfile, err := os.CreateTemp("", "devices")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if err := os.WriteFile(tmpfile.Name(), []byte(tempContent), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	originalReadFile := readFile
	readFile = func(name string) ([]byte, error) {
		if name == "/proc/bus/input/devices" {
			return os.ReadFile(tmpfile.Name())
		}
		return originalReadFile(name)
	}
	defer func() {
		readFile = originalReadFile
	}()

	tests := []struct {
		name       string
		deviceType string
		want       string
		wantErr    bool
	}{
		{
			name:       "Find mouse by handler when name has no 'mouse'",
			deviceType: "mouse",
			want:       "/dev/input/event0",
		},
		{
			name:       "Find keyboard by handler when name has no 'keyboard'",
			deviceType: "keyboard",
			want:       "/dev/input/event1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindInputDevice(tt.deviceType)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindInputDevice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FindInputDevice() = %v, want %v", got, tt.want)
			}
		})
	}
}
