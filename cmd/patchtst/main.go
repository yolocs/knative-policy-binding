package main

import "log"

// func main() {
// 	p := &duckv1.WithPod{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "foo",
// 			Namespace: "foo-ns",
// 		},
// 		Spec: duckv1.WithPodSpec{
// 			Template: duckv1.PodSpecable{
// 				Spec: corev1.PodSpec{
// 					Containers: []corev1.Container{
// 						{
// 							Name:  "container-0",
// 							Image: "image0",
// 							Env: []corev1.EnvVar{
// 								{
// 									Name:  "K_POLICY_DECIDER_1",
// 									Value: "url",
// 								},
// 								{
// 									Name:  "OTHER",
// 									Value: "val",
// 								},
// 							},
// 						},
// 						{
// 							Name:  "container-1",
// 							Image: "image1",
// 							Env: []corev1.EnvVar{
// 								{
// 									Name:  "K_POLICY_DECIDER",
// 									Value: "url",
// 								},
// 								{
// 									Name:  "OTHER",
// 									Value: "val",
// 								},
// 							},
// 						},
// 					},
// 					Volumes: []corev1.Volume{
// 						{
// 							Name: "vol-0",
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	orig := p.DeepCopy()
// 	change(p)

// 	patchBytes, err := duck.CreateBytePatch(orig, p)
// 	if err != nil {
// 		log.Println(err)
// 		os.Exit(1)
// 	}

// 	log.Println(string(patchBytes))
// }

// func change(ps *duckv1.WithPod) {
// 	spec := ps.Spec.Template.Spec
// 	for j, ev := range spec.Volumes {
// 		if ev.Name == "vol-0" {
// 			spec.Volumes = append(spec.Volumes[:j], spec.Volumes[j+1:]...)
// 		}
// 	}

// 	// for i, c := range spec.Containers {
// 	// 	if c.Name == "container-0" {
// 	// 		spec.Containers = append(spec.Containers[:i], spec.Containers[i+1:]...)
// 	// 	}
// 	// }

// 	for i, c := range spec.Containers {
// 		for j, ev := range c.Env {
// 			if ev.Name == "K_POLICY_DECIDER" {
// 				spec.Containers[i].Env = append(spec.Containers[i].Env[:j], spec.Containers[i].Env[j+1:]...)
// 				break
// 			}
// 		}
// 	}

// 	ps.Spec.Template.Spec = spec
// }

func main() {
	str := []string{"1", "2", "3", "4", "5"}
	del := []string{"3", "4", "5"}

	for i, s := range str {
		if toDelete(s, del) {
			str = append(str[:i], str[i+1:]...)
		}
	}

	log.Println(str)
}

func toDelete(target string, envs []string) bool {
	for _, v := range envs {
		if v == target {
			return true
		}
	}
	return false
}