package models

import (
	"testing"
	"time"

	"github.com/go-softwarelab/common/pkg/to"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
)

var mockWhen = time.Date(2023, 7, 25, 15, 4, 5, 0, time.UTC)

func TestAddHistoryNote(t *testing.T) {
	tests := map[string]struct {
		prevHistory datatypes.JSON
		what        string
		when        *time.Time
		attrs       map[string]any
		expected    string
	}{
		"Add initial note": {
			prevHistory: nil,
			what:        "Initial note",
			expected:    `{"notes":[{"what":"Initial note"}]}`,
		},
		"Add initial note with when": {
			prevHistory: nil,
			what:        "Initial note",
			when:        to.Ptr(mockWhen),
			expected:    `{"notes":[{"what":"Initial note","when":"2023-07-25T15:04:05Z"}]}`,
		},
		"Add initial note with attrs": {
			prevHistory: nil,
			what:        "Initial note",
			attrs:       map[string]any{"key": "value"},
			expected:    `{"notes":[{"what":"Initial note","key":"value"}]}`,
		},
		"Add initial note with when and attrs": {
			prevHistory: nil,
			what:        "Initial note",
			when:        to.Ptr(mockWhen),
			attrs:       map[string]any{"key": "value"},
			expected:    `{"notes":[{"what":"Initial note","when":"2023-07-25T15:04:05Z","key":"value"}]}`,
		},
		"Add subsequent note": {
			prevHistory: []byte(`{"notes":[{"what":"Initial note"}]}`),
			what:        "Subsequent note",
			expected:    `{"notes":[{"what":"Initial note"},{"what":"Subsequent note"}]}`,
		},
		"Add a note on malformed history": {
			prevHistory: []byte(`NOT A JSON AT ALL`),
			what:        "Subsequent note",
			expected:    `{"notes":[{"what":"Subsequent note"}]}`,
		},
		"Add a note on JSON without 'notes' key": {
			prevHistory: []byte(`{"foo":"bar"}`),
			what:        "Subsequent note",
			expected:    `{"notes":[{"what":"Subsequent note"}]}`,
		},
		"Add a note on JSON with empty 'notes' array": {
			prevHistory: []byte(`{"notes":[]}`),
			what:        "Subsequent note",
			expected:    `{"notes":[{"what":"Subsequent note"}]}`,
		},
		"Add a note on JSON with 'notes' as a string": {
			prevHistory: []byte(`{"notes":"foo"}`),
			what:        "Subsequent note",
			expected:    `{"notes":[{"what":"Subsequent note"}]}`,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			model := ProvenTxReq{
				History: test.prevHistory,
			}

			// when:
			model.AddNote(test.when, test.what, test.attrs)

			// then:
			require.JSONEq(t, test.expected, string(model.History))
		})
	}
}
