package pods

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

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

func Test_podParse(t *testing.T) {
	type args struct {
		valuesYaml string
		v          *podStruct
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := podParse(tt.args.valuesYaml, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("podParse() error = %v, wantErr %v", err, tt.wantErr)
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
			name: "Checking pod state - 2 container (positive)",
			fields: fields{
				checkName: "check1",
				namespace: "ns1",
				valuesYaml: `---
values:
  checksEnabled:
    state:
    - podName: pod-1-
      desiredState: Running`,
			},
			args: args{
				// Doc/example: https://gianarb.it/blog/unit-testing-kubernetes-client-in-go
				kubeClientSet: fake.NewSimpleClientset(&v1.PodList{
					Items: []v1.Pod{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name:        "pod-1-123abc-abc",
								Namespace:   "ns1",
								Annotations: map[string]string{},
							},
							Status: v1.PodStatus{
								Phase: "Running",
							},
						},
						{
							ObjectMeta: metav1.ObjectMeta{
								Name:        "pod-2222222-123abc-abc",
								Namespace:   "ns1",
								Annotations: map[string]string{},
							},
							Status: v1.PodStatus{
								Phase: "Running",
							},
						},
					},
				}),
			},
			want: Results{
				DidPass: true,
				Message: "* Pod pod-1-123abc-abc is in Running state\n",
			},
		},
		{
			name: "Checking pod state - 1 container (positive)",
			fields: fields{
				checkName: "check1",
				namespace: "ns1",
				valuesYaml: `---
values:
  checksEnabled:
    state:
    - podName: pod-1-
      desiredState: Running`,
			},
			args: args{
				// Doc/example: https://gianarb.it/blog/unit-testing-kubernetes-client-in-go
				kubeClientSet: fake.NewSimpleClientset(&v1.PodList{
					Items: []v1.Pod{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name:        "pod-1-123abc-abc",
								Namespace:   "ns1",
								Annotations: map[string]string{},
							},
							Status: v1.PodStatus{
								Phase: "Running",
							},
						},
					},
				}),
			},
			want: Results{
				DidPass: true,
				Message: "* Pod pod-1-123abc-abc is in Running state\n",
			},
		},
		{
			name: "Checking pod state - no containers found (negative)",
			fields: fields{
				checkName: "check1",
				namespace: "ns1",
				valuesYaml: `---
values:
  checksEnabled:
    state:
    - podName: pod-1-
      desiredState: Running`,
			},
			args: args{
				// Doc/example: https://gianarb.it/blog/unit-testing-kubernetes-client-in-go
				kubeClientSet: fake.NewSimpleClientset(&v1.PodList{
					Items: []v1.Pod{},
				}),
			},
			want: Results{
				DidPass: false,
				Message: "* Did not find pod: pod-1-\n",
			},
		},
		{
			name: "Checking pod state - found one container but not another (negative)",
			fields: fields{
				checkName: "check1",
				namespace: "ns1",
				valuesYaml: `---
values:
  checksEnabled:
    state:
    - podName: pod-1-
      desiredState: Running
    - podName: pod-2-
      desiredState: Running`,
			},
			args: args{
				// Doc/example: https://gianarb.it/blog/unit-testing-kubernetes-client-in-go
				kubeClientSet: fake.NewSimpleClientset(&v1.PodList{
					Items: []v1.Pod{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name:        "pod-1-123abc-abc",
								Namespace:   "ns1",
								Annotations: map[string]string{},
							},
							Status: v1.PodStatus{
								Phase: "Running",
							},
						},
					},
				}),
			},
			want: Results{
				DidPass: false,
				Message: "* Did not find pod: pod-1-\n",
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
				t.Errorf("inputs.GeneralCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}
