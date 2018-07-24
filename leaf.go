package behaviortree


// A basic node that takes a predicate function
// that returns true or false based on state
// Node return success or failure
type PredicateLeafNode struct {
  BasicNode
  handler func(state interface{}) bool
}

func (n *PredicateLeafNode) Update(state interface{}, messages []interface{}) []interface{} {
  if n.handler(state) {
    n.Status = Success
  } else {
    n.Status = Failure
  }
  return messages
}

func NewPredicateLeafNode(handler func(state interface{})bool) *PredicateLeafNode {
  n := new(PredicateLeafNode)
  n.handler = handler
  return n
}

// A leaf node that runs a handler in a goroutine
// handler receives ticks though a channel
// and should send over the status channel
type GoroutineLeafNode struct {
  BasicNode
  tickChannel chan interface{}
  statusChannel chan Status
  handler func(<-chan interface{}, chan<- Status)
}

func (n *GoroutineLeafNode) Initiate() {
  n.tickChannel = make(chan interface{})
  n.statusChannel = make(chan Status)
  go n.handler(n.tickChannel, n.statusChannel)
}

func (n *GoroutineLeafNode) Update(state interface{}, messages []interface{}) []interface{} {
  // TODO send messages to goroutine
  n.tickChannel <- state
  n.Status = <-n.statusChannel
  return messages
}

func (n *GoroutineLeafNode) Terminate() {
  close(n.tickChannel)
}

func NewGoroutineLeafNode(handler func(<-chan interface{}, chan<- Status)) *GoroutineLeafNode {
  n := new(GoroutineLeafNode)
  n.handler = handler
  return n
}
