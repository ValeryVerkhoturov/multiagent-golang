package multiagent

import (
	"fmt"
	"sync"
)

// Task represents a unit of work to be executed by an Agent
type Task struct {
	Name        string        // Name uniquely identifies the task
	Description string        // Description provides details for the task
	Agent       *Agent        // Agent is the executor assigned to run this task
	DependsOn   []*Task       // DependsOn lists other tasks that must be completed before this one can start
	output      *string       // Output holds the result produced by the task. Nill, if not executed
	Mu          sync.RWMutex  // Mu is used to safely access and update the task's output
	done        chan struct{} // done is a channel that signals when the task has finished execution
}

func (t *Task) GetOutput() *string {
	var output *string
	t.Mu.RLock()
	output = t.output
	t.Mu.RUnlock()
	return output
}

// Agent represents an executor that performs work using a specific function
type Agent struct {
	Name     string                                   // Name uniquely identifies the agent
	Function func(input string, tasks []*Task) string // Function is the work function the agent executes, taking current task description and executed depends on tasks
	Crew     *Crew                                    // Crew points back to the Crew the agent is part of (optional association)
}

// Crew manages a collection of tasks and agents.
type Crew struct {
	Tasks  map[string]*Task  // Tasks maps task names to Task objects for easy lookup.
	Agents map[string]*Agent // Agents maps agent names to Agent objects.
}

// NewCrew creates and returns a new Crew instance with initialized maps for tasks and agents.
func NewCrew() *Crew {
	return &Crew{
		Tasks:  make(map[string]*Task),
		Agents: make(map[string]*Agent),
	}
}

// AddAgent adds an Agent to the Crew.
// The agent is stored in the Crew's Agents map using the agent's name as the key.
func (c *Crew) AddAgent(agent *Agent) {
	c.Agents[agent.Name] = agent
}

// AddTask adds a Task to the Crew, assigns it to an agent, and establishes its dependencies.
// agentName: The name of the agent that will execute the task.
// dependsOnTaskNames: A list of names of tasks that must complete before this task can run.
func (c *Crew) AddTask(task *Task, agentName string, dependsOnTaskNames []string) error {
	// Look up the agent by name.
	agent, ok := c.Agents[agentName]
	if !ok {
		return fmt.Errorf("no agent with name: %v", agentName)
	}
	// Assign the agent to the task.
	task.Agent = agent

	// Initialize the DependsOn slice to hold dependency tasks.
	task.DependsOn = make([]*Task, 0, len(dependsOnTaskNames))

	// Initialize the done channel to signal task completion.
	task.done = make(chan struct{})

	// Loop over each dependency task name provided.
	for _, dependsOnTaskName := range dependsOnTaskNames {
		// Retrieve the dependency task from the Crew's task map.
		dependsOnTask, ok := c.Tasks[dependsOnTaskName]
		if !ok {
			return fmt.Errorf("no depends on task with name: %v", dependsOnTaskName)
		}
		// Add the dependency to the task's list.
		task.DependsOn = append(task.DependsOn, dependsOnTask)
	}

	// Store the task in the Crew's Tasks map using the task's name as the key.
	c.Tasks[task.Name] = task
	return nil
}

// Kickoff executes all tasks concurrently while respecting dependency order.
// Each task is launched as a goroutine and will wait for all its dependencies to signal completion before execution.
func (c *Crew) Kickoff() {
	var wg sync.WaitGroup // WaitGroup is used to wait for all tasks to complete.

	// Iterate over every task in the Crew.
	for _, task := range c.Tasks {
		wg.Add(1) // Increment WaitGroup counter for each task.
		// Launch the task in a new goroutine.
		go func(t *Task) {
			// Wait for all dependencies to complete.
			for _, dep := range t.DependsOn {
				<-dep.done // Blocks until the dependency task's done channel is closed.
			}

			// Execute the task's assigned function using the task's description as input.
			result := t.Agent.Function(t.Description, t.DependsOn)

			// Lock the task to safely update its Output field.
			t.Mu.Lock()
			t.output = &result
			t.Mu.Unlock()

			// Signal that this task is complete by closing its done channel.
			close(t.done)
			// Decrement the WaitGroup counter.
			wg.Done()
		}(task)
	}

	// Wait until all goroutines (tasks) have finished executing.
	wg.Wait()
}
