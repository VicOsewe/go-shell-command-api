package usecases_test

import (
	"testing"

	"github.com/VicOsewe/go-shell-command-api/usecases"
)

func TestExecuteCommand(t *testing.T) {
	type args struct {
		command string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				command: "pwd",
			},
			wantErr: false,
		},
		{
			name: "sad case - empty command",
			args: args{
				command: "",
			},
			wantErr: true,
		},
		{
			name: "sad case - invalid command",
			args: args{
				command: "pw",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := usecases.ExecuteCommand(tt.args.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) == 0 {
				t.Errorf("ExecuteCommand() = expected a command response.")
			}
		})
	}
}
