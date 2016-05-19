### go-tube
A simple command-line utility to provide status updates on all London Underground lines.
#### Configuration
`go-tube` requires a TfL App ID and API key; these must be specified in a `config.json` file kept in the same directory as the compiled binary. The structure is simple:

    {
        "appId": "sampleAppIdInHere",
        "apiKey": "sampleApiKeyInHere"
    }

An AppID and API key can be obtained by registering with [TfL's API portal](https://api-portal.tfl.gov.uk/)

#### Command-line flags
`go-tube` currently accepts one command-line flag:

    -modes     A comma-separated list of transport modes to check the status of

where valid modes go-tube currently recognises are:

    tube
    dlr
