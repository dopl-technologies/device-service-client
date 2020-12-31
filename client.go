package client

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"

	dtprotos "github.com/dopl-technologies/api-protos-go"
)

// Client device service client
type Client struct {
	client dtprotos.DeviceServiceClient
	conn   *grpc.ClientConn
}

// New creates a new client that connects to the
// given address
func New(address string) (Interface, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Client{
		client: dtprotos.NewDeviceServiceClient(conn),
		conn:   conn,
	}, nil
}

// Create creates a new device
func (c *Client) Create(info *dtprotos.DeviceInfo) (*dtprotos.Device, error) {
	req := &dtprotos.CreateDeviceRequest{
		Info: info,
	}
	res, err := c.client.Create(context.Background(), req)
	if err != nil {
		return nil, err
	}

	device := res.GetDevice()
	if device == nil {
		return nil, fmt.Errorf("unexpected device create response. Status is ok but device is nil")
	}

	return device, nil
}

// Get gets a device by its id
func (c *Client) Get(id uint64) (*dtprotos.Device, error) {
	req := &dtprotos.GetDeviceRequest{
		DeviceID: id,
	}
	res, err := c.client.Get(context.Background(), req)
	if err != nil {
		return nil, err
	}

	device := res.GetDevice()
	if device == nil {
		return nil, fmt.Errorf("unexpected get device response. Status is ok but device is nil")
	}

	return device, nil
}

// Update updates the device
func (c *Client) Update(id uint64, info *dtprotos.DeviceInfo) (*dtprotos.Device, error) {
	req := &dtprotos.UpdateDeviceRequest{
		DeviceID: id,
		Info:     info,
	}
	res, err := c.client.Update(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return res.GetDevice(), nil
}

// List lists devices
func (c *Client) List() ([]*dtprotos.Device, error) {
	req := &dtprotos.ListDevicesRequest{}
	stream, err := c.client.List(context.Background(), req)
	if err != nil {
		return nil, err
	}

	// Build the list
	var devices []*dtprotos.Device
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		devices = append(devices, res.GetDevice())
	}

	return devices, nil
}

// Delete deletes a device by its id
func (c *Client) Delete(id uint64) error {
	req := &dtprotos.DeleteDeviceRequest{
		DeviceID: id,
	}
	_, err := c.client.Delete(context.Background(), req)
	return err
}

// Close closes the connection
func (c *Client) Close() {
	c.conn.Close()
}
