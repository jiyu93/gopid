package gopid

import (
	"sync"
)

// PID controller
type PID struct {
	Kp          float64 // Proportion
	Ki          float64 // Integral
	Kd          float64 // Derivative
	TargetValue float64
	prevError   float64
	lastError   float64
	sumError    float64
	locker      sync.Mutex
}

// NewPID create a new PID object
func NewPID(p, i, d, target float64) *PID {
	return &PID{
		Kp:          p,
		Ki:          i,
		Kd:          d,
		TargetValue: target,
	}
}

// CalcIncPID calculate increase PID output
func (p *PID) CalcIncPID(currentValue float64) float64 {
	p.locker.Lock()
	defer p.locker.Unlock()
	errorValue := p.TargetValue - currentValue
	output := p.Kp * (errorValue + p.Ki*p.lastError + p.Kd*p.prevError)
	p.prevError = p.lastError
	p.lastError = errorValue
	return output
}

// CalcLocPID calculate locative PID output
func (p *PID) CalcLocPID(currentValue float64) float64 {
	p.locker.Lock()
	defer p.locker.Unlock()
	errorValue := p.TargetValue - currentValue
	p.sumError += errorValue
	dError := (errorValue - p.prevError)
	p.prevError = errorValue
	return (p.Kp*errorValue + p.Ki*p.sumError + p.Kd*dError)
}
