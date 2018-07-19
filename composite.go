package behaviortree

type CompositeNode struct {
  BasicNode
  children []Node
}

func (n *CompositeNode) Terminate() {
  for _, child := range n.children {
    child.Terminate()
  }
}

// Generic function for (memory) sequential an selector nodes
func compositeUpdate(n *CompositeNode, currentIndex int, endCondition Status) (Status, int) {
  for ; currentIndex<len(n.children); currentIndex++ {
    status := Tick(n.children[currentIndex])
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
  n.status, _ = compositeUpdate(&n.CompositeNode, 0, Failure)
}

// Create a new selector node with the given children
func NewSelectorNode(children[]Node) *SelectorNode{
  n := new(SelectorNode)
  n.children = children
  return n
}

// A node that runs all childs until one fails
type SequentialNode struct {
  CompositeNode
}

func (n *SequentialNode) Update() {
  n.status, _ = compositeUpdate(&n.CompositeNode, 0, Success)
}

// Create a new sequential node with the given children
func NewSequentialNode(children[]Node) *SequentialNode{
  n := new(SequentialNode)
  n.children = children
  return n
}

// A node that runs all children
// Success or Failure is defined by
// the number of children that fail or succeed
type ParallelNode struct {
  CompositeNode
  minimumSuccesses int
  minimumFailures int
}

func (n *ParallelNode) Update() {
  totalFailures := 0
  totalSuccesses := 0
  for _, child := range n.children {
    status := Tick(child)
    if status == Success {
      totalSuccesses++
    } else if status == Failure {
      totalFailures++
    }
  }
  if totalSuccesses >= n.minimumSuccesses {
    n.status = Success;
  } else if totalFailures >= n.minimumFailures {
    n.status = Failure;
  } else {
    n.status = Running
  }
}

// Create a new parallel node with the given children
// minSucc and minFail set the boundaries for success/failure of this node
func NewParallelNodeBounded(children[]Node, minSucc int, minFail int) *ParallelNode{
  n := new(ParallelNode)
  n.children = children
  n.minimumSuccesses = minSucc
  n.minimumFailures = minFail
  return n
}

// Create a new parallel node with the given children
// success or failure is either triggered by one or all nodes
func NewParallelNodeAll(children[]Node, successOnAll bool, failOnAll bool) *ParallelNode{
  n := new(ParallelNode)
  n.children = children
  if successOnAll {
    n.minimumSuccesses = len(children)
  } else {
    n.minimumSuccesses = 1
  }
  if failOnAll {
    n.minimumFailures = len(children)
  } else {
    n.minimumFailures = 1
  }
  return n
}

type MemoryNode struct {
  currentIndex int
}

func (n *MemoryNode) Initiate() {
  n.currentIndex = 0
}

// Like SequentialNode, but remembers its position
type SequentialMemoryNode struct {
  CompositeNode
  MemoryNode
}

func (n *SequentialMemoryNode) Update() {
  n.status, n.currentIndex = compositeUpdate(&n.CompositeNode, n.currentIndex, Success)
}

// Create a new sequential memory node with the given children
func NewSequentialMemoryNode(children[]Node) *SequentialMemoryNode{
  n := new(SequentialMemoryNode)
  n.children = children
  return n
}

// Like SelectorNode, but remembers its position
type SelectorMemoryNode struct {
  CompositeNode
  MemoryNode
}

func (n *SelectorMemoryNode) Update() {
  n.status, n.currentIndex = compositeUpdate(&n.CompositeNode, n.currentIndex, Failure)
}

// Create a new selector node with the given children
func NewSelectorMemoryNode(children[]Node) *SelectorMemoryNode{
  n := new(SelectorMemoryNode)
  n.children = children
  return n
}

