package minimal

import (
	"reflect"
	"testing"

	"github.com/adylanrff/bitcask/pkg/types"
)

func TestHandler_Get(t *testing.T) {
	type args struct {
		key types.Key
	}
	tests := []struct {
		name    string
		args    args
		want    types.Value
		wantErr bool
	}{
		{
			name: "empty data should return empty value",
			args: args{key: "test"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := NewHandler("test", nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handler.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got, err := h.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handler.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
