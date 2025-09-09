package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID    int       `json:"id"`
	Text  string    `json:"text"`
	Done  bool      `json:"done"`
	When  time.Time `json:"when"` // created time (optional)
}

func todoFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// fallback to current dir
		return "todo.json"
	}
	return filepath.Join(home, ".todo.json")
}

func loadTasks() ([]Task, error) {
	path := todoFilePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []Task{}, nil
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	if len(b) == 0 {
		return []Task{}, nil
	}
	if err := json.Unmarshal(b, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func saveTasks(tasks []Task) error {
	path := todoFilePath()
	b, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, 0644)
}

func nextID(tasks []Task) int {
	max := 0
	for _, t := range tasks {
		if t.ID > max {
			max = t.ID
		}
	}
	return max + 1
}

func cmdAdd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: todo add \"task text\"")
	}
	text := strings.Join(args, " ")
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	t := Task{
		ID:   nextID(tasks),
		Text: text,
		Done: false,
		When: time.Now(),
	}
	tasks = append(tasks, t)
	if err := saveTasks(tasks); err != nil {
		return err
	}
	fmt.Printf("added %d: %s\n", t.ID, t.Text)
	return nil
}

func cmdList(_ []string) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		fmt.Println("no tasks")
		return nil
	}
	fmt.Println("ID  Done  Task")
	fmt.Println("-------------------------------")
	for _, t := range tasks {
		done := " "
		if t.Done {
			done = "x"
		}
		fmt.Printf("%-3d  [%s]   %s\n", t.ID, done, t.Text)
	}
	return nil
}

func parseID(arg string) (int, error) {
	id, err := strconv.Atoi(arg)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid id: %s", arg)
	}
	return id, nil
}

func cmdDone(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: todo done <id>")
	}
	id, err := parseID(args[0])
	if err != nil {
		return err
	}
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	found := false
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Done = true
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("task %d not found", id)
	}
	if err := saveTasks(tasks); err != nil {
		return err
	}
	fmt.Printf("marked %d done\n", id)
	return nil
}

func cmdRm(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: todo rm <id>")
	}
	id, err := parseID(args[0])
	if err != nil {
		return err
	}
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	newTasks := tasks[:0]
	found := false
	for _, t := range tasks {
		if t.ID == id {
			found = true
			continue
		}
		newTasks = append(newTasks, t)
	}
	if !found {
		return fmt.Errorf("task %d not found", id)
	}
	if err := saveTasks(newTasks); err != nil {
		return err
	}
	fmt.Printf("removed %d\n", id)
	return nil
}

func printHelp() {
	fmt.Println("Simple TODO CLI")
	fmt.Println("Usage:")
	fmt.Println("  todo add \"task text\"   Add a new task")
	fmt.Println("  todo list               List tasks")
	fmt.Println("  todo done <id>          Mark task done")
	fmt.Println("  todo rm <id>            Remove task")
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}
	cmd := os.Args[1]
	args := os.Args[2:]

	var err error
	switch cmd {
	case "add":
		err = cmdAdd(args)
	case "list":
		err = cmdList(args)
	case "done":
		err = cmdDone(args)
	case "rm":
		err = cmdRm(args)
	case "help", "--help", "-h":
		printHelp()
		return
	default:
		fmt.Printf("unknown command: %s\n\n", cmd)
		printHelp()
		return
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
