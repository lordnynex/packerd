package packerd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"sync"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-openapi/strfmt"
	ps "github.com/mitchellh/go-ps"
	"github.com/tompscanlan/packerd/models"
)

func VarsToEnv(vars models.Variables) []string {

	env := os.Environ()

	// delete proxy vars from the env
	for i, variable := range env {
		match, err := regexp.MatchString("(?i)_proxy=", variable)
		if err != nil {
			log.Errorln(err)
		}
		if match {
			env = append(env[:i], env[i+1:]...)
		}
	}
	for _, variable := range vars {
		env = append(env, fmt.Sprintf("%s=%s", *variable.Key, *variable.Value))
	}
	return env
}

// BuildRequestToString returns the string representation of a build request
func BuildRequestToString(br models.Buildrequest) string {

	b, err := json.Marshal(br)
	if err != nil {
		msg := fmt.Sprintf("failed to marshal to json: %s", err.Error())
		return msg
	}

	return string(b[:])
}

func ProccessNameExists(name string) bool {

	procs, err := ps.Processes()
	if err != nil {
		log.Errorf("ProccessByName: %s", err)
	}
	for _, p := range procs {
		if p.Executable() == name {
			return true
		}
	}
	return false
}

// StreamToLog sends all the read contents to the global logger
func StreamToLog(reader io.Reader) {

	b := bufio.NewScanner(reader)

	for b.Scan() {
		log.WithFields(log.Fields{
			"function": "StreamToLog",
		}).Infoln(b.Text())
	}

	if err := b.Err(); err != nil {
		log.WithFields(log.Fields{
			"function": "StreamToLog",
		}).Errorln("error reading:", err)
	}
}

// StreamToString appends all read data to the given string
func StreamToString(reader io.Reader, s *string) {
	eol := "\n"
	var mutex = &sync.Mutex{}

	b := bufio.NewScanner(reader)
	for b.Scan() {
		mutex.Lock()
		*s = *s + b.Text() + eol
		mutex.Unlock()
		runtime.Gosched()
	}

	if err := b.Err(); err != nil {
		mutex.Lock()
		*s = *s + fmt.Sprintf("%v%s", err, eol)
		mutex.Unlock()
	}
}

// RunCmd is a helper to run commands and stream the output to the daemon logs and the stage log for end user consumption
func RunCmd(command string, args []string, dir string, env []string, stage *models.Buildstage) error {

	log.Debugf("RunCmd: running command [%s %v]", command, args)
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	//cmd.Env = []string{"PACKER_LOG=1"}
	for _, e := range cmd.Env {
		log.Println("precmd ", command, e)
	}
	cmd.Env = append(cmd.Env, env...)
	for _, e := range cmd.Env {
		log.Println("cmd ", command, e)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		stage.Status = "failed"
		log.Errorf("RunCmd %s: %s", stage.Name, err)
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		stage.Status = "failed"
		log.Errorf("RunCmd %s: %s", stage.Name, err)
		return err
	}

	stage.Status = "started"
	stage.Start = strfmt.DateTime(time.Now())

	if err := cmd.Start(); err != nil {
		stage.Status = "failed"
		log.Errorf("RunCmd %s: %s", stage.Name, err)
		return err
	}

	go StreamToString(stdout, &stage.Log)

	temp, err := ioutil.ReadAll(stderr)
	if err != nil {
		log.Errorf("RunCmd %s: %s", stage.Name, err)
	}
	errlog := string(temp[:])

	log.Debugf("command stdout: %s", stage.Log)
	log.Debugf("command stderr: %s", errlog)

	if err := cmd.Wait(); err != nil {
		stage.Status = "failed"
		log.Errorf("RunCmd %s: %s", stage.Name, err)
		return err
	}

	stage.Status = "complete"
	stage.End = strfmt.DateTime(time.Now())

	return nil
}

// GetDirFreeBytes takes a directory and returns how many bytes are available for storage there
func GetDirPercentFull(dir string) (error, int64) {
	statfs := syscall.Statfs_t{}
	err := syscall.Statfs(dir, &statfs)
	if err != nil {
		log.Error(err)
		return err, 0
	}
	percentfull := 1.0 - (float64(statfs.Bavail) / float64(statfs.Blocks))

	return err, int64(percentfull * 100)
}

func ParseDockerImageSha(packerlog string) []string {
	set := make(map[string]struct{})

	re := regexp.MustCompile(`(?:Imported Docker image: *sha256:)([A-Fa-f0-9]{64})`)
	found := re.FindAllStringSubmatch(packerlog, -1)

	// uniq the list of found matches
	for _, list := range found {
		if list[1] != "" {
			str := list[1]
			set[str] = struct{}{}
		}

	}
	matches := make([]string, 0)
	for k := range set {
		matches = append(matches, k)
	}
	return matches
}

func ParseDockerArtifactUrl(packerlog string) []string {
	set := make(map[string]struct{})

	re := regexp.MustCompile(`(?:docker-tag.*Repository:\s*)([/:\.\w]*)[\s]*`)
	found := re.FindAllStringSubmatch(packerlog, -1)

	// uniq the list of found matches
	for _, list := range found {
		if list[1] != "" {
			str := list[1]
			set[str] = struct{}{}
		}

	}
	matches := make([]string, 0)
	for k := range set {
		matches = append(matches, k)
	}

	return matches
}
