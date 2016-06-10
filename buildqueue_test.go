package packerd

import (
	//"fmt"
	"io/ioutil"
	"testing"

	"github.kdc.capitalone.com/kbs316/packerd/models"
)

func TestNewBuildQueue(t *testing.T) {

	bq := NewBuildQueue()
	br := new(models.Buildrequest)

	id, _ := bq.Add(br)
	if n := len(bq); n != 1 {
		t.Errorf("expected 1 key, found %d", n)
	}

	if n := len(bq); n != 1 {
		t.Errorf("expected 1 key, found %d", n)
	}
	bq.Delete(id)
	if n := len(bq); n != 0 {
		t.Errorf("expected 0 key, found %d", n)
	}

	bq.Delete(id)
	if n := len(bq); n != 0 {
		t.Errorf("expected 0 key, found %d", n)
	}

	badId := "asd123"
	err := bq.Delete(badId)
	if err == nil {
		t.Errorf("expected to hit an error, but didn't")
	}
}

func TestNewLookupBuildQueue(t *testing.T) {

	status := "New Test"
	giturl := "http://github.kdc.capitalone.com/kbs316/packerd"
	badId := "asd123"
	nonExistId := "b7aa6044-bbde-415c-89e2-7833d4e544dc"

	bq := NewBuildQueue()
	br := new(models.Buildrequest)
	br.Status = status
	br.Giturl = &giturl

	id, err := bq.Add(br)
	if err != nil {
		t.Errorf("got an unexpected error: %q", err)
	}

	if n := len(bq); n != 1 {
		t.Errorf("expected 1 key, found %d", n)
	}

	found, err := bq.LookUp(id)
	if err != nil {
		t.Errorf("got an unexpected error: %q", err)
	}
	if found.ID != id {
		t.Errorf("lookup was for %q, but found %q", id, found.ID)
	}
	if *found.Giturl != giturl {
		t.Errorf("lookup giturl should have been %q, but found %q", id, found.Giturl)
	}

	_, err = bq.LookUp(badId)
	if err == nil {
		t.Errorf("missing an expected error")
	}

	err = bq.Delete(nonExistId)
	if err != nil {
		t.Errorf("got an unexpected error: %q", err)
	}

	found, err = bq.LookUp(nonExistId)
	if err == nil {
		t.Errorf("missing an expected error")
	}

}

func TestStoreLoadBuildQueue(t *testing.T) {

	bq := NewBuildQueue()
	var ids []string
	recs := 1000

	for n := 1; n <= recs; n++ {
		br := new(models.Buildrequest)
		id, _ := bq.Add(br)
		ids = append(ids, id)

		if l := len(bq); n != l {
			t.Errorf("expected %d key, found %d", n, l)
		}
	}

	dir, err := ioutil.TempDir("", "packerd_test")
	if err != nil {
		t.Errorf("failed to find a temp dir: %v", err)
	}

	bqerr := bq.Store(dir + "/marshal.json")
	if bqerr != nil {
		t.Errorf("failed to store json: %s", *bqerr.Message)
	}

	for _, id := range ids {
		bq.Delete(id)
	}
	if l := len(bq); l != 0 {
		t.Errorf("expected %d key, found %d", 0, l)
	}

	bqerr = bq.Load(dir + "/marshal.json")
	if bqerr != nil {
		t.Errorf("failed to load json: %s", *bqerr.Message)
	}

	if l := len(bq); l != recs {
		t.Errorf("expected %d key, found %d", recs, l)
	}
}

func BenchmarkLotsOfRequests(b *testing.B) {

	bq := NewBuildQueue()
	var ids []string

	for n := 0; n < b.N; n++ {
		br := new(models.Buildrequest)
		id, _ := bq.Add(br)
		ids = append(ids, id)
	}
}
func BenchmarkDeleteLotsOfRequests(b *testing.B) {

	bq := NewBuildQueue()
	var ids []string

	for n := 0; n < b.N; n++ {
		br := new(models.Buildrequest)
		id, _ := bq.Add(br)
		ids = append(ids, id)
	}

	b.ResetTimer()
	for _, id := range ids {
		bq.Delete(id)

	}
}
