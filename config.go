package rhapsody


// DefaultPath is the path of config file for default
const DefaultPath = "./resources/application.conf"

// YAML is the suffix of yaml file
const YAML = "yaml"

// JSON is the suffix of json file
const JSON = "json"

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