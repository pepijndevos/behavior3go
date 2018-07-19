package behaviortree

import (
  "fmt"
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

func init() {
  NodeTypeRegister["Priority"] = func(root ProjectNode, nodes map[string]ProjectNode)Node {
    children := make([]Node, len(root.Children))
    for idx, child := range root.Children {
      fmt.Printf("Making %s\n", child)
      children[idx], _ = MakeNode(child, nodes)
    }
    return NewSelectorNode(children)
  }
}

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