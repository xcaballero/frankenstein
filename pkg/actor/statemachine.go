package actor

import "context"

type stateMachine struct {
	state   string
	actionc chan func()
}

func NewStateMachine() StateMachine {
	return &stateMachine{
		state:   "initial",
		actionc: make(chan func()),
	}
}

type StateMachine interface {
	Run(context.Context) error
}

/* func (sm *stateMachine) Run(cancel <-chan struct{}) {
	for {
		select {
		case f := <-sm.actionc:
			f()
		case <-cancel:
			return
		}
	}
} */

func (sm *stateMachine) Run(ctx context.Context) error {
	for {
		select {
		case f := <-sm.actionc:
			f()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
