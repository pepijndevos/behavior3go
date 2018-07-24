package behaviortree

// A leaf node that runs a handler in a goroutine
// handler receives ticks though a channel
// and should send over the status channel
type GoroutineLeafNode struct {
  BasicNode
  tickChannel chan struct{}
  statusChannel chan Status
  handler func(<-chan struct{}, chan<- Status)
}

func (n *GoroutineLeafNode) Initiate() {
  n.tickChannel = make(chan struct{})
  n.statusChannel = make(chan Status)
  go n.handler(n.tickChannel, n.statusChannel)
}

func (n *GoroutineLeafNode) Update() {
  n.tickChannel <- struct{}{}
  n.Status = <-n.statusChannel
}

func (n *GoroutineLeafNode) Terminate() {
  close(n.tickChannel)
}

func NewGoroutineLeafNode(handler func(<-chan struct{}, chan<- Status)) *GoroutineLeafNode {
  n := new(GoroutineLeafNode)
  n.handler = handler
  return n
}
