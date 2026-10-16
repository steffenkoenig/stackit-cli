## stackit beta vpn bgp-filter delete

Deletes a VPN BGP filter

### Synopsis

Deletes a VPN BGP filter. A filter which is still attached to a tunnel cannot be deleted.

```
stackit beta vpn bgp-filter delete FILTER_ID [flags]
```

### Examples

```
  Delete the BGP filter with the ID "yyy" of the gateway with the ID "xxx"
  $ stackit beta vpn bgp-filter delete yyy --gateway-id xxx
```

### Options

```
      --gateway-id string   Gateway ID
  -h, --help                Help for "stackit beta vpn bgp-filter delete"
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

