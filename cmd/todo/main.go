package main

import (
	"flag"
	"fmt"
	todo "go-to-do"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

const filename = "todo.json"

func main() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("***%s****", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	log := zerolog.New(output).With().Timestamp().Logger()
	// log.Info().Msg("Welcome to TODO cli app")

	// Parsing command line flags
	task := flag.String("task", "", "Task to be included in todo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "item to be completed")

	flag.Parse()

	l := &todo.List{}

	if err := l.Load(filename); err != nil {
		log.Fatal().Err(err).Msgf("can't read from file %s", filename)
	}

	switch {
	case *list:
		fmt.Print(l)
	case *task != "":
		l.Add(*task)
		if err := l.Save(filename); err != nil {
			log.Fatal().Err(err).Msgf("can't save to file %s", filename)
		}
	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			log.Fatal().Err(err)
		}
		if err := l.Save(filename); err != nil {
			log.Fatal().Err(err).Msgf("can't save to file %s", filename)
		}
	default:
		log.Fatal().Msg("invalid option")
	}
}
