package ctrl

import (
	"fmt"
)

const (

	actionsPath       = "actions"
	gitPath           = "git"
	blueprintPath     = "blueprint"
)

type ServiceLoc struct {
	rootDir string
	actionsDir string
	gitDir string
	blueprintDir string
}

func NewServeLoc(dir string) *ServiceLoc {
	if dir == "" {
		dir = "/usr/local/var/sticky"
	}

	return &ServiceLoc{
		rootDir: dir,
		actionsDir: fmt.Sprintf("%s/%s", dir, actionsPath),
		gitDir: fmt.Sprintf("%s/%s", dir, gitPath),
		blueprintDir: fmt.Sprintf("%s/%s", dir, blueprintPath),
	}
}

func (sl *ServiceLoc) ActionFolderPath() string {
	return sl.actionsDir
}

func (sl *ServiceLoc) GitRootPath() string {
	return sl.gitDir
}

func (sl *ServiceLoc) BlueprintDir() string {
	return sl.blueprintDir
}

func (sl *ServiceLoc) ActionPath(name string) string {
	return fmt.Sprintf("%s/%s", sl.actionsDir, name)
}

func (sl *ServiceLoc) GitActionPath(name string) string {
	return fmt.Sprintf("%s/%s.git", sl.gitDir, name)
}

func (sl *ServiceLoc) BlueprintActionPath(name string) string {
	return fmt.Sprintf("%s/%s", sl.blueprintDir, name)
}


