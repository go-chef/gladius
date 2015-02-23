package policy

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/jawher/mow.cli"
)

type Lock struct {
	Name       string   `json:"name"`
	RevisionId string   `json:"revision_id"`
	RunList    []string `json:"run_list"`
	CookBooks  map[string]CookLock
}

type CookLock struct {
	Version                 string `json:"version"`
	CacheKey                string `json:"cache_key"`
	DottedDecimalIdentifier string `json:"dotted_decimal_identifier"`
	Identifier              string `json:"identifier"`
	Origin                  string `json:"origin"`
	Source                  string `json:"source"`

	SourceOptions struct {
		Path           string `json:"path"`
		Artifactserver string `json:"artifactserver"`
		Version        string `json:"version"`
	} `json:"source_options"`

	ScmInfo struct {
		Published                  bool     `json:"published"`
		Remote                     string   `json:"remote"`
		Revision                   string   `json:"revision"`
		Scm                        string   `json:"scm"`
		SynchronizedRemoteBranches []string `json:"synchronized_remote_branches"`
		WorkingTreeClean           bool     `json:"working_tree_clean"`
	} `json:"scm_info"`
}

func Update(cmd *cli.Cmd) {
	pFile := cmd.StringArg("l lockfile", "./Policyfile.lock.json", "Path to Policyfile Lock")
	force := cmd.BoolOpt("f force", false, "force update all deps")

	spew.Dump(cmd)
	spew.Dump(pFile)
	spew.Dump(force)
}
