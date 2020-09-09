// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package emailtest

import (
	"github.com/ONSdigital/dp-frontend-feedback-controller/email"
	"sync"
)

var (
	lockSenderMockSend sync.RWMutex
)

// Ensure, that SenderMock does implement Sender.
// If this is not the case, regenerate this file with moq.
var _ email.Sender = &SenderMock{}

// SenderMock is a mock implementation of email.Sender.
//
//     func TestSomethingThatUsesSender(t *testing.T) {
//
//         // make and configure a mocked email.Sender
//         mockedSender := &SenderMock{
//             SendFunc: func(from string, to []string, msg []byte) error {
// 	               panic("mock out the Send method")
//             },
//         }
//
//         // use mockedSender in code that requires email.Sender
//         // and then make assertions.
//
//     }
type SenderMock struct {
	// SendFunc mocks the Send method.
	SendFunc func(from string, to []string, msg []byte) error

	// calls tracks calls to the methods.
	calls struct {
		// Send holds details about calls to the Send method.
		Send []struct {
			// From is the from argument value.
			From string
			// To is the to argument value.
			To []string
			// Msg is the msg argument value.
			Msg []byte
		}
	}
}

// Send calls SendFunc.
func (mock *SenderMock) Send(from string, to []string, msg []byte) error {
	if mock.SendFunc == nil {
		panic("SenderMock.SendFunc: method is nil but Sender.Send was just called")
	}
	callInfo := struct {
		From string
		To   []string
		Msg  []byte
	}{
		From: from,
		To:   to,
		Msg:  msg,
	}
	lockSenderMockSend.Lock()
	mock.calls.Send = append(mock.calls.Send, callInfo)
	lockSenderMockSend.Unlock()
	return mock.SendFunc(from, to, msg)
}

// SendCalls gets all the calls that were made to Send.
// Check the length with:
//     len(mockedSender.SendCalls())
func (mock *SenderMock) SendCalls() []struct {
	From string
	To   []string
	Msg  []byte
} {
	var calls []struct {
		From string
		To   []string
		Msg  []byte
	}
	lockSenderMockSend.RLock()
	calls = mock.calls.Send
	lockSenderMockSend.RUnlock()
	return calls
}