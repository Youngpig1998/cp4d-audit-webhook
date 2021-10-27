package main

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/api/admission/v1beta1"
	"k8s.io/klog"
)

type Content struct {
	Container string
	Volume    string
}

// const (
// 	podsContainerPatch string = `[
// 		{"op":"add","path":"/spec/volumes/-","value":` + content.Volume + `},
// 		{"op":"add","path":"/spec/containers/-","value":` + content.Container + `}
// 	]`
// )

var podsContainerPatch string

func (c *Content) addPodsContainerPatch() {
	// flag.StringVar(&c.Container, "container-json-content", c.Container, ""+
	// 	"File containing the sidecar container json content).")
	// flag.StringVar(&c.Volume, "volume-json-content", c.Volume, ""+
	// 	"File containing the sidecar volume json content")
	podsContainerPatch = `[
			{"op":"add","path":"/spec/volumes/-","value":` + c.Volume + `},
			{"op":"add","path":"/spec/containers/-","value":` + c.Container + `}
		]`
	fmt.Println(podsContainerPatch)
}

func addSidecar(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	klog.V(2).Info("calling add sidecar")
	podResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	if ar.Request.Resource != podResource {
		klog.Errorf("expect resource to be %s", podResource)
		return nil
	}

	raw := ar.Request.Object.Raw
	pod := corev1.Pod{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &pod); err != nil {
		klog.Error(err)
		return toAdmissionResponse(err)
	}
	klog.V(2).Info(fmt.Sprintf("Adding sidecar for pod: %s\n", pod.Name))
	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true
	var msg string
	if v, ok := pod.Labels["cp4d-audit"]; ok {
		if v == "yes" {
			reviewResponse.Patch = []byte(podsContainerPatch)
			pt := v1beta1.PatchTypeJSONPatch
			reviewResponse.PatchType = &pt
		} else {
			reviewResponse.Allowed = true
			msg = msg + "cp4d-audit is not set as yes"
		}
	} else {
		reviewResponse.Allowed = true
		msg = msg + "cp4d-audit label not found"
	}
	if !reviewResponse.Allowed {
		reviewResponse.Result = &metav1.Status{Message: strings.TrimSpace(msg)}
	}
	return &reviewResponse
}
