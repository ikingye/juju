// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package lxd_client

import (
	"fmt"
	"net"

	"github.com/juju/errors"
	"github.com/juju/utils/arch"
	"github.com/lxc/lxd/shared"

	"github.com/juju/juju/network"
)

// AliveStatuses are the LXD statuses that indicate a container is "alive".
var AliveStatuses = []string{
	//StatusOK,
	//StatusPending,
	StatusStarting,
	StatusStarted,
	StatusRunning,
	//StatusThawed,
}

// InstanceSpec holds all the information needed to create a new LXD
// container.
type InstanceSpec struct {
	// Name is the "name" of the instance.
	Name string

	// Profiles are the names of the container profiles to apply to the
	// new container, in order.
	Profiles []string

	// Ephemeral indicates whether or not the container should be
	// destroyed when the LXD host is restarted.
	Ephemeral bool

	// Metadata is the instance metadata.
	Metadata map[string]string

	// Disks
	// Networks
	// Metadata
	// Tags
}

func (spec InstanceSpec) info(namespace string) *shared.ContainerState {
	name := spec.Name
	if namespace != "" {
		name = namespace + "-" + name
	}

	return &shared.ContainerState{
		Architecture:    0,
		Config:          map[string]string{},
		Devices:         shared.Devices{},
		Ephemeral:       false,
		ExpandedConfig:  map[string]string{},
		ExpandedDevices: shared.Devices{},
		Name:            name,
		Profiles:        []string{},
		//Status:          ContainerStatus{},
		Userdata: []byte{}, // from spec.Metadata
	}
}

// Summary builds an InstanceSummary based on the spec and returns it.
func (spec InstanceSpec) Summary(namespace string) InstanceSummary {
	info := spec.info(namespace)
	return newInstanceSummary(info)
}

// InstanceHardware describes the hardware characteristics of a LXC container.
type InstanceHardware struct {
	// Architecture is the CPU architecture.
	Architecture string

	// NumCores is the number of CPU cores.
	NumCores uint

	// MemoryMB is the memory allocation for the container.
	MemoryMB uint

	//RootDisk uint64
}

// InstanceSummary captures all the data needed by Instance.
type InstanceSummary struct {
	// Name is the "name" of the instance.
	Name string

	// Status holds the status of the instance at a certain point in time.
	Status string

	// Hardware describes the instance's hardware characterstics.
	Hardware InstanceHardware

	// Metadata is the instance metadata.
	Metadata map[string]string

	// Addresses
	Addresses []network.Address
}

func newInstanceSummary(info *shared.ContainerState) InstanceSummary {
	archStr, _ := shared.ArchitectureName(info.Architecture)
	archStr = arch.NormaliseArch(archStr)

	var numCores uint = 0 // default to all
	if raw := info.Config["limits.cpus"]; raw != "" {
		fmt.Sscanf(raw, "%d", &numCores)
	}

	var mem uint = 0 // default to all
	if raw := info.Config["limits.memory"]; raw != "" {
		fmt.Sscanf(raw, "%d", &mem)
	}

	var addrs []network.Address
	for _, info := range info.Status.Ips {
		addr := network.NewAddress(info.Address)

		// Ignore loopback devices.
		// TODO(ericsnow) Move the loopback test to a network.Address method?
		ip := net.ParseIP(addr.Value)
		if ip != nil && ip.IsLoopback() {
			continue
		}

		addrs = append(addrs, addr)
	}

	// TODO(ericsnow) Factor this out into a function.
	statusStr := info.Status.Status
	for status, code := range allStatuses {
		if info.Status.StatusCode == code {
			statusStr = status
			break
		}
	}

	return InstanceSummary{
		Name:      info.Name,
		Status:    statusStr,
		Metadata:  unpackMetadata(info.Userdata),
		Addresses: addrs,
		Hardware: InstanceHardware{
			Architecture: archStr,
			NumCores:     numCores,
			MemoryMB:     mem,
		},
	}
}

// Instance represents a single realized LXD container.
type Instance struct {
	InstanceSummary

	// spec is the InstanceSpec used to create this instance.
	spec *InstanceSpec
}

func newInstance(info *shared.ContainerState, spec *InstanceSpec) *Instance {
	summary := newInstanceSummary(info)
	return NewInstance(summary, spec)
}

// NewInstance builds an instance from the provided summary and spec
// and returns it.
func NewInstance(summary InstanceSummary, spec *InstanceSpec) *Instance {
	if spec != nil {
		// Make a copy.
		val := *spec
		spec = &val
	}
	return &Instance{
		InstanceSummary: summary,
		spec:            spec,
	}
}

// Status returns a string identifying the status of the instance.
func (gi Instance) Status() string {
	return gi.InstanceSummary.Status
}

// CurrentStatus returns a string identifying the status of the instance.
func (gi Instance) CurrentStatus(client *Client) (string, error) {
	// TODO(ericsnow) Do this a better way?

	inst, err := client.Instance(gi.Name)
	if err != nil {
		return "", errors.Trace(err)
	}
	return inst.Status(), nil
}

// Metadata returns the user-specified metadata for the instance.
func (gi Instance) Metadata() map[string]string {
	// TODO*ericsnow) return a copy?
	return gi.InstanceSummary.Metadata
}

// packMetadata composes the provided data into the format required
// by the API.
func packMetadata(data map[string]string) []byte {
	// TODO(ericsnow) finish!
	return nil
}

// unpackMetadata decomposes the provided data from the format used
// in the API.
func unpackMetadata(data []byte) map[string]string {
	if data == nil {
		return nil
	}

	// TODO(ericsnow) finish!
	return nil
}
