package gateway

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/awcullen/machinist/cloud"
	"github.com/awcullen/machinist/opcua"
	"github.com/dgraph-io/badger/v2"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
)

const (
	compactionInterval = 5 * time.Minute
)

// Start initializes the Gateway and connects to the local broker for configuration details.
func Start(config cloud.Config) (*Gateway, error) {
	log.Printf("[gateway] Gateway '%s' starting.\n", config.DeviceID)
	g := &Gateway{
		workers:    &sync.WaitGroup{},
		projectID:  config.ProjectID,
		region:     config.Region,
		registryID: config.RegistryID,
		deviceID:   config.DeviceID,
		privateKey: config.PrivateKey,
		algorithm:  config.Algorithm,
		storePath:  config.StorePath,
		deviceMap:  make(map[string]cloud.Device, 16),
		configCh:   make(chan *Config, 4),
		startTime:  time.Now().UTC(),
	}
	g.ctx, g.cancel = context.WithCancel(context.Background())

	// ensure storePath directory exists
	dbPath := filepath.Join(g.storePath, g.deviceID)
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
	g.db = db

	// read gateway config from store, if available
	config2 := new(Config)
	err = g.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("config"))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return unmarshalConfig(val, config2)
		})
	})
	if err == nil {
		g.name = config2.Name
		g.location = config2.Location
		g.configCh <- config2
	}

	g.configTopic = fmt.Sprintf("/devices/%s/config", config.DeviceID)
	g.commandsTopic = fmt.Sprintf("/devices/%s/commands/#", config.DeviceID)
	g.eventsTopic = fmt.Sprintf("/devices/%s/events", config.DeviceID)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(cloud.MqttBrokerURL)
	opts.SetProtocolVersion(cloud.ProtocolVersion)
	opts.SetClientID(fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", config.ProjectID, config.Region, config.RegistryID, config.DeviceID))
	opts.SetTLSConfig(&tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert})
	opts.SetCredentialsProvider(g.onProvideCredentials)
	opts.SetConnectRetry(true)
	opts.SetMaxReconnectInterval(30 * time.Second)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		client.Subscribe(g.configTopic, 1, g.onConfig)
		client.Subscribe(g.commandsTopic, 1, g.onCommand)
	})
	g.broker = mqtt.NewClient(opts)
	g.broker.Connect()

	// start all service workers
	g.startWorkers(g.ctx, g.cancel, g.workers)
	return g, nil
}

// Gateway stores the connected devices.
type Gateway struct {
	sync.Mutex
	ctx           context.Context
	cancel        context.CancelFunc
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
	deviceMap     map[string]cloud.Device
	startTime     time.Time
	configTime    time.Time
	configCh      chan *Config
}

// Name ...
func (g *Gateway) Name() string {
	return g.name
}

// Location ...
func (g *Gateway) Location() string {
	return g.location
}

// DeviceID ...
func (g *Gateway) DeviceID() string {
	return g.deviceID
}

// ProjectID ...
func (g *Gateway) ProjectID() string {
	return g.projectID
}

// Region ...
func (g *Gateway) Region() string {
	return g.region
}

// RegistryID ...
func (g *Gateway) RegistryID() string {
	return g.registryID
}

// PrivateKey ...
func (g *Gateway) PrivateKey() string {
	return g.privateKey
}

// Algorithm ...
func (g *Gateway) Algorithm() string {
	return g.algorithm
}

// StorePath ...
func (g *Gateway) StorePath() string {
	return g.storePath
}

// Devices ...
func (g *Gateway) Devices() map[string]cloud.Device {
	return g.deviceMap
}

// Stop disconnects from the local broker and stops each device of the gateway.
func (g *Gateway) Stop() error {
	log.Printf("[gateway] Gateway '%s' stopping.\n", g.deviceID)
	g.cancel()
	g.broker.Disconnect(3000)
	g.workers.Wait()
	g.db.Close()
	return nil
}

// startServices starts all services.
func (g *Gateway) startWorkers(ctx context.Context, cancel context.CancelFunc, workers *sync.WaitGroup) {
	// start northbound service
	g.startNorthboundWorker(ctx, cancel, workers)
	// start southbound service
	g.startSouthboundWorker(ctx, cancel, workers)
	// start db compaction task
	g.startGCWorker(ctx, cancel, workers)
}

// startNorthboundWorker starts a task to handle communication to Google Clout IoT Core.
func (g *Gateway) startNorthboundWorker(ctx context.Context, cancel context.CancelFunc, workers *sync.WaitGroup) {
	workers.Add(1)
	go func() {
		defer cancel()
		defer workers.Done()
		// loop till gateway cancelled
		for {
			select {
			case <-ctx.Done():
				return
			}
		}
	}()
}

// startSouthboundWorker starts a task to handle starting and stopping devices.
func (g *Gateway) startSouthboundWorker(ctx context.Context, cancel context.CancelFunc, workers *sync.WaitGroup) {
	workers.Add(1)
	go func() {
		defer cancel()
		defer workers.Done()
		// loop till gateway cancelled
		for {
			select {
			case <-ctx.Done():
				g.stopDevices()
				return
			case config := <-g.configCh:
				g.stopDevices()
				g.startDevices(config)
			}
		}
	}()
}

// stopDevices stops each device of the gateway.
func (g *Gateway) stopDevices() {
	g.Lock()
	defer g.Unlock()
	for k, v := range g.deviceMap {
		if err := v.Stop(); err != nil {
			log.Println(errors.Wrapf(err, "[gateway] Error stopping device '%s'", k))
		}
	}
	g.deviceMap = make(map[string]cloud.Device, 16)
	return
}

// startDevices starts each device of the gateway.
func (g *Gateway) startDevices(config *Config) {
	g.Lock()
	defer g.Unlock()
	// map the config into devices
	for _, conf := range config.Devices {
		n := cloud.Config{
			ProjectID:  g.projectID,
			Region:     g.region,
			RegistryID: g.registryID,
			DeviceID:   conf.DeviceID,
			PrivateKey: conf.PrivateKey,
			Algorithm:  conf.Algorithm,
			StorePath:  g.storePath,
		}
		switch conf.Kind {
		case "opcua":
			// start opcua device
			d, err := opcua.Start(n)
			if err != nil {
				log.Println(errors.Wrapf(err, "[gateway] Error starting opcua device '%s'", n.DeviceID))
				continue
			}
			g.deviceMap[n.DeviceID] = d
		default:
			log.Println(errors.Errorf("[gateway] Error starting device '%s': kind unknown", n.DeviceID))
		}
	}
	g.configTime = time.Now().UTC()
	return
}

// onProvideCredentials returns a username and password that is valid for this device.
func (g *Gateway) onProvideCredentials() (string, string) {
	// With Google Cloud IoT Core, the username field is ignored (but must not be empty), and the
	// password field is used to transmit a JWT to authorize the device
	now := time.Now()
	exp := now.Add(cloud.TokenDuration)
	jwt, err := cloud.CreateJWT(g.projectID, g.privateKey, g.algorithm, now, exp)
	if err != nil {
		log.Println(errors.Wrap(err, "[gateway] Error creating jwt token"))
		return "ignored", ""
	}
	return "ignored", jwt
}

// onConfig handles messages to the 'config' topic of the MQTT broker.
func (g *Gateway) onConfig(client mqtt.Client, msg mqtt.Message) {
	g.Lock()
	defer g.Unlock()
	var isDuplicate bool
	err := g.db.View(func(txn *badger.Txn) error {
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
		log.Println(errors.Wrap(err, "[gateway] Error unmarshalling config"))
		return
	}
	// store for next start
	err = g.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("config"), msg.Payload())
	})
	if err != nil {
		log.Println(errors.Wrap(err, "[gateway] Error storing config"))
	}
	g.name = config2.Name
	g.location = config2.Location
	// restart devices
	g.configCh <- config2
}

func (g *Gateway) onCommand(client mqtt.Client, msg mqtt.Message) {
	switch {
	default:
		log.Printf("[gateway] Gateway '%s' received unknown command '%s':\n", g.deviceID, msg.Topic())
	}
}

// startGCWorker starts a task to periodically compact the database, removing deleted values.
func (g *Gateway) startGCWorker(ctx context.Context, cancel context.CancelFunc, workers *sync.WaitGroup) {
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
					err := g.db.RunValueLogGC(0.5)
					if err != nil {
						break
					}
				}
			}
		}
	}()
}
