package stacks

import (
	"fmt"
	"github.com/stickysh/sticky/pkg/fs"
	"log"
	"os"
	"os/exec"
)

type stackFile struct {
	name   string
	output string
	exectuable bool
	remove bool
}

type goStack struct {
	filenames []stackFile
}

func NewGoV1() *goStack {
	return &goStack{
		filenames: []stackFile{
			{
				name:   "action",
				remove: true,
			},
			{
				name: "config.json",
			},
			{
				name: "params.json",
			}},
	}
}

func (gs *goStack) Build(wd string) (map[string][]byte, error) {
	//
	cmdParams := []string{"build", "-o", fmt.Sprintf("%s/%s", wd, "action"), fmt.Sprintf("%s/%s", wd, "main.go")}
	cmd := exec.Command("go", cmdParams...)
	cmd.Dir = wd

	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}


	files := make(map[string][]byte, 0)

	for i := 0; i < len(gs.filenames); i++ {
		filename := gs.filenames[i].name
		path := fmt.Sprintf("%s/%s", wd, filename)
		content := fs.ReadFile(path)
		if content == nil {
			continue
		}

		
		if gs.filenames[i].remove {
			os.Remove(path)
		}



		files[filename] = content
	}

	return files, nil
}

func (gs *goStack) Deploy(path string) error {
	return nil
}
