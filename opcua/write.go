package opcua

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	ua "github.com/awcullen/opcua"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

// onWriteRequest handles protobuf messages published to the 'commands/write-requests' topic of the MQTT broker.
// WriteResponse is published to the 'events/write-requests' topic
func (d *device) onWriteRequest(client mqtt.Client, msg mqtt.Message) {
	req := new(WriteRequest)
	err := json.Unmarshal(msg.Payload(), req)
	if err != nil {
		log.Println(errors.Wrapf(err, "[opcua] Device '%s' error parsing WriteRequest", d.deviceID))
		return
	}

	res, err := d.Write(d.ctx, req)
	if err != nil {
		log.Println(errors.Wrapf(err, "[opcua] Device '%s' error browsing", d.deviceID))
		return
	}

	payload, err := proto.Marshal(res)
	if err != nil {
		log.Println(errors.Wrapf(err, "[opcua] Device '%s' error marshalling response", d.deviceID))
		return
	}
	if tok := client.Publish(fmt.Sprintf("/devices/%s/events/write-response", d.DeviceID()), 1, false, payload); tok.Wait() && tok.Error() != nil {
		log.Println(errors.Wrapf(tok.Error(), "[opcua] Device '%s' error publishing to broker", d.deviceID))
		return
	}
}

// Write service grpc implementation.
func (d *device) Write(ctx context.Context, req *WriteRequest) (*WriteResponse, error) {
	// prepare request
	nodes := make([]*ua.WriteValue, 0)
	for _, item := range req.NodesToWrite {
		nodes = append(nodes, &ua.WriteValue{
			NodeID:      ua.ParseNodeID(item.NodeId),
			AttributeID: item.AttributeId,
			IndexRange:  item.IndexRange,
			Value:       ua.NewDataValueVariant(toUaVariant(item.Value), 0, time.Time{}, 0, time.Time{}, 0),
		})
	}
	req2 := &ua.WriteRequest{
		NodesToWrite: nodes,
	}

	// do request
	res2, err := d.write(ctx, req2)
	if err != nil {
		log.Println(errors.Wrapf(err, "[opcua] Device '%s' Error writing", d.deviceID))
		res := &WriteResponse{
			RequestHandle: req.RequestHandle,
			StatusCode:    uint32(ua.BadConnectionClosed),
		}
		return res, nil
	}

	// prepare response
	results := make([]uint32, 0)
	for _, result := range res2.Results {
		results = append(results, uint32(result))
	}

	// return response
	res := &WriteResponse{
		RequestHandle: req.RequestHandle,
		StatusCode:    uint32(res2.ServiceResult),
		Results:       results,
	}
	return res, nil

	// return nil, status.Errorf(codes.Unimplemented, "method Write not implemented")
}

// Write sets the values of Attributes of one or more Nodes.
func (d *device) write(ctx context.Context, request *ua.WriteRequest) (*ua.WriteResponse, error) {
	response, err := d.request(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(*ua.WriteResponse), nil
}
