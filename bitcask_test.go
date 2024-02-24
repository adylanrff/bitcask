package bitcask

import (
	"os"
	"reflect"
	"testing"

	"github.com/adylanrff/bitcask/pkg/types"
)

func Test_Bitcask(t *testing.T) {
	type command struct {
		cmd   string
		key   string
		value string

		want    interface{}
		wantErr bool
	}

	type args struct {
		commands []command
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty get should return empty",
			args: args{
				commands: []command{
					{
						cmd:     "get",
						key:     "test",
						want:    types.Value(nil),
						wantErr: false,
					},
				},
			},
		},
		{
			name: "set one should return the set value",
			args: args{
				commands: []command{
					{
						cmd:     "set",
						key:     "test",
						value:   "value",
						wantErr: false,
					},
					{
						cmd:     "get",
						key:     "test",
						want:    types.Value("value"),
						wantErr: false,
					},
				},
			},
		},
		{
			name: "set multiple should return the correct set value",
			args: args{
				commands: []command{
					{
						cmd:     "set",
						key:     "test",
						value:   "value",
						wantErr: false,
					},
					{
						cmd:     "set",
						key:     "test2",
						value:   "value2",
						wantErr: false,
					},
					{
						cmd:     "get",
						key:     "test",
						want:    types.Value("value"),
						wantErr: false,
					},
					{
						cmd:     "get",
						key:     "test2",
						want:    types.Value("value2"),
						wantErr: false,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.RemoveAll("test")
			db, err := Open("test")
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			var outVal types.Value

			for _, command := range tt.args.commands {
				switch command.cmd {
				case "get":
					outVal, err = db.Get(types.Key(command.key))
					if (err != nil) != tt.wantErr {
						t.Errorf("Get() error = %v, wantErr %v", err, command.wantErr)
						return
					}
					if !reflect.DeepEqual(outVal, command.want) {
						t.Errorf("Get() = %#v, want %#v", outVal, command.want)
					}
				case "set":
					err = db.Put(types.Key(command.key), types.Value(command.value))
					if (err != nil) != tt.wantErr {
						t.Errorf("set() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
				}
			}
		})
	}
}
