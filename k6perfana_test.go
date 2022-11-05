package k6perfana

import "testing"

func TestIntegrationPerfana(t *testing.T) {
	perfanaConfig := new(K6Perfana)
	perfanaConfig.Completed = false
	perfanaConfig.Duration = "60"
	perfanaConfig.RampUp = "10"
	perfanaConfig.TestEnvironment = "experimental"
	perfanaConfig.SystemUnderTest = "core"
	perfanaConfig.Tags = []string{"k6"}
	perfanaConfig.Version = "v0.1.0"
	perfanaConfig.TestRunId = "#1"
	perfanaConfig.Workload = "load-10"
	perfanaConfig.CIBuildResultsUrl = "https://github.com/b4dc0d3rs/k6-perfana"

	startResponse, error := perfanaConfig.StartPerfana()

	if error != nil {
		t.Errorf("Failed, got error %s", error.Error())
	}

	if startResponse["statusCode"] != "200" {
		t.Logf(startResponse["config"])
		t.Logf(startResponse["body"])
		t.Errorf("Failed, got status code %s", startResponse["statusCode"])
	}

	stopResponse, error := perfanaConfig.StopPerfana()

	if error != nil {
		t.Errorf("Failed, got error %s", error.Error())
	}

	if stopResponse["statusCode"] != "200" {
		t.Logf(stopResponse["config"])
		t.Logf(stopResponse["body"])
		t.Errorf("Failed, got status code %s", stopResponse["statusCode"])
	}
}
