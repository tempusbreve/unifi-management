# unifi-management

Tools for managing devices on a UniFi network.

## Environment

Use environment variables `UNIFI_ENDPOINT`, `UNIFI_USERNAME`, and `UNIFI_PASSWORD` to connect.

## CLI

To build the CLI: `go build ./cmd/uncli`

CLI usage: ```sh
$ uncli list                    # show "known" devices
$ uncli block name [... name]   # block devices matching name(s) 
$ uncli unblock name [... name] # unblock devices matching name(s) 
```
