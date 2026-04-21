# unitd

A unit policy configuration DSL for [`systemd`](https://github.com/systemd/systemd).

## A vision for deterministic and validated configurations

### Problem

[`systemd`](https://github.com/systemd/systemd) has been a powerful system and service manager for Linux distributions for many years.
However, working with its configuration language and surrounding tooling can still be challenging.

Key issues include:

- Reliance on external tools to fetch and inject secrets from cloud providers.
- Dependence on orchestration tools such as Ansible for deployment automation.
- Incomplete developer tooling (e.g., LSPs and formatters), with limited support for advanced features such as references, composition, and large-scale deterministic configuration management.

These limitations motivate the design of a dedicated toolchain that addresses these concerns in a unified and coherent way.

### Approach

DevOps engineers familiar with cloud-native environments are often already experienced with Terraform.

Fortunately, HashiCorp has open-sourced the HCL (HashiCorp Configuration Language), which provides a strong foundation for designing declarative configuration systems.

Building on this idea, `unitd` explores a Terraform-like DSL tailored specifically for configuring `systemd`, aiming to improve:

- Determinism in configuration management
- Built-in validation and type safety
- First-class tooling support (LSP, formatting, refactoring)
- Reduced dependency on external automation layers

### Example

```hcl
systemd {
    version = "v260.1"
}

data "azurerm_key_vault" "example" {
  name                = "mykeyvault"
  resource_group_name = "some-resource-group"
}

service "db" {
  unit {
    description = "Database"
    after = [builtin.target.network]
  }

  service {
    exec_start = "/usr/bin/db-server"
    environment     = {
        SECRET = data.azurerm_key_vault.example.vault_uri
    }
  }

  install {
    wanted_by = [builtin.target.multi_user]
  }
}


service "nginx" {
  unit {
    description = "NGINX Web"
    after = [builtin.target.network, service.db]
  }

  service {
    exec_start      = "/usr/sbin/nginx -g 'daemon off;'"
    standard_output = "journal"
  }

  install {
    wanted_by = [builtin.target.multi_user]
  }
}

host "web01" {
  ssh {
    user = "root"
    host = "10.0.0.1"
  }

  enable  = [service.nginx, service.db]
  disable = [timer.nginx_reload]
}

host "web02" {
  ssh {
    name = "ssh-config-name"
  }

  enable  = [service.db]
  disable = [timer.nginx_reload]
}
```
