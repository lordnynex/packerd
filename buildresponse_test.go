package packerd

import (
	"testing"

	"github.com/tompscanlan/packerd/models"
)

func TestAddSingleBuildResponse(t *testing.T) {

	// add a build request first
	breq := new(models.Buildrequest)
	id, err := Builds.Add(breq)

	if err != nil {
		t.Errorf("unexpected err? %v", err)
	}

	bresp := new(models.Buildresponse)

	buildnumber, dupRespPointer := BuildResponses.Add(id, bresp)

	if dupRespPointer == nil {
		t.Error("returned pointer should not be nil")
	}

	merr, responses := BuildResponses.LookupResponses(id)
	if merr != nil {
		t.Error(merr.Error())
	}
	if responses == nil {
		t.Error("got back nil responses")
	}

	if responses[buildnumber].Buildrequestid != breq.ID {
		t.Error("ID's don't match")
	}

	if responses == nil {
		t.Errorf("Failed to look up just added response")
	}
}

func TestAddManyBuildResponse(t *testing.T) {

	// add a build request first
	breq := new(models.Buildrequest)
	id, err := Builds.Add(breq)

	if err != nil {
		t.Errorf("unexpected err? %v", err)
	}

	// add several responses for a single buildrequest
	count := 100
	for i := 1; i <= count; i++ {
		bresp := new(models.Buildresponse)

		buildnumber, dupRespPointer := BuildResponses.Add(id, bresp)
		merr, responses := BuildResponses.LookupResponses(id)
		if merr != nil {
			t.Error(merr.Error())
		}
		if responses[buildnumber].ID != bresp.ID {
			t.Error("ID's don't match")
		}

		if dupRespPointer == nil {
			t.Error("returned pointer should not be nil")
		}
	}

	merr, responses := BuildResponses.LookupResponses(id)
	if merr != nil {
		t.Error(merr.Error())
	}

	if len(responses) != count {
		t.Errorf("We added %d responses, but looked up %d", count, len(responses))
	}
}

func TestAddManyBuildResponsesForManyRequests(t *testing.T) {

	countReqs := 100

	for i := 1; i <= countReqs; i++ {

		// add a build request first
		breq := new(models.Buildrequest)
		id, err := Builds.Add(breq)

		if err != nil {
			t.Errorf("unexpected err? %v", err)
		}

		// add several responses for a single buildrequest
		count := 100
		for i := 1; i <= count; i++ {
			bresp := new(models.Buildresponse)

			buildnumber, dupRespPointer := BuildResponses.Add(id, bresp)
			merr, responses := BuildResponses.LookupResponses(id)
			if merr != nil {
				t.Error(merr.Error())
			}
			if responses[buildnumber].ID != bresp.ID {
				t.Error("ID's don't match")
			}
			if dupRespPointer == nil {
				t.Error("returned pointer should not be nil")
			}
		}

		merr, responses := BuildResponses.LookupResponses(id)
		if merr != nil {
			t.Error(merr.Error())
		}

		if len(responses) != count {
			t.Errorf("We added %d responses, but looked up %d", count, len(responses))
		}

	}
}

func TestAddBuildResponseStageToBadRequest(t *testing.T) {
	stage := new(models.Buildstage)
	stage.Name = "testing"

	err, resp := BuildResponses.AddStage("badid", 0, stage)
	if err == nil {
		t.Errorf("This should have failed")
	}
	if resp != nil {
		t.Errorf("This should have been nil")
	}
}

func TestAddBuildResponseStageOutOfBounds(t *testing.T) {
	breq := new(models.Buildrequest)
	id, _ := Builds.Add(breq)
	bresp := new(models.Buildresponse)
	buildnumber, _ := BuildResponses.Add(id, bresp)

	stage := new(models.Buildstage)
	stage.Name = "testing"

	nerr, resp := BuildResponses.AddStage(id, -1, stage)
	if nerr == nil {
		t.Errorf("This should have failed")
	}
	if resp != nil {
		t.Errorf("This should have been nil")
	}

	nerr, resp = BuildResponses.AddStage(id, buildnumber+1, stage)
	if nerr == nil {
		t.Errorf("This should have failed")
	}
	if resp != nil {
		t.Errorf("This should have been nil")
	}
}

func TestAddBuildResponseStageGood(t *testing.T) {
	breq := new(models.Buildrequest)
	id, _ := Builds.Add(breq)
	bresp := new(models.Buildresponse)
	buildnumber, _ := BuildResponses.Add(id, bresp)

	stage := new(models.Buildstage)
	stage.Name = "testing"

	nerr, resp := BuildResponses.AddStage(id, buildnumber, stage)
	if nerr != nil {
		t.Errorf("This should not have failed: %v", nerr)
	}
	if resp == nil {
		t.Errorf("This should not have been nil")
	}
}
