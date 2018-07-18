package stub

import (
	"context"
	"github.com/mvazquezc/python-api-hw/pkg/apis/ostack/v1alpha1"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
        appsv1 "k8s.io/api/apps/v1"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
        logrus.Infof("Inside Handler")
	switch o := event.Object.(type) {
	case *v1alpha1.PythonAPIHw:
                helloApiWorld := o
                dep := deploymentForHelloApi(helloApiWorld)
                err := sdk.Create(dep)
                logrus.Infof("Inside switch for my object")
		if err != nil && !errors.IsAlreadyExists(err) {
			logrus.Errorf("Failed to create deployment : %v", err)
			return err
		}
                svc := serviceForHelloApi(helloApiWorld)
                err = sdk.Create(svc)
                if err != nil && !errors.IsAlreadyExists(err) {
                        logrus.Errorf("Failed to create service : %v", err)
                        return err
                }
                // Ensure the deployment size is the same as the spec
                err = sdk.Get(dep)
                if err != nil {
                        logrus.Errorf("Failed to get deployment : %v", err)
                        return err
                }
                size := helloApiWorld.Spec.Size
                logrus.Infof("Size is set to %d, current replias %d", size, *dep.Spec.Replicas)
                if *dep.Spec.Replicas != size {
                        logrus.Infof("Need to update replicas from %d to %d", *dep.Spec.Replicas, size)
                        dep.Spec.Replicas = &size
                        err = sdk.Update(dep)
                        if err != nil {
                                logrus.Errorf("Failed to update deployment : %v", err)
                                return err
                        }
                }
	}
	return nil
}

// serviceForHelloApi returns a Service Object
func serviceForHelloApi(h *v1alpha1.PythonAPIHw) *corev1.Service {
        labels := map[string]string{
                 "app": "api-hello-world",
        }
        svc := &corev1.Service{
                TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
                ObjectMeta: metav1.ObjectMeta{
			Name:      h.Name,
			Namespace: h.Namespace,
		},
                Spec: corev1.ServiceSpec{
                        Type:     corev1.ServiceTypeLoadBalancer,
                        Selector: labels,
                        Ports: []corev1.ServicePort{
                                {
                                        Name: "http",
                                        Port: 5000,
                                },
                        },
                },
        }
        return svc
}

// deploymentForHelloApi returns a HelloApi Deployment Object
func deploymentForHelloApi(h *v1alpha1.PythonAPIHw) *appsv1.Deployment {
        labels := map[string]string{
                 "app": "api-hello-world",
        }
        replicas := h.Spec.Size
        dep := &appsv1.Deployment{
                TypeMeta: metav1.TypeMeta{
                        APIVersion: "apps/v1",
                        Kind:       "Deployment",
                },
                ObjectMeta: metav1.ObjectMeta{
                        Name:      h.Name,
                        Namespace: h.Namespace,
                },
                Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:   "quay.io/mavazque/hello-api",
						Name:    "api-hello-world",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 5000,
							Name:          "api-hello-world",
						}},
					}},
				},
			},
		},
	}
        return dep
}

