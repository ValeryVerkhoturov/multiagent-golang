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
	DependsOn   []*Task       // DependsOn lists other Tasks that must be completed before this one can start
	output      *string       // Output holds the result produced by the task. Nil if not executed
	mu          sync.RWMutex  // mu is used to safely access and update the task's output
	done        chan struct{} // done is a channel that signals when the task has finished execution
}

// GetOutput returns the output of the task, if available.
func (t *Task) GetOutput() *string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.output
}

// Agent represents an executor that performs work using a specific function
type Agent struct {
	Name     string                                   // Name uniquely identifies the agent
	Function func(input string, tasks []*Task) string // Function is the work function the agent executes, taking current task description and executed depends on Tasks
}

// Crew manages a collection of Tasks and Agents.
type Crew struct {
	Tasks  map[string]*Task  // Tasks maps task names to Task objects for easy lookup.
	Agents map[string]*Agent // Agents maps agent names to Agent objects.
}

// NewCrew creates and returns a new Crew instance with initialized maps for Tasks and Agents.
func NewCrew() *Crew {
	return &Crew{
		Tasks:  make(map[string]*Task),
		Agents: make(map[string]*Agent),
	}
}

// AddAgent adds an Agent to the Crew.
func (c *Crew) AddAgent(agent *Agent) {
	c.Agents[agent.Name] = agent
}

// AddTask adds a Task to the Crew, assigns it to an agent, and establishes its dependencies.
func (c *Crew) AddTask(task *Task, agentName string, dependsOnTaskNames []string) error {
	agent, ok := c.Agents[agentName]
	if !ok {
		return fmt.Errorf("no agent with name: %v", agentName)
	}
	task.Agent = agent
	task.DependsOn = make([]*Task, 0, len(dependsOnTaskNames))
	task.done = make(chan struct{})

	for _, dependsOnTaskName := range dependsOnTaskNames {
		dependsOnTask, ok := c.Tasks[dependsOnTaskName]
		if !ok {
			return fmt.Errorf("no depends on task with name: %v", dependsOnTaskName)
		}
		task.DependsOn = append(task.DependsOn, dependsOnTask)
	}

	c.Tasks[task.Name] = task
	return nil
}

// Kickoff executes all Tasks concurrently while respecting dependency order.
func (c *Crew) Kickoff() {
	var wg sync.WaitGroup

	for _, task := range c.Tasks {
		wg.Add(1)
		go func(t *Task) {
			defer wg.Done()
			for _, dep := range t.DependsOn {
				<-dep.done
			}
			result := t.Agent.Function(t.Description, t.DependsOn)
			t.mu.Lock()
			t.output = &result
			t.mu.Unlock()
			close(t.done)
		}(task)
	}

	wg.Wait()
}
