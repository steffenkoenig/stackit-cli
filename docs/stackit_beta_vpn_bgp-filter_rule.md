## stackit beta vpn bgp-filter rule

Provides functionality for VPN BGP filter rules

### Synopsis

Provides functionality for VPN BGP filter rules.

All non-empty match criteria of a rule are AND-combined. Rules within a filter are evaluated in
sequence order and the first matching rule decides the outcome. An implicit DENY follows the last
rule. A filter holds at most 10 rules.

```
stackit beta vpn bgp-filter rule [flags]
```

### Options

```
  -h, --help   Help for "stackit beta vpn bgp-filter rule"
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
* [stackit beta vpn bgp-filter rule create](./stackit_beta_vpn_bgp-filter_rule_create.md)	 - Creates a VPN BGP filter rule
* [stackit beta vpn bgp-filter rule delete](./stackit_beta_vpn_bgp-filter_rule_delete.md)	 - Deletes a VPN BGP filter rule
* [stackit beta vpn bgp-filter rule describe](./stackit_beta_vpn_bgp-filter_rule_describe.md)	 - Shows details of a VPN BGP filter rule
* [stackit beta vpn bgp-filter rule list](./stackit_beta_vpn_bgp-filter_rule_list.md)	 - Lists all rules of a VPN BGP filter
* [stackit beta vpn bgp-filter rule update](./stackit_beta_vpn_bgp-filter_rule_update.md)	 - Updates a VPN BGP filter rule

