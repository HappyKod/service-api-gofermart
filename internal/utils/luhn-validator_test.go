package utils

import "testing"

func TestValid(t *testing.T) {
	type args struct {
		number int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "проверка валидного когда",
			args: args{number: 4561261212345467},
			want: true,
		},
		{
			name: "проверка не валидного когда",
			args: args{number: 4561261212345464},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidLuhn(tt.args.number); got != tt.want {
				t.Errorf("ValidLuhn() = %v, want %v", got, tt.want)
			}
		})
	}
}
