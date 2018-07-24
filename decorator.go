package behaviortree

import "time"

type Decorator struct {
  Child Node
}

// A node that turns failure into success
// Do you want one for your life?
// Be carefull what you ask for,
// Because it also turns success into failure
type InverterNode struct {
  BasicNode
  Decorator
}

func (n *InverterNode) Update() {
  status := Tick(n.Child)
  switch status {
  case Success:
    n.Status = Failure
  case Failure:
    n.Status = Success
  default:
    n.Status = status
  }
}

func NewInverterNode(child Node) *InverterNode {
  n := new(InverterNode)
  n.Child = child
  return n
}

// Runs the child and always returns the same status
type WrapConstantNode struct {
  BasicNode
  Decorator
}

func (n *WrapConstantNode) Update() {
  Tick(n.Child)
}

func NewWrapConstantNode(status Status, child Node) *WrapConstantNode {
  n := new(WrapConstantNode)
  n.Child = child
  n.Status = status
  return n
}

// Runs the child until limit is reached
type RepeaterNode struct {
  BasicNode
  Decorator
  Counter int
  Limit int
}

func (n *RepeaterNode) Initiate() {
  n.Counter = 0
}

func (n *RepeaterNode) Update() {
  status := Tick(n.Child)
  if status != Running {
    n.Counter++
  }
  if n.Limit < 1 || n.Counter < n.Limit {
    n.Status = Running
  } else {
    n.Status = status
  }
}

func NewRepeaterNode(limit int, child Node) *RepeaterNode {
  n := new(RepeaterNode)
  n.Child = child
  n.Limit = limit
  return n
}

// Repeat Until the given status
type RepeatUntilNode struct {
  BasicNode
  Decorator
  Until Status
}

func (n *RepeatUntilNode) Update() {
  status := Tick(n.Child)
  if status == n.Until {
    n.Status = Success
  } else {
    n.Status = Running
  }
}

func NewRepeatUntilNode(until Status, child Node) *RepeatUntilNode {
  n := new(RepeatUntilNode)
  n.Child = child
  n.Until = until
  return n
}

type TimeoutNode struct {
  BasicNode
  Decorator
  Timeout time.Duration
  tchan <-chan time.Time
  Completion Status
}

func (n *TimeoutNode) Initiate() {
  n.tchan = time.After(n.Timeout)
}

func (n *TimeoutNode) Update() {
  select {
  case <-n.tchan:
    n.Status = n.Completion
  default:
    n.Status = Tick(n.Child)
  }
}

func NewTimeoutNode(timeout time.Duration, completion Status, child Node) *TimeoutNode {
  n := new(TimeoutNode)
  n.Child = child
  n.Timeout = timeout
  n.Completion = completion
  return n
}
