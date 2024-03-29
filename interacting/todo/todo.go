package todo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// item struct represents a todo item.
type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

// List represents a list of todo items.
type List []item

// CreateReport creates a string suitable for printing out on the commandline.
func (l *List) CreateReport(listPending bool, verbose bool) string {
	formatted := ""
	for k, t := range *l {
		if listPending && t.Done {
			continue
		}
		prefix := "  "
		if t.Done {
			prefix = "X "
		}
		// Adjust the item number k to print starting from 1 instead of 0.
		formatted += fmt.Sprintf("%s%d: %s\n", prefix, k+1, t.Task)

		if verbose {
			formatted += fmt.Sprintf("\tCreated:\t%s\n", t.CreatedAt)
			if t.Done {
				formatted += fmt.Sprintf("\tCompleted:\t%s\n", t.CompletedAt)
			}
		}
	}
	return formatted
}

// Add creates a new todo item and appends it to the list.
func (l *List) Add(task string) {
	t := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}
	*l = append(*l, t)
}

// Complete method marks a todo item as completed by setting Done = true and
// CompletedAt to the current time.
func (l *List) Complete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("Item %d does not exist", i)
	}
	ls[i-1].Done = true
	ls[i-1].CompletedAt = time.Now()

	return nil
}

// Save method encodes the List as JSON and saves it using the provided file
// name.
func (l *List) Save(filename string) error {
	js, err := json.Marshal(l)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, js, 0644)
}

// Get method opens the provided filename, decodes the JSON data and parses it
// into a List.
func (l *List) Get(filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if len(file) == 0 {
		return nil
	}
	return json.Unmarshal(file, l)
}
