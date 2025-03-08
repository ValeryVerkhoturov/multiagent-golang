package main

import (
	"fmt"
	"time"

	"github.com/ValeryVerkhoturov/multiagent-golang"
)

func createCrew() (*multiagent.Crew, error) {
	crew := multiagent.NewCrew()

	agents := []struct {
		Name     string
		Function func(input string, tasks []*multiagent.Task) string
	}{
		{
			Name: "Agent1",
			Function: func(input string, tasks []*multiagent.Task) string {
				time.Sleep(1 * time.Second)
				return "Processed by Agent1: " + input
			},
		},
		{
			Name: "Agent2",
			Function: func(input string, tasks []*multiagent.Task) string {
				time.Sleep(2 * time.Second)
				return "Processed by Agent2: " + input
			},
		},
	}

	for _, agent := range agents {
		crew.AddAgent(&multiagent.Agent{
			Name:     agent.Name,
			Function: agent.Function,
		})
	}

	tasks := []struct {
		Name         string
		Description  string
		AgentName    string
		Dependencies []string
	}{
		{"Task1", "Data for Task1", "Agent1", []string{}},
		{"Task2", "Data for Task2", "Agent2", []string{}},
		{"Task3", "Data for Task3", "Agent2", []string{}},
		{"Task4", "Data for Task4", "Agent1", []string{"Task1", "Task2"}},
		{"Task5", "Data for Task5", "Agent1", []string{"Task3", "Task4"}},
	}

	for _, task := range tasks {
		t := &multiagent.Task{
			Name:        task.Name,
			Description: task.Description,
		}
		if err := crew.AddTask(t, task.AgentName, task.Dependencies); err != nil {
			return nil, fmt.Errorf("error adding %s: %v", task.Name, err)
		}
	}

	return crew, nil
}

func main() {
	crew, err := createCrew()
	if err != nil {
		println(err.Error())
		return
	}

	println("Executing tasks...")
	start := time.Now()

	crew.Kickoff()

	elapsed := time.Since(start)
	fmt.Printf("All tasks completed in %s\n", elapsed)

	for _, task := range crew.Tasks {
		fmt.Printf("%s output: %s\n", task.Name, *task.GetOutput())
	}
}
