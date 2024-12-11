package main

//import (
//	"fmt"
//	"sync"
//)
//
//type Learner struct {
//	id           int
//	mutex        sync.Mutex
//	learnLog     map[int]string // Map of slot -> value
//	stateMachine []string       // State machine applying commands
//	nextExpected int            // Next expected slot for sequential application
//}
//
//func NewLearner(id int) *Learner {
//	return &Learner{
//		id:           id,
//		mutex:        sync.Mutex{},
//		learnLog:     make(map[int]string),
//		stateMachine: []string{},
//		nextExpected: 0,
//	}
//}
//
//func (l *Learner) Learn(proposal Proposal) {
//	l.mutex.Lock()
//	defer l.mutex.Unlock()
//
//	fmt.Printf("Learner %d: Learned proposal %d in slot %d with value '%s'\n",
//		l.id, proposal.Number, proposal.Slot, proposal.Value)
//
//	// Store the learned proposal
//	l.learnLog[proposal.Slot] = proposal.Value
//
//	// Apply sequentially to the state machine
//	for {
//		value, exists := l.learnLog[l.nextExpected]
//		if !exists {
//			break
//		}
//		l.stateMachine = append(l.stateMachine, value)
//		fmt.Printf("Learner %d: Applied value '%s' to state machine (slot %d)\n",
//			l.id, value, l.nextExpected)
//		l.nextExpected++
//	}
//}
