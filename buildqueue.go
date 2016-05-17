package packerd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	uuid "github.com/satori/go.uuid"
	"github.com/tompscanlan/packerd/models"
)

type BuildQueue map[string]*models.Buildrequest

var BuildQ = NewBuildQueue()

func NewBuildQueue() BuildQueue {
	bq := make(BuildQueue)
	return bq
}

func (bq *BuildQueue) LookUp(id string) (*models.Buildrequest, *models.Error) {

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

func (bq *BuildQueue) Add(br *models.Buildrequest) (string, *models.Error) {
	var id = uuid.NewV4()

	(*bq)[id.String()] = br

	return id.String(), nil

}

func (bq *BuildQueue) Update(id string, br *models.Buildrequest) (*models.Buildrequest, *models.Error) {
	request, ok := (*bq)[id]
	if !ok {
		msg := "no key for uuid: " + id
		return nil, &models.Error{Code: 100, Message: &msg}
	}
	request = br
	return request, nil
}

func (bq *BuildQueue) Delete(id string) *models.Error {

	uuid, err := uuid.FromString(id)
	if err != nil {
		msg := "invalid uuid string: " + err.Error()
		return &models.Error{Code: 100, Message: &msg}
	}

	delete(*bq, uuid.String())

	return nil
}

func (bq *BuildQueue) Store(filename string) *models.Error {
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

func (bq *BuildQueue) Load(filename string) *models.Error {
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
