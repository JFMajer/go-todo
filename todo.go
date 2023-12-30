package todo

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type List []item

func (l *List) String() string {
	formatted := ""
	for i, item := range *l {
		prefix := " "
		if item.Done {
			prefix = "X"
		}
		formatted += fmt.Sprintf("%s%d: %s\n", prefix, i+1, item.Task)
	}
	return formatted
}

// Add creates new item and appends it to the list
func (l *List) Add(task string) {
	t := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*l = append(*l, t)
}

// Complete marks an item as completed and sets CompletedAT
func (l *List) Complete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("item %d does not exist", i)
	}

	// Adjust index for 0 based index
	ls[i-1].Done = true
	ls[i-1].CompletedAt = time.Now()

	return nil
}

// Delete method deletes a ToDo item from the list
func (l *List) Delete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("item %d does not exist", i)
	}
	// Adjusting index for 0 based index
	*l = append(ls[:i-1], ls[i:]...)
	return nil
}

// Save list to a file as JSON
func (l *List) Save(filename string) error {
	data, err := json.Marshal(l)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Load list from file
func (l *List) Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			log.Info().Msgf("File %s not found. Creating a new empty file.", filename)
			*l = List{}
			return l.Save(filename)
		}
		return err
	}

	err = json.Unmarshal(data, l)
	if err != nil {
		return err
	}

	return nil
}
