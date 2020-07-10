package opcua

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/awcullen/machinist/cloud"
	ua "github.com/awcullen/opcua"
	"github.com/dgraph-io/badger/v2"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
	"github.com/reactivex/rxgo/v2"
	"google.golang.org/protobuf/proto"
)

const (
	maxIterationsPerTransaction = 1
	compactionInterval          = 5 * time.Minute
)

// Start the device.
func Start(config cloud.Config) (cloud.Device, error) {
	log.Printf("[opcua] Device '%s' starting.\n", config.DeviceID)
	d := &device{
		workers:    &sync.WaitGroup{},
		projectID:  config.ProjectID,
		region:     config.Region,
		registryID: config.RegistryID,
		deviceID:   config.DeviceID,
		privateKey: config.PrivateKey,
		algorithm:  config.Algorithm,
		storePath:  config.StorePath,
		configCh:   make(chan *Config, 4),
		commandsCh: make(chan *serviceOperation, 64),
		publishCh:  make(chan rxgo.Item),
	}

	// ensure storePath directory exists
	dbPath := filepath.Join(d.storePath, d.deviceID)
	if err := os.MkdirAll(dbPath, 0777); err != nil {
		return nil, err
	}
	// open the database
	db, err := badger.Open(badger.DefaultOptions(dbPath).WithTruncate(true).WithLogger(nil))
	if err != nil {
		os.Remove(dbPath)
		if err := os.MkdirAll(dbPath, 0777); err != nil {
			return nil, err
		}
		db, err = badger.Open(badger.DefaultOptions(dbPath).WithTruncate(true).WithLogger(nil))
		if err != nil {
			return nil, err
		}
	}
	d.db = db

	// read device config from store, if available
	config2 := new(Config)
	err = d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("config"))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return unmarshalConfig(val, config2)
		})
	})
	if err == nil {
		d.name = config2.Name
		d.location = config2.Location
		d.configCh <- config2
	}

	d.configTopic = fmt.Sprintf("/devices/%s/config", d.deviceID)
	d.commandsTopic = fmt.Sprintf("/devices/%s/commands/#", d.deviceID)
	d.eventsTopic = fmt.Sprintf("/devices/%s/events", d.deviceID)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(cloud.MqttBrokerURL)
	opts.SetProtocolVersion(cloud.ProtocolVersion)
	opts.SetClientID(fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", d.projectID, d.region, d.registryID, d.deviceID))
	opts.SetTLSConfig(&tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert})
	opts.SetCredentialsProvider(d.onProvideCredentials)
	opts.SetConnectRetry(true)
	opts.SetMaxReconnectInterval(30 * time.Second)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		client.Subscribe(d.configTopic, 1, d.onConfig)
		client.Subscribe(d.commandsTopic, 1, d.onCommand)
	})
	d.broker = mqtt.NewClient(opts)
	d.broker.Connect()

	// TODO: start grpc server?
	// lis, err := net.Listen("tcp", ":50051")
	// if err != nil {
	// 	log.Printf("failed to listen: %v", err)
	// }
	// d.grpcServer = grpc.NewServer()
	// RegisterBrowseServiceServer(d.grpcServer, d)
	// go d.grpcServer.Serve(lis)

	d.ctx, d.cancel = context.WithCancel(context.Background())
	d.observable = rxgo.FromEventSource(d.publishCh, rxgo.WithContext(d.ctx), rxgo.WithBackPressureStrategy(rxgo.Drop))

	// start all service workers
	d.startWorkers(d.ctx, d.cancel, d.workers)
	return d, nil
}

type device struct {
	sync.Mutex
	workers       *sync.WaitGroup
	name          string
	location      string
	projectID     string
	region        string
	registryID    string
	deviceID      string
	privateKey    string
	algorithm     string
	storePath     string
	configTopic   string
	commandsTopic string
	eventsTopic   string
	db            *badger.DB
	broker        mqtt.Client
	configCh      chan *Config
	commandsCh    chan *serviceOperation
	publishCh     chan rxgo.Item
	observable    rxgo.Observable
	ctx           context.Context
	cancel        context.CancelFunc
	// grpcServer    *grpc.Server
}

// Name ...
func (d *device) Name() string {
	return d.name
}

// Location ...
func (d *device) Location() string {
	return d.location
}

// DeviceID ...
func (d *device) DeviceID() string {
	return d.deviceID
}

// Kind ...
func (d *device) Kind() string {
	return "opcua"
}

// ProjectID ...
func (d *device) ProjectID() string {
	return d.projectID
}

// Region ...
func (d *device) Region() string {
	return d.region
}

// RegistryID ...
func (d *device) RegistryID() string {
	return d.registryID
}

// PrivateKey ...
func (d *device) PrivateKey() string {
	return d.privateKey
}

// Algorithm ...
func (d *device) Algorithm() string {
	return d.algorithm
}

// Stop the device.
func (d *device) Stop() error {
	log.Printf("[opcua] Device '%s' stopping.\n", d.deviceID)

	// if d.grpcServer != nil {
	// 	d.grpcServer.GracefulStop()
	// }

	d.cancel()
	d.broker.Disconnect(3000)
	d.workers.Wait()
	d.db.Close()

	return nil
}

// Observe returns a new channel to receive PublishResponse messages
func (d *device) Observe(opts ...rxgo.Option) <-chan rxgo.Item {
	return d.observable.Observe(opts...)
}

// startWorkers starts all service workers.
func (d *device) startWorkers(ctx context.Context, cancel context.CancelFunc, workers *sync.WaitGroup) {
	// start northbound worker
	d.startNorthboundWorker(ctx, cancel, workers)
	// start southbound worker
	d.startSouthboundWorker(ctx, cancel, workers)
	// start db compaction worker
	d.startGCWorker(ctx, cancel, workers)
}

// startNorthboundWorker starts a task to periodically check the database for PublishResponses and publishes to Google Cloud.
// As each publish is acknowleged, the PublishResponses are deleted.
func (d *device) startNorthboundWorker(ctx context.Context, cancel context.CancelFunc, workers *sync.WaitGroup) {
	workers.Add(1)
	go func() {
		defer cancel()
		defer workers.Done()

		var lastTs uint64
		// loop till device cancelled
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if !d.broker.IsConnectionOpen() {
					time.Sleep(500 * time.Millisecond)
					continue
				}
				var more bool
				err := d.db.View(func(txn *badger.Txn) error {
					opts := badger.DefaultIteratorOptions
					it := txn.NewIterator(opts)
					defer it.Close()
					for it.Seek([]byte{0}); it.ValidForPrefix([]byte{0}); it.Next() {
						item := it.Item()
						ts := binary.BigEndian.Uint64(item.Key())
						if ts < lastTs {
							log.Printf("[opcua] Device '%s' error published out of order.\n", d.deviceID)
						}
						lastTs = ts
						payload, err := it.Item().ValueCopy(nil)
						if err != nil {
							return err
						}
						t0 := time.Now()
						tok2 := d.broker.Publish(d.eventsTopic, byte(1), false, payload)
						if tok2.Wait() && tok2.Error() != nil {
							more = false
							return tok2.Error()
						}
						// if debug
						log.Printf("[opcua] Device '%s' published %d bytes in %d ms.\n", d.deviceID, len(payload), time.Now().Sub(t0).Milliseconds())
						err = d.db.Update(func(txn2 *badger.Txn) error {
							return txn2.Delete(item.KeyCopy(nil))
						})
						if err != nil {
							return err
						}
					}
					return nil
				})
				if err != nil {
					log.Println(errors.Wrapf(err, "[opcua] Device '%s' error publishing to broker", d.deviceID))
				}
				if !more {
					time.Sleep(500 * time.Millisecond)
				}
			}
		}
	}()
}

// startGCWorker starts a task to periodically compact the database, removing deleted values.
func (d *device) startGCWorker(ctx context.Context, cancel context.CancelFunc, workers *sync.WaitGroup) {
	workers.Add(1)
	go func() {
		defer cancel()
		defer workers.Done()
		ticker := time.NewTicker(compactionInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for {
					err := d.db.RunValueLogGC(0.5)
					if err != nil {
						break
					}
				}
			}
		}
	}()
}

// startSouthboundWorker opens a opcua session with the server and creates a subscription and monitored items.
// The function then starts three polling workers, a command worker, and a config worker
func (d *device) startSouthboundWorker(ctx context.Context, cancel context.CancelFunc, workers *sync.WaitGroup) {
	workers.Add(1)
	go func() {
		defer cancel()
		defer workers.Done()

		// wait for initial config
		var config *Config
		select {
		case <-ctx.Done():
			return
		case config = <-d.configCh:
		}

		for {
			select {
			case <-ctx.Done():
				return
			default:
				log.Printf("[opcua] Device '%s' opening channel to '%s'.\n", d.deviceID, config.EndpointURL)
				ch, err := ua.NewClient(ctx, config.EndpointURL, ua.WithInsecureSkipVerify())
				if err != nil {
					log.Println(errors.Wrapf(err, "[opcua] Device '%s' error opening channel to '%s'", d.deviceID, config.EndpointURL))
					time.Sleep(1 * time.Second)
					continue
				}
				req := &ua.CreateSubscriptionRequest{
					RequestedPublishingInterval: config.PublishInterval,
					RequestedMaxKeepAliveCount:  30,
					RequestedLifetimeCount:      30 * 3,
					PublishingEnabled:           true,
				}
				res, err := ch.CreateSubscription(ctx, req)
				if err != nil {
					log.Println(errors.Wrapf(err, "[opcua] Device '%s' error creating subscription", d.deviceID))
					ch.Abort(context.Background())
					time.Sleep(1 * time.Second)
					continue
				}

				// create the monitored items for the new subscription
				itemsToCreate := make([]*ua.MonitoredItemCreateRequest, 0, len(config.Metrics))
				for i, m := range config.Metrics {
					itemsToCreate = append(itemsToCreate, &ua.MonitoredItemCreateRequest{
						ItemToMonitor: &ua.ReadValueID{
							NodeID:      ua.ParseNodeID(m.NodeID),
							AttributeID: ua.AttributeIDValue,
						},
						MonitoringMode: ua.MonitoringModeReporting,
						RequestedParameters: &ua.MonitoringParameters{
							ClientHandle:     uint32(i),
							QueueSize:        m.QueueSize,
							DiscardOldest:    true,
							SamplingInterval: m.SamplingInterval,
						},
					})
				}
				req2 := &ua.CreateMonitoredItemsRequest{
					SubscriptionID:     res.SubscriptionID,
					TimestampsToReturn: ua.TimestampsToReturnBoth,
					ItemsToCreate:      itemsToCreate,
				}
				res2, err := ch.CreateMonitoredItems(ctx, req2)
				if err != nil {
					log.Println(errors.Wrapf(err, "[opcua] Device '%s' error creating items", d.deviceID))
					ch.Abort(context.Background())
					time.Sleep(1 * time.Second)
					continue
				}

				// check each item was successfully created.
				for i, item := range req2.ItemsToCreate {
					if res2.Results[i].StatusCode.IsBad() {
						log.Println(errors.Wrapf(err, "[opcua] Device '%s' error creating item '%s'", d.deviceID, item.ItemToMonitor.NodeID))
					}
				}

				// prepare a list of queues to be shared by the polling workers
				items := make([]*MonitoredItem, 0, len(config.Metrics))
				for _, m := range config.Metrics {
					items = append(items, NewMonitoredItem(m.MetricID, ua.MonitoringModeReporting, m.SamplingInterval, m.QueueSize, res2.Timestamp))
				}

				ctx1, cancel1 := context.WithCancel(ctx)
				workers1 := &sync.WaitGroup{}

				// start three polling workers in parallel to ensure no delays
				for i := 0; i < 3; i++ {
					d.startPollingWorker(ctx1, cancel1, workers1, ch, items, int(config.QueueSize), config.PublishTTL)
				}

				// start processing requests from command topics
				d.startCommandsWorker(ctx1, cancel1, workers1, ch)

				// start monitoring config channel for changes
				d.startConfigWorker(ctx1, cancel1, workers1, &config)

				workers1.Wait()

				log.Printf("[opcua] Device '%s' closing channel to '%s'.\n", d.deviceID, config.EndpointURL)
				err = ch.Close(context.Background())
				if err != nil {
					log.Println(errors.Wrapf(err, "[opcua] Device '%s' error closing opcua channel to '%s'", d.deviceID, config.EndpointURL))
					ch.Abort(context.Background())
					time.Sleep(1 * time.Second)
					continue
				}
			}
		}
	}()
}

// startPollingWorker starts a worker to poll the opcua server for data changes of monitored items of a subscription.
// At the sampling interval, each metric value is resampled and stored in a queue.
// At the publishing interval, all metric values are batched in a PublishResponse structure,
// stored in the database, and published to observers
func (d *device) startPollingWorker(ctx context.Context, cancel context.CancelFunc, workers *sync.WaitGroup,
	ch *ua.Client, items []*MonitoredItem, maxNotificationsPerItem int, publishTTL time.Duration) {
	workers.Add(1)
	go func() {
		defer workers.Done()
		defer cancel()
		req := &ua.PublishRequest{
			RequestHeader:                ua.RequestHeader{TimeoutHint: 60000},
			SubscriptionAcknowledgements: []*ua.SubscriptionAcknowledgement{},
		}
		size := len(items) * maxNotificationsPerItem
		for {
			res, err := ch.Publish(ctx, req)
			if err != nil {
				// log.Println(errors.Wrapf(err, "[opcua] Device '%s' error publishing ua", d.deviceID))
				return
			}
			// loop thru all the notifications.
			for _, n := range res.NotificationMessage.NotificationData {
				switch o := n.(type) {
				case *ua.DataChangeNotification:
					for _, z := range o.MonitoredItems {
						items[z.ClientHandle].Enqueue(z.Value)
					}
				}
			}
			tn := res.NotificationMessage.PublishTime

			metrics := make([]*MetricValue, 0, size)
			for _, item := range items {
				dvs, _ := item.Notifications(tn, maxNotificationsPerItem)
				for _, dv := range dvs {
					metrics = append(metrics, &MetricValue{Name: item.metricID, Value: toVariant(dv.InnerVariant()), Timestamp: dv.ServerTimestamp().UnixNano() / 1000000, StatusCode: uint32(dv.StatusCode())})
				}
			}
			pubresponse := &PublishResponse{Timestamp: tn.UnixNano() / 1000000, Metrics: metrics}
			// send to observers
			d.publishCh <- rxgo.Item{V: pubresponse}
			// store in db
			key := make([]byte, 8)
			binary.BigEndian.PutUint64(key, uint64(pubresponse.Timestamp))
			payload, err := proto.Marshal(pubresponse)
			if err != nil {
				log.Println(errors.Wrapf(err, "[opcua] Device '%s' error marshalling publish response", d.deviceID))
			} else {
				// compress payload would save about 50%
				// var b bytes.Buffer
				// w := zlib.NewWriter(&b)
				// w.Write(payload)
				// w.Close()
				// log.Println("Compressed size:", b.Len())
				err = d.db.Update(func(txn *badger.Txn) error {
					e := badger.NewEntry(key, payload)
					e.WithTTL(publishTTL)
					return txn.SetEntry(e)
				})
				if err != nil {
					log.Println(errors.Wrapf(err, "[opcua] Device '%s' error storing publish response", d.deviceID))
				}
			}
			req = &ua.PublishRequest{
				RequestHeader: ua.RequestHeader{TimeoutHint: 60000},
				SubscriptionAcknowledgements: []*ua.SubscriptionAcknowledgement{
					{
						SubscriptionID: res.SubscriptionID,
						SequenceNumber: res.NotificationMessage.SequenceNumber,
					},
				},
			}
		}
	}()
}

// startCommandsWorker starts a worker that monitors commands channel for opcua requests, execs the request, and returns response to response channel.
func (d *device) startCommandsWorker(ctx context.Context, cancel context.CancelFunc, workers *sync.WaitGroup, ch *ua.Client) {
	workers.Add(1)
	go func() {
		defer workers.Done()
		defer cancel()
		for {
			select {
			// a service operation is an opcua request and private channel for response.
			case op := <-d.commandsCh:
				req := op.request
				// skip if opcua request already timed out
				header := req.Header()
				if time.Now().After(header.Timestamp.Add(time.Duration(header.TimeoutHint) * time.Millisecond)) {
					continue
				}
				res, err := ch.Request(ctx, req)
				if err != nil {
					continue
				}
				op.responseCh <- res
			case <-ctx.Done():
				return
			}
		}
	}()
}

// startConfigWorker monitors config channel for update, installs config, and causes caller to restart
func (d *device) startConfigWorker(ctx context.Context, cancel context.CancelFunc, workers *sync.WaitGroup, config **Config) {
	workers.Add(1)
	go func() {
		defer workers.Done()
		defer cancel()
		select {
		case *config = <-d.configCh:
			return
		case <-ctx.Done():
			return
		}
	}()
}

// onCommand handles messages to the 'commands/#' topic of the MQTT broker.
func (d *device) onCommand(client mqtt.Client, msg mqtt.Message) {
	switch {
	case strings.HasSuffix(msg.Topic(), "browse-request"):
		d.onBrowseRequest(client, msg)

	case strings.HasSuffix(msg.Topic(), "write-request"):
		d.onWriteRequest(client, msg)
	default:
		log.Printf("[opcua] Device '%s' received unknown command '%s':\n", d.deviceID, msg.Topic())
	}
}

// onConfig handles messages to the 'config' topic of the MQTT broker.
func (d *device) onConfig(client mqtt.Client, msg mqtt.Message) {
	d.Lock()
	defer d.Unlock()
	var isDuplicate bool
	err := d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("config"))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			isDuplicate = bytes.Equal(val, msg.Payload())
			return nil
		})
	})
	if isDuplicate {
		return
	}
	config2 := new(Config)
	err = unmarshalConfig(msg.Payload(), config2)
	if err != nil {
		log.Println(errors.Wrapf(err, "[opcua] Device '%s' error unmarshalling config", d.deviceID))
		return
	}
	// store for next start
	err = d.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("config"), msg.Payload())
	})
	if err != nil {
		log.Println(errors.Wrapf(err, "[opcua] Device '%s' error storing config", d.deviceID))
	}
	d.name = config2.Name
	d.location = config2.Location
	// restart device
	d.configCh <- config2
}

// request sends the opcua service requests to the current opcua session and waits for a response or timeout.
func (d *device) request(ctx context.Context, req ua.ServiceRequest) (ua.ServiceResponse, error) {
	header := req.Header()
	header.Timestamp = time.Now()
	// header.RequestHandle = ch.getNextRequestHandle()
	// header.AuthenticationToken = ch.authenticationToken
	if header.TimeoutHint == 0 {
		header.TimeoutHint = defaultTimeoutHint
	}
	var operation = &serviceOperation{request: req, responseCh: make(chan ua.ServiceResponse, 1)}
	d.commandsCh <- operation
	ctx, cancel := context.WithTimeout(ctx, time.Duration(req.Header().TimeoutHint)*time.Millisecond)
	select {
	case res := <-operation.responseCh:
		if sr := res.Header().ServiceResult; sr != ua.Good {
			cancel()
			return nil, sr
		}
		cancel()
		return res, nil
	case <-ctx.Done():
		cancel()
		return nil, ua.BadRequestTimeout
	}
}

// onProvideCredentials returns a username and password that is valid for this device.
func (d *device) onProvideCredentials() (string, string) {
	// With Google Cloud IoT Core, the username field is ignored (but must not be empty), and the
	// password field is used to transmit a JWT to authorize the device
	now := time.Now()
	exp := now.Add(cloud.TokenDuration)
	jwt, err := cloud.CreateJWT(d.projectID, d.privateKey, d.algorithm, now, exp)
	if err != nil {
		log.Println(errors.Wrapf(err, "[opcua] Device '%s' error creating jwt token", d.deviceID))
		return "ignored", ""
	}
	return "ignored", jwt
}

type serviceOperation struct {
	request    ua.ServiceRequest
	responseCh chan ua.ServiceResponse
}
