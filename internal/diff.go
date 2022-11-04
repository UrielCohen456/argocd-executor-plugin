package argocd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/argoproj/argo-cd/v2/controller"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/settings"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/argoproj/argo-cd/v2/util/argo"
	"github.com/argoproj/argo-cd/v2/util/io"
	"github.com/argoproj/gitops-engine/pkg/sync/hook"
	"github.com/argoproj/gitops-engine/pkg/sync/ignore"
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Most of this is copied from the Argo CD CLI code. We should refactor this to be shared.
// https://github.com/argoproj/argo-cd/blob/master/cmd/argocd/commands/app.go

type objKeyLiveTarget struct {
	key    kube.ResourceKey
	live   *unstructured.Unstructured
	target *unstructured.Unstructured
}

type resourceInfoProvider struct {
	namespacedByGk map[schema.GroupKind]bool
}

// Infer if obj is namespaced or not from corresponding live objects list. If corresponding live object has namespace then target object is also namespaced.
// If live object is missing then it does not matter if target is namespaced or not.
func (p *resourceInfoProvider) IsNamespaced(gk schema.GroupKind) (bool, error) {
	return p.namespacedByGk[gk], nil
}

// liveObjects deserializes the list of live states into unstructured objects
func liveObjects(resources []*v1alpha1.ResourceDiff) ([]*unstructured.Unstructured, error) {
	objs := make([]*unstructured.Unstructured, len(resources))
	for i, resState := range resources {
		obj, err := resState.LiveObject()
		if err != nil {
			return nil, err
		}
		objs[i] = obj
	}
	return objs, nil
}

func getRefreshType(refresh bool, hardRefresh bool) *string {
	if hardRefresh {
		refreshType := string(v1alpha1.RefreshTypeHard)
		return &refreshType
	}

	if refresh {
		refreshType := string(v1alpha1.RefreshTypeNormal)
		return &refreshType
	}

	return nil
}

func groupObjsByKey(localObs []*unstructured.Unstructured, liveObjs []*unstructured.Unstructured, appNamespace string) (map[kube.ResourceKey]*unstructured.Unstructured, error) {
	namespacedByGk := make(map[schema.GroupKind]bool)
	for i := range liveObjs {
		if liveObjs[i] != nil {
			key := kube.GetResourceKey(liveObjs[i])
			namespacedByGk[schema.GroupKind{Group: key.Group, Kind: key.Kind}] = key.Namespace != ""
		}
	}
	localObs, _, err := controller.DeduplicateTargetObjects(appNamespace, localObs, &resourceInfoProvider{namespacedByGk: namespacedByGk})
	if err != nil {
		return nil, fmt.Errorf("failed to deduplicate target objects: %w", err)
	}
	objByKey := make(map[kube.ResourceKey]*unstructured.Unstructured)
	for i := range localObs {
		obj := localObs[i]
		if !(hook.IsHook(obj) || ignore.Ignore(obj)) {
			objByKey[kube.GetResourceKey(obj)] = obj
		}
	}
	return objByKey, nil
}

func groupObjsForDiff(resources *application.ManagedResourcesResponse, objs map[kube.ResourceKey]*unstructured.Unstructured, items []objKeyLiveTarget, argoSettings *settings.Settings, appName string) ([]objKeyLiveTarget, error) {
	resourceTracking := argo.NewResourceTracking()
	for _, res := range resources.Items {
		var live = &unstructured.Unstructured{}
		err := json.Unmarshal([]byte(res.NormalizedLiveState), &live)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal normalized live state: %w", err)
		}

		key := kube.ResourceKey{Name: res.Name, Namespace: res.Namespace, Group: res.Group, Kind: res.Kind}
		if key.Kind == kube.SecretKind && key.Group == "" {
			// Don't bother comparing secrets, argo-cd doesn't have access to k8s secret data
			delete(objs, key)
			continue
		}
		if local, ok := objs[key]; ok || live != nil {
			if local != nil && !kube.IsCRD(local) {
				err = resourceTracking.SetAppInstance(local, argoSettings.AppLabelKey, appName, "", v1alpha1.TrackingMethod(argoSettings.GetTrackingMethod()))
				if err != nil {
					return nil, fmt.Errorf("failed to set app instance: %w", err)
				}
			}

			items = append(items, objKeyLiveTarget{key, live, local})
			delete(objs, key)
		}
	}
	for key, local := range objs {
		if key.Kind == kube.SecretKind && key.Group == "" {
			// Don't bother comparing secrets, argo-cd doesn't have access to k8s secret data
			delete(objs, key)
			continue
		}
		items = append(items, objKeyLiveTarget{key, nil, local})
	}
	return items, nil
}

// GetDiff gets a diff between two unstructured objects to stdout using an external diff utility
func GetDiff(live *unstructured.Unstructured, target *unstructured.Unstructured) (string, error) {
	tempDir, err := os.MkdirTemp("", "argocd-diff")
	if err != nil {
		return "", err
	}
	defer func() {
		err = os.RemoveAll(tempDir)
		if err != nil {
			fmt.Printf("Failed to delete temp dir %s: %v\n", tempDir, err)
		}
	}()
	targetFile, err := os.CreateTemp(tempDir, "target")
	if err != nil {
		return "", err
	}
	defer io.Close(targetFile)
	targetData := []byte("")
	if target != nil {
		targetData, err = yaml.Marshal(target)
		if err != nil {
			return "", err
		}
	}
	_, err = targetFile.Write(targetData)
	if err != nil {
		return "", err
	}
	liveFile, err := os.CreateTemp(tempDir, "live")
	if err != nil {
		return "", err
	}
	defer io.Close(liveFile)
	liveData := []byte("")
	if live != nil {
		liveData, err = yaml.Marshal(live)
		if err != nil {
			return "", err
		}
	}
	_, err = liveFile.Write(liveData)
	if err != nil {
		return "", err
	}
	cmd := exec.Command("diff", liveFile.Name(), targetFile.Name())
	out, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 1 {
				return string(out), nil
			}
			return "", fmt.Errorf("diff command failed with exit code %d: %v", exitError.ExitCode(), err)
		}
		return "", fmt.Errorf("diff command failed: %v", err)
	}
	return "", nil
}
