# NFS Exporter README

## Introduction

The NFS Exporter is a Prometheus exporter designed to monitor the status of Network File Systems (NFS). It checks specified NFS servers for exported paths by invoking the `showmount` command and exposes this information to Prometheus for monitoring purposes.

## Installation

Ensure that your system has the following dependencies installed:
- Go language environment (version 1.16+)
- `showmount` command (usually included in the nfs-utils package)

To build and install the NFS Exporter from source:
```shell
go build -o nfs_exporter main.go
```

Alternatively, you can download pre-built binaries directly or use Docker for deployment.

## Usage

When starting the NFS Exporter, you can configure its behavior via command-line flags. Below is a list of commonly used parameters:

- `-web.telemetry-path`: Path under which to expose metrics. Default is `/metrics`.
- `-web.listen-address`: Address on which to expose metrics and web interface. Default is `:9689`.
- `-nfs.executable-path`: Path to the executable file for querying NFS. Default is `/usr/sbin/showmount`.
- `-nfs.uri`: A comma-separated list of NFS URIs, formatted as `<address>:<mount_path>`.

For example, to start the exporter with default settings, run the following command:
```shell
./nfs_exporter --nfs.uri "192.168.2.22:/mnt,192.168.2.22:/opt"
```

### Using Docker

You can deploy the NFS Exporter using Docker. First, pull the image from Docker Hub:

```shell
docker pull guobaocai/nfs_exporter:v0.0.1
```

Then, run the container with the desired configuration. For instance:
```shell
docker run -d \
  --name nfs_exporter \
  -p 9689:9689 \
  guobaocai/nfs_exporter:v0.0.1 \
  --nfs.uri "192.168.2.22:/mnt,192.168.2.22:/opt"
```

Make sure to replace `192.168.2.22:/mnt,192.168.2.22:/opt` with your actual NFS server addresses and mount paths.

## Configuring Prometheus

To enable Prometheus to scrape data from the NFS Exporter, add a corresponding job configuration in your Prometheus configuration file, such as:

```yaml
scrape_configs:
  - job_name: 'nfs_exporter'
    static_configs:
      - targets: ['localhost:9689']
```

Make sure to replace `localhost:9689` with the actual address where your NFS Exporter is running.

## Metrics

The NFS Exporter exposes a single metric named `nfs_up`, which indicates whether the last query of an NFS server was successful. The metric labels include `mount_path` and `nfs_address`.

### Example Query

To check if all monitored NFS mount points are up, you could use a PromQL query like:
```promql
nfs_up == 1
```

This will return a result set with all instances where the NFS mount point is accessible.

## License

[Specify the license under which the NFS Exporter is distributed.]

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md] for details on our code of conduct, and the process for submitting pull requests to us.

## Authors

- Your Name - Initial work - [Your GitHub Profile]

See also the list of [contributors] who participated in this project.

## Acknowledgments

- Mention any frameworks, libraries or other tools that were used in creating this exporter.

Please adjust the sections `[Specify the license under which the NFS Exporter is distributed.]`, `[CONTRIBUTING.md]`, `[contributors]`, and `[Your GitHub Profile]` with the appropriate information relevant to your project. This README now includes instructions for deploying the NFS Exporter using Docker.