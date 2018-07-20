package behaviortree

import (
  "fmt"
  "time"
  "io"
  "encoding/json"
)

type Project struct {
  Name string
  Data struct {
    Trees []struct {
      Title string
      Root string
      Nodes map[string] ProjectNode
    }
  }
}

type ProjectNode struct {
  Id string
  Name string
  Properties map[string]interface{}
  Child string
  Children []string
}

var NodeTypeRegister = make(map[string]func(ProjectNode, map[string]ProjectNode)Node)

func ReadProject(file io.Reader) (*Project, error) {
  var pr Project
  dec := json.NewDecoder(file)
  err := dec.Decode(&pr)
  if err != nil {
    return nil, err
  } else {
    return &pr, nil
  }
}

func MakeTrees(pr *Project) []Node {
  trees := make([]Node, len(pr.Data.Trees))
  for idx, tree := range pr.Data.Trees {
    node, _ := MakeNode(tree.Root, tree.Nodes)
    trees[idx] = node
  }
  return trees
}

func MakeNode(root string, nodes map[string]ProjectNode) (Node, bool) {
  node := nodes[root]
  fn, ok := NodeTypeRegister[node.Name]
  if ok {
    return fn(node, nodes), true
  } else {
    fmt.Printf("No constructor for %s\n", node.Name)
    return NewConstantNode(Failure), false
  }
}

func init() {
  // Composite nodes
  NodeTypeRegister["Priority"] = func(root ProjectNode, nodes map[string]ProjectNode)Node {
    children := make([]Node, len(root.Children))
    for idx, child := range root.Children {
      children[idx], _ = MakeNode(child, nodes)
    }
    return NewSelectorNode(children)
  }

  NodeTypeRegister["MemPriority"] = func(root ProjectNode, nodes map[string]ProjectNode)Node {
    children := make([]Node, len(root.Children))
    for idx, child := range root.Children {
      children[idx], _ = MakeNode(child, nodes)
    }
    return NewSelectorMemoryNode(children)
  }

  NodeTypeRegister["Sequence"] = func(root ProjectNode, nodes map[string]ProjectNode)Node {
    children := make([]Node, len(root.Children))
    for idx, child := range root.Children {
      children[idx], _ = MakeNode(child, nodes)
    }
    return NewSequentialNode(children)
  }

  NodeTypeRegister["MemSequence"] = func(root ProjectNode, nodes map[string]ProjectNode)Node {
    children := make([]Node, len(root.Children))
    for idx, child := range root.Children {
      children[idx], _ = MakeNode(child, nodes)
    }
    return NewSequentialMemoryNode(children)
  }

  NodeTypeRegister["ParallelSequence"] = func(root ProjectNode, nodes map[string]ProjectNode)Node {
    children := make([]Node, len(root.Children))
    for idx, child := range root.Children {
      children[idx], _ = MakeNode(child, nodes)
    }
    // This WILL panic if you don't set bounds
    minSuccess := int(root.Properties["minSuccess"].(float64))
    minFail    := int(root.Properties["minFail"].(float64))
    return NewParallelNodeBounded(children, minSuccess, minFail)
  }

  NodeTypeRegister["ParallelTactic"] = func(root ProjectNode, nodes map[string]ProjectNode)Node {
    children := make([]Node, len(root.Children))
    for idx, child := range root.Children {
      children[idx], _ = MakeNode(child, nodes)
    }
    // This WILL panic if you don't set bounds
    minSuccess := int(root.Properties["minSuccess"].(float64))
    minFail    := int(root.Properties["minFail"].(float64))
    return NewParallelMemoryNodeBounded(children, minSuccess, minFail)
  }

  // Decorator nodes
  NodeTypeRegister["Inverter"] = func(root ProjectNode, nodes map[string]ProjectNode)Node {
    child, _ := MakeNode(root.Child, nodes)
    return NewInverterNode(child)
  }

  NodeTypeRegister["FailerDec"] = func(root ProjectNode, nodes map[string]ProjectNode)Node {
    child, _ := MakeNode(root.Child, nodes)
    return NewWrapConstantNode(Failure, child)
  }

  NodeTypeRegister["SucceederDec"] = func(root ProjectNode, nodes map[string]ProjectNode)Node {
    child, _ := MakeNode(root.Child, nodes)
    return NewWrapConstantNode(Success, child)
  }

  NodeTypeRegister["Repeat"] = func(root ProjectNode, nodes map[string]ProjectNode)Node {
    child, _ := MakeNode(root.Child, nodes)
    //TODO extract limit if provided
    return NewRepeaterNode(-1, child)
  }

  NodeTypeRegister["RepeatUntilSuccess"] = func(root ProjectNode, nodes map[string]ProjectNode)Node {
    child, _ := MakeNode(root.Child, nodes)
    return NewRepeatUntilNode(Success, child)
  }

  NodeTypeRegister["RepeatUntilFailure"] = func(root ProjectNode, nodes map[string]ProjectNode)Node {
    child, _ := MakeNode(root.Child, nodes)
    return NewRepeatUntilNode(Failure, child)
  }

  // Utility nodes
  NodeTypeRegister["Failer"] = func(root ProjectNode, nodes map[string]ProjectNode)Node {
    return NewConstantNode(Failure)
  }

  NodeTypeRegister["Succeeder"] = func(root ProjectNode, nodes map[string]ProjectNode)Node {
    return NewConstantNode(Success)
  }

  NodeTypeRegister["Sleep"] = func(root ProjectNode, nodes map[string]ProjectNode)Node {
    fmt.Println(root.Properties)
    ms := time.Duration(root.Properties["ms"].(float64))*time.Millisecond
    return NewTimeoutNode(ms, Success, NewConstantNode(Running))
  }
}

