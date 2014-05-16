package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/go-chef/gladius/chef"
)

// VERSION is the gladius version
const VERSION = "0.0.1"

type Config struct {
	Client *chef.ClientConfig
}

func main() {
	usage := `usage: gladius <action> <object> [<names>]
	[-s|--server=<serverUri>] 
	[-k|--key=<keyFile>]
	[-c|--client=<clientName>]
	[-f|--file=<configFile>]
	[-o|--output=<formatter>] 
	[--version] [--help] 
  [--]

options:
	-h, --help 	   
	--version   
	-s, --server uri  The server uri to connect to. Include org if you use an org
	-k, --key    file The keyfile to read (this is your chef-client key)
	-c, --client name The chef Client name to use when talking to a server
	-f, --format formatter The Formatter to use (json,txt)

Actions: 
	solve     Run solver or execute a solve on a server.
	upload    Upload the object to the server.
 	download  Download an object frmo the server.
	edit      Edit an object in $EDITOR.
	show      Output the contents of the object.
	list      List objects.

Objects: 
	cookbook, role, run_list, environment, data_bag

See 'gladius help command' for more info on that command.
`

	args, _ := docopt.Parse(usage, nil, true, "git version 1.7.4.4", true)

	// Load/Validate the Config
	// Setup Chef Connection/config
	// Dispatch

	// switch on the action
	switch args["<action>"].(string) {
	case "download":
		// Dispatch to subcomand switch
	//	download(object, args)
	case "upload":
		fmt.Println("Not Implemented")
	case "solve":
		fmt.Println("Not Implemented")
	case "edit":
		fmt.Println("Not Implemented")
	case "show":
		fmt.Println("Not Implemented")
	case "list":
		fmt.Println("Not Implemented")
	case "help", "":
	}
}

/* old download crap
client, err := chef.Connect()
if err != nil {
	fmt.Println("Error:", err)
	os.Exit(1)
}
client.SSLNoVerify = true

// TODO: make this an arg
env, _, err := client.GetEnvironment("_default")

// TODO: make this an arg
runlist_json := `{ "run_list": [ "base"] }`
rl := strings.NewReader(runlist_json)
if err != nil {
	fmt.Println("error:", err)
	os.Exit(1)
}

// TODO: make this an arg
endpoint := fmt.Sprintf("/environments/%s/cookbook_versions", env.Name)

// send the request to the chef server to solve the run_list
fmt.Println("Requesting solve for runlist", rl)
resp, err := client.Post(endpoint, nil, rl)
if err != nil {
	fmt.Println("error:", err)
	os.Exit(1)
}

cookbooks := map[string]chef.CookbookVersion{}
json.NewDecoder(resp.Body).Decode(&cookbooks)

// PigThrusters Engage
Enqueue(cookbooks, client)
*/
