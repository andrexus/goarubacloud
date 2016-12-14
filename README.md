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
    server_name := "myServerNamePro"
	os_template_id := 481
	admin_password := "myPassword"

	createRequest := goarubacloud.NewCloudServerProCreateRequest(server_name, admin_password, os_template_id)
	err := createRequest.AddVirtualDisk(20)
	if err != nil{
		log.Println("[ERROR]", err)
	}
	
	createRequest.SetCPUQuantity(2)
	createRequest.SetRAMQuantity(4)
	createRequest.SetNote("created by goarubacloud")

	cloudServer, resp, err := client.CloudServers.Create(createRequest)
	if err != nil {
		fmt.Println("Something bad happened:", err)
	}
```

To create a new Cloud server **SMART**:

```go
    server_name := "myServerNameSmart"
	os_template_id := 482
	admin_password := "myPassword"
	createRequest := goarubacloud.NewCloudServerSmartCreateRequest(goarubacloud.MEDIUM, server_name, admin_password, os_template_id)
	
	cloudServer, resp, err := client.CloudServers.Create(createRequest)
	if err != nil {
		fmt.Println("Something bad happened:", err)
	}
```

## Contributing

Pull requests are appreciated!
