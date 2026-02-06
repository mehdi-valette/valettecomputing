package blog

import "testing"

func TestMakeSlug(t *testing.T) {
	type data struct {
		input  string
		output string
	}

	testData := []data{
		{"The event loop", "the-event-loop"},
		{"L'avant", "lavant"},
		{"L'Arc-En-Ciel", "larc-en-ciel"},
		{"ÿŷ", "y0"},
		{"Why? Why not!", "why-why-not"},
		{"L'Étrange noël de monsieur Jack!", "letrange-noel-de-monsieur-jack"},
		{"Ça c'est pas bien", "ca-cest-pas-bien"},
	}

	for _, test := range testData {
		result := makeSlug(test.input)

		if test.output != result {
			t.Errorf("expected \"%s\" to become \"%s\", got \"%s\"", test.input, test.output, result)
		}
	}
}
