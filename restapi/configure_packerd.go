package restapi

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	"github.com/tompscanlan/packerd"

	"github.com/tompscanlan/packerd/models"
	"github.com/tompscanlan/packerd/restapi/operations"
	"github.com/tompscanlan/packerd/restapi/operations/command"
	"github.com/tompscanlan/packerd/restapi/operations/informational"
)

// This file is safe to edit. Once it exists it will not be overwritten

func configureFlags(api *operations.PackerdAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }

	var LogFlags struct {
		LogLevel func(string) `short:"l" long:"log-level" description:"debug, info, warning, error, fatal and panic"`
		NoColor  func(bool)   `long:"nocolor" description:"set to disable color logging" default:"false" optional:"true" optional-value:"true" `
	}

	LogFlags.NoColor = func(nocolor bool) {
		log.SetFormatter(&log.TextFormatter{
			ForceColors:   !nocolor,
			DisableColors: nocolor,
		})

		return
	}

	LogFlags.LogLevel = func(s string) {
		switch s {
		case "debug":
			log.SetLevel(log.DebugLevel)
		default:
			log.SetLevel(packerd.DefaultLogLevel)
		}
		return
	}
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		{"log", "logging configuration", &LogFlags},
	}
}

func configureAPI(api *operations.PackerdAPI) http.Handler {
	fmt.Println("Starting the dispatcher")

	packerd.StartDispatcher(packerd.WorkerCount)

	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// s.api.Logger = log.Printf

	log.SetOutput(os.Stdout)
	api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.InformationalGetHealthHandler = informational.GetHealthHandlerFunc(func(params informational.GetHealthParams) middleware.Responder {
		h := GetHealth()

		health := informational.NewGetHealthOK()
		health.SetPayload(h)
		log.Debugf("%v", health.Payload)

		return health
	})

	api.InformationalGetBuildStageByNameHandler = informational.GetBuildStageByNameHandlerFunc(func(params informational.GetBuildStageByNameParams) middleware.Responder {
		//return middleware.NotImplemented("operation informational.GetBuildStageByName has not yet been implemented")
		err, stage := packerd.BuildResponses.LookupStage(params.ID, int(params.Buildnumber), params.Stagename)
		if err != nil {
			perr := models.NewError(4, "no stage with that name")
			return informational.NewGetBuildStageByNameBadRequest().WithPayload(perr)
		}
		return informational.NewGetBuildStageByNameOK().WithPayload(stage)
	})
	api.InformationalGetBuildListByIDHandler = informational.GetBuildListByIDHandlerFunc(func(params informational.GetBuildListByIDParams) middleware.Responder {
		br, bqerr := packerd.Builds.LookUp(params.ID)
		if bqerr != nil {
			var err = models.NewError(4, "non-existiant build request")
			return informational.NewGetBuildListByIDBadRequest().WithPayload(err)
		}
		var links []*models.Link

		cleanbr := br.Sanitize()

		//		responses := packerd.BuildResponses.LookupResponses(br.ID)

		var buildResponseLink = new(models.Link)
		buildResponseLink.Href = fmt.Sprintf("/build/responses/%s", br.ID)
		buildResponseLink.Rel = "buildresponses"
		links = append(links, buildResponseLink)

		err, responses := packerd.BuildResponses.LookupResponses(br.ID)
		if err != nil {
			log.Error(err)
			var perr = models.NewError(4, "non-existiant build request")

			return informational.NewGetBuildListByIDBadRequest().WithPayload(perr)
		}
		for i, _ := range responses {
			link := new(models.Link)
			link.Href = fmt.Sprintf("/build/responses/%s/%d", br.ID, i)
			link.Rel = "buildresponse"
			links = append(links, link)

			link = new(models.Link)
			link.Href = fmt.Sprintf("/build/responses/%s/%d/stages", br.ID, i)
			link.Rel = "buildstages"
			links = append(links, link)
		}
		cleanbr.Responselinks = links

		return informational.NewGetBuildListByIDOK().WithPayload(cleanbr)
	})
	api.InformationalGetBuildListHandler = informational.GetBuildListHandlerFunc(func(params informational.GetBuildListParams) middleware.Responder {
		var links []*models.Link

		log.Debugf("number of builds: %d", len(packerd.Builds))
		if len(packerd.Builds) <= 0 {
			return informational.NewGetBuildListNoContent()
		}

		for key, _ := range packerd.Builds {
			var link = new(models.Link)
			link.Rel = "buildrequest"
			link.Href = fmt.Sprintf("/build/queue/%s", key)
			links = append(links, link)
		}
		return informational.NewGetBuildListOK().WithPayload(links)
	})
	api.InformationalGetBuildStagesByIDBuildNumberHandler = informational.GetBuildStagesByIDBuildNumberHandlerFunc(func(params informational.GetBuildStagesByIDBuildNumberParams) middleware.Responder {
		var links []*models.Link

		err, response := packerd.BuildResponses.LookupResponse(params.ID, int(params.Buildnumber))
		if err != nil {
			log.Error(err)
			return informational.NewGetBuildResponseByIDAndBuildNumberBadRequest()
		}

		for _, stage := range response.Buildstages {
			var link = new(models.Link)
			link.Rel = "buildstage"
			link.Href = fmt.Sprintf("/build/responses/%s/%d/stages/%s", params.ID, params.Buildnumber, stage.Name)
			links = append(links, link)
		}
		return informational.NewGetBuildStagesByIDBuildNumberOK().WithPayload(links)
		//return middleware.NotImplemented("operation informational.GetBuildStagesByIDBuildNumber has not yet been implemented")
	})

	api.InformationalGetBuildResponseByIDAndBuildNumberHandler = informational.GetBuildResponseByIDAndBuildNumberHandlerFunc(func(params informational.GetBuildResponseByIDAndBuildNumberParams) middleware.Responder {
		err, response := packerd.BuildResponses.LookupResponse(params.ID, int(params.Buildnumber))
		if err != nil {
			log.Error(err)
			return informational.NewGetBuildResponseByIDAndBuildNumberBadRequest()
		}
		return informational.NewGetBuildResponseByIDAndBuildNumberOK().WithPayload(response)

	})
	api.InformationalGetBuildResponseByIDHandler = informational.GetBuildResponseByIDHandlerFunc(func(params informational.GetBuildResponseByIDParams) middleware.Responder {
		err, responses := packerd.BuildResponses.LookupResponses(params.ID)
		if err != nil {
			return informational.NewGetBuildResponseByIDBadRequest()
		}
		if responses == nil {
			perr := models.NewError(4, "non-existiant build request")
			return informational.NewGetBuildResponseByIDBadRequest().WithPayload(perr)
		}
		var links []*models.Link
		log.Debugf("making links for %v", responses)
		for index, _ := range responses {
			var link = new(models.Link)
			link.Rel = "buildresponse"
			link.Href = fmt.Sprintf("/build/responses/%s/%d", params.ID, index)
			links = append(links, link)
		}
		return informational.NewGetBuildResponseByIDOK().WithPayload(links)

	})
	api.CommandRunBuildHandler = command.RunBuildHandlerFunc(func(params command.RunBuildParams) middleware.Responder {

		// add the request to our list of builds
		id, bqerr := packerd.Builds.Add(params.Buildrequest)
		if bqerr != nil {
			log.Error(bqerr)
			return command.NewRunBuildBadRequest().WithPayload(bqerr)
		}
		params.Buildrequest.ID = id

		dir, err := ioutil.TempDir("", "packerd")
		if err != nil {

			log.Error(err)
			bqerr = models.NewError(5, err.Error())
			return command.NewRunBuildBadRequest().WithPayload(bqerr)
		}
		params.Buildrequest.Localpath = dir
		log.Debugf("Using %s as local working dir", dir)

		// signal the dispatcher to send build to the worker
		packerd.BuildSendWorkerChan <- params.Buildrequest

		var links []*models.Link

		// offer a link to the build, and the responses
		var link = new(models.Link)
		link.Rel = "buildrequest"
		link.Href = fmt.Sprintf("/build/queue/%s", id)
		links = append(links, link)

		link = new(models.Link)
		link.Href = fmt.Sprintf("/build/responses/%s", id)
		link.Rel = "buildresponses"
		links = append(links, link)

		link = new(models.Link)
		link.Rel = "health"
		link.Href = fmt.Sprintf("/health")
		links = append(links, link)

		return command.NewRunBuildAccepted().WithPayload(links)
	})

	api.ServerShutdown = func() {
		bqerr := packerd.Builds.Store("/var/cache/packerd-data.json")
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
	return handler
}
