service "nginx" {
  unit {
    description = "NGINX Web"
    after = ["network.target"]
  }

  service {
    exec_start = "/usr/sbin/nginx -g 'daemon off;'"
    standard_output = "journal"
  }

  install {
    wanted_by = ["multi-user.target"]
  }
}
