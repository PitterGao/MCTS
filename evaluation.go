package Mcts

var D = [8][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}
var Board [10][11]int
var TmpChess1 [10][10]int
var TmpChess2 [10][10]int
var head int
var tail int

type node1 [10][10]struct {
	kb int
	kr int
	qb int
	qr int
}

var Map node1

type node2 [110]struct {
	x    int
	y    int
	step int
}

var que node2

func Queen(x int, y int) {
	tx := 0
	ty := 0
	head = 0
	tail = 0
	que[tail].x = x
	que[tail].y = y
	que[tail].step = 0
	tail++
	TmpChess2[x][y] = 1
	for head < tail {
		for k := 0; k < 8; k++ {
			tx = que[head].x
			ty = que[head].y
			for {
				tx = tx + D[k][0]
				ty = ty + D[k][1]
				if tx >= 0 && tx < 10 && ty >= 0 && ty < 10 && TmpChess2[tx][ty] == -1 {
					continue
				}
				if tx >= 0 && tx < 10 && ty >= 0 && ty < 10 && TmpChess2[tx][ty] == 0 {
					TmpChess2[tx][ty] = -1
					que[tail].x = tx
					que[tail].y = ty
					que[tail].step = que[head].step + 1
					tail++
				} else {
					break
				}
			}
		}
		head++
	}
	for k := 1; k < tail; k++ {
		if Board[x][y] == -1 && Map[que[k].x][que[k].y].qb > que[k].step {
			Map[que[k].x][que[k].y].qb = que[k].step
		}
		if Board[x][y] == 1 && Map[que[k].x][que[k].y].qr > que[k].step {
			Map[que[k].x][que[k].y].qr = que[k].step
		}
	}
}

func King(x int, y int) {
	tx := 0
	ty := 0
	head = 0
	tail = 0
	que[tail].x = x
	que[tail].y = y
	que[tail].step = 0
	tail++
	for head < tail {
		for k := 0; k < 8; k++ {
			tx = que[head].x + D[k][0]
			ty = que[head].y + D[k][1]
			if tx < 0 || tx > 9 || ty < 0 || ty > 9 || TmpChess1[tx][ty] != 0 {
				continue
			} else {
				TmpChess1[tx][ty] = 1
				que[tail].x = tx
				que[tail].y = ty
				que[tail].step = que[head].step + 1
				tail++
			}
		}
		head++
	}
	for k := 1; k < tail; k++ {
		if Board[x][y] == -1 && Map[que[k].x][que[k].y].kb > que[k].step {
			Map[que[k].x][que[k].y].kb = que[k].step
		}
		if Board[x][y] == 1 && Map[que[k].x][que[k].y].kr > que[k].step {
			Map[que[k].x][que[k].y].kr = que[k].step
		}
	}

}

func evaluation(n *Node) int {
	queen := 0
	king := 0

	for i := 0; i < 100; i++ {
		x := i / 10
		y := i % 10
		if n.State.Board[i] == 0 {
			Board[x][y] = 0
		} else if n.State.Board[i] == 2 {
			Board[x][y] = 2
		} else if n.State.Board[i] == 1 {
			Board[x][y] = 1
		} else if n.State.Board[i] == -1 {
			Board[x][y] = -1
		}
	}

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			Map[i][j].kr = 10000
			Map[i][j].qr = 10000
			Map[i][j].kb = 10000
			Map[i][j].qb = 10000
		}
	}

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if Board[i][j] == 1 || Board[i][j] == -1 {
				for k := 0; k < 10; k++ {
					for t := 0; t < 10; t++ {
						if Board[k][t] == 0 {
							TmpChess1[k][t] = 0
							TmpChess2[k][t] = 0
						} else {
							TmpChess1[k][t] = 1
							TmpChess2[k][t] = 1
						}
					}
				}
				Queen(i, j)
				King(i, j)
			}
		}
	}

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if Board[i][j] == 0 {
				if Map[i][j].kb < Map[i][j].kr {
					king++
				}
				if Map[i][j].kb > Map[i][j].kr {
					king--
				}
				if Map[i][j].qb < Map[i][j].qr {
					queen++
				}
				if Map[i][j].qb > Map[i][j].qr {
					queen--
				}
			}
		}
	}
	return queen + king
}
