package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// TODO - split out URL properly
const statusURL = "https://api.tfl.gov.uk/line/mode/"

// Express valid modes of transport as a map, instead of a plain slice;
// this seems a little odd at first glance, but checking presnce of user
// supplied modes against the map is much simpler than doing the same with
// a slice
var ValidModes = map[string]struct{}{
	"tube": {},
	"dlr":  {},
}

// TubeStatusResult represents the top-level object returned in the JSON
// response; it contains a list of status objects, one for each specified
// mode line
type TubeStatusResult struct {
	Name         string
	LineStatuses []*Status
}

type Status struct {
	Description string `json:"statusSeverityDescription"`
}

// GetStatus returns a list of issues with the supplied transport modes
func GetStatus(modes []string) (*[]TubeStatusResult, error) {
	if len(modes) == 0 {
		// Provide default, of currently valid modes
		for key := range ValidModes {
			modes = append(modes, key)
		}
	}
	// TODO - better way of doing this
	q := "?app_id=" + config.AppId + "&app_key=" + config.ApiKey
	resp, err := http.Get(statusURL + strings.Join(modes, ",") + "/status" + q)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Tube status query failed: %s", resp.Status)
	}

	var result []TubeStatusResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
