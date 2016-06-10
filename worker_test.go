package packerd

import (
	"strings"
	"testing"

	"github.kdc.capitalone.com/kbs316/packerd/models"
)

func TestWorker(t *testing.T) {
	wq := make(chan chan *models.Buildrequest, 1)
	w := NewWorker(1, wq)

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

	var dir, status, log string
	wq := make(chan chan *models.Buildrequest, 1)
	w := NewWorker(1, wq)

	var args []string
	err := w.RunCmd("true", args, dir, &status, &log)

	if err != nil {
		t.Errorf("worker run failed: %v", err)
	}

	if strings.Contains(status, "Failed") {
		t.Errorf("worker run failed %d", w.State)
	}
}

func TestArgs(t *testing.T) {

	var dir, status, log string
	wq := make(chan chan *models.Buildrequest, 1)
	w := NewWorker(1, wq)

	args := []string{"a", "b", "-var", "\"asd=foo\""}
	err := w.RunCmd("echo", args, dir, &status, &log)

	if err != nil {
		t.Errorf("worker run failed: %v", err)
	}

	if strings.Contains(status, "Failed") {
		t.Errorf("worker run failed %d", w.State)
	}
}

func TestPacker(t *testing.T) {

	var status string
	packerconf := "testing/simple.json"
	wq := make(chan chan *models.Buildrequest, 1)
	w := NewWorker(1, wq)
	br := new(models.Buildrequest)

	varA := new(models.Variable)
	keyA := "empty"
	valA := "moo"

	varA.Key = &keyA
	varA.Value = &valA

	vars := append(br.Buildvars, varA)

	br.Templatepath = packerconf

	br.Buildvars = vars
	err := w.RunPackerValidate(br)

	if err != nil {
		t.Errorf("worker run failed: %v: %s", err, br.Buildlog)
	}

	if strings.Contains(status, "Failed") {
		t.Errorf("worker run failed %d", w.State)
	}
}
