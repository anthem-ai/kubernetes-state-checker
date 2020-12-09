package deployments

import (
	"reflect"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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

func Test_deploymentParse(t *testing.T) {
	type args struct {
		valuesYaml string
		v          *deploymentStruct
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
			if err := deploymentParse(tt.args.valuesYaml, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("deploymentParse() error = %v, wantErr %v", err, tt.wantErr)
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
			name: "Checking deployment envars (positive)",
			fields: fields{
				checkName: "check1",
				namespace: "ns1",
				// The spacing is real finicky.  yaml can't have tabs.  All spacing must be spaces
				valuesYaml: `---
values:
  # The service name to act on
  deploymentName: deployment1
  checksEnabled:
    containers:
    - name: container1
      env:
      - name: foo
        value: bar
      - name: foo2
        value: bar2
    - name: container2
      env:
      - name: foo
        value: bar`,
			},
			args: args{
				// Doc/example: https://gianarb.it/blog/unit-testing-kubernetes-client-in-go
				kubeClientSet: fake.NewSimpleClientset(&appsv1.DeploymentList{
					Items: []appsv1.Deployment{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name:        "deployment1",
								Namespace:   "ns1",
								Annotations: map[string]string{},
							},
							Spec: appsv1.DeploymentSpec{
								Template: corev1.PodTemplateSpec{
									Spec: corev1.PodSpec{
										Containers: []corev1.Container{
											{
												Name: "container1",
												Env: []corev1.EnvVar{
													{
														Name:  "foo",
														Value: "bar",
													},
													{
														Name:  "foo2",
														Value: "bar2",
													},
													{
														Name:  "foo3",
														Value: "bar3",
													},
												},
											},
											{
												Name: "container2",
												Env: []corev1.EnvVar{
													{
														Name:  "foo",
														Value: "bar",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				}),
			},
			want: Results{
				DidPass: true,
				Message: `* Found all envars in Deployment: deployment1 | container: container1
* Found all envars in Deployment: deployment1 | container: container2
* Found the correct number of containers in this deployment
`,
			},
		},
		{
			name: "Checking the number of pods in a deployment (positive)",
			fields: fields{
				checkName: "check1",
				namespace: "ns1",
				// The spacing is real finicky.  yaml can't have tabs.  All spacing must be spaces
				valuesYaml: `---
values:
  # The service name to act on
  deploymentName: check-number-of-pods
  checksEnabled:
    containers:
    - name: pod-container1
    - name: pod-container2`,
			},
			args: args{
				// Doc/example: https://gianarb.it/blog/unit-testing-kubernetes-client-in-go
				kubeClientSet: fake.NewSimpleClientset(&appsv1.DeploymentList{
					Items: []appsv1.Deployment{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name:        "check-number-of-pods",
								Namespace:   "ns1",
								Annotations: map[string]string{},
							},
							Spec: appsv1.DeploymentSpec{
								Template: corev1.PodTemplateSpec{
									Spec: corev1.PodSpec{
										Containers: []corev1.Container{
											{
												Name: "pod-container1",
												Env: []corev1.EnvVar{
													{
														Name:  "pod",
														Value: "bar",
													},
													{
														Name:  "pod",
														Value: "bar2",
													},
												},
											},
											{
												Name: "pod-container2",
												Env: []corev1.EnvVar{
													{
														Name:  "pod",
														Value: "bar",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				}),
			},
			want: Results{
				DidPass: true,
				Message: `* Found the correct number of containers in this deployment
`,
			},
		},
		{
			name: "Checking the number of pods in a deployment 2 (positive)",
			fields: fields{
				checkName: "check1",
				namespace: "ns1",
				// The spacing is real finicky.  yaml can't have tabs.  All spacing must be spaces
				valuesYaml: `---
values:
  # The service name to act on
  deploymentName: check-number-of-pods
  checksEnabled:
    containers:
    - name: pod-container1
    - name: pod-container2`,
			},
			args: args{
				// Doc/example: https://gianarb.it/blog/unit-testing-kubernetes-client-in-go
				kubeClientSet: fake.NewSimpleClientset(&appsv1.DeploymentList{
					Items: []appsv1.Deployment{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name:        "check-number-of-pods",
								Namespace:   "ns1",
								Annotations: map[string]string{},
							},
							Spec: appsv1.DeploymentSpec{
								Template: corev1.PodTemplateSpec{
									Spec: corev1.PodSpec{
										Containers: []corev1.Container{
											{
												Name: "pod-container1",
												Env: []corev1.EnvVar{
													{
														Name:  "pod",
														Value: "bar",
													},
													{
														Name:  "pod",
														Value: "bar2",
													},
												},
											},
											{
												Name: "pod-container2",
												Env: []corev1.EnvVar{
													{
														Name:  "pod",
														Value: "bar",
													},
												},
											},
											// THis test has more pods in the deployments than what the user is looking for
											{
												Name: "pod-container3",
												Env: []corev1.EnvVar{
													{
														Name:  "pod",
														Value: "bar",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				}),
			},
			want: Results{
				DidPass: true,
				Message: `* Found the correct number of containers in this deployment
`,
			},
		},
		{
			name: "Checking the number of pods in a deployment 1 (negative)",
			fields: fields{
				checkName: "check1",
				namespace: "ns1",
				// The spacing is real finicky.  yaml can't have tabs.  All spacing must be spaces
				valuesYaml: `---
values:
  # The service name to act on
  deploymentName: check-number-of-pods
  checksEnabled:
    containers:
    - name: pod-container1
    - name: pod-container2`,
			},
			args: args{
				// Doc/example: https://gianarb.it/blog/unit-testing-kubernetes-client-in-go
				kubeClientSet: fake.NewSimpleClientset(&appsv1.DeploymentList{
					Items: []appsv1.Deployment{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name:        "check-number-of-pods",
								Namespace:   "ns1",
								Annotations: map[string]string{},
							},
							Spec: appsv1.DeploymentSpec{
								Template: corev1.PodTemplateSpec{
									Spec: corev1.PodSpec{
										Containers: []corev1.Container{
											{
												Name: "pod-container1",
												Env: []corev1.EnvVar{
													{
														Name:  "pod",
														Value: "bar",
													},
													{
														Name:  "pod",
														Value: "bar2",
													},
												},
											},
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
				Message: "",
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
