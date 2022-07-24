package Mcts

import (
	"errors"
	amazonsChess "github.com/PitterGao/Regulation"
	"log"
)

type Mcts struct {
	tree   *Node
	deepth int
}

func NewMcts(tree *Node) (*Mcts, error) {
	if tree == nil {
		return nil, errors.New("error")
	}
	return &Mcts{tree: tree, deepth: 2}, nil
}

func AI(s *amazonsChess.State) amazonsChess.ChessMove {
	root := &Node{
		Children:    make([]*Node, 0),
		State:       s,
		FullyExpand: false,
	}
	m, err := NewMcts(root)
	if err != nil {
		log.Fatal(err)
	}
	bestChild := m.getBestChild(s, 1000)
	return bestChild.Move
}

func (m *Mcts) getBestChild(s *amazonsChess.State, loop int) *Node {
	for i := 0; i < loop; i++ {

		//t := time.Now()
		m.tree = m.search(s)
		//elapsed := time.Now().Sub(t)
		//fmt.Println("该函数执行完成耗时：", elapsed)

	}
	bestChild := m.tree.MostVisitedChild()
	return bestChild
}

func (m *Mcts) search(s *amazonsChess.State) *Node {
	Wins := 0.0

	node := m.Select()
	node = m.expand(node)
	Wins = m.rollout(s, node)
	m.tree = m.backupdate(Wins, node)

	return m.tree
}

func (m *Mcts) Select() *Node {
	currentDeepth := 0
	n := m.tree
	for len(n.Children) > 0 {
		n = n.SelectMostPromisingNode()
		currentDeepth++
		if currentDeepth >= m.deepth {
			break
		}
	}
	return n
}

func (m *Mcts) expand(n *Node) *Node {
	if n.State.GameOver() == 0 {
		if !n.FullyExpand {
			move := n.State.GetValid()
			for _, moves := range move {
				State, err := n.State.StateMove(moves)
				if err != nil {
					log.Fatal(err)
				}
				newNode := &Node{
					Parent:      n,
					State:       State,
					Move:        moves,
					FullyExpand: false,
				}
				n.Children = append(n.Children, newNode)
			}
			if len(move) == len(n.Children) {
				n.FullyExpand = true
			}
		}
	}
	//for k, nodes := range n.Children {
	//	n.Children[k].Eva = evaluation(nodes)
	//}
	//sort.Slice(n.Children, func(i, j int) bool {
	//	return n.Children[i].Eva > n.Children[j].Eva
	//})
	//if len(n.Children) > 400 {
	//	n.Children = n.Children[:400]
	//}

	return n
}

func (m *Mcts) rollout(s *State, n *node) float64 {
	res := node{
		State: n.State,
	}
	var err error
	for {
		if res.State.GameOver() != 0 {
			if n.State.CurrentPlayer == -res.State.GameOver() {
				return 1
			} else {
				return -1
			}
		}
		res.State, err = res.State.RandomMove()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (m *Mcts) backupdate(Wins float64, n *node) *node {
	n.Visits++
	n.Q += 1.0 * (Wins - n.Q) / float64(n.Visits)
	for n.Parent != nil {
		Wins = -Wins
		n = n.Parent
		n.Visits++
		n.Q += 1.0 * (Wins - n.Q) / float64(n.Visits)
	}
	return n
}
