package integration_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

const delay = 5 * time.Second

func TestMain(m *testing.M) {
	log.Printf("wait %s for service availability...", delay)
	time.Sleep(delay)

	status := godog.TestSuite{
		Name:                "integration",
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:    "pretty",
			Paths:     []string{"features"},
			Randomize: 0,
		},
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
