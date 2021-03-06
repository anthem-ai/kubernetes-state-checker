package services

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
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

func TestNew(t *testing.T) {
	type args struct {
		valuesYaml string
		checkName  string
		namespace  string
	}
	tests := []struct {
		name string
		args args
		want inputs
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.valuesYaml, tt.args.checkName, tt.args.namespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inputs_GeneralCheck(t *testing.T) {
	type fields struct {
		valuesYaml string
		checkName  string
		namespace  string
	}
	type args struct {
		kubeClientSet kubernetes.Interface
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Results
	}{
		// TODO: Add test cases.
		{
			name: "Checking clusterIP (positive)",
			fields: fields{
				checkName:  "check1",
				namespace:  "ns1",
				valuesYaml: "values:\n  serviceName: service1\n  port: 20014\n  checksEnabled:\n    clusterIP: true",
			},
			args: args{
				// Doc/example: https://gianarb.it/blog/unit-testing-kubernetes-client-in-go
				kubeClientSet: fake.NewSimpleClientset(&v1.ServiceList{
					Items: []v1.Service{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name:        "service1",
								Namespace:   "ns1",
								Annotations: map[string]string{},
							},
							Spec: v1.ServiceSpec{
								ClusterIP: "1.1.1.1",
							},
						},
					},
				}),
			},
			want: Results{
				DidPass: true,
				Message: "* ClusterIP Found\n",
			},
		},
		{
			name: "Checking clusterIP (negative)",
			fields: fields{
				checkName:  "check1",
				namespace:  "ns1",
				valuesYaml: "values:\n  serviceName: service1\n  port: 20014\n  checksEnabled:\n    clusterIP: true",
			},
			args: args{
				// Doc/example: https://gianarb.it/blog/unit-testing-kubernetes-client-in-go
				kubeClientSet: fake.NewSimpleClientset(&v1.ServiceList{
					Items: []v1.Service{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name:        "service1",
								Namespace:   "ns1",
								Annotations: map[string]string{},
							},
							Spec: v1.ServiceSpec{
								// No ClusterIP and this will fail the test
								// ClusterIP: "1.1.1.1",
							},
						},
					},
				}),
			},
			want: Results{
				DidPass: false,
				Message: "* No ClusterIP Found\n",
			},
		},
		// TODO: Fix this test: The endpoint k8s fake info doesnt show up.  The endpoint data is in the second k8s call.  Not sure why
		// {
		// 	name: "Checking Endpoints (positive)",
		// 	fields: fields{
		// 		checkName:  "check2",
		// 		namespace:  "ns1",
		// 		valuesYaml: "values:\n  serviceName: service1\n  port: 20014\n  checksEnabled:\n    endpoints: true",
		// 	},
		// 	args: args{
		// 		// Doc/example: https://gianarb.it/blog/unit-testing-kubernetes-client-in-go
		// 		kubeClientSet: fake.NewSimpleClientset(&v1.EndpointsList{
		// 			Items: []v1.Endpoints{
		// 				{
		// 					ObjectMeta: metav1.ObjectMeta{
		// 						Name: "service1",
		// 					},
		// 					Subsets: []v1.EndpointSubset{
		// 						{
		// 							Addresses: []v1.EndpointAddress{
		// 								{
		// 									IP: "1.1.1.1",
		// 								},
		// 							},
		// 						},
		// 					},
		// 				},
		// 			},
		// 		}, &v1.ServiceList{
		// 			Items: []v1.Service{
		// 				{
		// 					ObjectMeta: metav1.ObjectMeta{
		// 						Name:        "service1",
		// 						Namespace:   "ns1",
		// 						Annotations: map[string]string{},
		// 					},
		// 					Spec: v1.ServiceSpec{
		// 						ClusterIP: "1.1.1.1",
		// 					},
		// 				},
		// 			},
		// 		}),
		// 	},
		// 	want: Results{
		// 		DidPass: true,
		// 		Message: "* ClusterIP Found\n",
		// 	},
		// },
		{
			name: "Checking ports (positive)",
			fields: fields{
				checkName:  "check3",
				namespace:  "ns1",
				valuesYaml: "values:\n  serviceName: service1\n  port: 20014\n  checksEnabled:\n    ports: true",
			},
			args: args{
				// Doc/example: https://gianarb.it/blog/unit-testing-kubernetes-client-in-go
				kubeClientSet: fake.NewSimpleClientset(&v1.ServiceList{
					Items: []v1.Service{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name:        "service1",
								Namespace:   "ns1",
								Annotations: map[string]string{},
							},
							Spec: v1.ServiceSpec{
								ClusterIP: "1.1.1.1",
								Ports: []v1.ServicePort{
									{
										Port: 20014,
									},
								},
							},
						},
					},
				}),
			},
			want: Results{
				DidPass: true,
				Message: "* Port found: 20014\n",
			},
		},
		{
			name: "Checking ports (negative)",
			fields: fields{
				checkName:  "check4",
				namespace:  "ns1",
				valuesYaml: "values:\n  serviceName: service1\n  port: 20014\n  checksEnabled:\n    ports: true",
			},
			args: args{
				// Doc/example: https://gianarb.it/blog/unit-testing-kubernetes-client-in-go
				kubeClientSet: fake.NewSimpleClientset(&v1.ServiceList{
					Items: []v1.Service{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name:        "service1",
								Namespace:   "ns1",
								Annotations: map[string]string{},
							},
							Spec: v1.ServiceSpec{
								ClusterIP: "1.1.1.1",
								Ports: []v1.ServicePort{
									{
										// This is the incorrect port that will fail the test
										Port: 20013,
									},
								},
							},
						},
					},
				}),
			},
			want: Results{
				DidPass: false,
				Message: "* Port NOT found: 20014\n",
			},
		},
		{
			name: "Checking ports (negative) - port doesnt exist",
			fields: fields{
				checkName:  "check3",
				namespace:  "ns1",
				valuesYaml: "values:\n  serviceName: service1\n  port: 20014\n  checksEnabled:\n      clusterIP: true\n      endpoints: true\n      hostPort: false\n      ports: true",
			},
			args: args{
				// Doc/example: https://gianarb.it/blog/unit-testing-kubernetes-client-in-go
				kubeClientSet: fake.NewSimpleClientset(&v1.ServiceList{
					Items: []v1.Service{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name:        "service1",
								Namespace:   "ns1",
								Annotations: map[string]string{},
							},
							Spec: v1.ServiceSpec{
								ClusterIP: "1.1.1.1",
								Ports: []v1.ServicePort{
									{
										Port: 20015,
									},
								},
							},
						},
					},
				}, &v1.EndpointsList{
					Items: []v1.Endpoints{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name:        "service1",
								Namespace:   "ns1",
								Annotations: map[string]string{},
							},
							Subsets: []v1.EndpointSubset{
								{
									Addresses: []v1.EndpointAddress{
										{
											IP: "2.2.2.2",
										},
									},
								},
							},
						},
					},
				}),
			},
			want: Results{
				DidPass: false,
				Message: "* ClusterIP Found\n* Endpoint found: 2.2.2.2\n* Port NOT found: 20014\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := inputs{
				valuesYaml: tt.fields.valuesYaml,
				checkName:  tt.fields.checkName,
				namespace:  tt.fields.namespace,
			}
			if got := i.GeneralCheck(tt.args.kubeClientSet); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%v: inputs.GeneralCheck() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
