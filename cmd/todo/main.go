package main

import (
	"bufio"
	"flag"
	"fmt"
	todo "go-to-do"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var filename = "todo.json"

func main() {
	// Check if environment variable defines different file name
	if os.Getenv("TODO_FILE") != "" {
		filename = os.Getenv("TODO_FILE")
	}

	// Logging configuratgion
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

	// Parsing command line flags
	add := flag.Bool("add", false, "Add a new task to the list")
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
	case *add:
		task, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			log.Fatal().Err(err)
		}
		l.Add(task)
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

func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	// If arguments are not present, read input from STDIN
	s := bufio.NewScanner(r)
	s.Scan()

	if err := s.Err(); err != nil {
		return "", err
	}

	if len(s.Text()) == 0 {
		return "", fmt.Errorf("task cannot be empty")
	}

	return s.Text(), nil

}
