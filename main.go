package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
)

// Text template for outputting status updates
const templ = `Tube issues:
{{range .}}{{.Name}}: {{.LineStatuses | getDescription}}
{{end}}
`

var (
	// The issues report
	report = template.Must(template.New("statusList").
		Funcs(template.FuncMap{"getDescription": getDescription}).
		Parse(templ))
	// The configuration object, read from file
	config = Configuration{}
)

// Config is a struct containing the TfL app ID and API key
type Configuration struct {
	AppId  string
	ApiKey string
}

// ModesFlag is a flag value automatically split up into a slice, delimited
// by commas. It holds a list of user-supplied modes to check
type ModesFlag struct {
	modes []string
}

// String returns the strng representation of the receiving ModesFlag
func (s *ModesFlag) String() string {
	return fmt.Sprint(s.modes)
}

// Set sets the receiving ModesFlag based on the comma-seperated list of values
// received from the command line. If the stribg cannot be delimited, an error is
// thrown.
func (s *ModesFlag) Set(value string) error {
	modes := strings.Split(value, ",")
	for _, item := range modes {
		s.modes = append(s.modes, item)
	}

	// TODO
	return nil
}

// getDescription returns the description from the first Status in the slice; there
// is only ever one entry in the slice anyway, so this function provides a nicer way of
// getting the description in the template
func getDescription(lineStatuses []*Status) string {
	return lineStatuses[0].Description
}

// Read in the config file to obtain the API app ID and key
func init() {
	f, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Printf("Unable to read config.json\n", err)
		os.Exit(1)
	}
	json.Unmarshal(f, &config)
}

func main() {
	var modes ModesFlag
	flag.Var(&modes, "modes", "Modes of travel to check")
	flag.Parse()

	// Validate any modes supplied on the command line
	for _, mode := range modes.modes {
		if _, ok := ValidModes[mode]; ok != true {
			fmt.Printf("Invalid transport mode specified (%s)\n", mode)
			os.Exit(1)
		}
	}

	result, err := GetStatus(modes.modes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to retrieve URL; %v\n", err)
	}

	if err := report.Execute(os.Stdout, result); err != nil {
		log.Fatal(err)
	}
}
