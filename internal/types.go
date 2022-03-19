package types

// wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
//  "github.com/argoproj/argo-workflows/v3/pkg/plugins/executor"

type ArgocdPluginSpec struct {
	ServerUrl string 					`json:"serverUrl,omitempty"` 
	Actions  	[]ActionSpec		`json:"actions,omitempty"`
}

// All possible actions that can be taken
type ActionSpec struct {
	// A sync action 
	Sync			SyncAction		`json:"sync,omitempty"`
}

// An action that triggers an argocd sync
type SyncAction struct {
	// The project the application requested resides in
	Project		string		`json:"project,omitempty"`

	// The applications inside the given project 
	Apps			[]string		`json:"apps,omitempty"`
}
