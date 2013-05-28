package riakpbc

import (
	"github.com/bmizerany/assert"
	"os/exec"
	"testing"
)

type Farm struct {
	Animal string `json:"animal" riak:"index"`
}

func setupIndexing(t *testing.T) {
	cmd := exec.Command("search-cmd", "install", "riakpbctestbucket")
	err := cmd.Run()
	if err != nil {
		t.Error(err.Error())
	}
}

func teardownIndexing(t *testing.T) {
	cmd := exec.Command("search-cmd", "uninstall", "riakpbctestbucket")
	err := cmd.Run()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestMapReduce(t *testing.T) {
	riak := setupConnection(t)
	setupData(t, riak)

	twoLevelQuery := "{\"inputs\":[[\"riakpbctestbucket\",\"testkey\"]],\"query\":[{\"map\":{\"language\":\"javascript\",\"keep\":false,\"name\":\"Riak.mapValuesJson\"}},{\"reduce\":{\"language\":\"javascript\",\"keep\":true,\"name\":\"Riak.reduceMax\"}}]}"
	reduced, err := riak.MapReduce(twoLevelQuery, "application/json")
	if err != nil {
		t.Error(err.Error())
	}
	assert.T(t, string(reduced) == "[{\"data\":\"is awesome!\"}]")

	teardownData(t, riak)
}

func TestIndex(t *testing.T) {
	riak := setupConnection(t)
	if _, err := riak.StoreObject("farm", "chicken", &Farm{Animal: "chicken"}); err != nil {
		t.Error(err.Error())
	}
	if _, err := riak.StoreObject("farm", "hen", &Farm{Animal: "hen"}); err != nil {
		t.Error(err.Error())
	}
	if _, err := riak.StoreObject("farm", "rooster", &Farm{Animal: "rooster"}); err != nil {
		t.Error(err.Error())
	}

	data, err := riak.Index("farm", "animal_bin", "chicken", "", "")
	if err != nil {
		t.Log("In order for this test to pass storage_backend must be set to riak_kv_eleveldb_backend in app.config")
		t.Error(err.Error())
	}
	assert.T(t, len(data.GetKeys()) > 0)

	if _, err := riak.DeleteObject("farm", "chicken"); err != nil {
		t.Error(err.Error())
	}
	if _, err := riak.DeleteObject("farm", "hen"); err != nil {
		t.Error(err.Error())
	}
	if _, err := riak.DeleteObject("farm", "rooster"); err != nil {
		t.Error(err.Error())
	}
}

//func TestSearch(t *testing.T) {
//	riak := setupConnection(t)
//	setupIndexing(t)
//	setupData(t, riak)
//
//	data, err := riak.Search("*awesome*", "data")
//	if err != nil {
//		t.Log("In order for this test to pass riak_search may need to be enabled in app.config")
//		t.Error(err.Error())
//	}
//	assert.T(t, data.GetNumFound() > 0)
//
//	teardownData(t, riak)
//	teardownIndexing(t)
//}