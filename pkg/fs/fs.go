package fs

import (
	"io/ioutil"
	"log"
	"os"

	_ "github.com/stickysh/sticky/statik" // Generated FS
	statikFS "github.com/rakyll/statik/fs"
)


func WriteFile(filepath string, data []byte) error {
	fd, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = fd.Write(data)

	return err
}

func ReadFile(filePath string) []byte {

	_, err := os.Stat(filePath)
	if err != nil {
		return nil
	}

	// Copy handler
	srcFd, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer srcFd.Close()

	data, err := ioutil.ReadAll(srcFd)
	if err != nil {
		return nil
	}

	return data
}

func ReadStaticFile(path string) ([]byte, error) {
	sfs, err := statikFS.New()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(path)
	fd, err := sfs.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	buf, err := ioutil.ReadAll(fd)
	if err != nil {
		log.Fatal(err)
	}

	return buf, nil
}
