package rady

const (
	// DefaultPath is the path of config file for default
	DefaultPath = "./resources/application.conf"

	DefaultConfType = YAML

	// YAML is the suffix of yaml file
	YAML = "yaml"

	// JSON is the suffix of json file
	JSON = "json"
)

/*
CONF is a tag of Boot to define the path and type of config file

Usage:

	type Root struct {
		CONF `path:"./resources/app.conf" type:"yaml"`
	}

	func main() {
		CreateApplication(new(Root)).Run()
	}
*/
type CONF struct {
}
