package main

var config_template = `port: 8080
prometheus:
  switch_on: true
  port: 8081
service_name: {{.Package.Name}}
register:
  switch_on: true
  register_path: /microservice/service/
  timeout: 1s
  heart_beat: 10
  register_name: etcd
  register_addr: 127.0.0.1:2379
log:
  level: debug
  path: ./logs/
  chan_size: 10000
  console_log: true
limit:
  switch_on: true
  qps: 10
trace:
  switch_on: true
  report_addr: http://127.0.0.1:9412/api/v1/spans
  sample_type: const
  sample_rate: 1
`
