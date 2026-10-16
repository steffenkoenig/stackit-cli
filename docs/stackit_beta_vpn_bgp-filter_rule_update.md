## stackit beta vpn bgp-filter rule update

Updates a VPN BGP filter rule

### Synopsis

Updates a VPN BGP filter rule.

The rule is replaced as a whole, so every match criterion which should be kept has to be passed
again. Criteria which are not passed are removed from the rule.

```
stackit beta vpn bgp-filter rule update RULE_ID [flags]
```

### Examples

```
  Turn the rule with the ID "zzz" into a rule which permits every route
  $ stackit beta vpn bgp-filter rule update zzz --gateway-id xxx --filter-id yyy --action PERMIT --sequence 10

  Change the rule with the ID "zzz" to deny the routes 10.0.0.0/16 and 192.168.0.0/16
  $ stackit beta vpn bgp-filter rule update zzz --gateway-id xxx --filter-id yyy --action DENY --sequence 10 --match-prefixes 10.0.0.0/16,192.168.0.0/16
```

### Options

```
      --action string                           The action to take if a route matches all criteria of the rule (one of: [PERMIT, DENY])
      --filter-id string                        BGP filter ID
      --gateway-id string                       Gateway ID
  -h, --help                                    Help for "stackit beta vpn bgp-filter rule update"
      --match-as-path-contains-any int64Slice   Match routes whose AS-PATH contains any of the given ASNs (default [])
      --match-communities strings               Match routes carrying any of the given BGP standard communities in the 'asn:value' format (RFC 1997)
      --match-first-asn int                     Match routes whose first ASN in the AS-PATH, i.e. the immediate neighbor, equals the given ASN
      --match-max-prefix-length int32           Match routes with a subnet mask length of at most this value
      --match-min-prefix-length int32           Match routes with a subnet mask length of at least this value
      --match-peer string                       Match routes advertised by the BGP neighbor with this exact IPv4 address
      --match-prefixes strings                  Match routes whose prefix equals one of the given IPv4 networks (default [])
      --sequence int32                          The evaluation order of the rule. Lower numbers are evaluated first and the number must be unique within the filter.
      --set-local-preference int32              The BGP LOCAL_PREF to set on a permitted route. Higher values are preferred during best-path selection, the BGP default is 100. Ignored for DENY rules.
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

