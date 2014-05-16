// This is mainly a toy meant to explore a worker model in Go, but it solves a runlist and pulls the set of
// Cookbooks down from a chef-server quite nicely
package gladius

import (
	_ "encoding/json"
	"fmt"
	"github.com/go-chef/authentication"
	"github.com/go-chef/chefobjects"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	_ "strings"
)

const (
	//  TODO: both of these should be config || detected
	cookRoot   = "./cookbooks"
	maxWorkers = 20
)

// download_item represents the data returned by the api for each file in a cookbook
type download_item struct {
	Item     struct{ ChefObject }
	CookPath string
}

// makedir is a wrapper on MkdirAll so we don't always have to check return/err
func makedir(dir string, mode os.FileMode) {
	err := os.MkdirAll(dir, mode)
	if err != nil {
		fmt.Println("Error making dir:", dir, err)
		os.Exit(1)
	}
}

// worker Is a download worker pulling download_items from a server
func worker(id int, queue chan *download_item, client *chef.Chef) {
	fmt.Println("Started Worker ", id)
	for {
		job := <-queue
		if job == nil {
			break
		}
		fmt.Printf("Worker %d, writing %s/%s \n", id, job.CookPath, job.Item.Path)
		filePath := fmt.Sprintf("%s/%s", job.CookPath, job.Item.Path)
		makedir(filepath.Dir(filePath), 0755)
		out, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Error Openinig file:", filePath, err)
			// TODO: Should this exit, or how should we handle this?
			os.Exit(1)
		}

		// pull and assemble the path to this object
		u, err := url.Parse(job.Item.Url)
		qs, err := url.ParseQuery(u.RawQuery)
		// uggggh.
		v := map[string]string{
			"AWSAccessKeyId": qs.Get("AWSAccessKeyId"),
			"Expires":        qs.Get("Expires"),
			"Signature":      qs.Get("Signature"),
		}

		resp, err := client.Get(u.Path, v)
		if err != nil {
			fmt.Println("Got Error Getting file from server:", filePath, err)
		}

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			// TODO: should these? bubble instead of exit
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		// explicitly close these in the loop, otherwise we will max file descriptors in osx
		out.Close()
		resp.Body.Close()
	}
}

// enqueue places every chef.CookbookItem in the map of cookbooks onto the queue
func Enqueue(cooks map[string]chef.CookbookVersion, client *chef.Chef) {

	queue := make(chan *download_item)

	// TODO: prob want to be smarter about this
	runtime.GOMAXPROCS(maxWorkers)

	// spawn workers
	for i := 0; i < maxWorkers; i++ {
		go worker(i, queue, client)
	}

	for _, cook := range cooks {
		cookbook_path := fmt.Sprintf("%s/%s", cookRoot, cook.Name)
		fmt.Sprintf("Adding %s to queue\n", cook.Name)
		objects := [][]struct{ chef.CookbookItem }{
			cook.Files,
			cook.Templates,
			cook.Definitions,
			cook.Resources,
			cook.Providers,
			cook.Libraries,
		}

		for obj := range objects {
			for _, item := range objects[obj] {
				queue <- &download_item{item, cookbook_path}
			}
		}
	}

	for n := 0; n < maxWorkers; n++ {
		// work complete
		queue <- nil
	}
}
