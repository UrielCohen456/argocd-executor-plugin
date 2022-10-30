package argocd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/util/io"
	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/argoproj/argo-workflows/v3/pkg/plugins/executor"
	"k8s.io/utils/pointer"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
)

type ApiExecutor struct{
	apiClient apiclient.Client
}

func NewApiExecutor(apiClient apiclient.Client) ApiExecutor {
	return ApiExecutor{apiClient: apiClient}
}

func (e ApiExecutor) Execute(args executor.ExecuteTemplateArgs) executor.ExecuteTemplateReply {
	pluginJSON, err := args.Template.Plugin.MarshalJSON()
	if err != nil {
		err = fmt.Errorf("failed to marshal plugin to JSON from workflow spec: %w", err)
		log.Println(err.Error())
		return errorResponse(err)
	}

	plugin := &PluginSpec{}
	err = json.Unmarshal(pluginJSON, plugin)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal plugin JSON to plugin struct: %w", err)
		log.Println(err.Error())
		return errorResponse(err)
	}

	closer, appClient, err := e.apiClient.NewApplicationClient()
	if err != nil {
		err = fmt.Errorf("failed to initialize Application API client: %w", err)
		log.Println(err.Error())
		return errorResponse(err)
	}
	defer io.Close(closer)

	totalActionGroups := len(plugin.ArgoCD.Actions)
	actionGroupCount := 0

	for i, actionGroup := range plugin.ArgoCD.Actions {
		err = runActionsParallel(actionGroup, appClient)
		if err != nil {
			return failedResponse(wfv1.Progress(fmt.Sprintf("%d/%d", actionGroupCount, totalActionGroups)), fmt.Errorf("action group %d of %d failed: %w", i+1, totalActionGroups, err))
		}
		actionGroupCount += 1
	}

	return executor.ExecuteTemplateReply{
		Node: &wfv1.NodeResult{
			Phase:    wfv1.NodeSucceeded,
			Message:  "Actions completed",
			Progress: wfv1.Progress(fmt.Sprintf("%d/%d", actionGroupCount, totalActionGroups)),
		},
	}
}

// runActionsParallel runs the given group of actions in parallel and returns aggregated errors, if any.
func runActionsParallel(actionGroup []ActionSpec, appClient application.ApplicationServiceClient) error {
	wg := sync.WaitGroup{}
	errChan := make(chan error, len(actionGroup))
	for i, action := range actionGroup {
		action := action
		wg.Add(1)
		go func(actionNum int) {
			defer wg.Done()
			if action.App != nil && action.App.Sync != nil {
				err := syncAppsParallel(*action.App.Sync, action.Timeout, appClient)
				if err != nil {
					log.Println(err)
					errChan <- fmt.Errorf("parallel item %d of %d failed: failed to sync Application(s): %w", actionNum+1, len(actionGroup), err)
				}
			}
		}(i)
	}
	go func() {
		wg.Wait()
		close(errChan)
	}()
	var actionErrors []string
	for err := range errChan {
		actionErrors = append(actionErrors, err.Error())
	}
	if len(actionErrors) > 0 {
		return errors.New(strings.Join(actionErrors, ", "))

	}
	return nil
}

// syncAppsParallel loops over the apps in a SyncAction and syncs them in parallel. It waits for all responses and then
// aggregates any errors.
func syncAppsParallel(action SyncAction, timeout string, appClient application.ApplicationServiceClient) error {
	ctx, _, err := durationStringToContext(timeout)
	if err != nil {
		return fmt.Errorf("failed get action context: %w", err)
	}
	wg := sync.WaitGroup{}
	errChan := make(chan error, len(action.Apps))
	for _, app := range action.Apps {
		app := app
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := appClient.Sync(ctx, &application.ApplicationSyncRequest{
				Name:         pointer.String(app.Name),
				AppNamespace: pointer.String(app.Namespace),
				SyncOptions:  &application.SyncOptions{Items: action.Options},
			})
			if err != nil {
				errChan <- fmt.Errorf("failed to sync app %q: %w", app.Name, err)
			}
		}()
	}
	go func() {
		wg.Wait()
		close(errChan)
	}()
	var syncErrors []string
	for err := range errChan {
		syncErrors = append(syncErrors, err.Error())
	}
	if len(syncErrors) > 0 {
		return  errors.New(strings.Join(syncErrors, ", "))
	}
	return nil
}

// durationStringToContext parses a duration string and returns a context and cancel function. If timeout is empty, the
// context is context.Background().
func durationStringToContext(timeout string) (ctx context.Context, cancel func(), err error) {
	ctx = context.Background()
	cancel = func() {}
	if timeout != "" {
		duration, err := time.ParseDuration(timeout)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse timeout: %w", err)
		}
		ctx, cancel = context.WithTimeout(ctx, duration)
	}
	return ctx, cancel, nil
}

func errorResponse(err error) executor.ExecuteTemplateReply {
	return executor.ExecuteTemplateReply{
		Node: &wfv1.NodeResult{
			Phase: wfv1.NodeError,
			Message: err.Error(),
			Progress: wfv1.ProgressZero,
		},
	}
}

func failedResponse(progress wfv1.Progress, err error) executor.ExecuteTemplateReply {
	return executor.ExecuteTemplateReply{
		Node: &wfv1.NodeResult{
			Phase: wfv1.NodeFailed,
			Message: err.Error(),
			Progress: progress,
		},
	}
}
