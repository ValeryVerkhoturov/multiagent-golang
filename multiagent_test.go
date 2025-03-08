package multiagent

import (
	"testing"
)

func TestCrewExecution(t *testing.T) {
	// Create a new Crew instance
	crew := NewCrew()

	// Define a simple agent function
	agentFunc := func(input string, tasks []*Task) string {
		return "Completed: " + input
	}

	// Create an agent and add it to the crew
	agent := &Agent{Name: "Agent1", Function: agentFunc}
	crew.AddAgent(agent)

	// Create tasks
	task1 := &Task{Name: "Task1", Description: "First Task"}
	task2 := &Task{Name: "Task2", Description: "Second Task"}
	task3 := &Task{Name: "Task3", Description: "Third Task"}

	// Add tasks to the crew with dependencies
	if err := crew.AddTask(task1, "Agent1", nil); err != nil {
		t.Fatalf("Failed to add Task1: %v", err)
	}
	if err := crew.AddTask(task2, "Agent1", []string{"Task1"}); err != nil {
		t.Fatalf("Failed to add Task2: %v", err)
	}
	if err := crew.AddTask(task3, "Agent1", []string{"Task2"}); err != nil {
		t.Fatalf("Failed to add Task3: %v", err)
	}

	// Execute all tasks
	crew.Kickoff()

	// Validate execution order
	if output := task1.GetOutput(); output == nil || *output != "Completed: First Task" {
		t.Errorf("Unexpected output for Task1: got %v", output)
	}
	if output := task2.GetOutput(); output == nil || *output != "Completed: Second Task" {
		t.Errorf("Unexpected output for Task2: got %v", output)
	}
	if output := task3.GetOutput(); output == nil || *output != "Completed: Third Task" {
		t.Errorf("Unexpected output for Task3: got %v", output)
	}
}
