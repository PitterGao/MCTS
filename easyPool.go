package Mcts

import (
	"errors"
	amazonsChess "github.com/PitterGao/Regulation"
	"sync"
)

type task struct {
	rollout    func(s *amazonsChess.State, n *Node) float64
	backupdate func(Wins float64, n *Node)
	n          *Node
	s          *amazonsChess.State
}

type goPool struct {
	taskChannel    chan *task
	maxWaitingNum  int
	blockThreshold int
	passThreshold  int
	blockState     bool
	openState      bool
	waitingTask    int
	wg             sync.WaitGroup

	waitingtaskLock sync.Mutex
}

func NewPool(maxWaitingNum, blockThreshold, passThreshold int) *goPool {
	pool := &goPool{
		maxWaitingNum:  maxWaitingNum,
		taskChannel:    make(chan *task, maxWaitingNum),
		blockThreshold: blockThreshold, //最大线程数
		passThreshold:  passThreshold,  //最小线程数
		blockState:     false,          //线程池阻塞
		openState:      true,           //线程池打开
		waitingTask:    0,              //任务队列任务数
	}
	go pool.scheduler()
	return pool
}

func (p *goPool) Add(rollout func(s *amazonsChess.State, n *Node) float64, backupdate func(Wins float64, n *Node), s *amazonsChess.State, n *Node) error {
	if !p.openState {
		return errors.New("pool closed")
	}

	for p.blockState { //任务队列阻塞就不添加
	}

	p.taskChannel <- &task{
		rollout:    rollout,
		backupdate: backupdate,
		n:          n,
		s:          s,
	}

	p.waitingtaskLock.Lock()
	p.waitingTask++

	if p.waitingTask >= p.blockThreshold {
		p.blockState = true
	}
	p.waitingtaskLock.Unlock()

	return nil
}

func (p *goPool) Close() {

	p.openState = false
	p.wg.Wait()
}

func (p *goPool) scheduler() {
	for {
		if !p.openState {
			if cap(p.taskChannel) == 0 {
				p.wg.Wait()
				return
			}
		}
		select {
		case task := <-p.taskChannel:
			p.wg.Add(1)
			go func() {
				//wins := task.rollout(s, n)
				//task.backupdate(wins, n)
				task.backupdate(task.rollout(task.s, task.n), task.n)
				p.wg.Done()
			}() //go worker
			p.waitingtaskLock.Lock()
			p.waitingTask--
			if p.passThreshold >= p.waitingTask {
				p.blockState = false
			}

			p.waitingtaskLock.Unlock()

		}
	}
}

//func (m *Mcts) threadpool(s *State, n *Node) {
//	var p *goPool
//	p = new(goPool)
//	p = NewPool(30, 20, 10, s, n)
//	err := p.add(m.Rollout, m.Backupdate, n)
//	if err != nil {
//		log.Fatal(err)
//	}
//	p.Close()
//
//}
