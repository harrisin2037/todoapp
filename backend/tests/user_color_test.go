package tests

import (
	"testing"

	"github.com/harrisin2037/todoapp/internal/utils"
)

func TestGenerateColor(t *testing.T) {
	type args struct {
		id uint
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test GenerateColor",
			args: args{id: 1},
			want: "hsl(4, 70%, 80%)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.GenerateColor(tt.args.id); got != tt.want {
				t.Errorf("GenerateColor() = %v, want %v", got, tt.want)
			}
		})
	}
}
