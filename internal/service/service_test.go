package service

import "testing"

func TestConvertText(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "morse to text",
			input: ".--. .-. .. .-- . -",
			want:  "ПРИВЕТ",
		},
		{
			name:  "uppercase text to morse",
			input: "ПРИВЕТ",
			want:  ".--. .-. .. .-- . -",
		},
		{
			name:  "lowercase text to morse",
			input: "Привет",
			want:  ".--. .-. .. .-- . -",
		},
		{
			name:  "digits to morse",
			input: "123",
			want:  ".---- ..--- ...--",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertText(tt.input); got != tt.want {
				t.Fatalf("ConvertText(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestConvertTextRoundTrip(t *testing.T) {
	const text = "ПРИВЕТ"

	if got := ConvertText(ConvertText(text)); got != text {
		t.Fatalf("ConvertText(ConvertText(%q)) = %q, want %q", text, got, text)
	}
}
