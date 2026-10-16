## stackit beta vpn gateway update

Updates a vpn gateway

### Synopsis

Updates a vpn gateway.

Only the attributes given as flags are changed; every other attribute keeps its current value.
The routing type of a gateway can only be set at creation and is therefore not updatable. The same
applies to the predefined network prefix.

```
stackit beta vpn gateway update GATEWAY_ID [flags]
```

### Examples

```
  Rename the vpn gateway with the ID "xxx" to "yyy"
  $ stackit beta vpn gateway update xxx --name yyy

  Update the labels of the vpn gateway with the ID "xxx"
  $ stackit beta vpn gateway update xxx --labels foo=bar,x=y

  Change the advertised BGP routes of the vpn gateway with the ID "xxx"
  $ stackit beta vpn gateway update xxx --bgp-local-asn 64512 --bgp-override-advertised-routes 10.0.0.0/24,10.1.0.0/24
```

### Options

```
      --availability-zone-tunnel-1 string            Availability Zone of Tunnel 1
      --availability-zone-tunnel-2 string            Availability Zone of Tunnel 2
      --bgp-local-asn int                            ASN for private use (reserved by IANA), both 16Bit and 32Bit ranges are valid (RFC 6996)
      --bgp-override-advertised-routes stringArray   A list of IPv4 Prefixes to advertise via BGP
  -h, --help                                         Help for "stackit beta vpn gateway update"
      --labels stringToString                        Labels in key=value format, separated by commas (default [])
      --name string                                  Gateway name
      --network-routing-table-id string              Custom routing table ID for the vpn gateway
      --plan-id string                               Plan ID
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, (one of: [json, pretty, none, yaml])
  -p, --project-id string      Project ID
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, (one of: [debug, info, warning, error]) (default "info")
```

### SEE ALSO

* [stackit beta vpn gateway](./stackit_beta_vpn_gateway.md)	 - Provides functionality for VPN gateway

