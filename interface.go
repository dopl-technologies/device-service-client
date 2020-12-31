package client

import (
	dtprotos "github.com/dopl-technologies/api-protos-go"
)

// Interface device service client interface
type Interface interface {
	Create(*dtprotos.DeviceInfo) (*dtprotos.Device, error)

	Get(uint64) (*dtprotos.Device, error)

	Update(id uint64, info *dtprotos.DeviceInfo) (*dtprotos.Device, error)

	List() ([]*dtprotos.Device, error)

	Delete(id uint64) error

	Close()
}
