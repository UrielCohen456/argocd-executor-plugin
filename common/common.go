package common

import (
	"io/ioutil"
	"os"
	"strings"
)

// Returns the namespace the pod runs in.
func Namespace() string {
	// This way assumes you've set the POD_NAMESPACE environment variable using the downward API.
	if ns, ok := os.LookupEnv("POD_NAMESPACE"); ok {
		return ns
	}

	// Fall back to the namespace associated with the service account token, if available
	if data, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
		if ns := strings.TrimSpace(string(data)); len(ns) > 0 {
			return ns
		}
	}

	return "default"
}
