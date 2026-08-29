## stackit beta vpn bgp-filter describe

Shows details of a VPN BGP filter

### Synopsis

Shows details of a VPN BGP filter.

```
stackit beta vpn bgp-filter describe FILTER_ID [flags]
```

### Examples

```
  Show details of the BGP filter with the ID "yyy" of the gateway with the ID "xxx"
  $ stackit beta vpn bgp-filter describe yyy --gateway-id xxx

  Show details of the BGP filter with the ID "yyy" of the gateway with the ID "xxx" in JSON format
  $ stackit beta vpn bgp-filter describe yyy --gateway-id xxx --output-format json
```

### Options

```
      --gateway-id string   Gateway ID
  -h, --help                Help for "stackit beta vpn bgp-filter describe"
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

