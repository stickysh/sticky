package main

import (
	"fmt"
	"github.com/edkvm/ctrl"
	"github.com/edkvm/ctrl/packing"
	"log"
	"os"
	"path/filepath"


	"github.com/urfave/cli"
)


func main() {
	//currentPath, err := filepath.Abs(filepath.Dir(os.Args[0]))

	app := cli.NewApp()

	app.Name = "ctrl"

	app.Commands = []cli.Command{
		{
			Name: "deploy",
			Usage: "copy action to the runner, if the action does not exist it will be created",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "stack, s",
					Value: "node10",
				},
			},
			Action: func(c *cli.Context) error {
				wdPath, err := os.Getwd()
				if err != nil {
					log.Fatal(err)
				}

				if c.NArg() > 0 {
					pathArg := c.Args().First()
					wdPath, err = filepath.Abs(pathArg)
					if err != nil {
						log.Fatal(err)
					}
				}

				log.Println("deploying action from: ", wdPath)

				stackName := c.String("stack")
				pk, err := packing.BuildPack(stackName, wdPath)
				if err != nil {
					return err
				}

				err = pk.Deploy()
				if err != nil {
					log.Println(err)
				}

				return nil
			},

		},
		{
			Name: "init",
			Usage: "no usage yet",
			Action: func(c *cli.Context) error {
				log.Println("init not implemented")
				return nil
			},
		},
		{
			Name: "ls",
			Usage: "list all the available actions",
			Action: func(c *cli.Context) error {
				ar := ctrl.NewActionProvider()

				fmt.Println(ar.List())
				return nil
			},
		},
		{
			Name: "run",
			Usage: "runs the specified function",
			Action: func(c *cli.Context) error{

				actionName := c.Args().First()


				args := c.Args()[1:]
				ar := ctrl.NewActionProvider()

				if !ar.ActionExists(actionName) {
					log.Fatal(fmt.Sprintf("action does not exists: %s", actionName))
				}

				fr := ctrl.NewAction(actionName)

				log.Println("running action:", actionName)
				log.Println("params:", args)
				// Parse params
				params := fr.ParamsToJSON([]string(args))

				payload, env := fr.EncodePayload(params)
				result := ar.ExecuteAction(actionName, payload, env)


				log.Println(result)

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

