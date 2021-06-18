// This file was generated by counterfeiter
package fakes

import (
	"autoscaler/metricsgateway"
	"sync"

	"code.cloudfoundry.org/go-loggregator/v8/rpc/loggregator_v2"
)

type FakeEmitter struct {
	AcceptStub        func(envelope *loggregator_v2.Envelope)
	acceptMutex       sync.RWMutex
	acceptArgsForCall []struct {
		envelope *loggregator_v2.Envelope
	}
	EmitStub        func(envelope *loggregator_v2.Envelope) error
	emitMutex       sync.RWMutex
	emitArgsForCall []struct {
		envelope *loggregator_v2.Envelope
	}
	emitReturns struct {
		result1 error
	}
	StartStub        func() error
	startMutex       sync.RWMutex
	startArgsForCall []struct{}
	startReturns     struct {
		result1 error
	}
	StopStub         func()
	stopMutex        sync.RWMutex
	stopArgsForCall  []struct{}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeEmitter) Accept(envelope *loggregator_v2.Envelope) {
	fake.acceptMutex.Lock()
	fake.acceptArgsForCall = append(fake.acceptArgsForCall, struct {
		envelope *loggregator_v2.Envelope
	}{envelope})
	fake.recordInvocation("Accept", []interface{}{envelope})
	fake.acceptMutex.Unlock()
	if fake.AcceptStub != nil {
		fake.AcceptStub(envelope)
	}
}

func (fake *FakeEmitter) AcceptCallCount() int {
	fake.acceptMutex.RLock()
	defer fake.acceptMutex.RUnlock()
	return len(fake.acceptArgsForCall)
}

func (fake *FakeEmitter) AcceptArgsForCall(i int) *loggregator_v2.Envelope {
	fake.acceptMutex.RLock()
	defer fake.acceptMutex.RUnlock()
	return fake.acceptArgsForCall[i].envelope
}

func (fake *FakeEmitter) Emit(envelope *loggregator_v2.Envelope) error {
	fake.emitMutex.Lock()
	fake.emitArgsForCall = append(fake.emitArgsForCall, struct {
		envelope *loggregator_v2.Envelope
	}{envelope})
	fake.recordInvocation("Emit", []interface{}{envelope})
	fake.emitMutex.Unlock()
	if fake.EmitStub != nil {
		return fake.EmitStub(envelope)
	}
	return fake.emitReturns.result1
}

func (fake *FakeEmitter) EmitCallCount() int {
	fake.emitMutex.RLock()
	defer fake.emitMutex.RUnlock()
	return len(fake.emitArgsForCall)
}

func (fake *FakeEmitter) EmitArgsForCall(i int) *loggregator_v2.Envelope {
	fake.emitMutex.RLock()
	defer fake.emitMutex.RUnlock()
	return fake.emitArgsForCall[i].envelope
}

func (fake *FakeEmitter) EmitReturns(result1 error) {
	fake.EmitStub = nil
	fake.emitReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeEmitter) Start() error {
	fake.startMutex.Lock()
	fake.startArgsForCall = append(fake.startArgsForCall, struct{}{})
	fake.recordInvocation("Start", []interface{}{})
	fake.startMutex.Unlock()
	if fake.StartStub != nil {
		return fake.StartStub()
	}
	return fake.startReturns.result1
}

func (fake *FakeEmitter) StartCallCount() int {
	fake.startMutex.RLock()
	defer fake.startMutex.RUnlock()
	return len(fake.startArgsForCall)
}

func (fake *FakeEmitter) StartReturns(result1 error) {
	fake.StartStub = nil
	fake.startReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeEmitter) Stop() {
	fake.stopMutex.Lock()
	fake.stopArgsForCall = append(fake.stopArgsForCall, struct{}{})
	fake.recordInvocation("Stop", []interface{}{})
	fake.stopMutex.Unlock()
	if fake.StopStub != nil {
		fake.StopStub()
	}
}

func (fake *FakeEmitter) StopCallCount() int {
	fake.stopMutex.RLock()
	defer fake.stopMutex.RUnlock()
	return len(fake.stopArgsForCall)
}

func (fake *FakeEmitter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.acceptMutex.RLock()
	defer fake.acceptMutex.RUnlock()
	fake.emitMutex.RLock()
	defer fake.emitMutex.RUnlock()
	fake.startMutex.RLock()
	defer fake.startMutex.RUnlock()
	fake.stopMutex.RLock()
	defer fake.stopMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeEmitter) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ metricsgateway.Emitter = new(FakeEmitter)
