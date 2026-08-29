## stackit beta vpn connection update

Updates a VPN connection

### Synopsis

Updates a VPN connection.

Only the attributes given as flags are changed; every other attribute keeps its current value.
Pre-shared keys are never returned by the API, so they are only changed when passed explicitly.

```
stackit beta vpn connection update CONNECTION_ID [flags]
```

### Examples

```
  Rename the VPN connection with the ID "zzz"
  $ stackit beta vpn connection update zzz --gateway-id xxx --display-name my-connection

  Attach the BGP filter "yyy" to the BGP session of tunnel 1
  $ stackit beta vpn connection update zzz --gateway-id xxx --tunnel1-bgp-inbound-filter-id yyy

  Disable the VPN connection with the ID "zzz"
  $ stackit beta vpn connection update zzz --gateway-id xxx --enabled=false
```

### Options

```
      --display-name string                            A user friendly name for the connection.
      --enabled                                        Enable the connection (default true)
      --gateway-id string                              Required: Gateway ID
  -h, --help                                           Help for "stackit beta vpn connection update"
      --labels stringToString                          Map of custom labels. Key and values must be a string with max 63 chars, start/end with alphanumeric. (example: foo=bar) (default [])
      --local-subnets strings                          Defaults to 0.0.0.0/0 for Route-based VPN configurations. Mandatory for Policy-based.
      --remote-subnets strings                         Defaults to 0.0.0.0/0 for Route-based VPN configurations. Mandatory for Policy-based.
      --static-routes strings                          Use this for route-based VPN.
      --tunnel1-bgp-inbound-filter-id string           Tunnel 1 BGP Inbound Filter ID.
                                                       ID of the BGP filter to apply to incoming routes from this tunnel's BGP neighbor. The filter must exist in the same gateway and the tunnel must have a BGP remote ASN.
      --tunnel1-bgp-remote-asn int                     Tunnel 1 BGP Remote ASN.
                                                       ASN for private use (reserved by IANA), both 16Bit and 32Bit ranges are valid (RFC 6996).
      --tunnel1-peering-local-address string           Tunnel 1 Peering Local Address
      --tunnel1-peering-remote-address string          Tunnel 1 Peering Remote Address
      --tunnel1-phase1-dh-groups strings               Tunnel 1 Phase 1 DH Groups.
                                                       The Diffie-Hellman Group. Required, except if AEAD algorithms are selected. (multiple of: [modp1024, modp2048, ecp256, ecp384, modp2048s256]) (default [])
      --tunnel1-phase1-encryption-algorithms strings   Tunnel 1 Phase 1 Encryption Algorithms (multiple of: [aes256, aes128gcm16, aes256gcm16]) (default [])
      --tunnel1-phase1-integrity-algorithms strings    Tunnel 1 Phase 1 Integrity Algorithms (multiple of: [sha1, sha2_256, sha2_384, sha2_512]) (default [])
      --tunnel1-phase1-rekey-time int                  Tunnel 1 Phase 1 Rekey Time.
                                                       Time to schedule a IKE re-keying (in seconds).
      --tunnel1-phase2-dh-groups strings               Tunnel 1 Phase 2 DH Groups (multiple of: [modp1024, modp2048, ecp256, ecp384, modp2048s256]) (default [])
      --tunnel1-phase2-dpd-action string               Tunnel 1 Phase 2 DPD Action.
                                                       Action to perform for this CHILD_SA on DPD timeout. (one of: [clear, restart])
      --tunnel1-phase2-encryption-algorithms strings   Tunnel 1 Phase 2 Encryption Algorithms (multiple of: [aes256, aes128gcm16, aes256gcm16]) (default [])
      --tunnel1-phase2-integrity-algorithms strings    Tunnel 1 Phase 2 Integrity Algorithms (multiple of: [sha1, sha2_256, sha2_384, sha2_512]) (default [])
      --tunnel1-phase2-rekey-time int                  Tunnel 1 Phase 2 Rekey Time.
                                                       Time to schedule a Child SA re-keying (in seconds).
      --tunnel1-phase2-start-action string             Tunnel 1 Phase 2 Start Action.
                                                       Action to perform after loading the connection configuration. (one of: [none, start])
      --tunnel1-pre-shared-key string                  Tunnel 1 Pre Shared Key.
                                                       Only changed when passed explicitly.
      --tunnel1-remote-address string                  Tunnel 1 Remote Address
      --tunnel2-bgp-inbound-filter-id string           Tunnel 2 BGP Inbound Filter ID.
                                                       ID of the BGP filter to apply to incoming routes from this tunnel's BGP neighbor. The filter must exist in the same gateway and the tunnel must have a BGP remote ASN.
      --tunnel2-bgp-remote-asn int                     Tunnel 2 BGP Remote ASN
      --tunnel2-peering-local-address string           Tunnel 2 Peering Local Address
      --tunnel2-peering-remote-address string          Tunnel 2 Peering Remote Address
      --tunnel2-phase1-dh-groups strings               Tunnel 2 Phase 1 DH Groups.
                                                       The Diffie-Hellman Group. Required, except if AEAD algorithms are selected. (multiple of: [modp1024, modp2048, ecp256, ecp384, modp2048s256]) (default [])
      --tunnel2-phase1-encryption-algorithms strings   Tunnel 2 Phase 1 Encryption Algorithms (multiple of: [aes256, aes128gcm16, aes256gcm16]) (default [])
      --tunnel2-phase1-integrity-algorithms strings    Tunnel 2 Phase 1 Integrity Algorithms (multiple of: [sha1, sha2_256, sha2_384, sha2_512]) (default [])
      --tunnel2-phase1-rekey-time int                  Tunnel 2 Phase 1 Rekey Time.
                                                       Time to schedule a IKE re-keying (in seconds).
      --tunnel2-phase2-dh-groups strings               Tunnel 2 Phase 2 DH Groups (multiple of: [modp1024, modp2048, ecp256, ecp384, modp2048s256]) (default [])
      --tunnel2-phase2-dpd-action string               Tunnel 2 Phase 2 DPD Action.
                                                       Action to perform for this CHILD_SA on DPD timeout. (one of: [clear, restart])
      --tunnel2-phase2-encryption-algorithms strings   Tunnel 2 Phase 2 Encryption Algorithms (multiple of: [aes256, aes128gcm16, aes256gcm16]) (default [])
      --tunnel2-phase2-integrity-algorithms strings    Tunnel 2 Phase 2 Integrity Algorithms (multiple of: [sha1, sha2_256, sha2_384, sha2_512]) (default [])
      --tunnel2-phase2-rekey-time int                  Tunnel 2 Phase 2 Rekey Time.
                                                       Time to schedule a Child SA re-keying (in seconds).
      --tunnel2-phase2-start-action string             Tunnel 2 Phase 2 Start Action.
                                                       Action to perform after loading the connection configuration. (one of: [none, start])
      --tunnel2-pre-shared-key string                  Tunnel 2 Pre Shared Key.
                                                       Only changed when passed explicitly.
      --tunnel2-remote-address string                  Tunnel 2 Remote Address
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

* [stackit beta vpn connection](./stackit_beta_vpn_connection.md)	 - Provides functionality for VPN connections

