package opcua

import (
	"context"
	"fmt"
	"log"

	ua "github.com/awcullen/opcua"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

const (
	// defaultTimeoutHint is the default number of milliseconds before a request is cancelled. (15 sec)
	defaultTimeoutHint uint32 = 15000
	// defaultDiagnosticsHint is the default diagnostic hint that is sent in a request. (None)
	defaultDiagnosticsHint uint32 = 0x00000000
)

// onBrowseRequest handles messages to the 'commands/browse-requests' topic of the MQTT broker.
func (d *device) onBrowseRequest(client mqtt.Client, msg mqtt.Message) {
	req := new(BrowseRequest)
	err := proto.Unmarshal(msg.Payload(), req)
	if err != nil {
		log.Println(errors.Wrapf(err, "[opcua] Device '%s' error parsing browse-request", d.deviceID))
		return
	}

	res, err := d.Browse(d.ctx, req)
	if err != nil {
		log.Println(errors.Wrapf(err, "[opcua] Device '%s' error browsing", d.deviceID))
		return
	}

	payload, err := proto.Marshal(res)
	if err != nil {
		log.Println(errors.Wrapf(err, "[opcua] Device '%s' error marshalling response", d.deviceID))
		return
	}
	if tok := client.Publish(fmt.Sprintf("/devices/%s/events/browse-response", d.DeviceID()), 1, false, payload); tok.Wait() && tok.Error() != nil {
		log.Println(errors.Wrapf(tok.Error(), "[opcua] Device '%s' error publishing to broker", d.deviceID))
		return
	}
}

// Browse grpc service implementation.
func (d *device) Browse(ctx context.Context, req *BrowseRequest) (*BrowseResponse, error) {
	// prepare request
	nodes := make([]*ua.BrowseDescription, 0)
	for _, item := range req.NodesToBrowse {
		nodes = append(nodes, &ua.BrowseDescription{
			NodeID:          ua.ParseNodeID(item.NodeId),
			BrowseDirection: ua.BrowseDirection(item.BrowseDirection),
			ReferenceTypeID: ua.ParseNodeID(item.ReferenceTypeId),
			IncludeSubtypes: item.IncludeSubtypes,
			NodeClassMask:   item.NodeClassMask,
			ResultMask:      item.ResultMask,
		})
	}
	req2 := &ua.BrowseRequest{
		RequestedMaxReferencesPerNode: req.RequestedMaxReferencesPerNode,
		NodesToBrowse:                 nodes,
	}

	// do request
	res2, err := d.browse(ctx, req2)
	if err != nil {
		log.Println(errors.Wrapf(err, "[opcua] Device '%s' error browsing", d.deviceID))
		res := &BrowseResponse{
			RequestHandle: req.RequestHandle,
			StatusCode:    uint32(ua.BadConnectionClosed),
		}
		return res, nil
	}

	// prepare response
	results := make([]*BrowseResult, 0)
	for _, result := range res2.Results {
		refs := make([]*ReferenceDescription, 0)
		for _, ref := range result.References {
			refs = append(refs, &ReferenceDescription{
				ReferenceTypeId: ref.ReferenceTypeID.String(),
				IsForward:       ref.IsForward,
				NodeId:          ref.NodeID.String(),
				BrowseName:      ref.BrowseName.String(),
				DisplayName:     ref.DisplayName.Text,
				NodeClass:       uint32(ref.NodeClass),
				TypeDefinition:  ref.TypeDefinition.String(),
			})
		}
		results = append(results, &BrowseResult{
			StatusCode:        uint32(result.StatusCode),
			ContinuationPoint: []byte(result.ContinuationPoint),
			References:        refs,
		})
	}

	// return response
	res := &BrowseResponse{
		RequestHandle: req.RequestHandle,
		StatusCode:    uint32(res2.ServiceResult),
		Results:       results,
	}
	return res, nil
}

// Browse discovers the References of a specified Node.
func (d *device) browse(ctx context.Context, request *ua.BrowseRequest) (*ua.BrowseResponse, error) {
	response, err := d.request(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(*ua.BrowseResponse), nil
}

// onBrowseNextRequest handles messages to the 'commands/browse-next-requests' topic of the MQTT broker.
func (d *device) onBrowseNextRequest(client mqtt.Client, msg mqtt.Message) {
	req := new(BrowseNextRequest)
	err := proto.Unmarshal(msg.Payload(), req)
	if err != nil {
		log.Println(errors.Wrapf(err, "[opcua] Device '%s' error parsing browse-next-request", d.deviceID))
		return
	}

	res, err := d.BrowseNext(d.ctx, req)
	if err != nil {
		log.Println(errors.Wrapf(err, "[opcua] Device '%s' error browsing", d.deviceID))
		return
	}

	payload, err := proto.Marshal(res)
	if err != nil {
		log.Println(errors.Wrapf(err, "[opcua] Device '%s' error marshalling response", d.deviceID))
		return
	}
	if tok := client.Publish(fmt.Sprintf("/devices/%s/events/browse-next-response", d.DeviceID()), 1, false, payload); tok.Wait() && tok.Error() != nil {
		log.Println(errors.Wrapf(tok.Error(), "[opcua] Device '%s' error publishing to broker", d.deviceID))
		return
	}
}

// BrowseNext service implementation.
func (d *device) BrowseNext(ctx context.Context, req *BrowseNextRequest) (*BrowseNextResponse, error) {
	// prepare request
	cps := make([]ua.ByteString, 0)
	for _, item := range req.ContinuationPoints {
		cps = append(cps, ua.ByteString(item))
	}
	req2 := &ua.BrowseNextRequest{
		ReleaseContinuationPoints: req.ReleaseContinuationPoints,
		ContinuationPoints:        cps,
	}

	// do request
	res2, err := d.browseNext(ctx, req2)
	if err != nil {
		log.Println(errors.Wrapf(err, "[opcua] Device '%s' error browsing", d.deviceID))
		res := &BrowseNextResponse{
			RequestHandle: req.RequestHandle,
			StatusCode:    uint32(ua.BadConnectionClosed),
		}
		return res, nil
	}

	// prepare response
	results := make([]*BrowseResult, 0)
	for _, result := range res2.Results {
		refs := make([]*ReferenceDescription, 0)
		for _, ref := range result.References {
			refs = append(refs, &ReferenceDescription{
				ReferenceTypeId: ref.ReferenceTypeID.String(),
				IsForward:       ref.IsForward,
				NodeId:          ref.NodeID.String(),
				BrowseName:      ref.BrowseName.String(),
				DisplayName:     ref.DisplayName.Text,
				NodeClass:       uint32(ref.NodeClass),
				TypeDefinition:  ref.TypeDefinition.String(),
			})
		}
		results = append(results, &BrowseResult{
			StatusCode:        uint32(result.StatusCode),
			ContinuationPoint: []byte(result.ContinuationPoint),
			References:        refs,
		})
	}

	// return response
	res := &BrowseNextResponse{
		RequestHandle: req.RequestHandle,
		StatusCode:    uint32(res2.ServiceResult),
		Results:       results,
	}
	return res, nil
}

// BrowseNext requests the next set of Browse responses, when the information is too large to be sent in a single response.
func (d *device) browseNext(ctx context.Context, request *ua.BrowseNextRequest) (*ua.BrowseNextResponse, error) {
	response, err := d.request(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(*ua.BrowseNextResponse), nil
}
