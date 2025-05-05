package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
)

func TestAddHistoryNote(t *testing.T) {
	tests := map[string]struct {
		prevHistory datatypes.JSON
		what        string
		attrs       map[string]any
		expected    string
	}{
		"Add initial note": {
			prevHistory: nil,
			what:        "Initial note",
			expected:    `{"notes":[{"what":"Initial note","when":"2023-07-25T15:04:05Z"}]}`,
		},
		"Add initial note with attrs": {
			prevHistory: nil,
			what:        "Initial note",
			attrs:       map[string]any{"key": "value"},
			expected:    `{"notes":[{"what":"Initial note","key":"value","when":"2023-07-25T15:04:05Z"}]}`,
		},
		"Add subsequent note": {
			prevHistory: []byte(`{"notes":[{"what":"Initial note","when":"2023-07-24T15:04:05Z"}]}`),
			what:        "Subsequent note",
			expected:    `{"notes":[{"what":"Initial note","when":"2023-07-24T15:04:05Z"},{"what":"Subsequent note","when":"2023-07-25T15:04:05Z"}]}`,
		},
		"Add a note on malformed history": {
			prevHistory: []byte(`NOT A JSON AT ALL`),
			what:        "Subsequent note",
			expected:    `{"notes":[{"what":"Subsequent note","when":"2023-07-25T15:04:05Z"}]}`,
		},
		"Add a note on JSON without 'notes' key": {
			prevHistory: []byte(`{"foo":"bar"}`),
			what:        "Subsequent note",
			expected:    `{"notes":[{"what":"Subsequent note","when":"2023-07-25T15:04:05Z"}]}`,
		},
		"Add a note on JSON with empty 'notes' array": {
			prevHistory: []byte(`{"notes":[]}`),
			what:        "Subsequent note",
			expected:    `{"notes":[{"what":"Subsequent note","when":"2023-07-25T15:04:05Z"}]}`,
		},
		"Add a note on JSON with 'notes' as a string": {
			prevHistory: []byte(`{"notes":"foo"}`),
			what:        "Subsequent note",
			expected:    `{"notes":[{"what":"Subsequent note","when":"2023-07-25T15:04:05Z"}]}`,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			model := ProvenTxReq{
				History: test.prevHistory,
			}

			// and:
			mockWhen := time.Date(2023, 7, 25, 15, 4, 5, 0, time.UTC)

			// when:
			model.AddNote(mockWhen, test.what, test.attrs)

			// then:
			require.JSONEq(t, test.expected, string(model.History))
		})
	}
}
