package main

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"
	"sort"
	"time"

	"github.com/awcullen/machinist/gateway"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

// SysInfo is generic holder for passsing data back
type SysInfo struct {
	Hostname      string   `json:"hostname"`
	OS            string   `json:"os"`
	Arch          string   `json:"architecture"`
	CPUs          int      `json:"cpuCount"`
	GoVersion     string   `json:"goVersion"`
	NetRemoteAddr string   `json:"netRemoteAddress"`
	NetHost       string   `json:"netHost"`
	EnvVars       []string `json:"envVars"`
}

// Metrics are real time system counters
type Metrics struct {
	MemTotal     uint64  `json:"memTotal"`
	MemUsed      uint64  `json:"memUsed"`
	CPUPerc      float64 `json:"cpuPerc"`
	DiskTotal    uint64  `json:"diskTotal"`
	DiskFree     uint64  `json:"diskFree"`
	NetBytesSent uint64  `json:"netBytesSent"`
	NetBytesRecv uint64  `json:"netBytesRecv"`
}

// Routes is our exported class
type Routes struct {
	disableCORS bool
	gateway     *gateway.Gateway
}

//
// /api/info - Return system information and properties
//
func (r Routes) apiInfoRoute(resp http.ResponseWriter, req *http.Request) {
	// CORS is for wimps
	if r.disableCORS {
		resp.Header().Set("Access-Control-Allow-Origin", "*")
	}

	var info SysInfo

	// Grab various bits of infomation from where we can
	info.Hostname, _ = os.Hostname()
	info.GoVersion = runtime.Version()
	info.OS = runtime.GOOS
	info.Arch = runtime.GOARCH
	info.CPUs = runtime.NumCPU()
	info.NetRemoteAddr = req.RemoteAddr
	info.NetHost = req.Host
	info.EnvVars = os.Environ()

	// JSON-ify our info
	js, err := json.Marshal(info)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fire JSON result back down the internet tubes
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(js)
}

//
// /api/metrics - Return system metrics cpu, mem, etc
//
func (r Routes) apiMetricsRoute(resp http.ResponseWriter, req *http.Request) {
	// CORS is for wimps
	if r.disableCORS {
		resp.Header().Set("Access-Control-Allow-Origin", "*")
	}

	var metrics Metrics

	// Memory stuff
	memStats, err := mem.VirtualMemory()
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	metrics.MemTotal = memStats.Total
	metrics.MemUsed = memStats.Used

	// CPU / processor stuff
	cpuStats, err := cpu.Percent(time.Millisecond*1000, false)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	metrics.CPUPerc = cpuStats[0]

	// Disk and filesystem usage stuff
	diskStats, err := disk.Usage("/")
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	metrics.DiskTotal = diskStats.Total
	metrics.DiskFree = diskStats.Free

	// Network stuff
	netStats, err := net.IOCounters(false)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	metrics.NetBytesRecv = netStats[0].BytesRecv
	metrics.NetBytesSent = netStats[0].BytesSent

	// JSON-ify our metrics
	js, err := json.Marshal(metrics)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fire JSON result back down the internet tubes
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(js)
}

// GatewayConfig ...
type GatewayConfig struct {
	Name       string         `json:"name"`
	Location   string         `json:"location"`
	ProjectID  string         `json:"projectId"`
	Region     string         `json:"region"`
	RegistryID string         `json:"registryId"`
	DeviceID   string         `json:"deviceId"`
	StorePath  string         `json:"storePath"`
	Devices    []DeviceConfig `json:"devices"`
}

// DeviceConfig ...
type DeviceConfig struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	DeviceID string `json:"deviceId"`
	Kind     string `json:"kind"`
}

// GatewayMetrics are real time system counters
type GatewayMetrics struct {
	MemTotal     uint64  `json:"memTotal"`
	MemUsed      uint64  `json:"memUsed"`
	CPUPerc      float64 `json:"cpuPerc"`
	DiskTotal    uint64  `json:"diskTotal"`
	DiskFree     uint64  `json:"diskFree"`
	NetBytesSent uint64  `json:"netBytesSent"`
	NetBytesRecv uint64  `json:"netBytesRecv"`
}

//
// /api/gateway - Return gateway metrics
//
func (r Routes) apiHomeRoute(resp http.ResponseWriter, req *http.Request) {
	// CORS is for wimps
	if r.disableCORS {
		resp.Header().Set("Access-Control-Allow-Origin", "*")
	}

	g := r.gateway
	devices := []DeviceConfig{}
	for _, d := range g.Devices() {
		dc := DeviceConfig{
			Name:     d.Name(),
			Location: d.Location(),
			DeviceID: d.DeviceID(),
			Kind:     d.Kind(),
		}
		devices = append(devices, dc)
	}
	sort.Slice(devices, func(i, j int) bool {
		return devices[i].Name < devices[j].Name
	})
	config := &GatewayConfig{
		Name:       g.Name(),
		Location:   g.Location(),
		ProjectID:  g.ProjectID(),
		Region:     g.Region(),
		RegistryID: g.RegistryID(),
		DeviceID:   g.DeviceID(),
		StorePath:  g.StorePath(),
		Devices:    devices,
	}

	// JSON-ify our metrics
	js, err := json.Marshal(config)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	// Fire JSON result back out the internet tubes
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(js)
}
