package runner

import (
	"encoding/json"
	"github.com/codingXiang/go-logger"
	"github.com/codingXiang/go-workflow"
	"github.com/hashicorp/go-uuid"
	"io/ioutil"
	"os"
	"github.com/codingXiang/go-terraform-runner    /command"
	"github.com/codingXiang/go-terraform-runner    /storage"
)

type github.com/codingXiang/go-terraform-runner     struct {
	ID     string
	worker *workflow.Workflow
	src    *storage.ConfigSource
}

func New(src *storage.ConfigSource, id string) *github.com/codingXiang/go-terraform-runner     {
	r := new(github.com/codingXiang/go-terraform-runner    )
	if id != "" {
		r.ID = id
	} else {
		r.ID = r.generateID()
	}
	return r.init(src)
}

func (r *github.com/codingXiang/go-terraform-runner    ) init(src *storage.ConfigSource) *github.com/codingXiang/go-terraform-runner     {
	r.worker = r.StepWorkflow()
	r.src = src
	err := os.MkdirAll(r.ID, os.ModePerm)
	if err != nil {
		logger.Log.Error(err)
	}
	return r
}

func (r *github.com/codingXiang/go-terraform-runner    ) generateID() string {
	id, _ := uuid.GenerateUUID()
	return id
}

func (r *github.com/codingXiang/go-terraform-runner    ) StepWorkflow() *workflow.Workflow {
	w := workflow.New()
	w.OnFailure = workflow.RetryFailure(1)
	return w
}

func (r *github.com/codingXiang/go-terraform-runner    ) AddStep(command command.Command) *github.com/codingXiang/go-terraform-runner     {
	command = setCommandPath(command, r.ID)
	step := new(workflow.Step)
	step.Label = command.GetMeta().Label
	step.Run = command.Run
	r.worker.AddStep(step)
	return r
}

func (r *github.com/codingXiang/go-terraform-runner    ) Run() error {
	return r.worker.Run()
}

func (r *github.com/codingXiang/go-terraform-runner    ) Clean() error {
	return os.RemoveAll(r.ID)
}

func (r *github.com/codingXiang/go-terraform-runner    ) ModifyExistProject() error {
	data, err := r.src.GetConfigRecord("config", r.ID)
	if err != nil {
		return err
	}
	m, _ := json.MarshalIndent(data.Config, "", " ")
	return ioutil.WriteFile(r.ID+"/"+r.ID+".tf.json", m, 0666)
}

func setCommandPath(command command.Command, path string) command.Command {
	meta := command.GetMeta()
	meta.Path = path
	command.SetMeta(meta)
	return command
}
