package encoder

import "testing"

func TestBase62Encoder_Encode(t *testing.T) {
	encoder := NewBase62Encoder()

	tests := []struct {
		name string
		num  int64
		want string
	}{
		{
			name: "zero",
			num:  0,
			want: "0",
		},
		{
			name: "one",
			num:  1,
			want: "1",
		},
		{
			name: "nine",
			num:  9,
			want: "9",
		},
		{
			name: "ten",
			num:  10,
			want: "a",
		},
		{
			name: "thirty five",
			num:  35,
			want: "z",
		},
		{
			name: "thirty six",
			num:  36,
			want: "A",
		},
		{
			name: "sixty one",
			num:  61,
			want: "Z",
		},
		{
			name: "sixty two",
			num:  62,
			want: "10",
		},
		{
			name: "sixty three",
			num:  63,
			want: "11",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := encoder.Encode(tt.num)

			if got != tt.want {
				t.Fatalf("Encode(%d) = %s, want %s", tt.num, got, tt.want)
			}
		})
	}
}
