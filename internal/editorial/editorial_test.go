package editorial

import (
	"reflect"
	"testing"
)

func TestParseFeedback(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *Feedback
		wantErr bool
	}{
		{
			name: "well formed response",
			input: "<summary>A post about Go.</summary>\n" +
				"<tags>go, programming, system-design</tags>\n" +
				"<notes>\n- The intro is too long.\n- Consider splitting section 3.\n</notes>",
			want: &Feedback{
				Summary: "A post about Go.",
				Tags:    []string{"go", "programming", "system-design"},
				Notes:   []string{"The intro is too long.", "Consider splitting section 3."},
			},
		},
		{
			name:  "wrapped in a stray code fence",
			input: "```\n<summary>Short summary.</summary><tags>go</tags><notes></notes>\n```",
			want: &Feedback{
				Summary: "Short summary.",
				Tags:    []string{"go"},
				Notes:   nil,
			},
		},
		{
			name:  "notes contain quotes and code, no JSON escaping needed",
			input: `<summary>Fine.</summary><tags>go</tags><notes>- Don't use "self" as a receiver name, e.g. func (self *T) Foo().</notes>`,
			want: &Feedback{
				Summary: "Fine.",
				Tags:    []string{"go"},
				Notes:   []string{`Don't use "self" as a receiver name, e.g. func (self *T) Foo().`},
			},
		},
		{
			name:    "no recognizable tags at all",
			input:   "Sorry, I can't help with that.",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFeedback(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseFeedback() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFeedback() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
