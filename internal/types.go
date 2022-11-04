package argocd

// PluginSpec represents the `plugin` block of an Argo Workflows template.
type PluginSpec struct {
	ArgoCD ArgocdPluginSpec `json:"argocd,omitempty"`
}

// ArgocdPluginSpec describes the parameters of the Argo CD plugin.
type ArgocdPluginSpec struct {
	Actions [][]ActionSpec `json:"actions,omitempty"`
}

type ActionSpec struct {
	App     *AppActionSpec `json:"app,omitempty"`
	Timeout string         `json:"timeout,omitempty"`
}

// AppActionSpec describes all possible actions that can be taken by the plugin.
type AppActionSpec struct {
	// A sync action
	Sync *SyncAction `json:"sync,omitempty"`
	Diff *DiffAction `json:"diff,omitempty"`
}

type DiffAction struct {
	App         `json:"app,omitempty"`
	Revision    string `json:"revision,omitempty"`
	Refresh     bool   `json:"refresh,omitempty"`
	HardRefresh bool   `json:"hardRefresh,omitempty"`
}

// SyncAction describes an action that triggers an argocd sync.
type SyncAction struct {
	// Apps are the applications to be synced.
	Apps []App `json:"apps,omitempty"`
	// Options is a list of option=value pairs to configure the sync operation. https://argo-cd.readthedocs.io/en/stable/user-guide/sync-options/
	Options []string `json:"options,omitempty"`
}

// App specifies the app to be synced.
type App struct {
	// Namespace is the namespace in which the app is installed. If empty, assume the same namespace as the Argo CD
	// API server.
	Namespace string `json:"namespace,omitempty"`
	// Name is the name of the app (not prefixed with {namespace}/).
	Name string `json:"name,omitempty"`
}
