# Blocks

This document enumerates **all valid top-level blocks** supported by the system, their purpose, and how they participate in the compilation model.

Blocks are declarative. **No block has side effects unless referenced by a `host` block**.

---

## Terraform‑Based (Language Infrastructure)

These blocks exist to provide **configuration structure and reuse**, inspired by Terraform, but intentionally limited.

### `variable`

Declares an external input to the configuration.

- Typed
- Optional default
- Immutable
- Globally scoped

Used for:
- Environment names
- Paths
- Ports
- Feature flags

```hcl
variable "env" {
  type    = string
  default = "prod"
}
```

---

### `output`

Declares a computed value exposed to the user after evaluation.

- Read‑only
- No effect on system state
- Useful for debugging or introspection

```hcl
output "enabled_services" {
  value = host.web01.enable
}
```

---

### `data`

Reads **static local data** at compile time.

Allowed sources:
- Local files
- Inline literals

Disallowed:
- Network access
- Remote APIs
- Host access

```hcl
data "file" "nginx_tpl" {
  path = "nginx.service.tpl"
}
```

---

### `locals`

Defines pure computed values.

- Expressions only
- Deterministic
- No I/O

```hcl
locals {
  service_name = "nginx-${var.env}"
}
```

---

### `module` (optional / constrained)

Provides **structural reuse only**.

Constraints:
- No providers
- No state
- No lifecycle hooks
- Inputs = variables
- Outputs = outputs

Used only for:
- Grouping reusable service definitions
- Namespacing large configs

This block may be omitted entirely in early versions.

---

## Hosts

### `host`

Defines a **single reconciliation boundary**.

A host:
- Selects which units and policies apply
- Defines SSH connection metadata
- Is compiled and applied independently

```hcl
host "web01" {
  ssh {
    user = "root"
    host = "10.0.0.1"
  }

  enable  = [service.nginx, target.web_stack]
  disable = [timer.nginx_reload]
}
```

Rules:
- Hosts do not share state
- Hosts cannot reference each other
- A block not referenced by a host is inert

---

## systemd Units

All systemd unit blocks map **1:1 to native systemd unit types**.

Each unit block:
- Is reusable
- Is inert until enabled by a host
- Renders exactly one unit file

### Common Unit Sections

All unit types may contain:
- `unit {}`
- `install {}`

Additional sections depend on unit type.

---

### `service`

Maps to `.service` units.

Used for:
- Long‑running daemons
- One‑shot tasks

Additional sections:
- `service {}`
- `sandbox {}`
- `resources {}`
- `environment {}`

---

### `timer`

Maps to `.timer` units.

Used for:
- Scheduled execution
- Periodic maintenance tasks

Additional sections:
- `timer {}`

---

### `path`

Maps to `.path` units.

Used for:
- File or directory change triggers

Additional sections:
- `path {}`

---

### `target`

Maps to `.target` units.

Used for:
- Logical grouping of units
- High‑level enablement

---

### `socket`

Maps to `.socket` units.

Used for:
- Socket‑activated services

Additional sections:
- `socket {}`

---

### `mount`

Maps to `.mount` units.

Used for:
- Mounting filesystems

Additional sections:
- `mount {}`

---

### `automount`

Maps to `.automount` units.

Used for:
- On‑demand filesystem mounts

Additional sections:
- `automount {}`

---

### `swap`

Maps to `.swap` units.

Used for:
- Swap device management

Additional sections:
- `swap {}`

---

### `device`

Maps to `.device` units.

Used for:
- Hardware dependency ordering

---

### `busname`

Maps to `.busname` units.

Used for:
- D‑Bus name activation

Additional sections:
- `busname {}`

---

### `snapshot`

Maps to `.snapshot` units.

Used for:
- Capturing unit states

Typically informational; rarely enabled explicitly.

---

## Related Subsystems (systemd Satellites)

These blocks generate **configuration fragments** for systemd‑related subsystems.

They are reusable and applied **per host**.

---

## Journald

### `journal`

Manages `systemd‑journald` configuration.

Maps to:
- `/etc/systemd/journald.conf.d/*.conf`

Features:
- Stackable
- Ordered
- Conflict‑checked

```hcl
journal "server" {
  storage  = "persistent"
  compress = true
}
```

---

## Resolved (DNS)

### `dns`

Manages `systemd‑resolved` configuration.

Maps to:
- `/etc/systemd/resolved.conf.d/*.conf`

Features:
- Stackable
- Ordered
- Conflict‑checked

```hcl
dns "cloudflare" {
  servers = ["1.1.1.1", "1.0.0.1"]
}
```

---

## Timesync

### `time`

Manages `systemd‑timesyncd` configuration.

Maps to:
- `/etc/systemd/timesyncd.conf`

Rules:
- Only one `time` block may be applied per host

---

## sysusers

### `user`

Declares a system user.

Maps to:
- `sysusers.d`

```hcl
user "nginx" {
  uid   = 101
  shell = "/usr/sbin/nologin"
}
```

---

### `group`

Declares a system group.

Maps to:
- `sysusers.d`

---

## tmpfiles

### `file`

Declares a managed file.

Maps to:
- `tmpfiles.d`

---

### `directory`

Declares a managed directory.

Maps to:
- `tmpfiles.d`

---

## Optional / Future Extensions (Non‑binding)

These blocks are intentionally excluded from v1 but considered compatible:

- `logind` → systemd‑logind configuration
- `oomd` → systemd‑oomd policies
- `network` → systemd‑networkd (advanced, risky)
- `coredump` → systemd‑coredump config

These must follow the **policy block model** if added.

---

## Block Validity Rules

- Unknown blocks are compile‑time errors
- Blocks may not nest unless explicitly defined
- Block names are globally unique per type

---

## Summary

> Blocks describe **what exists**.
>
> Hosts decide **what applies**.
>
> systemd decides **how it runs**.

