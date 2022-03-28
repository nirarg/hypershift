package kubevirt

import (
	"github.com/openshift/hypershift/support/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	CloudConfigKey = "cloud.conf"
	Provider       = "azure"
)

func ReconcileCloudProviderServiceAccount(sa *corev1.ServiceAccount) error {
	// Nothing to do
	return nil
}

func ReconcileCloudControllerManagerDaemonSet(ds *appsv1.DaemonSet) error {
	ds.Labels = cloudControllerManagerLabels()
	ds.Spec = appsv1.DaemonSetSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: cloudControllerManagerLabels(),
		},
		UpdateStrategy: appsv1.DaemonSetUpdateStrategy{
			Type: appsv1.RollingUpdateDaemonSetStrategyType,
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: cloudControllerManagerLabels(),
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					util.BuildContainer(cloudControllerManagerContainer(), buildCloudControllerManagerContainer()),
				},
				Volumes: []corev1.Volume{
					{
						Name: "k8s-certs",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/etc/kubernetes/pki",
								Type: convertToPointer(corev1.HostPathDirectoryOrCreate),
							},
						},
					},
					{
						Name: "ca-certs",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/etc/ssl/certs",
								Type: convertToPointer(corev1.HostPathDirectoryOrCreate),
							},
						},
					},
					{
						Name: "kubeconfig",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/etc/kubernetes/controller-manager.conf",
								Type: convertToPointer(corev1.HostPathFileOrCreate),
							},
						},
					},
					{
						Name: "cloud-config",
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: "cloud-config",
							},
						},
					},
				},
			},
		},
	}
	return nil
}

func convertToPointer(src corev1.HostPathType) *corev1.HostPathType {
	return &src
}

func buildCloudControllerManagerContainer() func(c *corev1.Container) {
	return func(c *corev1.Container) {
		c.Args = []string{
			"--cloud-provider=kubevirt",
			"--cloud-config=/etc/cloud/cloud-config",
			"--use-service-account-credentials=true",
			"--kubeconfig=/etc/kubernetes/controller-manager.conf",
		}
		c.Image = "docker.io/dgonzalez/kubevirt-cloud-controller-manager:v0.0.7"
		c.Resources = corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceCPU: resource.MustParse("100m"),
			},
		}
		c.VolumeMounts = []corev1.VolumeMount{
			{
				Name:      "k8s-certs",
				MountPath: "/etc/kubernetes/pki",
				ReadOnly:  true,
			},
			{
				Name:      "ca-certs",
				MountPath: "/etc/ssl/certs",
				ReadOnly:  true,
			},
			{
				Name:      "kubeconfig",
				MountPath: "/etc/kubernetes/controller-manager.conf",
				ReadOnly:  true,
			},
			{
				Name:      "cloud-config",
				MountPath: "/etc/cloud",
				ReadOnly:  true,
			},
		}
	}
}

func cloudControllerManagerContainer() *corev1.Container {
	return &corev1.Container{
		Name: "kubevirt-cloud-controller-manager",
	}
}

func cloudControllerManagerLabels() map[string]string {
	return map[string]string{
		"k8s-app": "kubevirt-cloud-controller-manager",
	}
}
