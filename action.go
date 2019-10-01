package ctrl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"time"

	ctrlFS "github.com/edkvm/ctrl/fs"
	"github.com/oklog/ulid"
)

type ActionRepo struct {
	path string
}

func NewActionRepo() *ActionRepo {
	return &ActionRepo{
		ctrlFS.BuildActionRepoPath(),
	}
}

func (ar *ActionRepo) List() []string {

	items, err := ioutil.ReadDir(ar.path)
	if err != nil {
		return nil
	}

	dirList := make([]string, 0)
	for _, v := range items {
		dirList = append(dirList, v.Name())
	}

	return dirList

}

func (ar *ActionRepo) ActionExists(name string) bool {
	actionPath := ctrlFS.BuildActionPath(name)
	if _, err := os.Stat(actionPath); os.IsNotExist(err) {
		return false
	}

	return true
}

type Action struct {
	Name        string
	ExecId      string
	execName    string
	handlerPath string
	configPath  string
	paramsPath  string
	sockPath    string
	ctrlCh      chan struct{}
}

func NewAction(name string) *Action {
	execId := genULID()
	actionPath := ctrlFS.BuildActionPath(name)
	return &Action{
		Name:     name,
		ExecId:   execId,
		execName: "node",
		// Stack related
		handlerPath: fmt.Sprintf("%s/index.js", actionPath),
		configPath:  fmt.Sprintf("%s/config.json", actionPath),
		paramsPath:  fmt.Sprintf("%s/params.json", actionPath),
		sockPath:    fmt.Sprintf("%s/tmp/%s_%s.sock", actionPath, name, execId),
	}
}

func genULID() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	id, err := ulid.New(ulid.Timestamp(t), entropy)
	if err != nil {

	}

	return fmt.Sprintf("%s", id)
}

func (fr *Action) Execute(payload []byte) string {

	cmdParams := []string{
		fr.handlerPath,
		fr.sockPath,
	}

	// TODO: Add context to cmd
	cmd := exec.Command(fr.execName, cmdParams...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("failed to connect stdout: %v\n", err)
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("failed to connect stderr: %v\n", err)
	}
	defer stderr.Close()

	// TODO: Add Instrumentation
	if err := cmd.Start(); err != nil {
		log.Printf("failed to start cmd: %v\n", err)
	}

	outScanner := bufio.NewScanner(stdout)
	errScanner := bufio.NewScanner(stderr)
	go func() {
		for outScanner.Scan() {
			log.Printf("[%s] stdout: %v\n", fr.Name, outScanner.Text())
		}
	}()

	go func() {
		for errScanner.Scan() {
			log.Printf("[%s] stderr: %v\n", fr.Name, errScanner.Text())
		}
	}()

	inputCh := make(chan []byte, 0)
	outCh := make(chan []byte, 0)

	fr.openSock(inputCh, outCh)

	inputCh <- payload

	// Wait for result
	result := <-outCh

	if err := cmd.Wait(); err != nil {
		log.Print(err)
	}

	return string(result)
}

func (fr *Action) openSock(inputCh <-chan []byte, outCh chan []byte) {

	addr, err := net.ResolveUnixAddr("unix", fr.sockPath)
	if err != nil {
		log.Println("failed to resolve: %v", err)
		os.Exit(1)
	}

	sock, err := net.ListenUnix("unix", addr)
	if err != nil {
		log.Printf("failed to open listener: %v\n", err)
	}

	go func() {
		conn, err := sock.AcceptUnix()
		if err != nil {
			log.Printf("error start accept on conn: %v\n", err)
			return
		}
		defer conn.Close()
		if err != nil {
			log.Printf("error in accepting new connection: %v\n", err)
			return
		}

		// Wait for connection
		buf := make([]byte, 256)
		n, _, _ := conn.ReadFromUnix(buf)
		op := string(buf[:n])
		if op != "op|start" {
			log.Println(fmt.Sprintf("wrong handshake from client: %v\n", op))
			return
		}

		// TODO: Add system context
		select {
		case op := <-inputCh:
			conn.Write(op)
		}

		// Wait for function output
		result := make([]byte, 256)
		for {
			n, _, err = conn.ReadFromUnix(buf)
			if err != nil {
				log.Println(err)
				break
			}

			header := string(buf[:2])
			if header == "op" {
				op = string(buf[3:8])
				if op == "close" {
					break
				}
			}
			copy(result, buf)
		}

		outCh <- result

		sock.Close()
	}()

}

func (fr *Action) PayloadFromJSON(params map[string]interface{}) []byte {
	confDef := ctrlFS.ReadFile(fr.configPath)

	var config map[string]interface{}
	json.Unmarshal(confDef, &config)

	payload := make(map[string]interface{}, 0)

	payload["ctx"] = config
	payload["params"] = params

	buf, _ := json.Marshal(payload)

	return buf
}

func (fr *Action) PayloadFromString(args []string) []byte {
	confDef := ctrlFS.ReadFile(fr.configPath)
	paramDef := ctrlFS.ReadFile(fr.paramsPath)

	vals := make([]interface{}, len(args))
	for i, _ := range args {
		vals[i] = args[i]
	}

	var params map[string]interface{}
	json.Unmarshal(paramDef, &params)

	idx := 0
	for k, _ := range params {
		params[k] = vals[idx]
		idx = idx + 1
	}

	var config map[string]interface{}
	json.Unmarshal(confDef, &config)

	payload := make(map[string]interface{}, 0)

	payload["ctx"] = config
	payload["params"] = params

	buf, _ := json.Marshal(payload)

	return buf
}
