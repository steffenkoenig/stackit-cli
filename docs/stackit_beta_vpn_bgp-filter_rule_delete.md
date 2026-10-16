## stackit beta vpn bgp-filter rule delete

Deletes a VPN BGP filter rule

### Synopsis

Deletes a VPN BGP filter rule.

Since an implicit DENY follows the last rule of a filter, removing the last PERMIT rule of a filter
denies every route.

```
stackit beta vpn bgp-filter rule delete RULE_ID [flags]
```

### Examples

```
  Delete the rule with the ID "zzz" of the BGP filter with the ID "yyy"
  $ stackit beta vpn bgp-filter rule delete zzz --gateway-id xxx --filter-id yyy
```

### Options

```
      --filter-id string    BGP filter ID
      --gateway-id string   Gateway ID
  -h, --help                Help for "stackit beta vpn bgp-filter rule delete"
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

