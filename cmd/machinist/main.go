//go:generate statik -src=../../spa/dist

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"

	_ "expvar"
	_ "net/http/pprof"

	_ "github.com/awcullen/machinist/cmd/machinist/statik"

	"github.com/awcullen/machinist/cloud"
	"github.com/awcullen/machinist/gateway"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	command    = "machinist"
	configFile string
	debug      bool
	serverPort string
	version    string
	commit     string
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
	// mqtt.CRITICAL = log.New(os.Stderr, "", log.LstdFlags)
	// mqtt.ERROR = log.New(os.Stderr, "", log.LstdFlags)
	// mqtt.WARN = log.New(os.Stderr, "", log.LstdFlags)
	// mqtt.DEBUG = log.New(os.Stderr, "", log.LstdFlags)
	config := new(cloud.Config)
	in, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "[main] Error reading config file"))
	}
	err = yaml.Unmarshal(in, config)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "[main] Error unmarshalling from config file"))
	}
	g, err := gateway.Start(*config)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "[main] Error starting gateway"))
	}

	// Routing
	muxrouter := mux.NewRouter()
	routes := Routes{
		disableCORS: true,
		gateway:     g,
	}

	// API routes
	muxrouter.HandleFunc("/api/home", routes.apiHomeRoute)
	muxrouter.HandleFunc("/api/info", routes.apiInfoRoute)
	muxrouter.HandleFunc("/api/metrics", routes.apiMetricsRoute)

	// Handle static content
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	muxrouter.Handle("/", http.FileServer(statikFS))

	// EVERYTHING redirect to index.html if not found.
	muxrouter.NotFoundHandler = http.FileServer(&fsWrapper{statikFS})

	// Start server
	http.ListenAndServe(":"+serverPort, muxrouter)

	// time.Sleep(60 * time.Second)
	waitForSignal()
	g.Stop()
}

func setFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&debug, "debug", "d", false, "enable debug output")
	cmd.Flags().StringVarP(&configFile, "configFile", "c", "./config.yaml", "config file path")
	cmd.Flags().StringVarP(&serverPort, "serverPort", "p", "4000", "http server port")
}

func waitForSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	fmt.Println()
}

type fsWrapper struct {
	root http.FileSystem
}

func (f *fsWrapper) Open(name string) (http.File, error) {
	ret, err := f.root.Open(name)
	if !os.IsNotExist(err) || path.Ext(name) != "" {
		return ret, err
	}
	return f.root.Open("/index.html")
}
