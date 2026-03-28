package relay

import (
	"bytes"
	"testing"
)

func TestKeyboardRelay_ConvertEvent(t *testing.T) {
	tests := []struct {
		name       string
		event      InputEvent
		wantReport []byte
		wantErr    bool
	}{
		{
			name: "simple key press A",
			event: InputEvent{
				Type:  1,  // EV_KEY
				Code:  30, // KEY_A
				Value: 1,  // Press
			},
			wantReport: []byte{0, 0, 0x04, 0, 0, 0, 0, 0},
			wantErr:    false,
		},
		{
			name: "key release A",
			event: InputEvent{
				Type:  1,
				Code:  30,
				Value: 0, // Release
			},
			wantReport: []byte{0, 0, 0, 0, 0, 0, 0, 0},
			wantErr:    false,
		},
		{
			name: "left shift press",
			event: InputEvent{
				Type:  1,
				Code:  42, // Left Shift
				Value: 1,
			},
			wantReport: []byte{0x02, 0, 0, 0, 0, 0, 0, 0},
			wantErr:    false,
		},
		{
			name: "F1 press",
			event: InputEvent{
				Type:  1,
				Code:  59, // KEY_F1
				Value: 1,
			},
			wantReport: []byte{0, 0, 0x3A, 0, 0, 0, 0, 0},
			wantErr:    false,
		},
		{
			name: "F11 press",
			event: InputEvent{
				Type:  1,
				Code:  87, // KEY_F11
				Value: 1,
			},
			wantReport: []byte{0, 0, 0x44, 0, 0, 0, 0, 0},
			wantErr:    false,
		},
		{
			name: "F12 press",
			event: InputEvent{
				Type:  1,
				Code:  88, // KEY_F12
				Value: 1,
			},
			wantReport: []byte{0, 0, 0x45, 0, 0, 0, 0, 0},
			wantErr:    false,
		},
		{
			name: "Home press",
			event: InputEvent{
				Type:  1,
				Code:  102, // KEY_HOME
				Value: 1,
			},
			wantReport: []byte{0, 0, 0x4A, 0, 0, 0, 0, 0},
			wantErr:    false,
		},
		{
			name: "End press",
			event: InputEvent{
				Type:  1,
				Code:  107, // KEY_END
				Value: 1,
			},
			wantReport: []byte{0, 0, 0x4D, 0, 0, 0, 0, 0},
			wantErr:    false,
		},
		{
			name: "Page Up press",
			event: InputEvent{
				Type:  1,
				Code:  104, // KEY_PAGEUP
				Value: 1,
			},
			wantReport: []byte{0, 0, 0x4B, 0, 0, 0, 0, 0},
			wantErr:    false,
		},
		{
			name: "Page Down press",
			event: InputEvent{
				Type:  1,
				Code:  109, // KEY_PAGEDOWN
				Value: 1,
			},
			wantReport: []byte{0, 0, 0x4E, 0, 0, 0, 0, 0},
			wantErr:    false,
		},
		{
			name: "Up Arrow press",
			event: InputEvent{
				Type:  1,
				Code:  103, // KEY_UP
				Value: 1,
			},
			wantReport: []byte{0, 0, 0x52, 0, 0, 0, 0, 0},
			wantErr:    false,
		},
		{
			name: "KP Enter press",
			event: InputEvent{
				Type:  1,
				Code:  96, // KEY_KPENTER
				Value: 1,
			},
			wantReport: []byte{0, 0, 0x58, 0, 0, 0, 0, 0},
			wantErr:    false,
		},
		{
			name: "NumLock press",
			event: InputEvent{
				Type:  1,
				Code:  69, // KEY_NUMLOCK
				Value: 1,
			},
			wantReport: []byte{0, 0, 0x53, 0, 0, 0, 0, 0},
			wantErr:    false,
		},
		{
			name: "KP 5 press",
			event: InputEvent{
				Type:  1,
				Code:  76, // KEY_KP5
				Value: 1,
			},
			wantReport: []byte{0, 0, 0x5D, 0, 0, 0, 0, 0},
			wantErr:    false,
		},
		{
			name: "unmapped key returns nil",
			event: InputEvent{
				Type:  1,
				Code:  999, // unmapped
				Value: 1,
			},
			wantReport: nil,
			wantErr:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			k := new(KeyboardRelay)
			gotReport, err := k.convertEvent(test.event)
			if (err != nil) != test.wantErr {
				t.Errorf("KeyboardRelay.ConvertEvent() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !bytes.Equal(gotReport, test.wantReport) {
				t.Errorf("KeyboardRelay.ConvertEvent() = %v, want %v", gotReport, test.wantReport)
			}
		})
	}
}
