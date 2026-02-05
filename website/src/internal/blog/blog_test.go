package blog

import "testing"

func TestMakeSlug(t *testing.T) {
	type data struct {
		input  string
		output string
	}

	testData := []data{
		{"The event loop", "the-event-loop"},
		{"L'école", "l00cole"},
		{"L'Arc-En-Ciel", "l0arc-en-ciel"},
	}

	for _, test := range testData {
		if test.output != makeSlug(test.input) {
			t.Errorf("expected \"%s\" to become \"%s\"", test.input, test.output)
		}
	}
}
