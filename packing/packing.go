package packing

import (
	"bytes"
	"fmt"
	"github.com/stickysh/sticky"
	"github.com/stickysh/sticky/packing/stacks"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type StackConfig interface {
	Deploy(path string) error
	Build(wd string) (map[string][]byte, error)
}

var stacksList = map[string]StackConfig{
	"node10": stacks.NewNodev10(),
	"go": stacks.NewGoV1(),
}

type ActionPack struct {
	sl *ctrl.ServiceLoc
}

func NewActionPack(sl *ctrl.ServiceLoc) *ActionPack {
	return &ActionPack{
		sl: sl,
	}
}

func (ap *ActionPack) Create(name string) error {
	dir := ap.sl.GitActionPath(name)

	if st, err := os.Stat(dir); err == nil && st.IsDir() {
		return fmt.Errorf("function already exists")
	}

	// Create Bare git repo
	args := []string{"init", "--bare", fmt.Sprintf("%s.git", name) }
	command := exec.Command("/usr/bin/git", args...)
	command.Dir = ap.sl.GitRootPath()

	err := command.Run()

	return err
}

func (ap *ActionPack) Install(src string) error {
	panic("Implament ap.Install")
}

func (ap *ActionPack) Deploy(repoName string) error {
	list := strings.Split(repoName, ".git")

	if len(list) < 1 {
		return fmt.Errorf("wrong repo name")
	}

	actionName := list[0]
	wd := ap.sl.BlueprintActionPath(actionName)


	args := []string{"pull", "origin", "master" }
	dir := wd
	if _, err := os.Stat(wd); os.IsNotExist(err) {
		args = []string{"clone", ap.sl.GitActionPath(actionName) }
		dir = ap.sl.BlueprintDir()
	}
	// Clone from git
	// Create Bare git repo
	cmd := exec.Command("/usr/bin/git", args...)
	cmd.Dir = dir

	err := cmd.Run()
	if err != nil {
		return err
	}

	// Build in place
	pack, err := BuildPack("go", actionName, wd)
	if err != nil {
		return err
	}

	// Deploy
	err = pack.Deploy(ap.sl.ActionPath(actionName))
	if err != nil {
		return err
	}

	log.Println("method", "deploy", "action", actionName, "path", ap.sl.ActionPath(actionName))
	return nil
}

func (ap *ActionPack) Pack() {

}


type Pack struct {
	stack      StackConfig
	actionName string
	files      map[string][]byte
}

func BuildPack(stackName, actionName, wd string) (*Pack, error) {
	// TODO: Add more error handeling
	dirs := strings.Split(wd, "/")
	if len(dirs) < 2 {
		// TODO: return error, actionName is not absolute
	}

	pk := &Pack{
		stack:      stacksList[stackName],
		actionName: actionName,
		files:      make(map[string][]byte, 3),
	}

	files, err := pk.stack.Build(wd)
	if err != nil {
		return nil, err
	}

	pk.files = files

	log.Println("method", "build", "action", actionName, "wd", wd)
	return pk, nil
}

func (pk *Pack) Deploy(dst string) error {

	// Create tmp folder
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		err := os.MkdirAll(fmt.Sprintf("%s/tmp", dst), os.ModePerm)
		if err != nil {
			return err
		}
	}

	err := pk.stack.Deploy(dst)
	if err != nil {
		return nil
	}

	// Write Action files
	for name, _ := range pk.files {
		// TODO move to function
		file := pk.files[name]
		srcReader := bytes.NewReader(file)

		dstFd, err := os.OpenFile(fmt.Sprintf("%s/%s", dst, name),os.O_RDWR|os.O_CREATE|os.O_TRUNC,os.ModePerm)
		if err != nil {
			return err
		}
		defer dstFd.Close()

		_, err = io.Copy(dstFd, srcReader)
		if err != nil {
			// TODO: File didn't open, Report as (SystemError)
			return err
		}

	}

	return nil
}
