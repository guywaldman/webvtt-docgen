package vtt_parser

import (
	"strings"
	"testing"
)

func TestVttParser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       string
		expected    *VTTFile
		expectedErr bool
	}{
		{
			name: "basic",
			input: `WEBVTT

			1
			00:00:00.000 --> 00:00:05.000
			Speaker 1: This is the first cue.`,
			expected: &VTTFile{
				Cues: []*Cue{
					{
						Speaker: "Speaker 1",
						Start:   "00:00:00.000",
						End:     "00:00:05.000",
						Text:    "This is the first cue.",
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "cue consolidation",
			input: `WEBVTT

            1
            00:00:00.000 --> 00:00:0
            Speaker 1: This is the first cue.

            2
            00:00:05.000 --> 00:00:10.000
            Speaker 1: This is the second cue.

            3
            00:00:10.000 --> 00:00:15.000
            Speaker 2: This is the third cue.

			4
            00:00:17.000 --> 00:00:20.000
            Speaker 3: This is the third cue.
            `,
			expected: &VTTFile{
				Cues: []*Cue{
					{
						Speaker: "Speaker 1",
						Start:   "00:00:00.000",
						End:     "00:00:10.000",
						Text:    "This is the first cue.",
					},
					{
						Speaker: "Speaker 2",
						Start:   "00:00:10.000",
						End:     "00:00:15.000",
						Text:    "This is the third cue.",
					},
					{
						Speaker: "Speaker 3",
						Start:   "00:00:17.000",
						End:     "00:00:20.000",
						Text:    "This is the third cue.",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)

			actual, actualErr := ParseVTT(r)

			if (actualErr != nil) != tt.expectedErr {
				t.Errorf("ParseVTT(): actual error = %v, expected error = %v", actualErr, tt.expectedErr)
				return
			}

			for cue := range actual.Cues {
				if actual.Cues[cue].Speaker != tt.expected.Cues[cue].Speaker {
					t.Errorf("(Cue %d) Actual speaker = %s, expected speaker = %s", cue, actual.Cues[cue].Speaker, tt.expected.Cues[cue].Speaker)
				}
				if actual.Cues[cue].Start != tt.expected.Cues[cue].Start {
					t.Errorf("(Cue %d) Actual start = %s, expected start = %s", cue, actual.Cues[cue].Start, tt.expected.Cues[cue].Start)
				}
				if actual.Cues[cue].End != tt.expected.Cues[cue].End {
					t.Errorf("(Cue %d) Actual end = %s, expected end = %s", cue, actual.Cues[cue].End, tt.expected.Cues[cue].End)
				}
			}
		})
	}
}
