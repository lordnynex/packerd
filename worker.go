package packerd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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
	err := w.RunCmd("git", args, br.Localpath, &br.Status, &br.Buildlog)
	return err
}

func (w *Worker) RunPacker(br *models.Buildrequest) error {
	args := []string{"build", "-machine-readable"}

	if br.Buildonly != "" {
		args = append(args, fmt.Sprintf("-only=%s", br.Buildonly))
	}

	for _, v := range br.Buildvars {
		args = append(args, "-var", fmt.Sprintf("\"%s=%s\"", *v.Key, *v.Value))
	}

	// template must be last in command
	if br.Templatepath != "" {
		args = append(args, br.Templatepath)
	}

	err := w.RunCmd("packer", args, br.Localpath, &br.Status, &br.Buildlog)

	return err
}

func (w *Worker) RunBerks(br *models.Buildrequest) error {
	args := []string{"vendor", "provision/chef/vendor-cookbooks"}
	err := w.RunCmd("berks", args, br.Localpath, &br.Status, &br.Buildlog)

	return err
}

func (w *Worker) RunCmd(command string, args []string, dir string, status *string, fulllog *string) error {

	Logger.Printf("running command [%s %v]", command, args)
	cmd := exec.Command(command, args...)
	cmd.Dir = dir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		*status = fmt.Sprintf("Failed %s", command)
		Logger.Printf("%s run failed: %s", command, err)
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		*status = fmt.Sprintf("Failed %s", command)
		Logger.Printf("%s run failed: %s", command, err)
		return err
	}

	// stream packer to logs and build status
	go StreamToLog(stdout)
	go StreamToLog(stderr)
	go StreamToString(stdout, fulllog)

	*status = fmt.Sprintf("Started %s", command)

	if err := cmd.Start(); err != nil {
		*status = fmt.Sprintf("Failed %s", command)
		Logger.Printf("%s run failed: %s", command, err)
		return err
	}

	if err := cmd.Wait(); err != nil {
		*status = fmt.Sprintf("Failed %s", command)
		Logger.Printf("%s run failed: %s", command, err)
		return err
	}

	*status = fmt.Sprintf("Completed %s", command)

	return nil
}

func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerQueue <- w.BuildRequest
			Logger.Printf("worker%d: selecting", w.Id)

			select {
			case build := <-w.BuildRequest:
				Logger.Printf("worker%d: got build request %s", w.Id, BuildRequestToString(*build))

				w.RunGitClone(build)
				//br, brerr = BuildQ.Update(build.ID, br)

				if _, err := os.Stat(filepath.Join(build.Localpath, "Berksfile")); err == nil {
					w.RunBerks(build)
				}

				w.RunPacker(build)
				//br, brerr = BuildQ.Update(build.ID, br)

			case <-w.Done:
				Logger.Printf("worker%d: done", w.Id)
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

func StreamToLog(reader io.Reader) {
	b := bufio.NewScanner(reader)
	for b.Scan() {
		Logger.Println(b.Text())
	}

	if err := b.Err(); err != nil {
		Logger.Println("error reading:", err)
	}
}

func StreamToString(reader io.Reader, s *string) {
	b := bufio.NewScanner(reader)
	for b.Scan() {
		*s = *s + b.Text()
	}

	if err := b.Err(); err != nil {
		*s = *s + fmt.Sprintf("%v", err)
	}

}
