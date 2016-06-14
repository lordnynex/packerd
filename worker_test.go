package packerd

import (
	"strings"
	"testing"

	"github.com/tompscanlan/packerd/models"
)

func TestWorker(t *testing.T) {
	wq := make(chan chan *models.Buildrequest, 1)
	w := NewCommandWorker(1, wq)

	w.Start()
	if w.State != Started {
		t.Errorf("expected worker to be started, got %d", w.State)
	}

	w.Stop()
	if w.State != Stopping {
		t.Errorf("expected worker to be stopped, got %d", w.State)
	}

}

func TestTrue(t *testing.T) {

	var dir string
	bstage := new(models.Buildstage)

	wq := make(chan chan *models.Buildrequest, 1)
	w := NewCommandWorker(1, wq)

	var args []string
	err := RunCmd("true", args, dir, bstage)

	if err != nil {
		t.Errorf("worker run failed: %v", err)
	}

	if strings.Contains(bstage.Status, "Failed") {
		t.Errorf("worker run failed %d", w.State)
	}
}

func TestArgs(t *testing.T) {

	var dir string
	bstage := new(models.Buildstage)

	wq := make(chan chan *models.Buildrequest, 1)
	w := NewCommandWorker(1, wq)

	args := []string{"a", "b", "-var", "\"asd=foo\""}
	err := RunCmd("echo", args, dir, bstage)

	if err != nil {
		t.Errorf("worker run failed: %v", err)
	}

	if strings.Contains(bstage.Status, "Failed") {
		t.Errorf("worker run failed %d", w.State)
	}
}

func TestPacker(t *testing.T) {
	buildrequest := new(models.Buildrequest)
	Builds.Add(buildrequest)
	buildnumber, _ := BuildResponses.Add(buildrequest.ID, new(models.Buildresponse))

	packerconf := "testing/simple.json"
	wq := make(chan chan *models.Buildrequest, 1)
	w := NewCommandWorker(1, wq)

	varA := new(models.Variable)
	keyA := "empty"
	valA := "moo"

	varA.Key = &keyA
	varA.Value = &valA

	vars := append(buildrequest.Buildvars, varA)

	buildrequest.Templatepath = packerconf

	buildrequest.Buildvars = vars
	merr := w.RunPackerValidate(buildrequest, buildnumber)
	if merr != nil {
		t.Errorf("error running packer: %v", merr)
	}

	merr, resp := BuildResponses.LookupResponses(buildrequest.ID)

	if resp == nil {
		t.Error("worker got no response")
		return
	}

	if strings.Contains(resp[buildnumber].Status, "Failed") {
		t.Errorf("worker run failed %d", w.State)
	}
}
