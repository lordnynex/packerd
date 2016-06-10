package restapi

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.kdc.capitalone.com/kbs316/packerd"
	"github.kdc.capitalone.com/kbs316/packerd/models"

	"github.kdc.capitalone.com/kbs316/packerd/restapi/operations"
	"github.kdc.capitalone.com/kbs316/packerd/restapi/operations/command"
	"github.kdc.capitalone.com/kbs316/packerd/restapi/operations/informational"
)

// This file is safe to edit. Once it exists it will not be overwritten

func configureFlags(api *operations.PackerdAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.PackerdAPI) http.Handler {

	fmt.Println("Starting the dispatcher")

	packerd.StartDispatcher(5)

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

		for key, _ := range packerd.BuildQ {
			var link = new(models.Link)
			link.Rel = "status"
			link.Href = fmt.Sprintf("/queue/%s", key)
			links = append(links, link)
		}
		return informational.NewGetQueueOK().WithPayload(links)
	})

	api.InformationalGetQueueByIDHandler = informational.GetQueueByIDHandlerFunc(func(params informational.GetQueueByIDParams) middleware.Responder {

		br, bqerr := packerd.BuildQ.LookUp(params.ID)
		if bqerr != nil {
			var err = new(models.Error)
			err.Code = 4
			*err.Message = "non-existiant build request"

			return informational.NewGetQueueByIDBadRequest().WithPayload(err)
		}

		return informational.NewGetQueueByIDOK().WithPayload(br)

	})

	// /queue/{id}/buildlog
	api.InformationalGetPackerLogByIDHandler = informational.GetPackerLogByIDHandlerFunc(func(params informational.GetPackerLogByIDParams) middleware.Responder {

		request, err := packerd.BuildQ.LookUp(params.ID)

		if err.Message != nil {
			return informational.NewGetPackerLogByIDBadRequest().WithPayload(err)
		}
		return informational.NewGetPackerLogByIDOK().WithPayload(request.Status)

	})

	api.InformationalGetQueueTestLogByIDHandler = informational.GetQueueTestLogByIDHandlerFunc(func(params informational.GetQueueTestLogByIDParams) middleware.Responder {
		return middleware.NotImplemented("operation informational.GetQueueTestLogByID has not yet been implemented")
	})

	api.CommandRunBuildHandler = command.RunBuildHandlerFunc(func(params command.RunBuildParams) middleware.Responder {

		params.Buildrequest.Status = "Pending"
		id, bqerr := packerd.BuildQ.Add(params.Buildrequest)
		if bqerr != nil {
			log.Println(bqerr)
			return command.NewRunBuildBadRequest().WithPayload(bqerr)
		}
		log.Printf("added new build request %s", id)

		params.Buildrequest.ID = id

		log.Printf("build request: %v", params.Buildrequest)

		dir, err := ioutil.TempDir("", "packerd")
		if err != nil {
			log.Println(err)
			bqerr.Code = 400
			*bqerr.Message = err.Error()
			return command.NewRunBuildBadRequest().WithPayload(bqerr)
		}
		params.Buildrequest.Localpath = dir
		log.Printf("got safe local working dir: %s", params.Buildrequest.Localpath)

		packerd.WorkQueue <- params.Buildrequest
		log.Println("pushed a build request")

		var link = new(models.Link)
		link.Rel = "status"
		link.Href = fmt.Sprintf("/queue/%s", id)

		return command.NewRunBuildAccepted().WithPayload(link)

	})

	api.ServerShutdown = func() {
		bqerr := packerd.BuildQ.Store("/var/cache/packerd-data.json")
		if bqerr != nil {
			log.Printf("failed to store json: %s", *bqerr.Message)

		}
	}

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
