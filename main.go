package main

import (
	"fmt"
	"sync"
)

type StateMachine struct {
	state      map[int]string
	applyIndex int
	mu         sync.Mutex
}

func NewStateMachine() *StateMachine {
	return &StateMachine{
		state:      make(map[int]string),
		applyIndex: 0,
	}
}

func (sm *StateMachine) ApplyValue(index int, value string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if index == sm.applyIndex+1 {
		sm.state[index] = value
		sm.applyIndex = index
		fmt.Printf("StateMachine applied: [%d] -> %s\n", index, value)
	} else {
		fmt.Printf("Out of order log: [%d], expected: [%d]\n", index, sm.applyIndex+1)
	}
}

func (sm *StateMachine) DisplayState() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	fmt.Println("State Machine Current State:")
	for i := 1; i <= sm.applyIndex; i++ {
		fmt.Printf("[%d]: %s\n", i, sm.state[i])
	}
}

type Proposer struct {
	id           int
	proposalNum  int
	leader       bool
	acceptedLogs map[int]string
	mu           sync.Mutex
}

type Acceptor struct {
	id           int
	promisedNum  int
	acceptedNum  int
	acceptedVal  string
	acceptedLogs map[int]string
	mu           sync.Mutex
}

type PaxosState struct {
	proposers []*Proposer
	acceptors []*Acceptor
	quorum    int
}

func newPaxosSystem(proposerCount, acceptorCount int) *PaxosState {
	ps := &PaxosState{
		proposers: []*Proposer{},
		acceptors: []*Acceptor{},
		quorum:    (acceptorCount / 2) + 1,
	}
	for i := 1; i <= proposerCount; i++ {
		ps.proposers = append(ps.proposers, &Proposer{
			id:           i,
			proposalNum:  0,
			leader:       false,
			acceptedLogs: make(map[int]string),
		})
	}
	for i := 1; i <= acceptorCount; i++ {
		ps.acceptors = append(ps.acceptors, &Acceptor{
			id:           i,
			acceptedLogs: make(map[int]string),
		})
	}
	return ps
}

func (p *Proposer) RunPhase1(acceptors []*Acceptor, quorum int) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.proposalNum += 10 + p.id
	prepareNum := p.proposalNum
	promises := 0
	highestAcceptedNum := 0
	highestAcceptedVal := ""

	fmt.Printf("Proposer %d starts Phase 1 (Prepare Num: %d)\n", p.id, prepareNum)

	for _, acc := range acceptors {
		acc.mu.Lock()
		if prepareNum > acc.promisedNum {
			acc.promisedNum = prepareNum
			if acc.acceptedNum > highestAcceptedNum {
				highestAcceptedNum = acc.acceptedNum
				highestAcceptedVal = acc.acceptedVal
			}
			promises++
		}
		acc.mu.Unlock()
	}

	if promises >= quorum {
		p.leader = true
		if highestAcceptedNum > 0 {
			p.acceptedLogs[highestAcceptedNum] = highestAcceptedVal
		}
		fmt.Printf("Proposer %d becomes leader\n", p.id)
		return true
	}

	fmt.Printf("Proposer %d failed to become leader\n", p.id)
	return false
}

func (p *Proposer) ProposeValue(index int, value string, acceptors []*Acceptor, quorum int) bool {
	if !p.leader {
		fmt.Printf("Proposer %d is not the leader. Cannot propose value.\n", p.id)
		return false
	}

	acceptCount := 0
	for _, acc := range acceptors {
		acc.mu.Lock()
		if p.proposalNum >= acc.promisedNum {
			acc.acceptedNum = p.proposalNum
			acc.acceptedLogs[index] = value
			acc.acceptedVal = value
			acceptCount++
		}
		acc.mu.Unlock()
	}

	if acceptCount >= quorum {
		p.acceptedLogs[index] = value
		fmt.Printf("Proposer %d successfully proposed value '%s' at index %d\n", p.id, value, index)
		return true
	}

	fmt.Printf("Proposer %d failed to propose value '%s' at index %d\n", p.id, value, index)
	return false
}

func (p *Proposer) SimulateLeaderFailure() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.leader = false
	fmt.Printf("Leader Proposer %d has failed!\n", p.id)
}

func (p *Proposer) SynchronizeLogs(acceptors []*Acceptor) {
	for _, acc := range acceptors {
		acc.mu.Lock()
		for idx, val := range acc.acceptedLogs {
			if _, exists := p.acceptedLogs[idx]; !exists {
				p.acceptedLogs[idx] = val
			}
		}
		acc.mu.Unlock()
	}
	fmt.Printf("Proposer %d synchronized logs with acceptors.\n", p.id)
}

func main() {
	paxos := newPaxosSystem(2, 5)
	leaderElected := false
	var leader *Proposer

	// Leader Election Phase
	for _, proposer := range paxos.proposers {
		if proposer.RunPhase1(paxos.acceptors, paxos.quorum) {
			leader = proposer
			leaderElected = true
			break
		}
	}

	if !leaderElected {
		fmt.Println("No leader elected. Aborting.")
		return
	}

	// Propose Initial Values
	values := []string{"value1", "value2", "value3"}
	for i, value := range values {
		success := leader.ProposeValue(i+1, value, paxos.acceptors, paxos.quorum)
		if !success {
			fmt.Printf("Proposer %d failed to propose value %s\n", leader.id, value)
		}
	}

	// Simulate Leader Failure
	leader.SimulateLeaderFailure()

	// Elect a New Leader
	newLeaderElected := false
	for _, proposer := range paxos.proposers {
		if proposer.id != leader.id && proposer.RunPhase1(paxos.acceptors, paxos.quorum) {
			leader = proposer
			newLeaderElected = true
			break
		}
	}

	if newLeaderElected {
		fmt.Printf("New leader elected: Proposer %d\n", leader.id)

		// Synchronize logs with acceptors after new leader election
		leader.SynchronizeLogs(paxos.acceptors)
	} else {
		fmt.Println("No new leader elected. System is stalled.")
		return
	}

	// Apply Values to State Machine
	stateMachine := NewStateMachine()
	for i, value := range values {
		if leader.ProposeValue(i+1, value, paxos.acceptors, paxos.quorum) {
			stateMachine.ApplyValue(i+1, value)
		}
	}
	stateMachine.DisplayState()

}
