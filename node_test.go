package behaviortree

import (
  "testing"
)

// A debugging node that returns a fixed sequence of statusses
// Also print its name and status when updated
type ArrayLeafNode struct {
  BasicNode
  t *testing.T
  name string
  statuses []Status
  counter int
}

func (n *ArrayLeafNode) Update() {
  n.status = n.statuses[n.counter%len(n.statuses)]
  n.t.Logf("%s: %s", n.name, n.status)
  n.counter++
}

// Creat a new ArrayLeafNode with a given name
// and an array of statuses to loop through
func NewArrayLeafNode(t *testing.T, name string, statuses []Status) *ArrayLeafNode {
  n := new(ArrayLeafNode)
  n.t = t
  n.name = name
  n.statuses = statuses
  return n
}

func expectSequence(t *testing.T, node Node, statuses[]Status) {
  for idx, status := range statuses {
    t.Logf("---- tick ----")
    result := Tick(node)
    if status != result {
      t.Errorf("Status is %s, expected %s at index %d", result, status, idx)
    }
  }
}


func TestConstant(t *testing.T) {
  n := NewConstantNode(Success)
  status := Tick(*n)
  if status != Success {
    t.Errorf("Status is %s", status)
  }
}

func TestArrayLeaf(t *testing.T) {
  seq := []Status{Running, Success, Failure}
  n := NewArrayLeafNode(t, "name", seq)
  expectSequence(t, n, seq)
}

func TestSelector(t *testing.T) {
  seq := []Status{Success, Failure}
  ch := []Node{
    NewArrayLeafNode(t, "sel 1", seq),
    NewArrayLeafNode(t, "sel 2", seq),
  }
  n := NewSelectorNode(ch)

  expected := []Status{Success, Success, Success, Failure}
  expectSequence(t, n, expected)
}

func TestSequential(t *testing.T) {
  seq := []Status{Success, Failure}
  ch := []Node{
    NewArrayLeafNode(t, "seq 1", seq),
    NewArrayLeafNode(t, "seq 2", seq),
  }
  n := NewSequentialNode(ch)

  expected := []Status{Success, Failure, Failure, Failure, Success, Failure}
  expectSequence(t, n, expected)
}

func TestParallel(t *testing.T) {
  ch := []Node{
    NewArrayLeafNode(t, "par 1", []Status{Success, Failure}),
    NewArrayLeafNode(t, "par 2", []Status{Running, Success, Failure}),
  }
  n := NewParallelNodeAll(ch, true, false)

  expected := []Status{Running, Failure, Failure, Failure, Success}
  expectSequence(t, n, expected)
}

func TestSequentialMemory(t *testing.T) {
  seq := []Status{Running, Success, Failure}
  ch := []Node{
    NewArrayLeafNode(t, "memseq 1", seq),
    NewArrayLeafNode(t, "memseq 2", seq),
  }
  n := NewSequentialMemoryNode(ch)

  expected := []Status{Running, Running, Success, Failure, Running, Failure}
  expectSequence(t, n, expected)
}

func TestSelectorMemory(t *testing.T) {
  seq := []Status{Running, Success, Failure}
  ch := []Node{
    NewArrayLeafNode(t, "memsel 1", seq),
    NewArrayLeafNode(t, "memsel 2", seq),
  }
  n := NewSelectorMemoryNode(ch)

  expected := []Status{Running, Success, Running, Success, Running, Success}
  expectSequence(t, n, expected)
}
