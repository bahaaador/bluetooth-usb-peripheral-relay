package relay

import (
	"testing"
)

func TestMouseRelay_ConvertEvent(t *testing.T) {
	tests := []struct {
		name       string
		event      InputEvent
		wantReport []byte
		wantErr    bool
	}{
		{
			name: "left button press",
			event: InputEvent{
				Type:  1,   // EV_KEY
				Code:  272, // BTN_LEFT
				Value: 1,
			},
			wantReport: []byte{0x01, 0, 0, 0, 0},
			wantErr:    false,
		},
		{
			name: "mouse move right",
			event: InputEvent{
				Type:  2, // EV_REL
				Code:  0, // REL_X
				Value: 10,
			},
			wantReport: []byte{0, 10, 0, 0, 0},
			wantErr:    false,
		},
		{
			name: "mouse move down",
			event: InputEvent{
				Type:  2, // EV_REL
				Code:  1, // REL_Y
				Value: 5,
			},
			wantReport: []byte{0, 0, 5, 0, 0},
			wantErr:    false,
		},
		{
			name: "scroll wheel",
			event: InputEvent{
				Type:  2, // EV_REL
				Code:  8, // REL_WHEEL
				Value: 1,
			},
			wantReport: []byte{0, 0, 0, 1, 0},
			wantErr:    false,
		},
		{
			name: "horizontal scroll",
			event: InputEvent{
				Type:  2, // EV_REL
				Code:  6, // REL_HWHEEL
				Value: 1,
			},
			wantReport: []byte{0, 0, 0, 0, 1},
			wantErr:    false,
		},
		{
			name: "side button press (BTN_SIDE)",
			event: InputEvent{
				Type:  1,   // EV_KEY
				Code:  275, // BTN_SIDE
				Value: 1,
			},
			wantReport: []byte{0x08, 0, 0, 0, 0}, // bit 3
			wantErr:    false,
		},
		{
			name: "extra button press (BTN_EXTRA)",
			event: InputEvent{
				Type:  1,   // EV_KEY
				Code:  276, // BTN_EXTRA
				Value: 1,
			},
			wantReport: []byte{0x10, 0, 0, 0, 0}, // bit 4
			wantErr:    false,
		},
		{
			name: "forward button press (BTN_FORWARD maps to bit 4)",
			event: InputEvent{
				Type:  1,   // EV_KEY
				Code:  277, // BTN_FORWARD
				Value: 1,
			},
			wantReport: []byte{0x10, 0, 0, 0, 0}, // bit 4, same as BTN_EXTRA
			wantErr:    false,
		},
		{
			name: "back button press (BTN_BACK maps to bit 3)",
			event: InputEvent{
				Type:  1,   // EV_KEY
				Code:  278, // BTN_BACK
				Value: 1,
			},
			wantReport: []byte{0x08, 0, 0, 0, 0}, // bit 3, same as BTN_SIDE
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MouseRelay{}
			report, err := m.convertEvent(tt.event)

			if (err != nil) != tt.wantErr {
				t.Errorf("MouseRelay.convertEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(report) != len(tt.wantReport) {
				t.Errorf("MouseRelay.convertEvent() report length = %d, want %d", len(report), len(tt.wantReport))
				return
			}

			for i := range report {
				if report[i] != tt.wantReport[i] {
					t.Errorf("MouseRelay.convertEvent() report[%d] = %#x, want %#x", i, report[i], tt.wantReport[i])
				}
			}
		})
	}
}

func TestMouseRelay_ValidateEvent(t *testing.T) {
	tests := []struct {
		name  string
		event InputEvent
		want  bool
	}{
		{
			name: "valid button event (BTN_LEFT)",
			event: InputEvent{
				Type:  1,
				Code:  272, // BTN_LEFT
				Value: 1,
			},
			want: true,
		},
		{
			name: "valid button event (BTN_BACK)",
			event: InputEvent{
				Type:  1,
				Code:  278, // BTN_BACK
				Value: 1,
			},
			want: true,
		},
		{
			name: "invalid button event (out of range)",
			event: InputEvent{
				Type:  1,
				Code:  279,
				Value: 1,
			},
			want: false,
		},
		{
			name: "valid movement event",
			event: InputEvent{
				Type:  2,
				Code:  0, // REL_X
				Value: 10,
			},
			want: true,
		},
		{
			name: "valid horizontal wheel event",
			event: InputEvent{
				Type:  2,
				Code:  6, // REL_HWHEEL
				Value: 1,
			},
			want: true,
		},
		{
			name: "sync event",
			event: InputEvent{
				Type:  0,
				Code:  0,
				Value: 0,
			},
			want: false,
		},
		{
			name: "invalid event type",
			event: InputEvent{
				Type:  5,
				Code:  0,
				Value: 0,
			},
			want: false,
		},
	}

	m := &MouseRelay{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := m.validateEvent(test.event); got != test.want {
				t.Errorf("MouseRelay.validateEvent() = %v, want %v", got, test.want)
			}
		})
	}
}
