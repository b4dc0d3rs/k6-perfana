package k6perfana

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/perfana", new(K6Perfana))
}

type K6Perfana struct {
	Completed bool `json:"completed"`
	Duration string `json:"duration"`
	RampUp string `json:"rampUp"`
	TestEnvironment string `json:"testEnvironment"`
	SystemUnderTest string `json:"systemUnderTest"`
	Tags string `json:"tags"`
	Version string `json:"version"`
	TestRunId string `json:"testRunId"`
	Workload string `json:"workload"`
	CIBuildResultsUrl string `json:"CIBuildResultsUrl"`
}

func (perfanaSetup *K6Perfana) StartPerfana() (interface{}, error) {
	perfanaSetup.Completed = false
	perfanaSetup.Duration = os.Getenv("PERFANA_DURATION")
	perfanaSetup.RampUp = os.Getenv("PERFANA_RAMPUP")
	perfanaSetup.TestEnvironment = os.Getenv("PERFANA_TEST_ENVIRONMENT")
	perfanaSetup.SystemUnderTest = os.Getenv("PERFANA_SYSTEM_UNDER_TEST")
	perfanaSetup.Tags = os.Getenv("PERFANA_TAGS")
	perfanaSetup.Version = os.Getenv("PERFANA_VERSION")
	perfanaSetup.TestRunId = os.Getenv("PERFANA_TEST_RUN_ID")
	perfanaSetup.Workload = os.Getenv("PERFANA_WORKLOAD")
	perfanaSetup.CIBuildResultsUrl = os.Getenv("PERFANA_CI_BUILD_URL")

	return perfanaSetup.postToPerfana(), nil
}

func (perfanaConfig *K6Perfana) StopPerfana() (interface{}, error) {
	perfanaConfig.Completed = true
	return perfanaConfig.postToPerfana(), nil
}

func (perfanaConfig *K6Perfana) postToPerfana() []byte {
	PERFANA_URL := os.Getenv("PERFANA_URL")
	PERFANA_API_TOKEN := os.Getenv("PERFANA_API_TOKEN")

	reqBody, err := json.Marshal(perfanaConfig)
	if err != nil {
		print(err)
	}
	request, err := http.NewRequest("POST", PERFANA_URL + "/api/test", bytes.NewBuffer(reqBody))
	if err != nil {
		print(err)
	}
	request.Header.Set("Authorization", "Bearer " + PERFANA_API_TOKEN)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}

	return body;
}
