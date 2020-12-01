package ports

import (
	"testing"
)

func Test_doesPortExistParse(t *testing.T) {
	type args struct {
		valuesYaml string
		v          *doesPortExistStruct
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1 - success",
			args: args{
				valuesYaml: "values:\n  serviceName: healthapp-caregaps\n  port: 20014",
				v:          &doesPortExistStruct{},
			},
			wantErr: false,
		},
		{
			name: "test2 - no serviceName param present",
			args: args{
				valuesYaml: "values:\n  port: 20014",
				v:          &doesPortExistStruct{},
			},
			wantErr: true,
		},
		{
			name: "test3 - no port param present",
			args: args{
				valuesYaml: "values:\n  serviceName: healthapp-caregaps",
				v:          &doesPortExistStruct{},
			},
			wantErr: true,
		},
		{
			name: "test4 - invalid port range high",
			args: args{
				valuesYaml: "values:\n  serviceName: healthapp-caregaps\n  port: 65354",
				v:          &doesPortExistStruct{},
			},
			wantErr: true,
		},
		{
			name: "test4 - invalid port range high",
			args: args{
				valuesYaml: "values:\n  serviceName: healthapp-caregaps\n  port: -1",
				v:          &doesPortExistStruct{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := doesPortExistParse(tt.args.valuesYaml, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("doesPortExistParse() name = %v, error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}
