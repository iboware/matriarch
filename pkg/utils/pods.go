package utils

import v1 "k8s.io/api/core/v1"

// GetPodNames returns the pod names of the array of pods passed in
func GetPodNames(pods []v1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}
