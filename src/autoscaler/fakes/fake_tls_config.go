// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"crypto/tls"
	"sync"

	"code.cloudfoundry.org/app-autoscaler/src/autoscaler/eventgenerator/client"
	"google.golang.org/grpc/credentials"
)

type FakeTLSConfig struct {
	NewTLSStub        func(*tls.Config) credentials.TransportCredentials
	newTLSMutex       sync.RWMutex
	newTLSArgsForCall []struct {
		arg1 *tls.Config
	}
	newTLSReturns struct {
		result1 credentials.TransportCredentials
	}
	newTLSReturnsOnCall map[int]struct {
		result1 credentials.TransportCredentials
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeTLSConfig) NewTLS(arg1 *tls.Config) credentials.TransportCredentials {
	fake.newTLSMutex.Lock()
	ret, specificReturn := fake.newTLSReturnsOnCall[len(fake.newTLSArgsForCall)]
	fake.newTLSArgsForCall = append(fake.newTLSArgsForCall, struct {
		arg1 *tls.Config
	}{arg1})
	stub := fake.NewTLSStub
	fakeReturns := fake.newTLSReturns
	fake.recordInvocation("NewTLS", []interface{}{arg1})
	fake.newTLSMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeTLSConfig) NewTLSCallCount() int {
	fake.newTLSMutex.RLock()
	defer fake.newTLSMutex.RUnlock()
	return len(fake.newTLSArgsForCall)
}

func (fake *FakeTLSConfig) NewTLSCalls(stub func(*tls.Config) credentials.TransportCredentials) {
	fake.newTLSMutex.Lock()
	defer fake.newTLSMutex.Unlock()
	fake.NewTLSStub = stub
}

func (fake *FakeTLSConfig) NewTLSArgsForCall(i int) *tls.Config {
	fake.newTLSMutex.RLock()
	defer fake.newTLSMutex.RUnlock()
	argsForCall := fake.newTLSArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeTLSConfig) NewTLSReturns(result1 credentials.TransportCredentials) {
	fake.newTLSMutex.Lock()
	defer fake.newTLSMutex.Unlock()
	fake.NewTLSStub = nil
	fake.newTLSReturns = struct {
		result1 credentials.TransportCredentials
	}{result1}
}

func (fake *FakeTLSConfig) NewTLSReturnsOnCall(i int, result1 credentials.TransportCredentials) {
	fake.newTLSMutex.Lock()
	defer fake.newTLSMutex.Unlock()
	fake.NewTLSStub = nil
	if fake.newTLSReturnsOnCall == nil {
		fake.newTLSReturnsOnCall = make(map[int]struct {
			result1 credentials.TransportCredentials
		})
	}
	fake.newTLSReturnsOnCall[i] = struct {
		result1 credentials.TransportCredentials
	}{result1}
}

func (fake *FakeTLSConfig) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.newTLSMutex.RLock()
	defer fake.newTLSMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeTLSConfig) recordInvocation(key string, args []interface{}) {
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

var _ client.TLSConfig = new(FakeTLSConfig)
