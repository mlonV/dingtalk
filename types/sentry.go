package types

// sentry POST Webhook data struct
type SentryAlert struct {
	Culprit         string      `json:"culprit"`
	Event           Event       `json:"event"`
	ID              string      `json:"id"`
	Level           string      `json:"level"`
	Logger          interface{} `json:"logger"`
	Message         string      `json:"message"`
	Project         string      `json:"project"`
	ProjectName     string      `json:"project_name"`
	ProjectSlug     string      `json:"project_slug"`
	TriggeringRules []string    `json:"triggering_rules"`
	URL             string      `json:"url"`
}

type Event struct {
	Metrics         Metrics        `json:"_metrics"`
	Ref             int            `json:"_ref"`
	RefVersion      int            `json:"_ref_version"`
	Contexts        Contexts       `json:"contexts"`
	Culprit         string         `json:"culprit"`
	EventID         string         `json:"event_id"`
	Exception       Exception      `json:"exception"`
	Fingerprint     []string       `json:"fingerprint"`
	GroupingConfig  GroupingConfig `json:"grouping_config"`
	Hashes          []string       `json:"hashes"`
	ID              string         `json:"id"`
	KeyID           string         `json:"key_id"`
	Level           string         `json:"level"`
	Location        string         `json:"location"`
	Logger          string         `json:"logger"`
	Metadata        Metadata       `json:"metadata"`
	Modules         Modules        `json:"modules"`
	NodestoreInsert float64        `json:"nodestore_insert"`
	Platform        string         `json:"platform"`
	Project         int            `json:"project"`
	Received        float64        `json:"received"`
	Sdk             Sdk            `json:"sdk"`
	Tags            [][]string     `json:"tags"`
	Timestamp       float64        `json:"timestamp"`
	Title           string         `json:"title"`
	Type            string         `json:"type"`
	Version         string         `json:"version"`
}

type Metrics struct {
	BytesIngestedEvent int `json:"bytes.ingested.event"`
	BytesStoredEvent   int `json:"bytes.stored.event"`
}

type Contexts struct {
	Device  Device  `json:"device"`
	Os      Os      `json:"os"`
	Runtime Runtime `json:"runtime"`
}

type Device struct {
	Arch   string `json:"arch"`
	NumCPU int    `json:"num_cpu"`
	Type   string `json:"type"`
}

type Os struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
type Runtime struct {
	GoMaxprocs    int    `json:"go_maxprocs"`
	GoNumcgocalls int    `json:"go_numcgocalls"`
	GoNumroutines int    `json:"go_numroutines"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Version       string `json:"version"`
}

type Exception struct {
	Values []Value `json:"values"`
}

type Value struct {
	Stacktrace Stacktrace `json:"stacktrace"`
	Type       string     `json:"type"`
	Value      string     `json:"value"`
}

type Stacktrace struct {
	Frames []Frame `json:"frames"`
}

type Frame []struct {
	AbsPath     string   `json:"abs_path"`
	ContextLine string   `json:"context_line"`
	Filename    string   `json:"filename"`
	Function    string   `json:"function"`
	InApp       bool     `json:"in_app"`
	Lineno      int      `json:"lineno"`
	Module      string   `json:"module"`
	PostContext []string `json:"post_context"`
	PreContext  []string `json:"pre_context"`
}

type GroupingConfig struct {
	Enhancements string `json:"enhancements"`
	ID           string `json:"id"`
}

type Metadata struct {
	DisplayTitleWithTreeLabel bool   `json:"display_title_with_tree_label"`
	Filename                  string `json:"filename"`
	Function                  string `json:"function"`
	Type                      string `json:"type"`
	Value                     string `json:"value"`
}

type Modules struct {
	GithubComGetsentrySentryGo string `json:"github.com/getsentry/sentry-go"`
	GolangOrgXSys              string `json:"golang.org/x/sys"`
	GolangOrgXText             string `json:"golang.org/x/text"`
}

type Sdk struct {
	Integrations []string `json:"integrations"`
	Name         string   `json:"name"`
	Packages     Packages `json:"packages"`
	Version      string   `json:"version"`
}

type Packages []struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
