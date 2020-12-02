package services

import (
	"testing"
)

func Test_serviceParse(t *testing.T) {
	type args struct {
		valuesYaml string
		v          *serviceStruct
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1 - success",
			args: args{
				valuesYaml: "values:\n  serviceName: service1\n  port: 20014",
				v:          &serviceStruct{},
			},
			wantErr: false,
		},
		{
			name: "test2 - no serviceName param present",
			args: args{
				valuesYaml: "values:\n  port: 20014",
				v:          &serviceStruct{},
			},
			wantErr: true,
		},
		{
			name: "test3 - no port param present",
			args: args{
				valuesYaml: "values:\n  serviceName: service1",
				v:          &serviceStruct{},
			},
			wantErr: true,
		},
		{
			name: "test4 - invalid port range high",
			args: args{
				valuesYaml: "values:\n  serviceName: service1\n  port: 65354",
				v:          &serviceStruct{},
			},
			wantErr: true,
		},
		{
			name: "test4 - invalid port range high",
			args: args{
				valuesYaml: "values:\n  serviceName: service1\n  port: -1",
				v:          &serviceStruct{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := serviceParse(tt.args.valuesYaml, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("serviceParse() name = %v, error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}
