package main

import (
	"context"
	"errors"
	"flag"
	"io/ioutil"
	"log"

	"github.com/MasteryConnect/pipe/line"
	"github.com/MasteryConnect/pipe/message"

	"github.com/MasteryConnect/pipe/extras/csv"
	"github.com/MasteryConnect/pipe/extras/json"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
)

var (
	schemaCueFile string
	stopOnError   bool
)

func init() {
	flag.StringVar(&schemaCueFile, "schema", "", "[optional] the path to the cuelang schema file")
	flag.BoolVar(&stopOnError, "stop", false, "[optional] stop the pipeline if there is a schema validation error")
}

func main() {
	flag.Parse()
	cueCtx := cuecontext.New()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pipe := line.New().Add(
		csv.ReadStream(','),
		line.I(json.To),
	)

	if schemaCueFile != "" {
		schemaStr, err := ioutil.ReadFile(schemaCueFile)
		if err != nil {
			log.Fatal(err)
			return
		}
		schema := cueCtx.CompileBytes(schemaStr)

		if stopOnError {
			pipe.Add(stopOnCancel(ctx))
		}

		pipe.Add(newValidator(cueCtx, schema, cancel))
	}

	pipe.Add(line.Stdout).RunContext(ctx)
}

func newValidator(cueCtx *cue.Context, schema cue.Value, cancel func()) line.Tfunc {
	return line.I(func(m interface{}) (interface{}, error) {
		val := cueCtx.CompileString(message.String(m))
		err := schema.Subsume(val)
		if err != nil {
			cancel()
			return nil, err
		}
		return m, nil
	})
}

func stopOnCancel(ctx context.Context) line.Tfunc {
	return func(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
		for m := range in {
			select {
			case <-ctx.Done():
				errs <- errors.New("context cancelled")
				return
			default:
				out <- m
			}
		}
	}
}
