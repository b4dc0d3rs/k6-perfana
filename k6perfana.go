package k6perfana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/k6perfana", new(K6Perfana))
}

type K6Perfana struct {
	Completed bool `json:"completed"`
	Duration string `json:"duration"`
	RampUp string `json:"rampUp"`
	TestEnvironment string `json:"testEnvironment"`
	SystemUnderTest string `json:"systemUnderTest"`
	Tags []string `json:"tags"`
	Version string `json:"version"`
	TestRunId string `json:"testRunId"`
	Workload string `json:"workload"`
	CIBuildResultsUrl string `json:"CIBuildResultsUrl"`
}

func (perfanaConfig *K6Perfana) StartPerfana() (interface{}, error) {
	perfanaConfig.Completed = false

	perfanaConfig.Duration = os.Getenv("PERFANA_DURATION")
	if validationError := validateIfNilOrEmpty(perfanaConfig.Duration, "PERFANA_DURATION"); validationError != nil {
		return nil, validationError;
	}

	perfanaConfig.RampUp = os.Getenv("PERFANA_RAMPUP")
	if validationError := validateIfNilOrEmpty(perfanaConfig.RampUp, "PERFANA_RAMPUP"); validationError != nil {
		return nil, validationError;
	}

	perfanaConfig.TestEnvironment = os.Getenv("PERFANA_TEST_ENVIRONMENT")
	if validationError := validateIfNilOrEmpty(perfanaConfig.TestEnvironment, "PERFANA_TEST_ENVIRONMENT"); validationError != nil {
		return nil, validationError;
	}

	perfanaConfig.SystemUnderTest = os.Getenv("PERFANA_SYSTEM_UNDER_TEST")
	if validationError := validateIfNilOrEmpty(perfanaConfig.SystemUnderTest, "PERFANA_SYSTEM_UNDER_TEST"); validationError != nil {
		return nil, validationError;
	}

	perfanaConfig.Tags = strings.Split(os.Getenv("PERFANA_TAGS"), ",");
	if validationError := validateIfNilOrEmpty(perfanaConfig.Tags, "PERFANA_TAGS"); validationError != nil {
		return nil, validationError;
	}

	perfanaConfig.Version = os.Getenv("PERFANA_BUNDLE_VERSION")
	if validationError := validateIfNilOrEmpty(perfanaConfig.Version, "PERFANA_BUNDLE_VERSION"); validationError != nil {
		return nil, validationError;
	}

	perfanaConfig.TestRunId = os.Getenv("PERFANA_TEST_RUN_ID")
	if validationError := validateIfNilOrEmpty(perfanaConfig.TestRunId, "PERFANA_TEST_RUN_ID"); validationError != nil {
		return nil, validationError;
	}

	perfanaConfig.Workload = os.Getenv("PERFANA_WORKLOAD")
	if validationError := validateIfNilOrEmpty(perfanaConfig.Workload, "PERFANA_WORKLOAD"); validationError != nil {
		return nil, validationError;
	}

	perfanaConfig.CIBuildResultsUrl = os.Getenv("PERFANA_BUILD_URL")
	if validationError := validateIfNilOrEmpty(perfanaConfig.CIBuildResultsUrl, "PERFANA_BUILD_URL"); validationError != nil {
		return nil, validationError;
	}

	go perfanaConfig.scheduledPolling();

	return perfanaConfig.postToPerfana()
}

func (perfanaConfig *K6Perfana) scheduledPolling() {
	for perfanaConfig.Completed {
		time.Sleep(30 * time.Second)
		perfanaConfig.postToPerfana()
	}
}

func (perfanaConfig *K6Perfana) StopPerfana() (interface{}, error) {
	perfanaConfig.Completed = true
	return perfanaConfig.postToPerfana()
}

func validateIfNilOrEmpty(variable interface{}, variableName string) error {
	if variable == nil || variable == "" {
		return fmt.Errorf("Required environment variable `%s` is empty", variableName)
	}
	return nil
}

func (perfanaConfig *K6Perfana) postToPerfana() ([]byte, error) {
	PERFANA_URL := os.Getenv("PERFANA_URL")
	PERFANA_API_TOKEN := os.Getenv("PERFANA_API_TOKEN")

	reqBody, err := json.Marshal(perfanaConfig)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", PERFANA_URL + "/api/test", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer " + PERFANA_API_TOKEN)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, fmt.Errorf("Failed to login to Perfana, expected status was is 200 or 201, but got " + fmt.Sprint(resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil;
}
