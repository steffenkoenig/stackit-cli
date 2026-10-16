## stackit beta vpn bgp-filter rule create

Creates a VPN BGP filter rule

### Synopsis

Creates a VPN BGP filter rule.

All match criteria which are set must match for the rule to apply (logical AND). If no match
criterion is set, the rule matches every route. When the sequence is omitted, the rule is appended
after the last existing rule of the filter.

```
stackit beta vpn bgp-filter rule create [flags]
```

### Examples

```
  Permit every route received by the BGP filter with the ID "yyy" of the gateway with the ID "xxx"
  $ stackit beta vpn bgp-filter rule create --gateway-id xxx --filter-id yyy --action PERMIT

  Deny the routes 10.0.0.0/16 and 192.168.0.0/16 as the first rule of the filter
  $ stackit beta vpn bgp-filter rule create --gateway-id xxx --filter-id yyy --action DENY --sequence 5 --match-prefixes 10.0.0.0/16,192.168.0.0/16

  Permit the routes advertised by the peer 192.168.0.1 whose AS-PATH starts with the ASN 65000 and prefer them over other routes
  $ stackit beta vpn bgp-filter rule create --gateway-id xxx --filter-id yyy --action PERMIT --match-peer 192.168.0.1 --match-first-asn 65000 --set-local-preference 200
```

### Options

```
      --action string                           The action to take if a route matches all criteria of the rule (one of: [PERMIT, DENY])
      --filter-id string                        BGP filter ID
      --gateway-id string                       Gateway ID
  -h, --help                                    Help for "stackit beta vpn bgp-filter rule create"
      --match-as-path-contains-any int64Slice   Match routes whose AS-PATH contains any of the given ASNs (default [])
      --match-communities strings               Match routes carrying any of the given BGP standard communities in the 'asn:value' format (RFC 1997)
      --match-first-asn int                     Match routes whose first ASN in the AS-PATH, i.e. the immediate neighbor, equals the given ASN
      --match-max-prefix-length int32           Match routes with a subnet mask length of at most this value
      --match-min-prefix-length int32           Match routes with a subnet mask length of at least this value
      --match-peer string                       Match routes advertised by the BGP neighbor with this exact IPv4 address
      --match-prefixes strings                  Match routes whose prefix equals one of the given IPv4 networks (default [])
      --sequence int32                          The evaluation order of the rule. Lower numbers are evaluated first and the number must be unique within the filter. Defaults to the highest existing sequence plus 10.
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

