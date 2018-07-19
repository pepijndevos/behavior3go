package behaviortree

// Represents the status of a node
// can be one of Succes, Failure, Running
type Status int

func (s Status) String() string {
  switch s {
  case 0:
    return "Success"
  case 1:
    return "Failure"
  case 2:
    return "Running"
  default:
    return "Invalid"
  }
}

const (
  Success Status = iota
  Failure
  Running
)

// The basic Node interface
type Node interface {
  // Called when transitioning to Running
  Initiate()
  // Called every tick
  Update()
  // Called when transitioning from Running
  Terminate()
  // Get the current status
  Status() Status
}

// Calls Update on the node
// Also calls Initiate and Terminate when appropriate
func Tick(node Node) Status {
  if node.Status() != Running {
    node.Initiate()
  }

  node.Update()

  if node.Status() != Running {
    node.Terminate()
  }
  return node.Status()
}

// A basic node with a status
type BasicNode struct {
  status Status
}

func (n BasicNode) Initiate() {}
func (n BasicNode) Update() {}
func (n BasicNode) Terminate() {}
func (n BasicNode) Status() Status { return n.status }

// Create a new node that always returns the same status
func NewConstantNode(status Status) *BasicNode {
  n := new(BasicNode)
  n.status = status
  return n
}
