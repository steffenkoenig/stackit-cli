## stackit beta vpn bgp-filter rule list

Lists all rules of a VPN BGP filter

### Synopsis

Lists all rules of a VPN BGP filter, ordered by their sequence.

```
stackit beta vpn bgp-filter rule list [flags]
```

### Examples

```
  List all rules of the BGP filter with the ID "yyy" of the gateway with the ID "xxx"
  $ stackit beta vpn bgp-filter rule list --gateway-id xxx --filter-id yyy

  List up to 4 rules of the BGP filter with the ID "yyy"
  $ stackit beta vpn bgp-filter rule list --gateway-id xxx --filter-id yyy --limit 4
```

### Options

```
      --filter-id string    BGP filter ID
      --gateway-id string   Gateway ID
  -h, --help                Help for "stackit beta vpn bgp-filter rule list"
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

* [stackit beta vpn bgp-filter rule](./stackit_beta_vpn_bgp-filter_rule.md)	 - Provides functionality for VPN BGP filter rules

