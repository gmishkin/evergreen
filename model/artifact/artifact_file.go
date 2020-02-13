package artifact

const Collection = "artifact_files"

const (
	// strings for setting visibility
	Public  = "public"
	Private = "private"
	None    = "none"
	Signed  = "signed"
)

var ValidVisibilities = []string{Public, Private, None, Signed, ""}

// Entry stores groups of names and links (not content!) for
// files uploaded to the api server by a running agent. These links could
// be for build or task-relevant files (things like extra results,
// test coverage, etc.)
type Entry struct {
	TaskId          string `json:"task" bson:"task"`
	TaskDisplayName string `json:"task_name" bson:"task_name"`
	BuildId         string `json:"build" bson:"build"`
	Files           []File `json:"files" bson:"files"`
	Execution       int    `json:"execution" bson:"execution"`
}

// Params stores file entries as key-value pairs, for easy parameter parsing.
//  Key = Human-readable name for file
//  Value = link for the file
type Params map[string]string

// File is a pairing of name and link for easy storage/display
type File struct {
	// Name is a human-readable name for the file being linked, e.g. "Coverage Report"
	Name string `json:"name" bson:"name"`
	// Link is the link to the file, e.g. "http://fileserver/coverage.html"
	Link string `json:"link" bson:"link"`
	// Visibility determines who can see the file in the UI
	Visibility string `json:"visibility" bson:"visibility"`
	// When true, these artifacts are excluded from reproduction
	IgnoreForFetch bool `bson:"fetch_ignore,omitempty" json:"ignore_for_fetch"`
	//AwsKey is the key with which the file was uploaded to s3
	AwsKey string `json:"aws_key,omitempty" bson:"aws_key,omitempty"`
	//AwsSercret is the secret with which the file was uploaded to s3
	AwsSecret string `json:"aws_secret,omitempty" bson:"aws_secret,omitempty"`
}

// Array turns the parameter map into an array of File structs.
// Deprecated.
func (params Params) Array() []File {
	var files []File
	for name, link := range params {
		files = append(files, File{
			Name: name,
			Link: link,
		})
	}
	return files
}
