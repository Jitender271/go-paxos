package main

//import (
//	"fmt"
//	"sync"
//)
//
//type Proposal struct {
//	Number int
//	Slot   int    // Command slot (sequence number)
//	Value  string // Command value
//}
//
//type Acceptor struct {
//	mutex            sync.Mutex
//	id               int
//	promisedNumber   int
//	acceptedProposal *Proposal
//	recovered        bool
//}
//
//func (a *Acceptor) Prepare(proposalNumber int) (bool, *Proposal) {
//	a.mutex.Lock()
//	defer a.mutex.Unlock()
//
//	if proposalNumber > a.promisedNumber {
//		a.promisedNumber = proposalNumber
//		return true, a.acceptedProposal
//	}
//	return false, nil
//}
//
//func (a *Acceptor) Accept(proposal Proposal) bool {
//	a.mutex.Lock()
//	defer a.mutex.Unlock()
//
//	if proposal.Number >= a.promisedNumber {
//		a.promisedNumber = proposal.Number
//		a.acceptedProposal = &proposal
//		return true
//	}
//	return false
//}
//
//func (a *Acceptor) RecoverState(promisedNumber int, acceptedProposal *Proposal) {
//	a.mutex.Lock()
//	defer a.mutex.Unlock()
//	a.promisedNumber = promisedNumber
//	a.acceptedProposal = acceptedProposal
//	a.recovered = true
//}
//
//func NewAcceptor(id int) *Acceptor {
//	return &Acceptor{
//		mutex:            sync.Mutex{},
//		id:               id,
//		promisedNumber:   0,
//		recovered:        false,
//		acceptedProposal: nil,
//	}
//}
//
//func (a *Acceptor) PersistState() {
//	// Simulate persisting state to disk
//	fmt.Printf("Acceptor %d persists state: PromisedNumber=%d, AcceptedProposal=%v\n",
//		a.id, a.promisedNumber, a.acceptedProposal)
//}
//
//func (a *Acceptor) RecoverStateFromDisk() {
//	// Simulate recovering state
//	fmt.Printf("Acceptor %d recovers state from disk: PromisedNumber=%d, AcceptedProposal=%v\n",
//		a.id, a.promisedNumber, a.acceptedProposal)
//	a.recovered = true
//}
