## stackit beta vpn bgp-filter update

Updates a VPN BGP filter

### Synopsis

Updates the metadata of a VPN BGP filter. The rules of the filter are managed via 'stackit beta vpn bgp-filter rule'.

```
stackit beta vpn bgp-filter update FILTER_ID [flags]
```

### Examples

```
  Rename the BGP filter with the ID "yyy" of the gateway with the ID "xxx" to "my-filter"
  $ stackit beta vpn bgp-filter update yyy --gateway-id xxx --display-name my-filter
```

### Options

```
      --display-name string   A user friendly name for the filter
      --gateway-id string     Gateway ID
  -h, --help                  Help for "stackit beta vpn bgp-filter update"
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

* [stackit beta vpn bgp-filter](./stackit_beta_vpn_bgp-filter.md)	 - Provides functionality for VPN BGP filters

