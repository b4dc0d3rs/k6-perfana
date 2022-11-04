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
	Completed         bool     `json:"completed"`
	Duration          string   `json:"duration"`
	RampUp            string   `json:"rampUp"`
	TestEnvironment   string   `json:"testEnvironment"`
	SystemUnderTest   string   `json:"systemUnderTest"`
	Tags              []string `json:"tags"`
	Version           string   `json:"version"`
	TestRunId         string   `json:"testRunId"`
	Workload          string   `json:"workload"`
	CIBuildResultsUrl string   `json:"CIBuildResultsUrl"`
}

func (perfanaConfig *K6Perfana) StartPerfana() (map[string]string, error) {
	perfanaConfig.Completed = false

	variablesFailed := []string{}

	validateIfNilOrEmpty(variablesFailed, perfanaConfig.Duration, "PERFANA_URL")
	validateIfNilOrEmpty(variablesFailed, perfanaConfig.Duration, "PERFANA_TOKEN")

	perfanaConfig.Duration = os.Getenv("PERFANA_DURATION")
	variablesFailed = validateIfNilOrEmpty(variablesFailed, perfanaConfig.Duration, "PERFANA_DURATION")

	perfanaConfig.RampUp = os.Getenv("PERFANA_RAMPUP")
	variablesFailed = validateIfNilOrEmpty(variablesFailed, perfanaConfig.RampUp, "PERFANA_RAMPUP")

	perfanaConfig.TestEnvironment = os.Getenv("PERFANA_TEST_ENVIRONMENT")
	variablesFailed = validateIfNilOrEmpty(variablesFailed, perfanaConfig.TestEnvironment, "PERFANA_TEST_ENVIRONMENT")

	perfanaConfig.SystemUnderTest = os.Getenv("PERFANA_SYSTEM_UNDER_TEST")
	variablesFailed = validateIfNilOrEmpty(variablesFailed, perfanaConfig.SystemUnderTest, "PERFANA_SYSTEM_UNDER_TEST")

	perfanaConfig.Tags = strings.Split(os.Getenv("PERFANA_TAGS"), ",")
	variablesFailed = validateIfNilOrEmpty(variablesFailed, perfanaConfig.Tags, "PERFANA_TAGS")

	perfanaConfig.Version = os.Getenv("PERFANA_BUNDLE_VERSION")
	variablesFailed = validateIfNilOrEmpty(variablesFailed, perfanaConfig.Version, "PERFANA_BUNDLE_VERSION")

	perfanaConfig.TestRunId = os.Getenv("PERFANA_TEST_RUN_ID")
	variablesFailed = validateIfNilOrEmpty(variablesFailed, perfanaConfig.TestRunId, "PERFANA_TEST_RUN_ID")

	perfanaConfig.Workload = os.Getenv("PERFANA_WORKLOAD")
	variablesFailed = validateIfNilOrEmpty(variablesFailed, perfanaConfig.Workload, "PERFANA_WORKLOAD")

	perfanaConfig.CIBuildResultsUrl = os.Getenv("PERFANA_BUILD_URL")
	variablesFailed = validateIfNilOrEmpty(variablesFailed, perfanaConfig.CIBuildResultsUrl, "PERFANA_BUILD_URL")

	if len(variablesFailed) != 0 {
		return nil, fmt.Errorf("Required environment variables `%s` aren't valid", strings.Join(variablesFailed[:], ","))
	}

	go perfanaConfig.scheduledPolling()

	startResponse, startError := perfanaConfig.postToPerfana()
	return startResponse, startError
}

func (perfanaConfig *K6Perfana) scheduledPolling() {
	for perfanaConfig.Completed {
		time.Sleep(30 * time.Second)
		if perfanaConfig.Completed {
			perfanaConfig.postToPerfana()
		}
	}
}

func (perfanaConfig *K6Perfana) StopPerfana() (interface{}, error) {
	perfanaConfig.Completed = true
	stopResponse, stopError := perfanaConfig.postToPerfana()
	return stopResponse, stopError
}

func validateIfNilOrEmpty(failedVariables []string, variable interface{}, variableName string) []string {
	if variable == nil || variable == "" {
		return append(failedVariables, variableName)
	}
	return failedVariables
}

func (perfanaConfig *K6Perfana) postToPerfana() (map[string]string, error) {
	PERFANA_URL := os.Getenv("PERFANA_URL")
	PERFANA_TOKEN := os.Getenv("PERFANA_TOKEN")

	reqBody, err := json.Marshal(perfanaConfig)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", PERFANA_URL+"/api/test", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+PERFANA_TOKEN)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := make(map[string]string)
	response["perfanaPayload"] = bytes.NewBuffer(reqBody).String()
	response["statusCode"] = fmt.Sprint(resp.StatusCode)
	response["body"] = bytes.NewBuffer(body).String()

	return response, nil
}
