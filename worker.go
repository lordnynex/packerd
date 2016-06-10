package packerd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/Sirupsen/logrus"

	"github.kdc.capitalone.com/kbs316/packerd/models"
)

const (
	Stopped  = 0
	Stopping = 1
	Started  = 2
)

type Worker struct {
	Id           int
	Done         chan bool
	BuildRequest chan *models.Buildrequest
	WorkerQueue  chan chan *models.Buildrequest
	State        int
}

func NewWorker(id int, workerqueue chan chan *models.Buildrequest) *Worker {

	worker := new(Worker)

	worker.Id = id
	worker.Done = make(chan bool)
	worker.BuildRequest = make(chan *models.Buildrequest)
	worker.WorkerQueue = workerqueue
	worker.State = Stopped

	return worker
}

func (w *Worker) RunGitCheckout(br *models.Buildrequest) error {
	args := []string{"clone", br.Branch}
	err := w.RunCmd("git", args, br.Localpath, &br.Status, &br.Buildlog)
	return err
}

func (w *Worker) RunGitClone(br *models.Buildrequest) error {
	args := []string{"clone", *br.Giturl, br.Localpath}
	err := w.RunCmd("git", args, br.Localpath, &br.Status, &br.Buildlog)
	return err
}

func (w *Worker) RunPackerValidate(br *models.Buildrequest) error {
	args := []string{"validate"}
	err := w.RunPacker(args, br)
	return err
}

func (w *Worker) RunPackerBuild(br *models.Buildrequest) error {
	// args := []string{"build", "-machine-readable"}
	args := []string{"build"}
	err := w.RunPacker(args, br)
	return err
}

func (w *Worker) RunPacker(args []string, br *models.Buildrequest) error {
	if br.Buildonly != "" {
		args = append(args, fmt.Sprintf("-only=%s", br.Buildonly))
	}

	for _, v := range br.Buildvars {
		args = append(args, "-var", fmt.Sprintf("%s=%s", *v.Key, *v.Value))

	}

	// template must be last in command
	if br.Templatepath != "" {
		args = append(args, br.Templatepath)
	}

	err := w.RunCmd("packer", args, br.Localpath, &br.Status, &br.Buildlog)

	return err
}

func (w *Worker) RunBerks(br *models.Buildrequest) error {
	args := []string{"vendor"}
	err := w.RunCmd("berks", args, br.Localpath, &br.Status, &br.Buildlog)

	return err
}

func (w *Worker) RunCmd(command string, args []string, dir string, status *string, fulllog *string) error {

	log.Infof("running command [%s %v]", command, args)
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	//cmd.Env = []string{"PACKER_LOG=1"}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		*status = fmt.Sprintf("Failed %s", command)
		log.Errorf("%s run failed: %s", command, err)
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		*status = fmt.Sprintf("Failed %s", command)
		log.Errorf("%s run failed: %s", command, err)
		return err
	}

	// stream packer to logs and build status
	go StreamToLog(stdout)
	go StreamToLog(stderr)
	go StreamToString(stdout, fulllog)

	*status = fmt.Sprintf("Started %s", command)

	if err := cmd.Start(); err != nil {
		*status = fmt.Sprintf("Failed %s", command)
		log.Errorf("%s run failed: %s", command, err)
		return err
	}

	if err := cmd.Wait(); err != nil {
		*status = fmt.Sprintf("Failed %s", command)
		log.Errorf("%s run failed: %s", command, err)
		return err
	}

	*status = fmt.Sprintf("Completed %s", command)

	return nil
}

func (w *Worker) Start() {
	w.State = Started

	go func() {
		for {
			// put our channel onto the queue.
			// dispatcher will shift it off, and send work down it
			// Re-add after doing the job
			w.WorkerQueue <- w.BuildRequest

			select {
			case build := <-w.BuildRequest:
				log.Debugf("worker%d: got build request %s", w.Id, BuildRequestToString(*build))

				w.RunGitClone(build)
				if build.Branch != "" {
					w.RunGitCheckout(build)
				}

				if _, err := os.Stat(filepath.Join(build.Localpath, "Berksfile")); err == nil {
					w.RunBerks(build)
				}

				w.RunPackerBuild(build)

			case <-w.Done:
				log.Debugf("worker%d: done", w.Id)
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	w.State = Stopping
	go func() {
		w.Done <- true
		w.State = Stopped
	}()
}
