package restapi

import (
	"crypto/tls"
	"fmt"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"io/ioutil"
	"log"
	"os"
	//"path/filepath"

	uuid "github.com/satori/go.uuid"

	"os/exec"

	"github.com/tompscanlan/packerd/models"

	"github.com/tompscanlan/packerd/restapi/operations"
	"github.com/tompscanlan/packerd/restapi/operations/command"
	"github.com/tompscanlan/packerd/restapi/operations/informational"
)

// This file is safe to edit. Once it exists it will not be overwritten

func configureFlags(api *operations.PackerdAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.PackerdAPI) http.Handler {
	var logger = log.New(os.Stderr,
		"debug: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	var buildQueue = make(map[uuid.UUID]*models.Buildstatus)

	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.InformationalGetHealthHandler = informational.GetHealthHandlerFunc(func() middleware.Responder {
		var health = informational.NewGetHealthOK()
		var pl = new(models.Health)
		pl.Status = "OK"
		health.SetPayload(pl)
		return health
	})
	api.InformationalGetPackerLogByIDHandler = informational.GetPackerLogByIDHandlerFunc(func(params informational.GetPackerLogByIDParams) middleware.Responder {
		return middleware.NotImplemented("operation informational.GetPackerLogByID has not yet been implemented")
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
			*resp.Message = "invalid uuid"

			return informational.NewGetQueueByIDBadRequest().WithPayload(resp)
		}

		status, ok := buildQueue[id]
		if !ok {
			var err = new(models.Error)
			err.Code = 1
			*err.Message = "no git url"

			return informational.NewGetQueueByIDBadRequest().WithPayload(err)
		}

		status.Eta = 30
		status.Status = "pending"
		return informational.NewGetQueueByIDOK().WithPayload(status)

	})
	api.InformationalGetQueueTestLogByIDHandler = informational.GetQueueTestLogByIDHandlerFunc(func(params informational.GetQueueTestLogByIDParams) middleware.Responder {
		return middleware.NotImplemented("operation informational.GetQueueTestLogByID has not yet been implemented")
	})
	api.CommandRunBuildHandler = command.RunBuildHandlerFunc(func(params command.RunBuildParams) middleware.Responder {

		if *params.Buildrequest.Giturl == "" {
			var err = new(models.Error)
			err.Code = 1
			*err.Message = "no git url"

			return command.NewRunBuildBadRequest().WithPayload(err)
		}

		dir, err := ioutil.TempDir("", "packerd")
		if err != nil {
			logger.Fatal(err)
		}
		defer os.RemoveAll(dir) // clean up
		logger.Printf("gotcha url right here %s", *params.Buildrequest.Giturl)

		var cmd = "git"
		var cmdOut []byte

		var args = []string{"clone", *params.Buildrequest.Giturl, dir}
		if cmdOut, err = exec.Command(cmd, args...).Output(); err != nil {
			var resp = new(models.Error)
			resp.Code = 2
			*resp.Message = "failed to clone"

			logger.Println("There was an error running git clone command: ", err, cmdOut)
			return command.NewRunBuildBadRequest().WithPayload(resp)

		}
		logger.Println(cmdOut)
		var id = uuid.NewV4()
		buildQueue[id] = new(models.Buildstatus)

		var link = new(models.Link)
		link.Rel = "status"
		link.Href = fmt.Sprintf("/queue/%s", id)
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
