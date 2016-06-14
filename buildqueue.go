package packerd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"

	uuid "github.com/satori/go.uuid"
	"github.com/tompscanlan/packerd/models"
)

// BuildMap provides mapping from uuid to a specific build request
type BuildMap map[string]*models.Buildrequest

//Builds is the global list of all known builds
var Builds = NewBuildMap()

// NewBuildMap creates a new build map
func NewBuildMap() BuildMap {
	bq := make(BuildMap)
	return bq
}

// LookUp looks up a build when given a uuid as a string
func (bq *BuildMap) LookUp(id string) (*models.Buildrequest, *models.Error) {

	uuid, err := uuid.FromString(id)
	if err != nil {
		msg := "invalid uuid string: " + err.Error()
		return nil, &models.Error{Code: 100, Message: &msg}
	}

	request, ok := (*bq)[uuid.String()]
	if !ok {
		msg := "no key for uuid: " + uuid.String()
		return nil, &models.Error{Code: 100, Message: &msg}
	}
	return request, nil
}

// Add a new build request to the map
func (bq *BuildMap) Add(br *models.Buildrequest) (string, *models.Error) {
	var id = uuid.NewV4()

	// if lookup of random new uuid succeeds, we're adding a duplicate
	_, ok := (*bq)[id.String()]
	if ok {
		msg := "duplicate add for uuid: " + id.String()
		return id.String(), &models.Error{Code: 100, Message: &msg}
	}

	br.ID = id.String()
	(*bq)[id.String()] = br

	log.WithFields(log.Fields{
		"function": "BuildMap.Add",
	}).Debugf("Adding: %v", br)

	return id.String(), nil
}

// Delete a build
func (bq *BuildMap) Delete(id string) *models.Error {

	uuid, err := uuid.FromString(id)
	if err != nil {
		msg := "invalid uuid string: " + err.Error()
		return &models.Error{Code: 100, Message: &msg}
	}

	delete(*bq, uuid.String())

	return nil
}

// Store all the builds as JSON in a file
func (bq *BuildMap) Store(filename string) *models.Error {
	b, err := json.Marshal(bq)
	if err != nil {
		//panic(err)
		msg := fmt.Sprintf("failed to marshal to json: %s", err.Error())
		return &models.Error{Code: 100, Message: &msg}
	}

	err = ioutil.WriteFile(filename, b, 0644)
	if err != nil {
		//panic(err)
		msg := fmt.Sprintf("failed write file: %s", err.Error())
		return &models.Error{Code: 100, Message: &msg}
	}

	return nil
}

// Load builds from a JSON file
func (bq *BuildMap) Load(filename string) *models.Error {
	blob, err := ioutil.ReadFile(filename)
	if err != nil {
		//panic(err)
		msg := fmt.Sprintf("failed read file: %s", err.Error())
		return &models.Error{Code: 100, Message: &msg}
	}

	err = json.Unmarshal(blob, &bq)
	if err != nil {
		//panic(err)
		msg := fmt.Sprintf("failed to unmarshal: %s", err.Error())
		return &models.Error{Code: 100, Message: &msg}
	}

	return nil
}
