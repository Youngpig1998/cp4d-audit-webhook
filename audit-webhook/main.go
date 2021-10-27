package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
	"net/http"
	"os"

	// TODO: try this library to see if it generates correct json patch
	// https://github.com/mattbaird/jsonpatch
)

// toAdmissionResponse is a helper function to create an AdmissionResponse
// with an embedded error
func toAdmissionResponse(err error) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}

// admitFunc is the type we use for all of our validators and mutators
type admitFunc func(v1beta1.AdmissionReview) *v1beta1.AdmissionResponse

// serve handles the http portion of a request prior to handing to an admit
// function
func serve(w http.ResponseWriter, r *http.Request, admit admitFunc) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		klog.Errorf("contentType=%s, expect application/json", contentType)
		return
	}

	klog.V(2).Info(fmt.Sprintf("handling request: %s", body))

	// The AdmissionReview that was sent to the webhook
	requestedAdmissionReview := v1beta1.AdmissionReview{}

	// The AdmissionReview that will be returned
	responseAdmissionReview := v1beta1.AdmissionReview{}

	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(body, nil, &requestedAdmissionReview); err != nil {
		klog.Error(err)
		responseAdmissionReview.Response = toAdmissionResponse(err)
	} else {
		// pass to admitFunc
		responseAdmissionReview.Response = admit(requestedAdmissionReview)
	}

	// Return the same UID
	responseAdmissionReview.Response.UID = requestedAdmissionReview.Request.UID

	klog.V(2).Info(fmt.Sprintf("sending response: %v", responseAdmissionReview.Response))

	respBytes, err := json.Marshal(responseAdmissionReview)
	if err != nil {
		klog.Error(err)
	}
	if _, err := w.Write(respBytes); err != nil {
		klog.Error(err)
	}
}

func addSidecarToPods(w http.ResponseWriter, r *http.Request) {
	serve(w, r, addSidecar)
}


func ReadyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}


func main() {



	var config = Config{CertFile: "certs/tls.crt", KeyFile: "certs/tls.key"}
	volume := os.Getenv("VOLUME_PATCH")
	container := os.Getenv("CONTAINER_PATCH")
	var content = Content{Volume: volume, Container: container}
	config.addFlags()
	content.addPodsContainerPatch()
	flag.Parse()
	http.HandleFunc("/readyz", ReadyHandler)
	http.HandleFunc("/healthz", HealthHandler)
	http.HandleFunc("/add-sidecar", addSidecarToPods)
	fmt.Print("SERVER STARTING")
	server := &http.Server{
		Addr:      ":8081",
		TLSConfig: configTLS(config),
	}

	//server.ListenAndServe()
	server.ListenAndServeTLS("", "")
	fmt.Print("SERVER COMPLETED")


}
