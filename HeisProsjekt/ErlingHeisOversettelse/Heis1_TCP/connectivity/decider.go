package connectivity

import "sync"

type decisionType int

const (
	Master decisionType = 1
	Slave  decisionType = -1
	Alone  decisionType = 0
)

var (
	decisionState       [NR_OF_ELEVATORS]decisionType //arreay of length NR_OF_ELEVATORS, deafult values Alone
	decisionState_mutex sync.Mutex
)

func init() {
	// Initialize all elements to Alone (0) using a loop
	for i := range decisionState {
		decisionState[i] = Alone
	}
}

func Get_decision_type(id int) decisionType {
	decisionState_mutex.Lock()
	defer decisionState_mutex.Unlock()
	return decisionState[id]
}

func SetDecisionState() {
	//Check is this elevator is connected to anyuthing

}
