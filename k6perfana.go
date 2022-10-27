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

func (perfanaSetup *K6Perfana) StartPerfana() (interface{}, error) {
	fmt.Printf("Starting perfana")
	perfanaSetup.Completed = false

	perfanaSetup.Duration = os.Getenv("PERFANA_DURATION")
	if validationError := validateIfNilOrEmpty(perfanaSetup.Duration, "PERFANA_DURATION"); validationError != nil {
		return nil, validationError;
	}

	perfanaSetup.RampUp = os.Getenv("PERFANA_RAMPUP")
	if validationError := validateIfNilOrEmpty(perfanaSetup.RampUp, "PERFANA_RAMPUP"); validationError != nil {
		return nil, validationError;
	}

	perfanaSetup.TestEnvironment = os.Getenv("PERFANA_TEST_ENVIRONMENT")
	if validationError := validateIfNilOrEmpty(perfanaSetup.TestEnvironment, "PERFANA_TEST_ENVIRONMENT"); validationError != nil {
		return nil, validationError;
	}

	perfanaSetup.SystemUnderTest = os.Getenv("PERFANA_SYSTEM_UNDER_TEST")
	if validationError := validateIfNilOrEmpty(perfanaSetup.SystemUnderTest, "PERFANA_SYSTEM_UNDER_TEST"); validationError != nil {
		return nil, validationError;
	}

	perfanaSetup.Tags = strings.Split(os.Getenv("PERFANA_TAGS"), ",");
	if validationError := validateIfNilOrEmpty(perfanaSetup.Tags, "PERFANA_TAGS"); validationError != nil {
		return nil, validationError;
	}

	perfanaSetup.Version = os.Getenv("PERFANA_VERSION")
	if validationError := validateIfNilOrEmpty(perfanaSetup.Version, "PERFANA_VERSION"); validationError != nil {
		return nil, validationError;
	}

	perfanaSetup.TestRunId = os.Getenv("PERFANA_TEST_RUN_ID")
	if validationError := validateIfNilOrEmpty(perfanaSetup.TestRunId, "PERFANA_TEST_RUN_ID"); validationError != nil {
		return nil, validationError;
	}

	perfanaSetup.Workload = os.Getenv("PERFANA_WORKLOAD")
	if validationError := validateIfNilOrEmpty(perfanaSetup.Workload, "PERFANA_WORKLOAD"); validationError != nil {
		return nil, validationError;
	}

	perfanaSetup.CIBuildResultsUrl = os.Getenv("PERFANA_CI_BUILD_URL")
	if validationError := validateIfNilOrEmpty(perfanaSetup.CIBuildResultsUrl, "PERFANA_CI_BUILD_URL"); validationError != nil {
		return nil, validationError;
	}

	return perfanaSetup.postToPerfana()
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
		return nil, fmt.Errorf("Expected status was is 200 or 201, but got " + fmt.Sprint(resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil;
}
