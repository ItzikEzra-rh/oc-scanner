# OC Scanner

A Go-based tool for scanning OpenShift/Kubernetes resources using the `oc` CLI.

## Features

- **Concurrent Scanning**: Scans multiple resource types simultaneously using goroutines
- **Modular Architecture**: Easy to extend with new resource scanners
- **OpenShift Native**: Uses `oc` CLI for cluster interactions
- **Namespace Support**: Scan resources in specific namespaces

## Supported Resources

- **Pods**: Scan and list all pods in a namespace
- **Deployments**: Scan and list all deployments in a namespace

## Prerequisites

- [OpenShift CLI (`oc`)](https://docs.openshift.com/container-platform/latest/cli_reference/openshift_cli/getting-started-cli.html) installed and configured
- Access to an OpenShift cluster (logged in via `oc login`)
- Go 1.19+ (for building from source)

## Installation

### From Source

```bash
git clone https://github.com/ItzikEzra-rh/oc-scanner.git
cd oc-scanner
go build -o oc-scanner .
```

## Usage

### Basic Usage

```bash
# Scan pods in a specific namespace
./oc-scanner scan <namespace> pods

# Scan deployments in a specific namespace  
./oc-scanner scan <namespace> deployments

# Scan multiple resources simultaneously
./oc-scanner scan <namespace> pods deployments
```

### Examples

```bash
# Scan pods in the openshift-kmm namespace
./oc-scanner scan openshift-kmm pods

# Scan both pods and deployments in default namespace
./oc-scanner scan default pods deployments

# Development mode (using go run)
go run main.go scan openshift-kmm pods
```

## Architecture

The scanner follows a modular design:

```
oc-scanner/
├── main.go              # Main application logic
├── scanner/
│   ├── interface.go     # Scanner interface definition
│   ├── pods.go         # Pod scanner implementation
│   └── deployments.go  # Deployment scanner implementation
├── go.mod
└── go.sum
```

### Scanner Interface

```go
type Scanner interface {
    Scan() error
}
```

Each resource type implements this interface, making it easy to add new scanners.

## Adding New Scanners

1. Create a new file in the `scanner/` directory (e.g., `services.go`)
2. Implement the `Scanner` interface:

```go
package scanner

type ServiceScanner struct {
    Namespace string
}

func (s ServiceScanner) Scan() error {
    // Implementation here
    return nil
}
```

3. Add the scanner to the factory map in `main.go`:

```go
scannerMap := map[string]func(string) scanner.Scanner{
    "pods":        func(ns string) scanner.Scanner { return scanner.PodScanner{Namespace: ns} },
    "deployments": func(ns string) scanner.Scanner { return scanner.DeploymentScanner{Namespace: ns} },
    "services":    func(ns string) scanner.Scanner { return scanner.ServiceScanner{Namespace: ns} },
}
```

## Concurrency

The scanner uses Go's goroutines and `sync.WaitGroup` to scan multiple resource types concurrently:

- Each resource type runs in its own goroutine
- All scans complete in parallel, improving performance
- Proper synchronization ensures all scans complete before exit

## Error Handling

- Individual scanner errors don't stop other scanners
- Errors are logged with resource context
- Non-zero exit codes for critical failures

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o oc-scanner .
```

### Running Directly

```bash
go run main.go scan <namespace> <resource-types...>
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

[MIT License](LICENSE)

## Requirements

- Go 1.19+
- OpenShift CLI (`oc`)
- Active OpenShift cluster connection 