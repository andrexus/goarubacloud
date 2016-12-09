# Goarubacloud

Goarubacloud is a Go client library for accessing the Arubacloud API.

Inspired by [godo](https://github.com/digitalocean/godo) library.

Current version implements basic services:
- DataCenters (basic info about services used in selected datacenter)
- CloudServers (manage cloud servers)
- CloudServerActions (cloud servers actions like PowerOn, PowerOn, PowerCycle, Reinitialize, etc.)
- PurchasedIPs (manage IP addresses)
- VLans (manage VLANs)

## Usage

```go
import "github.com/andrexus/goarubacloud"
```

Create a new Arubacloud client, then use the exposed services to
access different parts of the Arubacloud API.

```go
import "github.com/andrexus/goarubacloud"

username := "AWI-XXXXX"
password:="XXXXX"

client := goarubacloud.NewClient(goarubacloud.Germany, username, password)
```

## Examples


To create a new Cloud server **PRO**:

```go
createRequest := &goarubacloud.CloudServerCreateRequestPro{
		Name:         "yourServerName",
		CPUQuantity:  1,
		RAMQuantity:  4,
		OSTemplateId: 481,
		Note:         "created with goarubaclient",
		AdministratorPassword:        "XXXXXXX",
		NetworkAdaptersConfiguration: []goarubacloud.NetworkAdapter{},
	}
	createRequest.VirtualDisks = append(createRequest.VirtualDisks, goarubacloud.CloudServerCreateVirtualDisk{
		VirtualDiskType: 0,
		Size:            10,
	})

	cloudServer, resp, err := client.CloudServers.Create(createRequest)
	if err != nil {
		fmt.Printf("Something bad happened: %s\n\n", err)
	}
```

To create a new Cloud server **SMART**:

```go
createRequest := &goarubacloud.CloudServerCreateRequestSmart{
		Name:         "yourServerName",
		OSTemplateId: 482,
		Note:         "created with goarubaclient",
		AdministratorPassword: "XXXXXXX",
		CloudServerSmartType:  goarubacloud.MEDIUM,
	}

	cloudServer, resp, err := client.CloudServers.Create(createRequest)
	if err != nil {
		fmt.Printf("Something bad happened: %s\n\n", err)
	}
```

## Contributing

Pull requests are appreciated!
