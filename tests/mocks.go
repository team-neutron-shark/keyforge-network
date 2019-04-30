package tests

import (
	"bytes"
	"net"
	"time"
)

type MockAddress struct {
}

func (m MockAddress) Network() string {
	return "tcp"
}

func (m MockAddress) String() string {
	return "0.0.0.0:1234"
}

type MockNetworkConnection struct {
	Buffer *bytes.Buffer
}

func NewMockNetworkConnection() *MockNetworkConnection {
	mockConnection := new(MockNetworkConnection)
	temp := make([]byte, 0)
	mockConnection.Buffer = bytes.NewBuffer(temp)
	return mockConnection
}

func (m MockNetworkConnection) Read(b []byte) (int, error) {
	return m.Buffer.Read(b)
}

func (m MockNetworkConnection) Write(b []byte) (int, error) {
	m.Buffer.Grow(len(b) + m.Buffer.Len())
	return m.Buffer.Write(b)
}

func (m MockNetworkConnection) Close() error {
	return nil
}

func (m MockNetworkConnection) LocalAddr() net.Addr {
	return MockAddress{}
}

func (m MockNetworkConnection) RemoteAddr() net.Addr {
	return MockAddress{}
}

func (m MockNetworkConnection) SetDeadline(t time.Time) error {
	return nil
}

func (m MockNetworkConnection) SetReadDeadline(t time.Time) error {
	return nil
}

func (m MockNetworkConnection) SetWriteDeadline(t time.Time) error {
	return nil
}
