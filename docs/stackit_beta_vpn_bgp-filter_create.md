## stackit beta vpn bgp-filter create

Creates a VPN BGP filter

### Synopsis

Creates a VPN BGP filter.

The filter is created without any rules. Since an empty filter denies every route, add rules with
'stackit beta vpn bgp-filter rule create' before attaching it to a tunnel.

```
stackit beta vpn bgp-filter create [flags]
```

### Examples

```
  Create a VPN BGP filter named "my-filter" for the gateway with the ID "xxx"
  $ stackit beta vpn bgp-filter create --gateway-id xxx --display-name my-filter
```

### Options

```
      --display-name string   A user friendly name for the filter
      --gateway-id string     Gateway ID
  -h, --help                  Help for "stackit beta vpn bgp-filter create"
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

