package packerd

import (
	"bytes"
	"os/exec"
	//"strings"
	//"runtime"

	"github.com/tompscanlan/packerd/models"
)

type Worker struct {
	Id           int
	Done         chan bool
	BuildRequest chan *models.Buildrequest
	WorkerQueue  chan chan *models.Buildrequest
}

func NewWorker(id int, workerqueue chan chan *models.Buildrequest) *Worker {

	worker := new(Worker)

	worker.Id = id
	worker.Done = make(chan bool)
	worker.BuildRequest = make(chan *models.Buildrequest)
	worker.WorkerQueue = workerqueue

	return worker
}

func (w *Worker) RunGitClone(br *models.Buildrequest) error {

	args := []string{"clone", *br.Giturl, br.Localpath}
	cmd := exec.Command("git", args...)
	br.Status = "Checking Out"

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		br.Status = "Failed"
		Logger.Printf("failed to git clone: %v", err)
		return err
	}

	Logger.Printf("git clone: done")

	br.Buildlog = stdout.String() + stderr.String()

	br.Status = "Checked Out"
	return nil
}

func (w *Worker) RunPacker(br *models.Buildrequest) {

	args := []string{"build", "-machine-readable"}
	if br.Templatepath != "" {
		args = append(args, br.Templatepath)
	}

	cmd := exec.Command("packer", args...)
	cmd.Dir = br.Localpath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	Logger.Printf("packer run starting: %v", BuildRequestToString(*br))
	br.Status = "Started"
	err := cmd.Run()
	if err != nil {
		br.Status = "Failed"
		Logger.Printf("packer run failed")
	} else {
		br.Status = "Done"
		Logger.Printf("packer run done")
	}
	br.Buildlog = br.Buildlog + stdout.String() + stderr.String()
}

func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerQueue <- w.BuildRequest
			Logger.Printf("worker%d: selecting", w.Id)

			select {
			case build := <-w.BuildRequest:
				Logger.Printf("worker%d: got build request %s", w.Id, BuildRequestToString(*build))
				//br, brerr := BuildQ.LookUp(build.ID)

				//if brerr != nil {
				//	Logger.Printf("worker%d: failed to look up %s", w.Id, build.ID)
				//}

				w.RunGitClone(build)
				//br, brerr = BuildQ.Update(build.ID, br)
				w.RunPacker(build)
				//br, brerr = BuildQ.Update(build.ID, br)

			case <-w.Done:
				Logger.Printf("worker%d: done")
				return
			}
		}
	}()
}
func (w *Worker) Stop() {

	go func() {
		w.Done <- true
	}()
}
