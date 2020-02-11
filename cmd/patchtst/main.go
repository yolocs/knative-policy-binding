package main

import (
	"encoding/json"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"

	jsonpatch "gomodules.xyz/jsonpatch/v2"
)

func main() {
	p1 := &duckv1.WithPod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "foo-ns",
		},
		Spec: duckv1.WithPodSpec{
			Template: duckv1.PodSpecable{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "curl-1",
							Image:   "governmentpaas/curl-ssl",
							Command: []string{"/bin/sleep", "3650d"},
							Env: []corev1.EnvVar{
								{
									Name:  "KEY1",
									Value: "VAL1",
								},
							},
						},
						{
							Name:    "curl-2",
							Image:   "governmentpaas/curl-ssl",
							Command: []string{"/bin/sleep", "3650d"},
							Env: []corev1.EnvVar{
								{
									Name:  "KEY2",
									Value: "VAL2",
								},
							},
						},
					},
				},
			},
		},
	}

	// p15 := &duckv1.WithPod{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name:      "foo",
	// 		Namespace: "foo-ns",
	// 	},
	// 	Spec: duckv1.WithPodSpec{
	// 		Template: duckv1.PodSpecable{
	// 			Spec: corev1.PodSpec{
	// 				Containers: []corev1.Container{
	// 					{
	// 						Name:    "curl-1",
	// 						Image:   "governmentpaas/curl-ssl",
	// 						Command: []string{"/bin/sleep", "3650d"},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	p2 := &duckv1.WithPod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "foo-ns",
		},
		Spec: duckv1.WithPodSpec{
			Template: duckv1.PodSpecable{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "curl-1",
							Image:   "governmentpaas/curl-ssl",
							Command: []string{"/bin/sleep", "3650d"},
							Env: []corev1.EnvVar{
								{
									Name:  "KEY3",
									Value: "VAL3",
								},
							},
						},
					},
				},
			},
		},
	}

	b1, _ := json.Marshal(p1)
	b2, _ := json.Marshal(p2)

	op, _ := jsonpatch.CreatePatch(b1, b2)
	for _, o := range op {
		log.Println(o.Json())
	}
}
