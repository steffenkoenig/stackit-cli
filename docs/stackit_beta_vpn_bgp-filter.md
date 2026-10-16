## stackit beta vpn bgp-filter

Provides functionality for VPN BGP filters

### Synopsis

Provides functionality for VPN BGP filters.

A BGP filter holds an ordered set of rules which are used to filter the routes a VPN gateway
receives via BGP. A route is evaluated against each rule in sequence order and the first match
decides the outcome. An implicit DENY is applied after the last rule, so a route which matches no
PERMIT rule is dropped. Attaching an empty filter to a tunnel therefore denies every route.

A filter is attached to a tunnel of a VPN connection via the '--tunnel1-bgp-inbound-filter-id' and
'--tunnel2-bgp-inbound-filter-id' flags of the 'stackit beta vpn connection create' command.

```
stackit beta vpn bgp-filter [flags]
```

### Options

```
  -h, --help   Help for "stackit beta vpn bgp-filter"
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

* [stackit beta vpn](./stackit_beta_vpn.md)	 - Provides functionality for VPN
* [stackit beta vpn bgp-filter create](./stackit_beta_vpn_bgp-filter_create.md)	 - Creates a VPN BGP filter
* [stackit beta vpn bgp-filter delete](./stackit_beta_vpn_bgp-filter_delete.md)	 - Deletes a VPN BGP filter
* [stackit beta vpn bgp-filter describe](./stackit_beta_vpn_bgp-filter_describe.md)	 - Shows details of a VPN BGP filter
* [stackit beta vpn bgp-filter list](./stackit_beta_vpn_bgp-filter_list.md)	 - Lists all VPN BGP filters of a gateway
* [stackit beta vpn bgp-filter rule](./stackit_beta_vpn_bgp-filter_rule.md)	 - Provides functionality for VPN BGP filter rules
* [stackit beta vpn bgp-filter update](./stackit_beta_vpn_bgp-filter_update.md)	 - Updates a VPN BGP filter

