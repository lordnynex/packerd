package restapi

import (
	"crypto/tls"
	"fmt"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"bytes"
	"io/ioutil"
	"log"
	"os"
	

	uuid "github.com/satori/go.uuid"

	"os/exec"

	"github.com/tompscanlan/packerd/models"

	"github.com/tompscanlan/packerd/restapi/operations"
	"github.com/tompscanlan/packerd/restapi/operations/command"
	"github.com/tompscanlan/packerd/restapi/operations/informational"
	"github.com/tompscanlan/packerd/worker"
)

var buildQueue = make(map[uuid.UUID]*models.Buildrequest)
var logger = log.New(os.Stderr,
	"debug: ",
	log.Ldate|log.Ltime|log.Lshortfile)

// This file is safe to edit. Once it exists it will not be overwritten

func configureFlags(api *operations.PackerdAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.PackerdAPI) http.Handler {

	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// /health
	api.InformationalGetHealthHandler = informational.GetHealthHandlerFunc(func() middleware.Responder {
		var health = informational.NewGetHealthOK()
		var pl = new(models.Health)
		pl.Status = "OK"

		return health.WithPayload(pl)
	})

	api.InformationalGetQueueHandler = informational.GetQueueHandlerFunc(func() middleware.Responder {
		var links []*models.Link

		for key, _ := range buildQueue {
			var link = new(models.Link)
			link.Rel = "status"
			link.Href = fmt.Sprintf("/queue/%s", key)
			links = append(links, link)
		}
		return informational.NewGetQueueOK().WithPayload(links)
	})

	api.InformationalGetQueueByIDHandler = informational.GetQueueByIDHandlerFunc(func(params informational.GetQueueByIDParams) middleware.Responder {

		id, err := uuid.FromString(params.ID)
		if err != nil {
			var resp = new(models.Error)
			resp.Code = 3
			*resp.Message = "invalid or missing uuid"

			return informational.NewGetQueueByIDBadRequest().WithPayload(resp)
		}

		status, ok := buildQueue[id]
		if !ok {
			var err = new(models.Error)
			err.Code = 4
			*err.Message = "non-existiant build"

			return informational.NewGetQueueByIDBadRequest().WithPayload(err)
		}

		status.Eta = 30
		status.Status = "pending"
		return informational.NewGetQueueByIDOK().WithPayload(status)

	})

	// /queue/{id}/buildlog
	api.InformationalGetPackerLogByIDHandler = informational.GetPackerLogByIDHandlerFunc(func(params informational.GetPackerLogByIDParams) middleware.Responder {

		request, err := idToBuildRequest(params.ID)

		if err.Message != nil {
			return informational.NewGetPackerLogByIDBadRequest().WithPayload(err)
		}
		return informational.NewGetPackerLogByIDOK().WithPayload(request.Status)

	})

	api.InformationalGetQueueTestLogByIDHandler = informational.GetQueueTestLogByIDHandlerFunc(func(params informational.GetQueueTestLogByIDParams) middleware.Responder {
		return middleware.NotImplemented("operation informational.GetQueueTestLogByID has not yet been implemented")
	})

	api.CommandRunBuildHandler = command.RunBuildHandlerFunc(func(params command.RunBuildParams) middleware.Responder {

		var error = new(models.Error)
		var id = uuid.NewV4()

		buildQueue[id] = params.Buildrequest

		if *buildQueue[id].Giturl == "" {
			error.Code = 1
			*error.Message = "no git url"
			return command.NewRunBuildBadRequest().WithPayload(error)
		}

		dir, err := ioutil.TempDir("", "packerd")
		if err != nil {
			logger.Println(err)
			error.Code = 400
			*error.Message = err.Error()
			return command.NewRunBuildBadRequest().WithPayload(error)
		}
		buildQueue[id].Localpath = dir

		//error = runGitClone(*buildQueue[id].Giturl, buildQueue[id].Localpath)
	
		
		//if *error.Message != "" {
		//return command.NewRunBuildBadRequest().WithPayload(error)
		//}

		var link = new(models.Link)
		link.Rel = "status"
		link.Href = fmt.Sprintf("/queue/%s", id)

		error = runPacker(buildQueue[id])
		//if *error.Message != "" {
		//	return command.NewRunBuildBadRequest().WithPayload(error)
		//}

		return command.NewRunBuildAccepted().WithPayload(link)
		//return middleware.NotImplemented(*params.Buildrequest.Giturl)
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {

	//	recovery := recover.New(&recover.Options{
	//		Log: log.Print,
	//	})
	//	return recovery(handler)
	return handler

}

func idToBuildRequest(id string) (*models.Buildrequest, *models.Error) {

	var error = new(models.Error)
	var bs = new(models.Buildrequest)

	if id == "" {
		error.Code = 5
		*error.Message = "missing uuid"
		return bs, error
	}

	uuid, err := uuid.FromString(id)
	if err != nil {
		error.Code = 3
		*error.Message = "invalid uuid"
		return bs, error
	}

	request, ok := buildQueue[uuid]
	if !ok {
		error.Code = 4
		*error.Message = "non-existiant build"

		return bs, error
	}

	return request, nil
}

func runPacker(br *models.Buildrequest) *models.Error {
	var error = new(models.Error)
	var bin = "packer"
	var args = []string{"build", "-machine-readable"}

	if br.Templatepath != "" {
		args = append(args, br.Templatepath)
	}

	cmd := exec.Command(bin, args...)
	cmd.Dir = br.Localpath

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		logger.Println(err)
		logger.Println("packer:" + stdout.String() + stderr.String())

		error.Code = 400
		//		error.Message = "Failed to clone: " + stderr.String()

		return error
	}
	logger.Println("packer:" + stdout.String() + stderr.String())

	return error
}

func runGitClone(url string, dir string) *models.Error {
	var error = new(models.Error)
	var bin = "git"
	var args = []string{"clone", url, dir}

	cmd := exec.Command(bin, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		logger.Println(err)

		error.Code = 400
		*error.Message = "Failed to clone: " + err.Error() + ":" + stderr.String()

		return error
	}
	logger.Println("git clone:" + stdout.String() + stderr.String())
	return error
}
