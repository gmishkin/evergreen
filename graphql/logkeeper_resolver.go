package graphql

import (
	"context"
	"fmt"

	"github.com/evergreen-ci/evergreen/rest/model"
	"github.com/evergreen-ci/plank"
)

// Task is the resolver for the task field.
func (r *logkeeperBuildResolver) Task(ctx context.Context, obj *plank.Build) (*model.APITask, error) {
	task, err := getTask(ctx, obj.TaskID, &obj.TaskExecution, r.sc.GetURL())
	if err != nil {
		return nil, ResourceNotFound.Send(ctx, fmt.Sprintf("Finding task %s for buildId: %s: %s", obj.TaskID, obj.ID, err.Error()))
	}
	return task, nil
}

// LogkeeperBuild returns LogkeeperBuildResolver implementation.
func (r *Resolver) LogkeeperBuild() LogkeeperBuildResolver { return &logkeeperBuildResolver{r} }

type logkeeperBuildResolver struct{ *Resolver }
