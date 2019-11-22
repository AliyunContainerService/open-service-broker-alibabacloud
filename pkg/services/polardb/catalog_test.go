package polardb

import (
	"testing"
)

func Test_dealServiceInstanceMetadataValueString(t *testing.T) {
	type args struct {
		parameter map[string]interface{}
		paramKey  string
		metadata  *ServiceInstanceMetadata
	}

	metadataNew := &ServiceInstanceMetadata{}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test0",
			args: args{
				parameter: map[string]interface{}{
					"Engine": "test",
				},
				metadata: metadataNew,
				paramKey: "Engine",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := dealServiceInstanceMetadataValueString(tt.args.parameter, tt.args.paramKey, tt.args.metadata); (err != nil) != tt.wantErr {
				t.Errorf("dealServiceInstanceMetadataValueString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
