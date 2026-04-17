service "nginx" {
  unit {
    description = "NGINX Web"

    after   = [builtin.target.network, service.db]
    wants   = [builtin.target.network_online]
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
