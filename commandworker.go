package packerd

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"

	"github.com/tompscanlan/packerd/models"
)

// Worker recieves work on the WorkerQueue and returns Buildrequests on the BuildRequest channel
// send true to the Done channel to stop the worker
type CommandWorker struct {
	ID    int
	Done  chan bool
	State int

	// BuildRequest will recieve build requests to handle
	Work chan *models.Buildrequest
	// BuildRequestQueue, send our receiving channel to this channel,
	// which will be used by dispatcher to send us work
	WorkQueue chan chan *models.Buildrequest
}

// NewCommandWorker makes a new worker
func NewCommandWorker(id int, buildq chan chan *models.Buildrequest) *CommandWorker {

	worker := new(CommandWorker)

	worker.ID = id
	worker.Done = make(chan bool)
	worker.Work = make(chan *models.Buildrequest)
	worker.WorkQueue = buildq

	worker.State = Stopped

	return worker
}

// Start gets the worker recieving from the global WorkQueue and running the needed stages
func (w *CommandWorker) Start() {
	w.State = Started

	go func() {
		for {
			// put our channel onto the dispatcher's queue.
			// dispatcher will shift it off, and send work down it
			// Re-add it after doing the job
			log.Debug("before listen select")

			w.WorkQueue <- w.Work
			log.Debugf("added worker%d to the build Q", w.ID)

			select {
			case build := <-w.Work:

				// new request must trigger a response, add it to the list of responses
				response := models.NewBuildresponse()
				buildnumber, _ := BuildResponses.Add(build.ID, response)
				build.Responses = append(build.Responses, response.ID)

				log.Debugf("worker%d: got build request %s", w.ID, BuildRequestToString(*build))
				log.Debugf("worker%d: using response %s/%d", build.ID, buildnumber)

				response.Status = models.RUNNING
				if err := w.RunGitClone(build, buildnumber); err != nil {
					log.Error(err.Error())
					response.Status = models.FAILED
					return
				}
				if err := w.RunGitCheckout(build, buildnumber); err != nil {
					log.Error(err.Error())
					response.Status = models.FAILED
					return
				}
				if err := w.RunBerks(build, buildnumber); err != nil {
					log.Error(err.Error())
					response.Status = models.FAILED
					return
				}
				if err, _ := w.RunPackerBuild(build, buildnumber); err != nil {
					log.Error(err.Error())
					response.Status = models.FAILED
					return
				}
				if err := w.RunKitchenTest(build, buildnumber); err != nil {
					log.Error(err.Error())
					response.Status = models.FAILED
					return
				}

				if err := w.RunDockerClean(response, buildnumber); err != nil {
					log.Error(err.Error())
					// not fatal to miss a cleanup
				}
				response.Status = models.COMPLETE

			case <-w.Done:
				log.Debugf("worker%d: done", w.ID)
				return
			}
		}
	}()
}

// Stop this worker
func (w *CommandWorker) Stop() {
	w.State = Stopping
	go func() {
		w.Done <- true
		w.State = Stopped
	}()
}

// RunDockerClean creates a new build stage, and puts the results of checking out the requested branch/tag/commit
func (w *CommandWorker) RunDockerClean(resp *models.Buildresponse, buildnumber int) error {
	stagename := "DockerClean"
	stage := new(models.Buildstage)
	stage.Name = stagename
	err, _ := BuildResponses.AddStage(resp.Buildrequestid, buildnumber, stage)

	args := []string{"rmi", "-f"}
	args = append(args, resp.Images...)
	err = RunCmd("docker", args, "", nil, stage)

	return err
}

// RunKitchenTest adds a new build stage and puts the results of test kitchen in it
func (w *CommandWorker) RunKitchenTest(req *models.Buildrequest, buildnumber int) error {
	stage := models.NewBuildstage("KitchenTest")

	err, resp := BuildResponses.AddStage(req.ID, buildnumber, stage)
	if err != nil {
		return err
	}

	gem := "/opt/chefdk/embedded/bin/gem"
	if _, err := os.Stat(gem); err == nil {
		args := []string{"install", "kitchen-docker"}
		err := RunCmd(gem, args, req.Localpath, VarsToEnv(req.Buildvars), resp)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(filepath.Join(req.Localpath, "gen-kitchen-dockerfile.sh")); err == nil {
		args := []string{"gen-kitchen-dockerfile.sh"}
		err := RunCmd("bash", args, req.Localpath, VarsToEnv(req.Buildvars), resp)
		if err != nil {
			return err
		}
	}
	if _, err := os.Stat(filepath.Join(req.Localpath, ".kitchen.yml")); err == nil {

		args := []string{"test"}
		err := RunCmd("kitchen", args, req.Localpath, VarsToEnv(req.Buildvars), resp)
		if err != nil {
			return err
		}
	}

	return nil
}

// RunGitCheckout creates a new build stage, and puts the results of checking out the requested branch/tag/commit
func (w *CommandWorker) RunGitCheckout(br *models.Buildrequest, buildnumber int) error {

	if br.Branch != "" {
		stagename := "GitCheckout"
		stage := new(models.Buildstage)
		stage.Name = stagename
		err, resp := BuildResponses.AddStage(br.ID, buildnumber, stage)
		if err != nil {
			return err
		}
		args := []string{"checkout", br.Branch}
		err = RunCmd("git", args, br.Localpath, VarsToEnv(br.Buildvars), resp)
		return err
	}
	return nil
}

// RunGitClone creates a new build stage and adds results of cloning the requested git url
func (w *CommandWorker) RunGitClone(br *models.Buildrequest, buildnumber int) error {
	stage := models.NewBuildstage("GitClone")
	err, resp := BuildResponses.AddStage(br.ID, buildnumber, stage)

	args := []string{"clone", *br.Giturl, br.Localpath}
	err = RunCmd("git", args, br.Localpath, VarsToEnv(br.Buildvars), resp)

	return err
}

// RunPackerValidate creates a new build stage and adds the results of validating the packer config
func (w *CommandWorker) RunPackerValidate(br *models.Buildrequest, buildnumber int) error {
	args := []string{"validate"}
	return w.RunPacker(args, buildnumber, br)
}

// RunPackerBuild creates a new build stage,
// runs a packer build in a directory with the given packer template,
// and puts the results to the stage
func (w *CommandWorker) RunPackerBuild(br *models.Buildrequest, buildnumber int) (error, []string) {
	args := []string{"build", "-color=false"}

	err := w.RunPacker(args, buildnumber, br)
	if err != nil {
		return err, nil
	}

	// record build images for later clean up, and remote artifacts
	err, response := BuildResponses.LookupResponse(br.ID, buildnumber)
	err, stage := BuildResponses.LookupStage(br.ID, buildnumber, "RunPacker")
	response.Images = ParseDockerImageSha(stage.Log)
	for _, a := range ParseDockerArtifactUrl(stage.Log) {
		var link = new(models.Link)
		link.Rel = "artifact"
		link.Href = a
		response.Artifacts = append(response.Artifacts, link)
	}

	return nil, nil
}

// RunPacker is a helper to  get the right arguments to packer in the right order
func (w *CommandWorker) RunPacker(args []string, buildnumber int, br *models.Buildrequest) error {
	stage := models.NewBuildstage("RunPacker")
	err, _ := BuildResponses.AddStage(br.ID, buildnumber, stage)
	if err != nil {
		log.Errorf("RunPacker: %v", err)
		return err
	}

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

	err = RunCmd("packer", args, br.Localpath, VarsToEnv(br.Buildvars), stage)

	//	err, response := BuildResponses.LookupResponse(br.ID, buildnumber)

	return err
}

// RunBerks creates a build stage, and runs berks vendor to pull needed chef cookbooks for preovisioning
func (w *CommandWorker) RunBerks(br *models.Buildrequest, buildnumber int) error {
	stage := models.NewBuildstage("RunBerksVendor")
	err, resp := BuildResponses.AddStage(br.ID, buildnumber, stage)
	if err != nil {
		return err
	}

	if _, err := os.Stat(filepath.Join(br.Localpath, "Berksfile")); err == nil {
		args := []string{"vendor"}
		err := RunCmd("berks", args, br.Localpath, VarsToEnv(br.Buildvars), resp)

		return err
	}
	return nil
}
