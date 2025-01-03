package main

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
)

const (
	nfsCmd    = "/usr/sbin/showmount"
	namespace = "nfs"
)

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last query of NFS successful.",
		[]string{"mount_path", "nfs_address"}, nil,
	)
	metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	listenAddress = kingpin.Flag("web.listen-address", "Address on which to expose metrics and web interface.").Default(":9689").String()
	nfsExcPath    = kingpin.Flag("nfs.executable-path", "Path to nfs executable.").Default(nfsCmd).String()
	nfsUri        = kingpin.Flag("nfs.uri", "NFS URIs.").Default("192.168.2.22:/mnt,192.168.2.22:/opt").String()
)

type Exporter struct {
	hostname string
	execpath string
	servers  []struct{ address, mountPath string }
}

// Describe all the metrics exported by NFS exporter. It implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
}

// Collect collects all the metrics
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}

	for _, server := range e.servers {
		wg.Add(1)
		go func(server struct{ address, mountPath string }) {
			defer wg.Done()
			cmd := exec.Command(e.execpath, "-e", server.address)
			output, err := cmd.Output()
			found := false
			if err == nil {
				lines := strings.Split(strings.TrimSpace(string(output)), "\n")
				for _, line := range lines {
					fields := strings.Fields(line)
					if len(fields) > 0 && fields[0] == server.mountPath {
						log.Infoln("Mount Path is matching NFS server:", server.mountPath, server.address)
						found = true
						break
					}
				}
			} else {
				log.Errorf("Exec Command %s %s failed: %v", e.execpath, server.address, err)
			}

			value := 0.0
			if found {
				value = 1.0
			}

			mutex.Lock()
			ch <- prometheus.MustNewConstMetric(
				up, prometheus.GaugeValue, value,
				server.mountPath, server.address,
			)
			mutex.Unlock()
		}(server)
	}

	wg.Wait()
}

func NewExporter(hostname, nfsExcPath string, nfsUris []string) (*Exporter, error) {
	servers := make([]struct{ address, mountPath string }, 0)
	for _, uri := range nfsUris {
		parts := strings.SplitN(uri, ":", 2)
		if len(parts) != 2 {
			log.Warnf("Invalid NFS URI format: %s", uri)
			continue
		}
		address, mountPath := parts[0], parts[1]
		servers = append(servers, struct{ address, mountPath string }{address, mountPath})
	}
	return &Exporter{
		hostname: hostname,
		execpath: nfsExcPath,
		servers:  servers,
	}, nil
}

func main() {
	kingpin.Version(version.Print("nfs_exporter"))
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()
	log.Infoln("Starting nfs_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())
	log.Infof("NFS URI: %s", *nfsUri)

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("While trying to get Hostname error happened: ", err)
	}

	// 解析nfsUri参数为数组
	nfsUris := strings.Split(*nfsUri, ",")

	// 创建一个Exporter实例
	exporter, err := NewExporter(hostname, *nfsExcPath, nfsUris)
	if err != nil {
		log.Fatalf("Failed to create exporters: %v", err)
	}

	// 注册Exporter
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>NFS Exporter v` + version.Version + `</title></head>
			<body>
			<h1>NFS Exporter v` + version.Version + `</h1>
			<p><a href='` + *metricsPath + `'>Metrics</a></p>
			</body>
			</html>`))
	})

	log.Infoln("Listening on", *listenAddress)
	err = http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}
