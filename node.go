package Mcts

import (
	amazonsChess "github.com/PitterGao/Regulation"
	"math"
)

type Node struct {
	Parent      *Node
	Children    []*Node
	Move        amazonsChess.ChessMove
	State       *amazonsChess.State
	Q           float64
	Visits      int
	FullyExpand bool
	P           float64
	Eva         int
}

func (n *Node) SelectMostPromisingNode() *Node {
	bestChild := n.Children[0]
	bestScore := n.Children[0].Q + math.Sqrt(2.0*math.Log(float64(n.Visits))/float64(n.Children[0].Visits))
	var score float64
	for _, child := range n.Children {
		Q := child.Q
		cVisits := float64(child.Visits)
		visits := float64(n.Visits)
		if cVisits == 0 {
			return child
		} else {
			score = Q + math.Sqrt(2.0*math.Log(visits)/cVisits)
		}
		if score > bestScore {
			bestChild = child
			bestScore = score
		}
	}
	return bestChild
}

func (n *Node) MostVisitedChild() *Node {
	var mostVisits int
	var mostVisited *Node
	for _, child := range n.Children {
		if child.Visits > mostVisits {
			mostVisits = child.Visits
			mostVisited = child
		}
	}
	if mostVisited != nil {
		return mostVisited
	}
	return nil
}

func (n *Node) Probably() {
	n.P = 1
}
