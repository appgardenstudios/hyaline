package e2e

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/go-github/v74/github"
)

// dispatchWorkflow triggers a GitHub Actions workflow and waits for it to complete
func dispatchWorkflow(t *testing.T, token string, repo string, workflowFileName string, inputs map[string]interface{}) {
	t.Helper()

	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		t.Fatalf("invalid repo format: %s", repo)
	}
	owner, repoName := parts[0], parts[1]

	client := github.NewClient(nil).WithAuthToken(token)
	ctx := context.Background()

	// Get the current latest run ID before dispatching
	runs, _, err := client.Actions.ListWorkflowRunsByFileName(ctx, owner, repoName, workflowFileName, &github.ListWorkflowRunsOptions{
		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	})
	if err != nil {
		t.Fatalf("failed to list workflow runs: %v", err)
	}

	var lastRunID int64
	if runs.GetTotalCount() > 0 && len(runs.WorkflowRuns) > 0 {
		lastRunID = runs.WorkflowRuns[0].GetID()
	}

	// Get the default branch
	repoInfo, _, err := client.Repositories.Get(ctx, owner, repoName)
	if err != nil {
		t.Fatalf("failed to get repository info: %v", err)
	}
	defaultBranch := repoInfo.GetDefaultBranch()

	// Dispatch the workflow
	t.Logf("Dispatching workflow %s with inputs: %v on branch %s", workflowFileName, inputs, defaultBranch)
	_, err = client.Actions.CreateWorkflowDispatchEventByFileName(ctx, owner, repoName, workflowFileName, github.CreateWorkflowDispatchEventRequest{
		Ref:    defaultBranch,
		Inputs: inputs,
	})
	if err != nil {
		t.Fatalf("failed to dispatch workflow: %v", err)
	}

	// Wait for the new run to appear and complete
	t.Logf("Waiting for workflow run to complete...")
	timeout := time.After(60 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatalf("timeout waiting for workflow to complete")
		case <-ticker.C:
			runs, _, err := client.Actions.ListWorkflowRunsByFileName(ctx, owner, repoName, workflowFileName, &github.ListWorkflowRunsOptions{
				ListOptions: github.ListOptions{
					PerPage: 1,
				},
			})
			if err != nil {
				t.Fatalf("failed to list workflow runs: %v", err)
			}

			// Check the newest run (at index 0) if it's newer than lastRunID
			if len(runs.WorkflowRuns) > 0 {
				run := runs.WorkflowRuns[0]
				if run.GetID() > lastRunID {
					status := run.GetStatus()
					conclusion := run.GetConclusion()
					t.Logf("Workflow run %d status: %s, conclusion: %s", run.GetID(), status, conclusion)

					if status == "completed" {
						if conclusion == "success" {
							t.Logf("Workflow completed successfully")
							return
						}
						t.Fatalf("workflow failed with conclusion: %s", conclusion)
					}
					// Still running, continue waiting
				}
			}
		}
	}
}
