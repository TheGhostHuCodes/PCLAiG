package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/TheGhostHuCodes/interacting/todo"
)

// Default filename
var todoFilename = ".todo.json"

func main() {
	// Check if the user defined the custom filename environment variable.
	if os.Getenv("TODO_FILENAME") != "" {
		todoFilename = os.Getenv("TODO_FILENAME")
	}

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool. Developed by tghc.\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2019.\n")
		fmt.Fprintln(
			flag.CommandLine.Output(),
			`New tasks can be added to the tool by using the --add flag and...
	1. following the flag with the task.
	2. piping the task(s) in using STDIN. Newlines seperate individual tasks.`)
		fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
	}
	// Parsing commandline flags
	add := flag.Bool("add", false, "Add task to the todo list")
	list := flag.Bool("list", false, "List all tasks")
	listPending := flag.Bool("list-pending", false, "List pending tasks only")
	complete := flag.Int("complete", 0, "Item to be completed")
	verbose := flag.Bool("verbose", false, "Show verbose output")
	flag.Parse()

	l := &todo.List{}

	if err := l.Get(todoFilename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *list:
		// List current todo items.
		fmt.Print(l.CreateReport(*listPending, *verbose))
	case *complete > 0:
		// Complete the given item.
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Save the new list
		if err := l.Save(todoFilename); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *add:
		// When any arguments (excluding flags) are provided, they will be used
		// as the new task.
		ts, err := getTasks(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, t := range ts {
			l.Add(t)
		}

		// Save the new list
		if err := l.Save(todoFilename); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		// Invalid flag provided.
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}

// getTasks function decides where to get the description for a new task from:
// arguments or STDIN. Returns an array of strings to add as tasks.
func getTasks(r io.Reader, args ...string) ([]string, error) {
	var tasks []string
	if len(args) > 0 {
		tasks = append(tasks, strings.Join(args, " "))
		return tasks, nil
	}

	s := bufio.NewScanner(r)
	for s.Scan() {
		if err := s.Err(); err != nil {
			return tasks, err
		}
		if len(s.Text()) == 0 {
			return tasks, fmt.Errorf("Task cannot be blank")
		}
		tasks = append(tasks, s.Text())
	}
	return tasks, nil
}
