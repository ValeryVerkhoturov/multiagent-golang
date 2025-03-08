# Multi-agents framework for LLM based on layered communication

![](/pic/paper.jpg)

## Installation

```shell
go get -u github.com/ValeryVerkhoturov/multiagent-golang
```

## Example 

### Code

```go
package main

import (
	"fmt"
	"time"

	"github.com/ValeryVerkhoturov/multiagent-golang"
)

func createCrew() (*multiagent.Crew, error) {
	crew := multiagent.NewCrew()

	agent1Name := "Agent1"
	agent2Name := "Agent2"

	task1Name := "Task1"
	task2Name := "Task2"
	task3Name := "Task3"
	task4Name := "Task4"
	task5Name := "Task5"

	crew.AddAgent(&multiagent.Agent{
		Name: agent1Name,
		Function: func(input string, tasks []*multiagent.Task) string {
			time.Sleep(1 * time.Second)
			return "Processed by Agent1: " + input
		},
	})

	crew.AddAgent(&multiagent.Agent{
		Name: agent2Name,
		Function: func(input string, tasks []*multiagent.Task) string {
			time.Sleep(2 * time.Second)
			return "Processed by Agent2: " + input
		},
	})

	task1 := &multiagent.Task{
		Name:        task1Name,
		Description: "Data for Task1",
	}
	if err := crew.AddTask(task1, agent1Name, []string{}); err != nil {
		return nil, fmt.Errorf("error adding Task1: %v", err)
	}

	task2 := &multiagent.Task{
		Name:        task2Name,
		Description: "Data for Task2",
	}
	if err := crew.AddTask(task2, agent2Name, []string{}); err != nil {
		return nil, fmt.Errorf("error adding Task2: %v", err)
	}

	task3 := &multiagent.Task{
		Name:        task3Name,
		Description: "Data for Task3",
	}
	if err := crew.AddTask(task3, agent2Name, []string{}); err != nil {
		return nil, fmt.Errorf("error adding Task3: %v", err)
	}

	task4 := &multiagent.Task{
		Name:        task4Name,
		Description: "Data for Task4",
	}
	if err := crew.AddTask(task4, agent1Name, []string{task1Name, task2Name}); err != nil {
		return nil, fmt.Errorf("error adding Task4: %v", err)
	}

	task5 := &multiagent.Task{
		Name:        task5Name,
		Description: "Data for Task5",
	}
	if err := crew.AddTask(task5, agent1Name, []string{task3Name, task4Name}); err != nil {
		return nil, fmt.Errorf("error adding Task5: %v", err)
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

```

### Output
```text
Executing tasks...
All tasks completed in 4.002615625s
Task3 output: Processed by Agent2: Data for Task3
Task4 output: Processed by Agent1: Data for Task4
Task5 output: Processed by Agent1: Data for Task5
Task1 output: Processed by Agent1: Data for Task1
Task2 output: Processed by Agent2: Data for Task2
```

