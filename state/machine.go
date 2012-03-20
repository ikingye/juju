// launchpad.net/juju/state
//
// Copyright (c) 2011-2012 Canonical Ltd.

package state

import (
	"fmt"
	"launchpad.net/juju/go/state/presence"
	"strconv"
	"strings"
	"time"
)

// Machine represents the state of a machine.
type Machine struct {
	st  *State
	key string
}

// Id returns the machine id.
func (m *Machine) Id() int {
	return machineId(m.key)
}

// AgentAlive returns whether the respective remote agent is alive.
func (m *Machine) AgentAlive() (bool, error) {
	return presence.Alive(m.st.zk, m.zkAgentPath())
}

// WaitAgentAlive blocks until the respective agent is alive.
func (m *Machine) WaitAgentAlive(timeout time.Duration) error {
	err := presence.WaitAlive(m.st.zk, m.zkAgentPath(), timeout)
	if err != nil {
		return fmt.Errorf("wait for machine agent %d failed: %v", m.Id(), err)
	}
	return nil
}

// SetAgentAlive signals that the respective agent is alive. By stopping
// the returned pinger the agent can notify others about the end of its
// work.
func (m *Machine) SetAgentAlive() (*presence.Pinger, error) {
	return presence.StartPinger(m.st.zk, m.zkAgentPath(), agentPingerPeriod)
}

// zkKey returns the ZooKeeper key of the machine.
func (m *Machine) zkKey() string {
	return m.key
}

// zkPath returns the ZooKeeper base path for the machine.
func (m *Machine) zkPath() string {
	return fmt.Sprintf("/machines/%s", m.key)
}

// zkAgentPath returns the ZooKeeper path for the machine agent.
func (m *Machine) zkAgentPath() string {
	return fmt.Sprintf("/machines/%s/agent", m.key)
}

// machineId returns the machine id corresponding to machineKey.
func machineId(machineKey string) (id int) {
	if machineKey == "" {
		panic("machineId: empty machine key")
	}
	i := strings.Index(machineKey, "-")
	var id64 int64
	var err error
	if i >= 0 {
		id64, err = strconv.ParseInt(machineKey[i+1:], 10, 32)
	}
	if i < 0 || err != nil {
		panic("machineId: invalid machine key: " + machineKey)
	}
	return int(id64)
}

// machineKey returns the machine key corresponding to machineId.
func machineKey(machineId int) string {
	return fmt.Sprintf("machine-%010d", machineId)
}
