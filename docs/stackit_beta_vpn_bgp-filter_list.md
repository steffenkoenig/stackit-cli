## stackit beta vpn bgp-filter list

Lists all VPN BGP filters of a gateway

### Synopsis

Lists all VPN BGP filters of a gateway.

```
stackit beta vpn bgp-filter list [flags]
```

### Examples

```
  List all BGP filters of the gateway with the ID "xxx"
  $ stackit beta vpn bgp-filter list --gateway-id xxx

  List up to 4 BGP filters of the gateway with the ID "xxx"
  $ stackit beta vpn bgp-filter list --gateway-id xxx --limit 4
```

### Options

```
      --gateway-id string   Gateway ID
  -h, --help                Help for "stackit beta vpn bgp-filter list"
      --limit int           Maximum number of entries to list
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

