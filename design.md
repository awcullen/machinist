## Implementation Details

### Google Cloud IoT Core

Machinist communicates using MQTT to IoT Core through Gateways and Devices.

Machinist subscribes to the following topics:
- /devices/gateway-id/config
- /devices/gateway-id/commands/#

Machinist publishes to the following topics:
- /devices/gateway-id/state
- /devices/gateway-id/events/#

When Machinist receives a config, it starts a machine for each device-id listed in config.

Machines subscribe to the following topics:
- /devices/device-id/config
- /devices/device-id/commands/#

Machine publish to the following topics:
- /devices/device-id/state
- /devices/device-id/events/#
