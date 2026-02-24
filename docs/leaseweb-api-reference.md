# Leaseweb API Reference

OpenAPI 3.0.3 — Developer API

## Introduction

### Error


### Authentication


## Cloud

### Public Clouds

- `POST /publicCloud/v1/autoScalingGroups` — Create Auto Scaling Group
- `GET /publicCloud/v1/autoScalingGroups` — Get Auto Scaling Group list
- `DELETE /publicCloud/v1/autoScalingGroups/{autoScalingGroupId}` — Delete Auto Scaling Group
- `PUT /publicCloud/v1/autoScalingGroups/{autoScalingGroupId}` — Update Auto Scaling Group
- `GET /publicCloud/v1/autoScalingGroups/{autoScalingGroupId}` — Get Auto Scaling Group details
- `POST /publicCloud/v1/autoScalingGroups/{autoScalingGroupId}/deregisterTargetGroup` — Deregister Target Group
- `GET /publicCloud/v1/autoScalingGroups/{autoScalingGroupId}/instances` — Get list of instances belonging to an Auto Scaling Group
- `POST /publicCloud/v1/autoScalingGroups/{autoScalingGroupId}/refreshImage` — Refresh auto scaling group image
- `POST /publicCloud/v1/autoScalingGroups/{autoScalingGroupId}/registerTargetGroup` — Register Target Group
- `GET /publicCloud/v1/equipments/{equipmentId}/expenses` — Get costs for a given month.
- `GET /publicCloud/v1/images` — List all available Images
- `POST /publicCloud/v1/images` — Create Custom Image
- `PUT /publicCloud/v1/images/{imageId}` — Update Custom Image
- `POST /publicCloud/v1/instance/{instanceId}/monitoring/enable` — Enable monitoring
- `GET /publicCloud/v1/instance/{instanceId}/monitoring/status` — Get service monitoring status
- `GET /publicCloud/v1/instanceTypes` — List instance types
- `POST /publicCloud/v1/instances` — Launch instance
- `GET /publicCloud/v1/instances` — Get instance list
- `GET /publicCloud/v1/instances/{instanceId}` — Get instance details
- `DELETE /publicCloud/v1/instances/{instanceId}` — Terminate instance
- `PUT /publicCloud/v1/instances/{instanceId}` — Update instance
- `PUT /publicCloud/v1/instances/{instanceId}/addToPrivateNetwork` — Add instance to Private Network
- `POST /publicCloud/v1/instances/{instanceId}/attachIso` — Attach ISO to a specific Instance
- `POST /publicCloud/v1/instances/{instanceId}/attachSecurityGroups` — Attach security groups to instance
- `POST /publicCloud/v1/instances/{instanceId}/cancelTermination` — Cancel instance termination
- `GET /publicCloud/v1/instances/{instanceId}/console` — Get console access
- `GET /publicCloud/v1/instances/{instanceId}/credentials` — List credentials stored for a specific Instance
- `POST /publicCloud/v1/instances/{instanceId}/credentials` — Store credentials for a specific Instance
- `DELETE /publicCloud/v1/instances/{instanceId}/credentials` — Delete all credentials associated with a specific Instance
- `GET /publicCloud/v1/instances/{instanceId}/credentials/{type}` — Get credentials by type for a specific Instance
- `GET /publicCloud/v1/instances/{instanceId}/credentials/{type}/{username}` — Get Instance credentials by type and username.
- `PUT /publicCloud/v1/instances/{instanceId}/credentials/{type}/{username}` — Update credentials for a given type and username
- `DELETE /publicCloud/v1/instances/{instanceId}/credentials/{type}/{username}` — Delete Instance credential for a given type and username
- `POST /publicCloud/v1/instances/{instanceId}/detachIso` — Detach ISO from a specific Instance
- `POST /publicCloud/v1/instances/{instanceId}/detachSecurityGroups` — Detach security groups from instance
- `GET /publicCloud/v1/instances/{instanceId}/instanceTypesUpdate` — List available instance types for update
- `GET /publicCloud/v1/instances/{instanceId}/ips` — List IP addresses associated with a specific Instance
- `GET /publicCloud/v1/instances/{instanceId}/ips/{ip}` — Get IP details for a specific Instance
- `PUT /publicCloud/v1/instances/{instanceId}/ips/{ip}` — Update the IP address for a specific Instance
- `POST /publicCloud/v1/instances/{instanceId}/ips/{ip}/null` — Null route IP address for a specific resource Instance
- `POST /publicCloud/v1/instances/{instanceId}/ips/{ip}/unnull` — Remove an IP null route for a specific Instance
- `GET /publicCloud/v1/instances/{instanceId}/metrics/cpu` — Get instance CPU metrics
- `GET /publicCloud/v1/instances/{instanceId}/metrics/datatraffic` — Get data traffic metrics for a specific Instance
- `GET /publicCloud/v1/instances/{instanceId}/notificationSettings/dataTraffic` — List the notification settings of a customer
- `GET /publicCloud/v1/instances/{instanceId}/notificationSettings/dataTraffic/{notificationSettingId}` — Get details of a notification Setting
- `POST /publicCloud/v1/instances/{instanceId}/notificationSettings/dataTraffic/{notificationSettingId}` — Create a notification setting
- `PUT /publicCloud/v1/instances/{instanceId}/notificationSettings/dataTraffic/{notificationSettingId}` — Update Notification Setting details
- `DELETE /publicCloud/v1/instances/{instanceId}/notificationSettings/dataTraffic/{notificationSettingId}` — Delete a notification setting setting
- `POST /publicCloud/v1/instances/{instanceId}/reboot` — Reboot a specific Instance
- `PUT /publicCloud/v1/instances/{instanceId}/reinstall` — Reinstall an Instance
- `GET /publicCloud/v1/instances/{instanceId}/reinstall/images` — List images available for reinstall
- `DELETE /publicCloud/v1/instances/{instanceId}/removeFromPrivateNetwork` — Remove instance from Private Network
- `POST /publicCloud/v1/instances/{instanceId}/resetPassword` — Reset the password for a specific Instance
- `GET /publicCloud/v1/instances/{instanceId}/securityGroups` — Get Instance Security Groups
- `GET /publicCloud/v1/instances/{instanceId}/snapshots` — List snapshots
- `POST /publicCloud/v1/instances/{instanceId}/snapshots` — Create instance snapshot
- `GET /publicCloud/v1/instances/{instanceId}/snapshots/{snapshotId}` — Get snapshot detail
- `PUT /publicCloud/v1/instances/{instanceId}/snapshots/{snapshotId}` — Restore instance snapshot
- `DELETE /publicCloud/v1/instances/{instanceId}/snapshots/{snapshotId}` — Delete instance snapshot
- `POST /publicCloud/v1/instances/{instanceId}/start` — Start a specific resource Instance
- `POST /publicCloud/v1/instances/{instanceId}/stop` — Stop a specific Instance
- `GET /publicCloud/v1/instances/{instanceId}/userData` — Get User Data supplied while launching the instance.
- `GET /publicCloud/v1/isos` — List available ISOs
- `POST /publicCloud/v1/loadBalancers` — Launch Load balancer
- `GET /publicCloud/v1/loadBalancers` — Get load balancer list
- `GET /publicCloud/v1/loadBalancers/{loadBalancerId}` — Get load balancer details
- `DELETE /publicCloud/v1/loadBalancers/{loadBalancerId}` — Delete load balancer
- `PUT /publicCloud/v1/loadBalancers/{loadBalancerId}` — Update load balancer
- `GET /publicCloud/v1/loadBalancers/{loadBalancerId}/ips` — List IP addresses associated with a specific Load Balancer
- `GET /publicCloud/v1/loadBalancers/{loadBalancerId}/ips/{ip}` — Get IP details for a specific Load Balancer
- `PUT /publicCloud/v1/loadBalancers/{loadBalancerId}/ips/{ip}` — Update the IP address for a specific Load Balancer
- `POST /publicCloud/v1/loadBalancers/{loadBalancerId}/ips/{ip}/null` — Null route IP address for a specific resource Load Balancer
- `POST /publicCloud/v1/loadBalancers/{loadBalancerId}/ips/{ip}/unnull` — Remove an IP null route for a specific Load Balancer
- `GET /publicCloud/v1/loadBalancers/{loadBalancerId}/listeners` — Get listener list
- `POST /publicCloud/v1/loadBalancers/{loadBalancerId}/listeners` — Create listener
- `GET /publicCloud/v1/loadBalancers/{loadBalancerId}/listeners/{listenerId}` — Get listener details
- `PUT /publicCloud/v1/loadBalancers/{loadBalancerId}/listeners/{listenerId}` — Update a listener
- `DELETE /publicCloud/v1/loadBalancers/{loadBalancerId}/listeners/{listenerId}` — Delete load balancer listener
- `GET /publicCloud/v1/loadBalancers/{loadBalancerId}/metrics/connections` — Get connections metrics
- `GET /publicCloud/v1/loadBalancers/{loadBalancerId}/metrics/connectionsPerSecond` — Get connections per second metrics
- `GET /publicCloud/v1/loadBalancers/{loadBalancerId}/metrics/cpu` — Get load balancer CPU metrics
- `GET /publicCloud/v1/loadBalancers/{loadBalancerId}/metrics/dataTransferred` — Get load balancer data transferred metrics
- `GET /publicCloud/v1/loadBalancers/{loadBalancerId}/metrics/dataTransferredPerSecond` — Get load balancer data transferred per second metrics
- `GET /publicCloud/v1/loadBalancers/{loadBalancerId}/metrics/datatraffic` — Get data traffic metrics for a specific Load Balancer
- `GET /publicCloud/v1/loadBalancers/{loadBalancerId}/metrics/requests` — Get load balancer requests metrics. Not available for listeners with TCP protocol
- `GET /publicCloud/v1/loadBalancers/{loadBalancerId}/metrics/requestsPerSecond` — Get load balancer requests per second metrics. Not available for listeners with TCP protocol
- `GET /publicCloud/v1/loadBalancers/{loadBalancerId}/metrics/responseCodes` — Get response codes metrics
- `GET /publicCloud/v1/loadBalancers/{loadBalancerId}/metrics/responseCodesPerSecond` — Get response codes per second metrics
- `POST /publicCloud/v1/loadBalancers/{loadBalancerId}/monitoring/enable` — Enable monitoring
- `GET /publicCloud/v1/loadBalancers/{loadBalancerId}/monitoring/status` — Get service monitoring status
- `POST /publicCloud/v1/loadBalancers/{loadBalancerId}/reboot` — Reboot a specific Load Balancer
- `POST /publicCloud/v1/loadBalancers/{loadBalancerId}/start` — Start a specific resource Load Balancer
- `POST /publicCloud/v1/loadBalancers/{loadBalancerId}/stop` — Stop a specific Load Balancer
- `GET /publicCloud/v1/marketApps` — Get marketplace apps
- `GET /publicCloud/v1/regions` — List regions
- `GET /publicCloud/v1/securityGroups` — Get Security Group list
- `POST /publicCloud/v1/securityGroups` — Create Security Group
- `GET /publicCloud/v1/securityGroups/{securityGroupId}` — Get Security Group details
- `PUT /publicCloud/v1/securityGroups/{securityGroupId}` — Update Security Group
- `DELETE /publicCloud/v1/securityGroups/{securityGroupId}` — Delete Security Group
- `POST /publicCloud/v1/securityGroups/{securityGroupId}/authorizeFirewallRules` — Add firewall rules to security group
- `GET /publicCloud/v1/securityGroups/{securityGroupId}/firewallRules` — Get Security Group Firewall Rules
- `POST /publicCloud/v1/securityGroups/{securityGroupId}/revokeFirewallRules` — Remove firewall rules from security group
- `GET /publicCloud/v1/targetGroups` — Get Target Group list
- `POST /publicCloud/v1/targetGroups` — Create Target Group
- `GET /publicCloud/v1/targetGroups/{targetGroupId}` — Get Target Group details
- `PUT /publicCloud/v1/targetGroups/{targetGroupId}` — Update Target Group
- `DELETE /publicCloud/v1/targetGroups/{targetGroupId}` — Delete Target Group
- `POST /publicCloud/v1/targetGroups/{targetGroupId}/deregisterTargets` — Deregister Targets
- `POST /publicCloud/v1/targetGroups/{targetGroupId}/registerTargets` — Register Targets
- `GET /publicCloud/v1/targetGroups/{targetGroupId}/targets` — Get Targets

### VPS

- `GET /publicCloud/v1/vps/` — Get VPS list
- `GET /publicCloud/v1/vps/isos` — List available ISOs
- `GET /publicCloud/v1/vps/{vpsId}` — Get VPS details
- `PUT /publicCloud/v1/vps/{vpsId}` — Update VPS details
- `POST /publicCloud/v1/vps/{vpsId}/attachIso` — Attach ISO to a specific VPS
- `GET /publicCloud/v1/vps/{vpsId}/console` — Get console access
- `GET /publicCloud/v1/vps/{vpsId}/credentials` — List credentials stored for a specific VPS
- `POST /publicCloud/v1/vps/{vpsId}/credentials` — Store credentials for a specific VPS
- `DELETE /publicCloud/v1/vps/{vpsId}/credentials` — Delete all credentials associated with a specific VPS
- `GET /publicCloud/v1/vps/{vpsId}/credentials/{type}` — Get credentials by type for a specific VPS
- `GET /publicCloud/v1/vps/{vpsId}/credentials/{type}/{username}` — Get VPS credential by type and username.
- `PUT /publicCloud/v1/vps/{vpsId}/credentials/{type}/{username}` — Update credentials for a given type and username
- `DELETE /publicCloud/v1/vps/{vpsId}/credentials/{type}/{username}` — Delete VPS credential
- `POST /publicCloud/v1/vps/{vpsId}/detachIso` — Detach ISO from a specific VPS
- `GET /publicCloud/v1/vps/{vpsId}/ips` — List IP addresses associated with a specific VPS
- `GET /publicCloud/v1/vps/{vpsId}/ips/{ip}` — Get IP details for a specific VPS
- `PUT /publicCloud/v1/vps/{vpsId}/ips/{ip}` — Update the IP address for a specific VPS
- `POST /publicCloud/v1/vps/{vpsId}/ips/{ip}/null` — Null route IP address for a specific resource VPS
- `POST /publicCloud/v1/vps/{vpsId}/ips/{ip}/unnull` — Remove an IP null route for a specific VPS
- `GET /publicCloud/v1/vps/{vpsId}/metrics/datatraffic` — Get data traffic metrics for a specific VPS
- `POST /publicCloud/v1/vps/{vpsId}/monitoring/enable` — Enable monitoring
- `GET /publicCloud/v1/vps/{vpsId}/monitoring/status` — Get service monitoring status
- `GET /publicCloud/v1/vps/{vpsId}/notificationSettings/dataTraffic` — List the notification settings of a customer
- `GET /publicCloud/v1/vps/{vpsId}/notificationSettings/dataTraffic/{notificationSettingId}` — Get details of a notification Setting
- `POST /publicCloud/v1/vps/{vpsId}/notificationSettings/dataTraffic/{notificationSettingId}` — Create a notification setting
- `PUT /publicCloud/v1/vps/{vpsId}/notificationSettings/dataTraffic/{notificationSettingId}` — Update Notification Setting details
- `DELETE /publicCloud/v1/vps/{vpsId}/notificationSettings/dataTraffic/{notificationSettingId}` — Delete a notification setting setting
- `POST /publicCloud/v1/vps/{vpsId}/reboot` — Reboot a specific VPS
- `PUT /publicCloud/v1/vps/{vpsId}/reinstall` — Reinstall a VPS
- `GET /publicCloud/v1/vps/{vpsId}/reinstall/images` — List images available for reinstall
- `POST /publicCloud/v1/vps/{vpsId}/resetPassword` — Reset the password for a specific VPS
- `GET /publicCloud/v1/vps/{vpsId}/snapshots` — List snapshots
- `POST /publicCloud/v1/vps/{vpsId}/snapshots` — Create VPS snapshot
- `GET /publicCloud/v1/vps/{vpsId}/snapshots/{snapshotId}` — Get snapshot detail
- `PUT /publicCloud/v1/vps/{vpsId}/snapshots/{snapshotId}` — Restore VPS snapshot
- `DELETE /publicCloud/v1/vps/{vpsId}/snapshots/{snapshotId}` — Delete VPS snapshot
- `POST /publicCloud/v1/vps/{vpsId}/start` — Start a specific VPS
- `POST /publicCloud/v1/vps/{vpsId}/stop` — Stop a specific VPS

### Virtual Servers

- `GET /cloud/v2/virtualServers` — List virtual servers
- `GET /cloud/v2/virtualServers/{virtualServerId}` — Inspect a virtual server
- `PUT /cloud/v2/virtualServers/{virtualServerId}` — Update a virtual server
- `PUT /cloud/v2/virtualServers/{virtualServerId}/credentials` — Update credentials of a virtual machine
- `GET /cloud/v2/virtualServers/{virtualServerId}/credentials/{credentialType}` — List credentials
- `GET /cloud/v2/virtualServers/{virtualServerId}/credentials/{credentialType}/{username}` — Inspect user credentials
- `GET /cloud/v2/virtualServers/{virtualServerId}/metrics/datatraffic` — Inspect datatraffic metrics
- `POST /cloud/v2/virtualServers/{virtualServerId}/powerOff` — Power off
- `POST /cloud/v2/virtualServers/{virtualServerId}/powerOn` — Power on
- `POST /cloud/v2/virtualServers/{virtualServerId}/reboot` — Reboot
- `POST /cloud/v2/virtualServers/{virtualServerId}/reinstall` — Reinstall
- `GET /cloud/v2/virtualServers/{virtualServerId}/snapshots` — List snapshots
- `POST /cloud/v2/virtualServers/{virtualServerId}/snapshots` — Create snapshot
- `GET /cloud/v2/virtualServers/{virtualServerId}/snapshots/{snapshotId}` — Get snapshot
- `DELETE /cloud/v2/virtualServers/{virtualServerId}/snapshots/{snapshotId}` — Delete snapshot
- `POST /cloud/v2/virtualServers/{virtualServerId}/snapshots/{snapshotId}/restore` — Restore snapshot
- `GET /cloud/v2/virtualServers/{virtualServerId}/templates` — List templates

### Private Clouds

- `GET /cloud/v2/privateClouds` — List private clouds
- `GET /cloud/v2/privateClouds/{[privateCloudId}/metrics/datatraffic` — Inspect datatraffic metrics
- `GET /cloud/v2/privateClouds/{privateCloudId}` — Inspect a private cloud
- `GET /cloud/v2/privateClouds/{privateCloudId}/credentials/{credentialType}` — List credentials
- `GET /cloud/v2/privateClouds/{privateCloudId}/credentials/{credentialType}/{username}` — Inspect user credentials
- `GET /cloud/v2/privateClouds/{privateCloudId}/metrics/bandwidth` — Inspect bandwidth metrics
- `GET /cloud/v2/privateClouds/{privateCloudId}/metrics/cpu` — Inspect CPU metrics
- `GET /cloud/v2/privateClouds/{privateCloudId}/metrics/memory` — Inspect memory metrics
- `GET /cloud/v2/privateClouds/{privateCloudId}/metrics/storage` — Inspect storage metrics

### Storage

- `GET /storage/v1/storageVMs` — 
- `GET /storage/v1/storageVMs/{storageVMId}/jobs/{jobId}` — 
- `GET /storage/v1/storageVMs/{storageVMId}/volumes` — 
- `POST /storage/v1/storageVMs/{storageVMId}/volumes/{volumeId}/grow` — 
- `GET /storage/v1/storages` — 

### Acronis Backup

- `GET /backup/v1/backup` — List backup items
- `GET /backup/v1/backup/{equipmentId}` — Inspect a backup item
- `GET /backup/v1/metrics/storage` — Get usage metrics

## Dedicated Services

### Dedicated Servers

- `GET /bareMetals/v2/controlPanels` — List control panels
- `GET /bareMetals/v2/hardwareMonitoring` — Show hardware monitoring data for all servers
- `GET /bareMetals/v2/operatingSystems` — List Operating Systems
- `GET /bareMetals/v2/operatingSystems/{operatingSystemId}` — Show an operating system
- `GET /bareMetals/v2/operatingSystems/{operatingSystemId}/controlPanels` — List control panels by Operating System
- `GET /bareMetals/v2/rescueImages` — Rescue Images
- `GET /bareMetals/v2/servers` — List servers
- `GET /bareMetals/v2/servers/{serverId}` — Get server
- `PUT /bareMetals/v2/servers/{serverId}` — Update server
- `POST /bareMetals/v2/servers/{serverId}/cancelActiveJob` — Cancel active job
- `GET /bareMetals/v2/servers/{serverId}/credentials` — List server credentials
- `POST /bareMetals/v2/servers/{serverId}/credentials` — Create new server credentials
- `GET /bareMetals/v2/servers/{serverId}/credentials/{type}` — List server credentials by type
- `DELETE /bareMetals/v2/servers/{serverId}/credentials/{type}/{username}` — Delete server credentials
- `GET /bareMetals/v2/servers/{serverId}/credentials/{type}/{username}` — Show server credentials
- `PUT /bareMetals/v2/servers/{serverId}/credentials/{type}/{username}` — Update server credentials
- `POST /bareMetals/v2/servers/{serverId}/expireActiveJob` — Expire active job
- `GET /bareMetals/v2/servers/{serverId}/hardwareInfo` — Show hardware information
- `GET /bareMetals/v2/servers/{serverId}/hardwareMonitoring` — Show hardware monitoring data for a server
- `POST /bareMetals/v2/servers/{serverId}/hardwareScan` — Launch hardware scan
- `POST /bareMetals/v2/servers/{serverId}/install` — Launch installation
- `POST /bareMetals/v2/servers/{serverId}/ipmiReset` — Launch IPMI reset
- `GET /bareMetals/v2/servers/{serverId}/ips` — List IPs
- `GET /bareMetals/v2/servers/{serverId}/ips/{ip}` — Show a server IP
- `PUT /bareMetals/v2/servers/{serverId}/ips/{ip}` — Update an IP
- `POST /bareMetals/v2/servers/{serverId}/ips/{ip}/null` — Null route an IP
- `POST /bareMetals/v2/servers/{serverId}/ips/{ip}/unnull` — Remove a null route
- `GET /bareMetals/v2/servers/{serverId}/jobs` — List jobs
- `GET /bareMetals/v2/servers/{serverId}/jobs/{jobId}` — Show a job
- `POST /bareMetals/v2/servers/{serverId}/jobs/{jobId}/retry` — Retry a job
- `DELETE /bareMetals/v2/servers/{serverId}/leases` — Delete a DHCP reservation
- `GET /bareMetals/v2/servers/{serverId}/leases` — List DHCP reservations
- `POST /bareMetals/v2/servers/{serverId}/leases` — Create a DHCP reservation
- `GET /bareMetals/v2/servers/{serverId}/metrics/bandwidth` — Show bandwidth metrics
- `GET /bareMetals/v2/servers/{serverId}/metrics/datatraffic` — Show datatraffic metrics
- `GET /bareMetals/v2/servers/{serverId}/networkInterfaces` — List network interfaces
- `POST /bareMetals/v2/servers/{serverId}/networkInterfaces/close` — Close all network interfaces
- `POST /bareMetals/v2/servers/{serverId}/networkInterfaces/open` — Open all network interfaces
- `GET /bareMetals/v2/servers/{serverId}/networkInterfaces/{networkTypeURL}` — Show a network interface
- `POST /bareMetals/v2/servers/{serverId}/networkInterfaces/{networkTypeURL}/close` — Close network interface
- `POST /bareMetals/v2/servers/{serverId}/networkInterfaces/{networkTypeURL}/open` — Open network interface
- `GET /bareMetals/v2/servers/{serverId}/notificationSettings/bandwidth` — List bandwidth notification settings
- `POST /bareMetals/v2/servers/{serverId}/notificationSettings/bandwidth` — Create a bandwidth notification setting
- `DELETE /bareMetals/v2/servers/{serverId}/notificationSettings/bandwidth/{notificationSettingId}` — Delete a bandwidth notification setting
- `GET /bareMetals/v2/servers/{serverId}/notificationSettings/bandwidth/{notificationSettingId}` — Show a bandwidth notification setting
- `PUT /bareMetals/v2/servers/{serverId}/notificationSettings/bandwidth/{notificationSettingId}` — Update a bandwidth notification setting
- `GET /bareMetals/v2/servers/{serverId}/notificationSettings/datatraffic` — List data traffic notification settings
- `POST /bareMetals/v2/servers/{serverId}/notificationSettings/datatraffic` — Create a data traffic notification setting
- `DELETE /bareMetals/v2/servers/{serverId}/notificationSettings/datatraffic/{notificationSettingId}` — Delete a data traffic notification setting
- `GET /bareMetals/v2/servers/{serverId}/notificationSettings/datatraffic/{notificationSettingId}` — Show a data traffic notification setting
- `PUT /bareMetals/v2/servers/{serverId}/notificationSettings/datatraffic/{notificationSettingId}` — Update a data traffic notification setting
- `GET /bareMetals/v2/servers/{serverId}/notificationSettings/ddos` — Inspect DDoS notification settings
- `PUT /bareMetals/v2/servers/{serverId}/notificationSettings/ddos` — Update DDoS notification settings
- `GET /bareMetals/v2/servers/{serverId}/nullRouteHistory` — Show null route history
- `POST /bareMetals/v2/servers/{serverId}/powerCycle` — Power cycle a server
- `GET /bareMetals/v2/servers/{serverId}/powerInfo` — Show power status
- `POST /bareMetals/v2/servers/{serverId}/powerOff` — Power off server
- `POST /bareMetals/v2/servers/{serverId}/powerOn` — Power on server
- `DELETE /bareMetals/v2/servers/{serverId}/privateNetworks/{privateNetworkId}` — Delete a server from a private network
- `PUT /bareMetals/v2/servers/{serverId}/privateNetworks/{privateNetworkId}` — Add a server to private network
- `POST /bareMetals/v2/servers/{serverId}/rescueMode` — Launch rescue mode

### Dedicated Network Equipments

- `GET /bareMetals/v2/networkEquipments` — List network equipment
- `GET /bareMetals/v2/networkEquipments/{networkEquipmentId}` — Get network equipment
- `PUT /bareMetals/v2/networkEquipments/{networkEquipmentId}` — Update network equipment
- `GET /bareMetals/v2/networkEquipments/{networkEquipmentId}/credentials` — List network equipment credentials
- `POST /bareMetals/v2/networkEquipments/{networkEquipmentId}/credentials` — Create new network equipment credentials
- `GET /bareMetals/v2/networkEquipments/{networkEquipmentId}/credentials/{type}` — List network equipment credentials by type
- `DELETE /bareMetals/v2/networkEquipments/{networkEquipmentId}/credentials/{type}/{username}` — Delete network equipment credentials
- `GET /bareMetals/v2/networkEquipments/{networkEquipmentId}/credentials/{type}/{username}` — Show network equipment credentials
- `PUT /bareMetals/v2/networkEquipments/{networkEquipmentId}/credentials/{type}/{username}` — Update network equipment credentials
- `GET /bareMetals/v2/networkEquipments/{networkEquipmentId}/ips` — List IPs
- `GET /bareMetals/v2/networkEquipments/{networkEquipmentId}/ips/{ip}` — Show a network equipment IP
- `PUT /bareMetals/v2/networkEquipments/{networkEquipmentId}/ips/{ip}` — Update an IP
- `POST /bareMetals/v2/networkEquipments/{networkEquipmentId}/ips/{ip}/null` — Null route an IP
- `POST /bareMetals/v2/networkEquipments/{networkEquipmentId}/ips/{ip}/unnull` — Remove a null route
- `GET /bareMetals/v2/networkEquipments/{networkEquipmentId}/nullRouteHistory` — Show null route history
- `POST /bareMetals/v2/networkEquipments/{networkEquipmentId}/powerCycle` — Power cycle a network equipment
- `GET /bareMetals/v2/networkEquipments/{networkEquipmentId}/powerInfo` — Show power status
- `POST /bareMetals/v2/networkEquipments/{networkEquipmentId}/powerOff` — Power off network equipment
- `POST /bareMetals/v2/networkEquipments/{networkEquipmentId}/powerOn` — Power on network equipment

### Dedicated Racks

- `GET /bareMetals/v2/privateRacks` — List dedicated racks
- `GET /bareMetals/v2/privateRacks/{privateRackId}` — Inspect a dedicated rack
- `PUT /bareMetals/v2/privateRacks/{privateRackId}` — Update a dedicated rack
- `GET /bareMetals/v2/privateRacks/{privateRackId}/credentials` — List credentials
- `POST /bareMetals/v2/privateRacks/{privateRackId}/credentials` — Create new credentials
- `GET /bareMetals/v2/privateRacks/{privateRackId}/credentials/{type}` — List credentials by type
- `GET /bareMetals/v2/privateRacks/{privateRackId}/credentials/{type}/{username}` — Inspect user credentials
- `PUT /bareMetals/v2/privateRacks/{privateRackId}/credentials/{type}/{username}` — Update user credentials
- `DELETE /bareMetals/v2/privateRacks/{privateRackId}/credentials/{type}/{username}` — Delete user credentials
- `GET /bareMetals/v2/privateRacks/{privateRackId}/ips` — List IPs
- `GET /bareMetals/v2/privateRacks/{privateRackId}/ips/{ip}` — Inspect an IP
- `PUT /bareMetals/v2/privateRacks/{privateRackId}/ips/{ip}` — Update an IP
- `POST /bareMetals/v2/privateRacks/{privateRackId}/ips/{ip}/null` — Null route an IP
- `POST /bareMetals/v2/privateRacks/{privateRackId}/ips/{ip}/unnull` — Remove a null route
- `GET /bareMetals/v2/privateRacks/{privateRackId}/metrics/bandwidth` — Inspect bandwidth metrics
- `GET /bareMetals/v2/privateRacks/{privateRackId}/metrics/datatraffic` — Inspect datatraffic metrics
- `GET /bareMetals/v2/privateRacks/{privateRackId}/notificationSettings/bandwidth` — List bandwidth notification settings
- `POST /bareMetals/v2/privateRacks/{privateRackId}/notificationSettings/bandwidth` — Create new bandwidth notification setting
- `GET /bareMetals/v2/privateRacks/{privateRackId}/notificationSettings/bandwidth/{notificationId}` — Inspect bandwidth notification settings
- `PUT /bareMetals/v2/privateRacks/{privateRackId}/notificationSettings/bandwidth/{notificationId}` — Update bandwidth notification settings
- `DELETE /bareMetals/v2/privateRacks/{privateRackId}/notificationSettings/bandwidth/{notificationId}` — Delete bandwidth notification settings
- `GET /bareMetals/v2/privateRacks/{privateRackId}/notificationSettings/datatraffic` — List datatraffic notification settings
- `POST /bareMetals/v2/privateRacks/{privateRackId}/notificationSettings/datatraffic` — Create datatraffic notification settings
- `GET /bareMetals/v2/privateRacks/{privateRackId}/notificationSettings/datatraffic/{notificationId}` — Inspect datatraffic notification settings
- `PUT /bareMetals/v2/privateRacks/{privateRackId}/notificationSettings/datatraffic/{notificationId}` — Update  datatraffic notification settings
- `DELETE /bareMetals/v2/privateRacks/{privateRackId}/notificationSettings/datatraffic/{notificationId}` — Delete datatraffic notification settings
- `GET /bareMetals/v2/privateRacks/{privateRackId}/notificationSettings/ddos` — Inspect DDos notification settings
- `PUT /bareMetals/v2/privateRacks/{privateRackId}/notificationSettings/ddos` — Update DDos notification settings
- `GET /bareMetals/v2/privateRacks/{privateRackId}/nullRouteHistory` — Inspect null route history

### Colocations

- `GET /bareMetals/v2/colocations` — List colocations
- `GET /bareMetals/v2/colocations/{colocationId}` — Inspect a colocation
- `PUT /bareMetals/v2/colocations/{colocationId}` — Update a colocation
- `GET /bareMetals/v2/colocations/{colocationId}/credentials` — List credentials
- `POST /bareMetals/v2/colocations/{colocationId}/credentials` — Create new credentials
- `GET /bareMetals/v2/colocations/{colocationId}/credentials/{type}` — List credentials by type
- `GET /bareMetals/v2/colocations/{colocationId}/credentials/{type}/{username}` — Inspect user credentials
- `PUT /bareMetals/v2/colocations/{colocationId}/credentials/{type}/{username}` — Update user credentials
- `DELETE /bareMetals/v2/colocations/{colocationId}/credentials/{type}/{username}` — Delete user credentials
- `GET /bareMetals/v2/colocations/{colocationId}/ips` — List IPs
- `GET /bareMetals/v2/colocations/{colocationId}/ips/{ip}` — Inspect an IP
- `PUT /bareMetals/v2/colocations/{colocationId}/ips/{ip}` — Update an IP
- `POST /bareMetals/v2/colocations/{colocationId}/ips/{ip}/null` — Null route an IP
- `POST /bareMetals/v2/colocations/{colocationId}/ips/{ip}/unnull` — Remove a null route
- `GET /bareMetals/v2/colocations/{colocationId}/metrics/bandwidth` — Inspect bandwidth metrics
- `GET /bareMetals/v2/colocations/{colocationId}/metrics/datatraffic` — Inspect datatraffic metrics
- `GET /bareMetals/v2/colocations/{colocationId}/networkInterfaces/public` — Inspect public network interface
- `POST /bareMetals/v2/colocations/{colocationId}/networkInterfaces/public/{action}` — Open or close public network interface
- `GET /bareMetals/v2/colocations/{colocationId}/notificationSettings/bandwidth` — List bandwidth notification settings
- `POST /bareMetals/v2/colocations/{colocationId}/notificationSettings/bandwidth` — Create new bandwidth notification settings
- `GET /bareMetals/v2/colocations/{colocationId}/notificationSettings/bandwidth/{notificationId}` — Inspect bandwidth notification settings
- `PUT /bareMetals/v2/colocations/{colocationId}/notificationSettings/bandwidth/{notificationId}` — Update bandwidth notification settings
- `DELETE /bareMetals/v2/colocations/{colocationId}/notificationSettings/bandwidth/{notificationId}` — Delete bandwidth notification settings
- `GET /bareMetals/v2/colocations/{colocationId}/notificationSettings/datatraffic` — List datatraffic notification settings
- `POST /bareMetals/v2/colocations/{colocationId}/notificationSettings/datatraffic` — Create new datatraffic notification settings
- `GET /bareMetals/v2/colocations/{colocationId}/notificationSettings/datatraffic/{notificationId}` — Inspect datatraffic notification settings
- `PUT /bareMetals/v2/colocations/{colocationId}/notificationSettings/datatraffic/{notificationId}` — Update datatraffic notification settings
- `DELETE /bareMetals/v2/colocations/{colocationId}/notificationSettings/datatraffic/{notificationId}` — Delete datatraffic notification settings
- `GET /bareMetals/v2/colocations/{colocationId}/notificationSettings/ddos` — Inspect DDoS notification settings
- `PUT /bareMetals/v2/colocations/{colocationId}/notificationSettings/ddos` — Update DDoS notifications settings
- `GET /bareMetals/v2/colocations/{colocationId}/nullRouteHistory` — Inspect null route history

### Datacenter Access

- `GET /datacenterAccess/v1/accessCards` — Get a list of access cards
- `POST /datacenterAccess/v1/accessCards` — Create an access card
- `GET /datacenterAccess/v1/accessCards/{id}` — Get an access card by User Id
- `PUT /datacenterAccess/v1/accessCards/{id}` — Update an access card
- `DELETE /datacenterAccess/v1/accessCards/{id}` — Delete an access card
- `GET /datacenterAccess/v1/visitRequests` — Get a list of access requests
- `POST /datacenterAccess/v1/visitRequests` — Create access request
- `GET /datacenterAccess/v1/visitRequests/{id}` — Get access request by ID

## Network

### IP management

- `GET /ipMgmt/v2/ips` — List IPs
- `GET /ipMgmt/v2/ips/{ip}` — Inspect an IP
- `PUT /ipMgmt/v2/ips/{ip}` — Update an IP
- `POST /ipMgmt/v2/ips/{ip}/nullRoute` — Null route an IP
- `DELETE /ipMgmt/v2/ips/{ip}/nullRoute` — Remove a null route
- `GET /ipMgmt/v2/ips/{ip}/nullRouted` — List null routed IPv6
- `GET /ipMgmt/v2/ips/{ip}/reverseLookup` — List reverse lookup records for an IPv6 range
- `PUT /ipMgmt/v2/ips/{ip}/reverseLookup` — Set or remove reverse lookup records for an IPv6 range
- `GET /ipMgmt/v2/nullRoutes` — Inspect null route history
- `GET /ipMgmt/v2/nullRoutes/{id}` — Inspect null route history
- `PUT /ipMgmt/v2/nullRoutes/{id}` — Update a null route

### Floating IPs

- `GET /floatingIps/v2/ranges` — List ranges
- `GET /floatingIps/v2/ranges/{rangeId}` — Inspect a range
- `GET /floatingIps/v2/ranges/{rangeId}/floatingIpDefinitions` — List range definitions
- `POST /floatingIps/v2/ranges/{rangeId}/floatingIpDefinitions` — Create a range definition
- `GET /floatingIps/v2/ranges/{rangeId}/floatingIpDefinitions/{floatingIpDefinitionId}` — Inspect a range definition
- `PUT /floatingIps/v2/ranges/{rangeId}/floatingIpDefinitions/{floatingIpDefinitionId}` — Update a range definition
- `DELETE /floatingIps/v2/ranges/{rangeId}/floatingIpDefinitions/{floatingIpDefinitionId}` — Remove a range definition
- `POST /floatingIps/v2/ranges/{rangeId}/floatingIpDefinitions/{floatingIpDefinitionId}/toggleAnchorIp` — Toggle the anchor IP with the passive anchor IP

### Private Networks

- `GET /bareMetals/v2/privateNetworks` — List private networks
- `POST /bareMetals/v2/privateNetworks` — Create a private network
- `GET /bareMetals/v2/privateNetworks/{privateNetworkId}` — Inspect a private network
- `PUT /bareMetals/v2/privateNetworks/{privateNetworkId}` — Update a private network
- `DELETE /bareMetals/v2/privateNetworks/{privateNetworkId}` — Delete a private network
- `GET /bareMetals/v2/privateNetworks/{privateNetworkId}/reservations` — Get DHCP reservations
- `POST /bareMetals/v2/privateNetworks/{privateNetworkId}/reservations` — Create a DHCP reservation
- `DELETE /bareMetals/v2/privateNetworks/{privateNetworkId}/reservations/{ip}` — Delete a DHCP reservation

### Remote Management

- `POST /bareMetals/v2/remoteManagement/changeCredentials` — Change OpenVPN credentials
- `GET /bareMetals/v2/remoteManagement/profiles` — List OpenVPN profiles
- `GET /bareMetals/v2/remoteManagement/profiles/lsw-rmvpn-{datacenter}.ovpn` — Inspect an OpenVPN profile

### Aggregation packs

- `GET /bareMetals/v2/aggregationPacks` — List aggregation packs
- `GET /bareMetals/v2/aggregationPacks/{aggregationPackId}` — Get aggregation pack

## CDN

### CDN

- `GET /cdn/v2/actions` — Get All Actions
- `GET /cdn/v2/actions/{action_id}` — Get Action By Id
- `POST /cdn/v2/certificates` — Create Certificate
- `GET /cdn/v2/certificates` — Get All Certificates
- `GET /cdn/v2/certificates/{certificate_id}` — Get Certificate By Id
- `DELETE /cdn/v2/certificates/{certificate_id}` — Delete Certificate By Id
- `POST /cdn/v2/debug` — Debug Url
- `POST /cdn/v2/distributions` — Create Distribution
- `GET /cdn/v2/distributions` — Get All Distributions
- `GET /cdn/v2/distributions/{distribution_id}` — Get Distribution By Id
- `PUT /cdn/v2/distributions/{distribution_id}` — Update Distribution By Id
- `DELETE /cdn/v2/distributions/{distribution_id}` — Delete Distribution By Id
- `POST /cdn/v2/distributions/{distribution_id}/domains` — Create Domain
- `GET /cdn/v2/distributions/{distribution_id}/domains` — Get All Domains
- `GET /cdn/v2/distributions/{distribution_id}/domains/{domain_id}` — Get Domain By Id
- `DELETE /cdn/v2/distributions/{distribution_id}/domains/{domain_id}` — Delete Domain By Id
- `POST /cdn/v2/distributions/{distribution_id}/policies` — Create Policy
- `GET /cdn/v2/distributions/{distribution_id}/policies` — Get All Policies
- `GET /cdn/v2/distributions/{distribution_id}/policies/{policy_id}` — Get Policy By Id
- `PUT /cdn/v2/distributions/{distribution_id}/policies/{policy_id}` — Update Policy By Id
- `DELETE /cdn/v2/distributions/{distribution_id}/policies/{policy_id}` — Delete Policy By Id
- `POST /cdn/v2/distributions/{distribution_id}/policies/{policy_id}/acls` — Create Acl
- `GET /cdn/v2/distributions/{distribution_id}/policies/{policy_id}/acls` — Get All Acls
- `GET /cdn/v2/distributions/{distribution_id}/policies/{policy_id}/acls/{acl_id}` — Get Acl By Id
- `DELETE /cdn/v2/distributions/{distribution_id}/policies/{policy_id}/acls/{acl_id}` — Delete Acl By Id
- `POST /cdn/v2/invalidations` — Do Invalidation
- `POST /cdn/v2/originGroups` — Create Origin Group
- `GET /cdn/v2/originGroups` — Get All Origin Groups
- `GET /cdn/v2/originGroups/{origin_group_id}` — Get Origin Group By Id
- `PUT /cdn/v2/originGroups/{origin_group_id}` — Update Origin Group By Id
- `DELETE /cdn/v2/originGroups/{origin_group_id}` — Delete Origin Group By Id
- `POST /cdn/v2/originGroups/{origin_group_id}/members` — Create Origin Group Member
- `GET /cdn/v2/originGroups/{origin_group_id}/members` — Get All Origin Group Members
- `GET /cdn/v2/originGroups/{origin_group_id}/members/{origin_group_member_id}` — Get Origin Group Member By Id
- `PUT /cdn/v2/originGroups/{origin_group_id}/members/{origin_group_member_id}` — Update Origin Group Member By Id
- `DELETE /cdn/v2/originGroups/{origin_group_id}/members/{origin_group_member_id}` — Delete Origin Group Member By Id
- `POST /cdn/v2/origins` — Create Origin
- `GET /cdn/v2/origins` — Get All Origins
- `GET /cdn/v2/origins/{origin_id}` — Get Origin By Id
- `PUT /cdn/v2/origins/{origin_id}` — Update Origin By Id
- `DELETE /cdn/v2/origins/{origin_id}` — Delete Origin By Id
- `GET /cdn/v2/origins/{origin_id}/resolve` — Validate If Orgin Can Be Resolved
- `GET /cdn/v2/policyTemplates` — Get All Templates
- `GET /cdn/v2/statistics/customer` — Get Statistics For Customer
- `GET /cdn/v2/statistics/customer/totals` — Get Total Statistics For Customer

## Hosting

### Domains

- `GET /hosting/v2/domains` — List domains
- `GET /hosting/v2/domains/{domainName}` — Inspect a domain
- `GET /hosting/v2/domains/{domainName}/available` — Check availability
- `GET /hosting/v2/domains/{domainName}/contacts` — List contacts
- `PUT /hosting/v2/domains/{domainName}/contacts` — Update all contacts
- `PUT /hosting/v2/domains/{domainName}/contacts/{type}` — Updating a contact
- `PUT /hosting/v2/domains/{domainName}/contacts/{type}/verify` — Contact Verification
- `GET /hosting/v2/domains/{domainName}/dnssec` — Inspect DNS Security (DNSSEC)
- `PUT /hosting/v2/domains/{domainName}/dnssec` — Update DNS Security (DNSSEC)
- `GET /hosting/v2/domains/{domainName}/health` — Get Domain Health
- `GET /hosting/v2/domains/{domainName}/locks` — Get Domain locks
- `POST /hosting/v2/domains/{domainName}/locks` — Domain locks
- `GET /hosting/v2/domains/{domainName}/nameservers` — List nameservers
- `PUT /hosting/v2/domains/{domainName}/nameservers` — Update nameservers
- `GET /hosting/v2/domains/{domainName}/overview` — Retrieve detailed information

### DNS

- `POST /hosting/v2/domains/{domainName}/keyRollover` — Key Rollover
- `GET /hosting/v2/domains/{domainName}/metrics/dnsQuery` — Show Dns Query Metrics
- `GET /hosting/v2/domains/{domainName}/resourceRecordSets` — List resource record sets
- `POST /hosting/v2/domains/{domainName}/resourceRecordSets` — Create a resource record set
- `PUT /hosting/v2/domains/{domainName}/resourceRecordSets` — Update all DNS records
- `DELETE /hosting/v2/domains/{domainName}/resourceRecordSets` — Delete all DNS records
- `GET /hosting/v2/domains/{domainName}/resourceRecordSets/export` — Export dns records as a bind file content
- `POST /hosting/v2/domains/{domainName}/resourceRecordSets/import` — Import dns records from bind file content
- `POST /hosting/v2/domains/{domainName}/resourceRecordSets/validateSet` — Validate a resource record set
- `GET /hosting/v2/domains/{domainName}/resourceRecordSets/{name}/{type}` — Inspect resource record set
- `PUT /hosting/v2/domains/{domainName}/resourceRecordSets/{name}/{type}` — Update a specific DNS record
- `DELETE /hosting/v2/domains/{domainName}/resourceRecordSets/{name}/{type}` — Delete a specific DNS record
- `POST /hosting/v2/domains/{domainName}/validateZone` — Validate zone

### Emails

- `GET /hosting/v2/domains/{domainName}/catchAll` — Retrieve catch all information
- `PUT /hosting/v2/domains/{domainName}/catchAll` — Change/create catch all
- `DELETE /hosting/v2/domains/{domainName}/catchAll` — Delete a catch all
- `GET /hosting/v2/domains/{domainName}/dkim` — Retrieve Dkim Record
- `POST /hosting/v2/domains/{domainName}/dkim` — Create Dkim Record
- `DELETE /hosting/v2/domains/{domainName}/dkim` — Delete Dkim Record
- `GET /hosting/v2/domains/{domainName}/emailAliases` — List email aliases
- `POST /hosting/v2/domains/{domainName}/emailAliases` — Create an email alias
- `GET /hosting/v2/domains/{domainName}/emailAliases/{source}/{destination}` — Inspect email alias
- `PUT /hosting/v2/domains/{domainName}/emailAliases/{source}/{destination}` — Update an email alias
- `DELETE /hosting/v2/domains/{domainName}/emailAliases/{source}/{destination}` — Delete an email alias
- `GET /hosting/v2/domains/{domainName}/forwards` — List domain forwards
- `GET /hosting/v2/domains/{domainName}/mailboxes` — List mailboxes
- `POST /hosting/v2/domains/{domainName}/mailboxes` — Create a mailbox
- `GET /hosting/v2/domains/{domainName}/mailboxes/{emailAddress}` — Inspect a mailbox
- `PUT /hosting/v2/domains/{domainName}/mailboxes/{emailAddress}` — Update a mailbox
- `DELETE /hosting/v2/domains/{domainName}/mailboxes/{emailAddress}` — Delete a mailbox
- `GET /hosting/v2/domains/{domainName}/mailboxes/{emailAddress}/autoResponder` — Inspect an auto responder
- `PUT /hosting/v2/domains/{domainName}/mailboxes/{emailAddress}/autoResponder` — Change an auto responder
- `DELETE /hosting/v2/domains/{domainName}/mailboxes/{emailAddress}/autoResponder` — Delete an auto responder
- `GET /hosting/v2/domains/{domainName}/mailboxes/{emailAddress}/forwards` — List forwards
- `POST /hosting/v2/domains/{domainName}/mailboxes/{emailAddress}/forwards` — Create a forward
- `GET /hosting/v2/domains/{domainName}/mailboxes/{emailAddress}/forwards/{destination}` — Inspect a forward
- `PUT /hosting/v2/domains/{domainName}/mailboxes/{emailAddress}/forwards/{destination}` — Update a forward
- `DELETE /hosting/v2/domains/{domainName}/mailboxes/{emailAddress}/forwards/{destination}` — Delete a forward

### Web hosting

- `GET /hosting/v2/webhosting/{domainName}/{equipmentId}/aliasdomain` — Get Alias Domains
- `POST /hosting/v2/webhosting/{domainName}/{equipmentId}/aliasdomain` — Add Alias Domain
- `PUT /hosting/v2/webhosting/{domainName}/{equipmentId}/aliasdomain/{aliasDomainName}` — Update Alias Domain
- `DELETE /hosting/v2/webhosting/{domainName}/{equipmentId}/aliasdomain/{aliasDomainName}` — Delete Alias Domain
- `GET /hosting/v2/webhosting/{domainName}/{equipmentId}/dotNet` — Get current selected dot net version
- `POST /hosting/v2/webhosting/{domainName}/{equipmentId}/dotNet` — Change dot net version
- `GET /hosting/v2/webhosting/{domainName}/{equipmentId}/subdomain` — Get Subdomains
- `POST /hosting/v2/webhosting/{domainName}/{equipmentId}/subdomain` — Add Subdomain
- `DELETE /hosting/v2/webhosting/{domainName}/{equipmentId}/subdomain/{subDomainName}` — Delete Sub Domain

### Traffic Policy

- `GET /hosting/v2/domains/{domainName}/trafficPolicies` — Get all traffic policy records
- `GET /hosting/v2/domains/{domainName}/trafficPolicies/{recordName}` — Get all traffic policies
- `POST /hosting/v2/domains/{domainName}/trafficPolicies/{recordName}` — Create traffic policies
- `PATCH /hosting/v2/domains/{domainName}/trafficPolicies/{recordName}` — Update traffic policies
- `DELETE /hosting/v2/domains/{domainName}/trafficPolicies/{recordName}` — Delete traffic policies

## Services

### Abuse Reports

- `GET /abuse/v1/reports` — List reports
- `GET /abuse/v1/reports/{reportId}` — Inspect a report
- `GET /abuse/v1/reports/{reportId}/messageAttachments/{fileId}` — Inspect a message attachment
- `GET /abuse/v1/reports/{reportId}/messages` — List report messages
- `POST /abuse/v1/reports/{reportId}/messages` — Create new message
- `GET /abuse/v1/reports/{reportId}/reportAttachments/{fileId}` — Inspect a report attachment
- `GET /abuse/v1/reports/{reportId}/resolutions` — List resolution options
- `POST /abuse/v1/reports/{reportId}/resolve` — Resolve a report

### Api Key management

- `GET /apiKeys/v1/keys` — List API keys
- `GET /apiKeys/v1/keys/{apiKeyId}` — Get API key
- `GET /apiKeys/v1/keys/{apiKeyId}/whiteListedIps` — List allowed IPs
- `POST /apiKeys/v1/keys/{apiKeyId}/whiteListedIps` — Add IP to allowed IP list
- `PUT /apiKeys/v1/keys/{apiKeyId}/whiteListedIps` — Replace allowed IP list
- `DELETE /apiKeys/v1/keys/{apiKeyId}/whiteListedIps` — Delete all whitelisted ips for an API key.
- `DELETE /apiKeys/v1/keys/{apiKeyId}/whiteListedIps/{ip}` — Delete allowed IP

### Invoices

- `GET /invoices/v1/invoices` — List Invoices
- `GET /invoices/v1/invoices/export/csv` — Export Invoices (CSV)
- `GET /invoices/v1/invoices/{Id}` — Inspect a invoice
- `GET /invoices/v1/invoices/{Id}/pdf` — Get Invoice PDF
- `GET /invoices/v1/proforma` — Pro Forma

### Ordering

- `GET /ordering/v1/products/dedicatedServers` — List available dedicated server configurations.
- `GET /ordering/v1/products/dedicatedServers/{dedicatedServerId}` — Get a single dedicated server and its price.
- `POST /ordering/v1/products/dedicatedServers/{dedicatedServerId}/order` — Dedicated Server ordering.
- `GET /ordering/v1/products/vps` — Get a list of available VPS products and their prices.
- `GET /ordering/v1/products/vps/{vpsId}` — Get a single VPS and its price.
- `POST /ordering/v1/products/vps/{vpsId}/order` — VPS ordering.

### Services

- `GET /services/v1/services` — List Services
- `GET /services/v1/services/cancellationReasons` — Get cancellation reasons
- `GET /services/v1/services/{Id}` — Inspect a service
- `PUT /services/v1/services/{Id}` — Update a service
- `POST /services/v1/services/{Id}/cancel` — Cancel a Service.
- `POST /services/v1/services/{Id}/uncancel` — Uncancel a Service.

### Orders

- `GET /account/v1/orders` — List of orders
- `GET /account/v1/orders/{Id}` — Inspect an order

### Aws Ec2 Wrapper

- `GET /aws/ec2/?Action=CreateListener` — CreateListener
- `GET /aws/ec2/?Action=CreateLoadBalancer` — CreateLoadBalancer
- `GET /aws/ec2/?Action=DescribeInstances` — DescribeInstances
- `GET /aws/ec2/?Action=DescribeLoadBalancers` — DescribeLoadBalancers
- `GET /aws/ec2/?Action=RebootInstances` — RebootInstances
- `GET /aws/ec2/?Action=RegisterTargets` — RegisterTargets
- `GET /aws/ec2/?Action=RunInstances` — RunInstances
- `GET /aws/ec2/?Action=StartInstances` — StartInstances
- `GET /aws/ec2/?Action=StopInstances` — StopInstances

