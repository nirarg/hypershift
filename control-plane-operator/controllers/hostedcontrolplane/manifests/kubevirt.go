package manifests

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func KubevirtCloudConfigSecret(ns string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cloud-config",
			Namespace: ns,
		},
	}
}

func KubevirtCloudProviderServiceAccount() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cloud-controller-manager",
			Namespace: "kube-system",
		},
	}
}
func KubevirtCloudProviderDaemonSet() *appsv1.DaemonSet {
	return &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kubevirt-cloud-controller-manager",
			Namespace: "kube-system",
		},
	}
}
