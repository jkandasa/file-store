package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToDuration(t *testing.T) {
	tests := []struct {
		name            string
		duration        string
		defaultDuration time.Duration
		expected        time.Duration
	}{
		{name: "valid_duration_1m", duration: "1m", defaultDuration: time.Minute * 2, expected: time.Minute * 1},
		{name: "valid_duration_3h30m5s", duration: "3h30m5s", defaultDuration: time.Minute * 2, expected: ((time.Hour * 3) + (time.Minute * 30) + (time.Second * 5))},
		{name: "invalid_duration_1m", duration: "1minute", defaultDuration: time.Minute * 2, expected: time.Minute * 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, ToDuration(test.duration, test.defaultDuration))
		})
	}
}

func TestStructToMap(t *testing.T) {
	sample := struct {
		Name  string
		Value int64
	}{
		Name:  "hi",
		Value: 43,
	}
	got := StructToMap(sample)
	expected := map[string]interface{}{"Name": "hi", "Value": int64(43)}
	assert.Equal(t, expected, got)
}
