# Machinist

Machinist is an application that makes it easy for customers to collect data from their OPC UA data servers and publish the data to Google Cloud IoT Core. Machinist connects to any OPC UA server on your factory's network. You can add metrics for collection at any time through a browser. If internet connectivity is interrupted, Machinist stores metric values for a configurable time period and forwards the data to Google when your connection is restored. Machinist is available for Windows, Linux and MacOS.

## A Walkthrough

A customer visits your website and chooses to trial the service. After accepting the license, the customer downloads two files. The first file 'config.yaml' is unique to the customer. It's contents:

    kind: gateway
    projectID: our-chassis-269114
    region: us-central1
    registryID: golang-iot-test-registry
    deviceID: gateway-92d6686e-92de-4c3f-bca2-dea6f054ab24
    privateKey: |
      -----BEGIN RSA PRIVATE KEY-----
      -----END RSA PRIVATE KEY-----
    algorithm: RS256
    storePath: ./data

The customer chooses a second file to download, depending on the preferred OS. Example 'machinist-windows-amd64-0.1.0.zip'.

The customer uncompresses the second file and places it in a directory along with the first file ‘config.yaml’. The customer then starts the Machinist application.

## Walkthrough, Part 2

The customer returns to the website and configures the machinist by adding machines and metrics.

### Adding a machine to the machinist

The customer adds a machine and specifies a name, opc ua endpoint url and various properties that all have reasonable defaults (sampling interval, publishing interval, etc.)

### Adding a metric to a machine

The customer adds a metric to a machine and specifies a name, opc ua nodeID and various properties that all have reasonable defaults. The customer could also browse through the namespace of the opc ua server and select the metrics of interest.

## Walkthrough, Part 3

The Machinist application is always listening for service requests, such as 'Browse', 'Read', 'Write', and 'Call' published by a client application to the Google Cloud IoT 'commands' topic. For example:
1) A client issues a 'Write' request and publishes it to the command topic of a particular device, i.e. '/devices/\<device-id\>/commands/write-request'. 
2) Machinist receives the request and writes to the OPC UA server on the factory network. 
3) The OPC UA server sends the 'Write' response back to Machinist.
4) Machinist forwards the response to Google Cloud IoT topic '/devices/\<device-id\>/events/write-response'.
5) The client reads the response and can check the status code.