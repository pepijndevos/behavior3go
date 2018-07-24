package behaviortree

type CompositeNode struct {
  BasicNode
  Children []Node
}

func (n *CompositeNode) Terminate() {
  for _, child := range n.Children {
    child.Terminate()
  }
}

// Generic function for (memory) sequential an selector nodes
func compositeUpdate(n *CompositeNode, currentIndex int, endCondition Status) (Status, int) {
  for ; currentIndex<len(n.Children); currentIndex++ {
    status := Tick(n.Children[currentIndex])
    if status == endCondition {
      continue
    } else {
      return status, currentIndex
    }
  }
  return endCondition, currentIndex
}

// A node that finds the first successfull child
type SelectorNode struct {
  CompositeNode
}

func (n *SelectorNode) Update() {
  n.Status, _ = compositeUpdate(&n.CompositeNode, 0, Failure)
}

// Create a new selector node with the given children
func NewSelectorNode(children[]Node) *SelectorNode{
  n := new(SelectorNode)
  n.Children = children
  return n
}

// A node that runs all childs until one fails
type SequentialNode struct {
  CompositeNode
}

func (n *SequentialNode) Update() {
  n.Status, _ = compositeUpdate(&n.CompositeNode, 0, Success)
}

// Create a new sequential node with the given children
func NewSequentialNode(children[]Node) *SequentialNode{
  n := new(SequentialNode)
  n.Children = children
  return n
}

// A node that runs all children
// Success or Failure is defined by
// the number of children that fail or succeed
type ParallelNode struct {
  CompositeNode
  MinimumSuccesses int
  MinimumFailures int
}

func (n *ParallelNode) Update() {
  totalFailures := 0
  totalSuccesses := 0
  for _, child := range n.Children {
    status := Tick(child)
    if status == Success {
      totalSuccesses++
    } else if status == Failure {
      totalFailures++
    }
  }
  if totalSuccesses >= n.MinimumSuccesses {
    n.Status = Success;
  } else if totalFailures >= n.MinimumFailures {
    n.Status = Failure;
  } else {
    n.Status = Running
  }
}

// Create a new parallel node with the given children
// minSucc and minFail set the boundaries for success/failure of this node
func NewParallelNodeBounded(children[]Node, minSucc int, minFail int) *ParallelNode{
  n := new(ParallelNode)
  n.Children = children
  n.MinimumSuccesses = minSucc
  n.MinimumFailures = minFail
  return n
}

// Create a new parallel node with the given children
// success or failure is either triggered by one or all nodes
func NewParallelNodeAll(children[]Node, successOnAll bool, failOnAll bool) *ParallelNode{
  n := new(ParallelNode)
  n.Children = children
  if successOnAll {
    n.MinimumSuccesses = len(children)
  } else {
    n.MinimumSuccesses = 1
  }
  if failOnAll {
    n.MinimumFailures = len(children)
  } else {
    n.MinimumFailures = 1
  }
  return n
}

// A node that runs all children
// Completed children are not re-run
// Success or Failure is defined by
// the number of children that fail or succeed
type ParallelMemoryNode struct {
  ParallelNode
  Completed []bool
  TotalFailures int
  TotalSuccesses int
}

func (n *ParallelMemoryNode) Initiate() {
  for i := range n.Completed {
    n.Completed[i] = false
  }
  n.TotalFailures = 0
  n.TotalSuccesses = 0
}

func (n *ParallelMemoryNode) Update() {
  for i, child := range n.Children {
    if !n.Completed[i] {
      status := Tick(child)
      if status != Running {
        n.Completed[i] = true
      }
      if status == Success {
        n.TotalSuccesses++
      } else if status == Failure {
        n.TotalFailures++
      }
    }
  }
  if n.TotalSuccesses >= n.MinimumSuccesses {
    n.Status = Success;
  } else if n.TotalFailures >= n.MinimumFailures {
    n.Status = Failure;
  } else {
    n.Status = Running
  }
}

// Create a new parallel node with the given children
// minSucc and minFail set the boundaries for success/failure of this node
func NewParallelMemoryNodeBounded(children[]Node, minSucc int, minFail int) *ParallelMemoryNode{
  n := new(ParallelMemoryNode)
  n.Children = children
  n.Completed = make([]bool, len(children))
  n.MinimumSuccesses = minSucc
  n.MinimumFailures = minFail
  return n
}

// Create a new parallel node with the given children
// success or failure is either triggered by one or all nodes
func NewParallelMemoryNodeAll(children[]Node, successOnAll bool, failOnAll bool) *ParallelMemoryNode{
  n := new(ParallelMemoryNode)
  n.Children = children
  n.Completed = make([]bool, len(children))
  if successOnAll {
    n.MinimumSuccesses = len(children)
  } else {
    n.MinimumSuccesses = 1
  }
  if failOnAll {
    n.MinimumFailures = len(children)
  } else {
    n.MinimumFailures = 1
  }
  return n
}
type MemoryNode struct {
  CurrentIndex int
}

func (n *MemoryNode) Initiate() {
  n.CurrentIndex = 0
}

// Like SequentialNode, but remembers its position
type SequentialMemoryNode struct {
  CompositeNode
  MemoryNode
}

func (n *SequentialMemoryNode) Update() {
  n.Status, n.CurrentIndex = compositeUpdate(&n.CompositeNode, n.CurrentIndex, Success)
}

// Create a new sequential memory node with the given children
func NewSequentialMemoryNode(children[]Node) *SequentialMemoryNode{
  n := new(SequentialMemoryNode)
  n.Children = children
  return n
}

// Like SelectorNode, but remembers its position
type SelectorMemoryNode struct {
  CompositeNode
  MemoryNode
}

func (n *SelectorMemoryNode) Update() {
  n.Status, n.CurrentIndex = compositeUpdate(&n.CompositeNode, n.CurrentIndex, Failure)
}

// Create a new selector node with the given children
func NewSelectorMemoryNode(children[]Node) *SelectorMemoryNode{
  n := new(SelectorMemoryNode)
  n.Children = children
  return n
}

