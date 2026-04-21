service "nginx" {
  unit {
    description = "NGINX Web"

    after = [builtin.target.network, service.db, instance.queue_workers["q1"]]
    wants = [builtin.target.network_online]
  }

  service {
    exec_start = "/usr/sbin/nginx -g 'daemon off;'"
    standard_output = "journal"
  }

  install {
    wanted_by = [builtin.target.multi_user]
  }
}

service "db" {
  unit {
    description = "Database"

    after = [builtin.target.network]
  }

  service {
    exec_start = "/usr/bin/db-server"
  }

  install {
    wanted_by = [builtin.target.multi_user]
  }
}

service "worker" {
  template = true
  for_each = {
    queue = "Queue processing"
    email = "Email sending"
  }

  unit {
    description = "Worker - ${each.value}"

    after = [builtin.target.network, service.db]
  }

  service {
    exec_start = "/usr/bin/worker --type ${each.key} --name ${self.instance}"
  }

  install {
    wanted_by = [builtin.target.multi_user]
  }
}

instance "queue_workers" {
  template  = service.worker["queue"]
  instances = ["q1", "q2"]
}
