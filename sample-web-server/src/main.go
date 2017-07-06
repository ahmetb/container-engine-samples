package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/compute/metadata"
)

const (
	httpAddr  = ":8080"
	httpsAddr = ":443"
)

var (
	tpl      = template.Must(template.ParseFiles("template.html"))
	hostname string
)

func init() {
	h, err := os.Hostname()
	if err != nil {
		panic(fmt.Errorf("cannot get hostname: %+v", err))
	}
	hostname = h
}

func main() {
	http.HandleFunc("/", htmlHandler)

	var addr string
	tlsDir := os.Getenv("TLS_DIR")
	if tlsDir != "" {
		addr = httpsAddr
	} else {
		addr = httpAddr
	}

	log.Printf("running on gce: %v", metadata.OnGCE())
	log.Println("starting web server on " + addr)
	if tlsDir == "" {
		log.Fatal(http.ListenAndServe(addr, http.DefaultServeMux))
	} else {
		log.Fatal(http.ListenAndServeTLS(addr, filepath.Join(tlsDir, "tls.crt"), filepath.Join(tlsDir, "tls.key"), nil))
	}
}

func htmlHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Connection","close")

	log.Println("handing request on " + r.URL.String())
	envs := os.Environ()
	env := make(map[string]string)
	for _, e := range envs {
		ss := strings.SplitN(e, "=", 2)
		env[ss[0]] = ss[1]
	}

	data := map[string]interface{}{
		"env":      env,
		"req":      r,
		"hostname": hostname,
		"onGCE":    metadata.OnGCE(),
	}

	// Fetch instance metadata if on GCE/GKE
	if metadata.OnGCE() {
		val := func(s string, _ error) string { return s }
		data["gce"] = map[string]string{
			"hostname":     val(metadata.Hostname()),
			"projectID":    val(metadata.ProjectID()),
			"instanceName": val(metadata.InstanceName()),
			"instanceID":   val(metadata.InstanceID()),
			"zone":         val(metadata.Zone()),
		}
	}

	tpl.Execute(w, data)
}
