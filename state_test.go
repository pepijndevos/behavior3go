package behaviortree

import (
  "testing"
  "log"
  "io/ioutil"
)

// Disable logging
func init() {
  log.SetOutput(ioutil.Discard)
}

func expectMessageSequence(t *testing.T, node Node, states []interface{}, expectedMessages [][]interface{}, expectedStatuses []Status) {
  for idx, state := range states {
    t.Logf("---- tick ----")
    status, messages := Tick(node, state, nil)
    if status != expectedStatuses[idx] {
      t.Errorf("Status is %s, expected %s at index %d", status, expectedStatuses[idx], idx)
    }
    if len(messages) != len(expectedMessages[idx]) {
      t.Errorf("Number of messages differ:\n%v\n%v\n", messages, expectedMessages[idx])
    } else {
      for midx := range messages {
        if messages[midx] != expectedMessages[idx][midx] {
          t.Errorf("Messages differ:\n%v\n%v\n", messages[midx], expectedMessages[idx][midx])
        }
      }
    }
  }
}

func NewModulusNode(t *testing.T, mod int) *PredicateLeafNode {
  return NewPredicateLeafNode(
    func (state interface{}) bool {
      //t.Logf("%d %% %d = %d == 0: %t", state.(int), mod, state.(int) % mod, state.(int) % mod == 0)
      t.Logf("%d %% %d", state.(int), mod)
      return state.(int) % mod == 0
    })
}

type MessageNode struct {
  BasicNode
  t *testing.T
  Message interface{}
}

func (n *MessageNode) Update(state interface{}, messages []interface{}) []interface{} {
  n.t.Logf("%s", n.Message)
  return append(messages, n.Message)
}

func NewMessageNode(t *testing.T, message interface{}) *MessageNode {
  n := new(MessageNode)
  n.Message = message
  n.Status = Success
  n.t = t
  return n
}

type ConcatNode struct {
  BasicNode
  t *testing.T
}

func (n *ConcatNode) Update(state interface{}, messages []interface{}) []interface{} {
  s := ""
  for _, m := range messages {
    s += m.(string)
  }
  n.t.Logf(s)
  return []interface{}{s}
}

func NewConcatNode(t *testing.T) *ConcatNode {
  n := new(ConcatNode)
  n.Status = Success
  n.t = t
  return n
}

func TestPredicate(t *testing.T) {
  n := NewPredicateLeafNode(
    func (state interface{}) bool {
      return true
    })
  sep := []Status{Success, Success}
  expectSequence(t, n, sep)
  m := NewPredicateLeafNode(
    func (state interface{}) bool {
      return false
    })
  seq := []Status{Failure, Failure}
  expectSequence(t, m, seq)
}

func TestMessage(t *testing.T) {
  n := NewMessageNode(t, "Hello")
  expectMessageSequence(t, n,
    []interface{}{1,2},
    [][]interface{}{{"Hello"}, {"Hello"}},
    []Status{Success,Success},
  )
}

func TestSequentialMessage(t *testing.T) {
  n := NewSequentialNode([]Node{
    NewMessageNode(t, "Hello"),
    NewMessageNode(t, "World"),
  })
  expectMessageSequence(t, n,
    []interface{}{1},
    [][]interface{}{{"Hello", "World"}},
    []Status{Success,Success},
  )
}

func TestFizzBuzz(t *testing.T) {
  n := NewSequentialNode([]Node{
    NewParallelNodeAll(false, true, []Node{
      NewSequentialNode([]Node{
        NewModulusNode(t, 3),
        NewMessageNode(t, "Fizz"),
      }),
      NewSequentialNode([]Node{
        NewModulusNode(t, 5),
        NewMessageNode(t, "Buzz"),
      }),
    }),
    NewConcatNode(t),
  })
  expectMessageSequence(t, n,
    []interface{}{1,2,3,4,5,6,7,8,9,10,11,12,13,14,15},
    [][]interface{}{{},{},{"Fizz"},{},{"Buzz"},{"Fizz"},{},{},{"Fizz"},{"Buzz"},{},{"Fizz"},{},{},{"FizzBuzz"}},
    []Status{Failure,Failure,Success,Failure,Success,Success,Failure,Failure,Success,Success,Failure,Success,Failure,Failure,Success},
  )
}
