package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var (
	binaryName = "todoTest"
)

func TestMain(m *testing.M) {

	fmt.Println("Building tools...")
	build := exec.Command("go", "build", "-o", binaryName)

	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Can't build tool %s: %s", binaryName, err)
	}

	fmt.Println("Running tests....")
	result := m.Run()
	fmt.Println("Cleaning up...")
	os.Remove(binaryName)
	os.Remove(filename)
	os.Exit(result)
}

func TestTodoCLI(t *testing.T) {
	task := "test task 1"
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	cmdPath := filepath.Join(dir, binaryName)

	t.Run("AddNewTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-task", task)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ListTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}

		expected := fmt.Sprintf(" 1: %s\n", task)
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})
}
