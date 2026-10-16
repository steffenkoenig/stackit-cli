## stackit beta vpn bgp-filter rule describe

Shows details of a VPN BGP filter rule

### Synopsis

Shows details of a VPN BGP filter rule.

```
stackit beta vpn bgp-filter rule describe RULE_ID [flags]
```

### Examples

```
  Show details of the rule with the ID "zzz" of the BGP filter with the ID "yyy"
  $ stackit beta vpn bgp-filter rule describe zzz --gateway-id xxx --filter-id yyy

  Show details of the rule with the ID "zzz" of the BGP filter with the ID "yyy" in JSON format
  $ stackit beta vpn bgp-filter rule describe zzz --gateway-id xxx --filter-id yyy --output-format json
```

### Options

```
      --filter-id string    BGP filter ID
      --gateway-id string   Gateway ID
  -h, --help                Help for "stackit beta vpn bgp-filter rule describe"
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

