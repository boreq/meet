package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/boreq/errors"
	"github.com/boreq/hydro/internal/fixture"
	"github.com/boreq/hydro/internal/wire"
	"github.com/boreq/hydro/ports/http/hydro"
	"github.com/stretchr/testify/require"
)

func TestHydro(t *testing.T) {
	s, cleanup := fixture.Service(t)
	defer cleanup()

	controllers, err := apiHydroListControllers(s)
	require.NoError(t, err)
	require.Empty(t, controllers)

	controllerAddress := "controller-address"

	err = apiHydroAddController(s, hydro.AddControllerRequest{
		Address: controllerAddress,
	})
	require.NoError(t, err)

	controllers, err = apiHydroListControllers(s)
	require.NoError(t, err)
	require.Len(t, controllers, 1)
	require.NotEmpty(t, controllers[0].UUID)
	require.Equal(t, controllerAddress, controllers[0].Address)
}

const (
	apiUrlHydroControllers = "/hydro/controllers"
)

func apiHydroAddController(s wire.ComponentTestService, r hydro.AddControllerRequest) error {
	response, err := apiPost(s, apiUrlHydroControllers, r)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("invalid status: '%s'", response.Status)
	}

	var controllers []hydro.Controller

	if err := json.NewDecoder(response.Body).Decode(&controllers); err != nil {
		return errors.Wrap(err, "json decoding failed")
	}

	return nil
}

func apiHydroListControllers(s wire.ComponentTestService) ([]hydro.Controller, error) {
	response, err := apiGet(s, apiUrlHydroControllers)
	if err != nil {
		return nil, err
	}

	var controllers []hydro.Controller

	if err := json.NewDecoder(response.Body).Decode(&controllers); err != nil {
		return nil, errors.Wrap(err, "json decoding failed")
	}

	return controllers, nil
}

func apiGet(s wire.ComponentTestService, url string) (*http.Response, error) {
	url = apiUrl(s, url)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create a request")
	}

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "client do failed")
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status: '%s'", response.Status)
	}

	return response, nil
}

func apiPost(s wire.ComponentTestService, url string, body interface{}) (*http.Response, error) {
	url = apiUrl(s, url)

	j, err := json.Marshal(body)
	if err != nil {
		return nil, errors.Wrap(err, "json marshal failed")
	}

	request, err := http.NewRequest("POST", url, bytes.NewReader(j))
	if err != nil {
		return nil, errors.Wrap(err, "could not create a request")
	}

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "client do failed")
	}

	return response, nil
}

func apiUrl(s wire.ComponentTestService, url string) string {
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}
	return "http://" + strings.Trim(s.Config.ServeAddress, "/") + "/api" + url
}
