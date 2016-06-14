package packerd

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/satori/go.uuid"

	"github.com/tompscanlan/packerd/models"
)

// BuildResponseMap is a mapping of uuid as a string to build responses
type BuildResponseMap map[string][]*models.Buildresponse

// BuildResponses is the global list of all responses to build requests
var BuildResponses = NewBuildResponseMap()

// NewBuildResponseMap creates a new map of uuid to build response
func NewBuildResponseMap() BuildResponseMap {
	br := make(BuildResponseMap)
	return br
}

// LookupResponses a response based on the build request uuid
func (buildresponses *BuildResponseMap) LookupResponses(requestid string) (error, []*models.Buildresponse) {
	l := log.WithFields(log.Fields{
		"function": "BuildResponse.LookupResponses",
	})

	responses, ok := (*buildresponses)[requestid]
	if ok {
		l.Debugf("got %v", responses)
		return nil, responses
	}
	l.Debugf("lookup missed for id %s", requestid)
	return errors.New(fmt.Sprintf("lookup missed for id %s", requestid)), nil
}

//  LookupResponse returns a specific build stage when called with a uuid of the build request and a stage name
func (buildresponses *BuildResponseMap) LookupResponse(id string, buildnumber int) (error, *models.Buildresponse) {

	err, responses := buildresponses.LookupResponses(id)
	if err != nil {
		log.Error(err)
		return err, nil
	}

	if buildnumber < 0 || buildnumber >= len(responses) {
		log.Error("buildnumber %d, out of range: 0-%d", buildnumber, len(responses))
		return errors.New("build number out of range"), nil
	}

	return nil, responses[buildnumber]
}

//  LookupStage returns a specific build stage when called with a uuid of the build request and a stage name
func (buildresponses *BuildResponseMap) LookupStage(id string, buildnumber int, stagename string) (error, *models.Buildstage) {

	err, response := buildresponses.LookupResponse(id, buildnumber)
	if err != nil {
		log.Error(err)
		return err, nil
	}

	for _, stage := range response.Buildstages {
		if stage.Name == stagename {
			return nil, stage
		}
	}

	return errors.New("No stage by that name"), nil
}

// Add a new build response for a given build request id
func (buildresponses *BuildResponseMap) Add(requestid string, buildresponse *models.Buildresponse) (int, *models.Buildresponse) {
	var respid = uuid.NewV4().String()
	l := log.WithFields(log.Fields{
		"function": "BuildResponseMap.Add",
	})
	buildresponse.ID = respid
	buildresponse.Buildrequestid = requestid
	l.Debugf("Adding %s", buildresponse)

	responses := (*buildresponses)[requestid]
	buildnumber := len(responses)
	responses = append(responses, buildresponse)
	(*buildresponses)[requestid] = responses

	return buildnumber, buildresponse

}

// AddStage puts a new stage onto an existing build response
func (buildresponses *BuildResponseMap) AddStage(id string, buildnumber int, stage *models.Buildstage) (error, *models.Buildstage) {
	log.Debugf("AddStage %s to id: %s, buildnumber: %d", stage, id, buildnumber)

	err, responses := buildresponses.LookupResponses(id)
	if err != nil {
		log.Error(err)
		return err, nil
	}

	if responses == nil {
		return errors.New("must add a response before adding a stage"), nil
	}
	log.Infof("buildnumber: %d, len of reponses: %d", buildnumber, len(responses))

	if buildnumber < 0 || buildnumber >= len(responses) {
		log.Error("buildnumber %d, out of range: 0-%d", buildnumber, len(responses))
		return errors.New("build number out of range"), nil
	}

	responses[buildnumber].Buildstages = append(responses[buildnumber].Buildstages, stage)

	return nil, stage
}

func (b *BuildResponseMap) String() string {
	var response []string

	for k, _ := range *b {
		responses := (*b)[k]
		for _, j := range responses {

			response = append(response, j.String())
		}
	}
	return strings.Join(response, ",")
}
