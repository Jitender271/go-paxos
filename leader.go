package main

//import (
//	"fmt"
//	"sync"
//)
//
//type Leader struct {
//	id           int
//	proposalNum  int
//	acceptors    []*Acceptor
//	learners     []*Learner
//	quorumSize   int
//	proposalsLog map[int]string // Map of slot -> value
//	nextSlot     int            // Next slot to propose
//	mutex        sync.Mutex
//	isLeader     bool
//}
//
//func (l *Leader) ElectLeader() bool {
//	prepareResponses := 0
//	for _, acceptor := range l.acceptors {
//		success, _ := acceptor.Prepare(l.proposalNum)
//		if success {
//			prepareResponses++
//		}
//	}
//	l.isLeader = prepareResponses >= l.quorumSize
//	return l.isLeader
//}
//
//func (l *Leader) Propose(value string) {
//	if !l.isLeader {
//		fmt.Printf("Node %d is not the leader. Cannot propose values.\n", l.id)
//		return
//	}
//
//	l.mutex.Lock()
//	defer l.mutex.Unlock()
//
//	l.proposalNum++
//	proposal := Proposal{
//		Number: l.proposalNum,
//		Slot:   l.nextSlot,
//		Value:  value,
//	}
//
//	acceptResponses := 0
//	for _, acceptor := range l.acceptors {
//		if acceptor.Accept(proposal) {
//			acceptResponses++
//		}
//	}
//
//	if acceptResponses >= l.quorumSize {
//		fmt.Printf("Leader %d: Proposal %d in slot %d with value '%s' was accepted!\n",
//			l.id, proposal.Number, proposal.Slot, proposal.Value)
//		l.proposalsLog[proposal.Slot] = proposal.Value
//		for _, learner := range l.learners {
//			learner.Learn(proposal)
//		}
//		l.nextSlot++
//	} else {
//		fmt.Printf("Leader %d: Proposal %d in slot %d failed to achieve quorum.\n", l.id, proposal.Number, proposal.Slot)
//	}
//}
//
//func (l *Leader) RecoverLeader() bool {
//	// Try to get accepted proposals from all acceptors
//	highestProposal := Proposal{}
//	prepareResponses := 0
//
//	for _, acceptor := range l.acceptors {
//		success, proposal := acceptor.Prepare(l.proposalNum)
//		if success {
//			prepareResponses++
//			if proposal != nil && proposal.Number > highestProposal.Number {
//				highestProposal = *proposal
//			}
//		}
//	}
//
//	if prepareResponses >= l.quorumSize {
//		fmt.Printf("Node %d is elected as the new leader!\n", l.id)
//		l.isLeader = true
//
//		// Recover the latest accepted proposal
//		if highestProposal.Value != "" {
//			l.nextSlot = highestProposal.Slot + 1
//			fmt.Printf("Leader %d recovers proposal in slot %d: '%s'\n",
//				l.id, highestProposal.Slot, highestProposal.Value)
//		}
//		return true
//	}
//
//	fmt.Printf("Node %d failed to recover leadership.\n", l.id)
//	return false
//}
