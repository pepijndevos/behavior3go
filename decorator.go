package behaviortree

import "time"

type Decorator struct {
  child Node
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
  status := Tick(n.child)
  switch status {
  case Success:
    n.status = Failure
  case Failure:
    n.status = Success
  default:
    n.status = status
  }
}

func NewInverterNode(child Node) *InverterNode {
  n := new(InverterNode)
  n.child = child
  return n
}

type TimeoutNode struct {
  BasicNode
  Decorator
  timeout time.Duration
  tchan <-chan time.Time
}

func (n *TimeoutNode) Initiate() {
  n.tchan = time.After(n.timeout)
}

func (n *TimeoutNode) Update() {
  select {
  case <-n.tchan:
    n.status = Failure
  default:
    n.status = Tick(n.child)
  }
}

func NewTimeoutNode(child Node, timeout time.Duration) *TimeoutNode {
  n := new(TimeoutNode)
  n.child = child
  n.timeout = timeout
  return n
}
