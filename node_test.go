package behaviortree

import (
  "testing"
  "log"
  "io/ioutil"
  "time"
)

// Disable logging
func init() {
  log.SetOutput(ioutil.Discard)
}

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

type PanicNode struct {
  BasicNode
}

func (n *PanicNode) Update() {
  panic("welp")
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

func TestPanic(t *testing.T) {
  n := new(PanicNode)
  status := Tick(n)
  if status != Failure {
    t.Errorf("Status is %s", status)
  }
}

func TestArrayLeaf(t *testing.T) {
  seq := []Status{Running, Success, Failure}
  n := NewArrayLeafNode(t, "name", seq)
  expectSequence(t, n, seq)
}

func TestGoroutineLeaf(t *testing.T) {
  n := NewGoroutineLeafNode(func(ticks <-chan struct{}, status chan<- Status) {
    i := 0
    for range ticks {
      t.Logf("go %d\n", i)
      if i < 2 {
        status <- Running
      } else {
        status <- Success
        close(status)
      }
      i++
    }
  })
  seq := []Status{Running, Running, Success}
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

func TestParallelMemory(t *testing.T) {
  ch := []Node{
    NewArrayLeafNode(t, "mempar 1", []Status{Success, Failure}),
    NewArrayLeafNode(t, "mempar 2", []Status{Running, Running, Success}),
  }
  n := NewParallelMemoryNodeAll(ch, true, false)

  expected := []Status{Running, Running, Success, Failure, Running, Success}
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

func TestInverter(t *testing.T) {
  seq := []Status{Running, Success, Failure}
  n := NewInverterNode(NewArrayLeafNode(t, "inv", seq))
  expected := []Status{Running, Failure, Success}
  expectSequence(t, n, expected)
}

func TestConstantWrap(t *testing.T) {
  seq := []Status{Failure}
  n := NewWrapConstantNode(Success, NewArrayLeafNode(t, "const", seq))
  expected := []Status{Success, Success}
  expectSequence(t, n, expected)
}

func TestRepeater(t *testing.T) {
  seq := []Status{Success}
  n := NewRepeaterNode(2, NewArrayLeafNode(t, "rep", seq))
  expected := []Status{Running, Running, Success}
  expectSequence(t, n, expected)
}

func TestUntilSuccess(t *testing.T) {
  seq := []Status{Failure, Failure, Success}
  n := NewUntilSuccessNode(NewArrayLeafNode(t, "suc", seq))
  expected := []Status{Running, Running, Success}
  expectSequence(t, n, expected)
}

func TestUntilFailure(t *testing.T) {
  seq := []Status{Success, Success, Failure}
  n := NewUntilFailureNode(NewArrayLeafNode(t, "fail", seq))
  expected := []Status{Running, Running, Success}
  expectSequence(t, n, expected)
}

func TestTimeout(t *testing.T) {
  seq := []Status{Running}
  n := NewTimeoutNode(time.Millisecond, Failure, NewArrayLeafNode(t, "tout", seq))
  expected := []Status{Running, Running}
  expectSequence(t, n, expected)
  time.Sleep(time.Millisecond)
  expected = []Status{Failure, Running}
  expectSequence(t, n, expected)
}
