package util

import (
	"testing"
	"time"
)

var z *struct{}

var testCase = []struct {
	Value   interface{}
	IsEmpty bool
}{
	{
		"", true,
	},
	{
		" ", true,
	},
	{
		" jakarta", false,
	},
	{
		int64(0), true,
	},
	{
		uint(0), true,
	},
	{
		float64(0), true,
	},
	{
		[]string{}, true,
	},
	{
		map[string]interface{}{}, true,
	},
	{
		false, true,
	},
	{
		nil, true,
	},
	{
		struct {
		}{}, false,
	},
	{
		new(struct{}), false,
	},
	{
		time.Now(), false,
	},
	{
		time.Time{}, true,
	},
	{
		z, true,
	},
}

func TestIsEmptyValue(t *testing.T) {

	for _, c := range testCase {
		ie := IsEmptyValue(c.Value)
		if ie == c.IsEmpty {
			t.Logf("expected %v, got actual: %v", c.IsEmpty, ie)
			continue
		}

		t.Errorf("expected %v, got actual: %v", c.IsEmpty, ie)
	}
}

func TestValidJSON(t *testing.T) {
	tCase := []struct {
		STR      string
		Expected bool
	}{
		{
			STR:      `test`,
			Expected: false,
		},
		{
			STR:      `{"test": "test"}`,
			Expected: true,
		},
	}

	for _, c := range tCase {
		ie := ValidJSON([]byte(c.STR))
		if ie == c.Expected {
			t.Logf("expected %v, got actual: %v", c.Expected, ie)
			continue
		}

		t.Errorf("expected %v, got actual: %v", c.Expected, ie)
	}

}
