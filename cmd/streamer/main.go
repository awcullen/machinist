package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/awcullen/machinist/opcua"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/proto"
)

var (
	command   = "streamer"
	projectID string
	subID     string
	deviceID  string
	metricID  string
	version   string
	commit    string
)

func main() {
	var cmd = &cobra.Command{
		Use:     command,
		Args:    cobra.NoArgs,
		Version: fmt.Sprintf("%s+%s", version, commit),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s %s\n", cmd.Use, cmd.Version)
			run()
		},
	}

	// Set the command line flags.
	setFlags(cmd)

	// Check for environment variables prefixed by the command
	// name and referenced by the flag name (eg: FOOCOMMAND_BAR)
	// and override flag values before running the main command.
	viper.SetEnvPrefix(command)
	viper.AutomaticEnv()
	viper.BindPFlags(cmd.Flags())

	cmd.Execute()
}

func run() {
	if projectID == "" {
		fmt.Println("Missing parameter 'project'.")
		return
	}
	if subID == "" {
		fmt.Println("Missing parameter 'sub'.")
		return
	}
	cctx, cancel := context.WithCancel(context.Background())
	pullMsgs(cctx, projectID, subID, deviceID, metricID)
	waitForSignal()
	cancel()
	time.Sleep(3 * time.Second)
}

func setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&projectID, "project", "p", "our-chassis-269114", "project id")
	cmd.Flags().StringVarP(&subID, "sub", "s", "test-sub12", "subscription id")
	cmd.Flags().StringVarP(&deviceID, "device", "d", "", "filter to this device id")
	cmd.Flags().StringVarP(&metricID, "metric", "m", "", "filter to this metric id")
}

func waitForSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	fmt.Println()
}

func pullMsgs(ctx context.Context, projectID, subID, deviceID, metricID string) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Println(errors.Wrap(err, "pubsub.NewClient"))
	}
	sub := client.Subscription(subID)
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()
		mDeviceID := msg.Attributes["deviceId"]
		if deviceID != "" && deviceID != mDeviceID {
			return
		}
		fmt.Println(msg.PublishTime.Format(time.StampMilli))
		notification := &opcua.PublishResponse{}
		proto.Unmarshal(msg.Data, notification)
		for _, m := range notification.Metrics {
			if metricID != "" && metricID != m.Name {
				continue
			}
			fmt.Printf("ts: %d, d: %s, m: %s, v: %+v, q: %#X\n", m.Timestamp, mDeviceID, m.Name, m.Value, m.StatusCode)
		}
	})
	if err != nil {
		log.Println(errors.Wrap(err, "Error receiving"))
	}
}
