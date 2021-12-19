package dockerfile

import (
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	log "github.com/sirupsen/logrus"
	"os"
)

func LookupEnvironment(file string) map[string]string {
	stages, _ := LoadDockerfile(file)

	result := make(map[string]string)

	for _, stage := range stages {
		log.Debugf("parsing stage: %s", stage.Name)

		for _, cmd := range stage.Commands {
			if env, ok := cmd.(*instructions.EnvCommand); ok {

				for _, e := range env.Env {
					log.Debugf("found environment value: %s=%s", e.Key, e.Value)
					result[e.Key] = e.Value
				}
			}
		}
	}

	return result
}

func LoadDockerfile(file string) (stages []instructions.Stage, metaArgs []instructions.ArgCommand) {
	fh, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	p, err := parser.Parse(fh)
	if err != nil {
		panic(err)
	}

	stages, metaArgs, err = instructions.Parse(p.AST)
	if err != nil {
		panic(err)
	}

	return
}
