package args

import (
	"fmt"
	settingsPackage "gitlab.com/mcsolutions/tools/gek/utils/settings"
	"os"
	"text/template"
)

const (
	route = `                        - match: { prefix: "/api.{{ .Name }}." }
                          route:
                            cluster: {{ .Name }}-{{ .Version }}-cluster
                            max_grpc_timeout: 0s
`
	cluster = `    - name: {{ .Name }}-{{ .Version }}-cluster
      connect_timeout: 0.25s
      type: logical_dns
      http2_protocol_options: {}
      lb_policy: round_robin
      load_assignment:
        cluster_name: {{ .Name }}-{{ .Version }}__cluster
        endpoints:
          - lb_endpoints:
              - endpoint:
                  address:
                    socket_address:
                      address: {{ .Name }}-{{ .Version }}
                      port_value: 9000
`
	service = `
  {{ .Name }}-{{ .Version }}:
    image: registry.gitlab.com/mcsolutions/find-psy/back/{{ .Name }}-{{ .Version }}:v0.0.1
    build:
      context: ../../back/{{ .Name }}-{{ .Version }}
      dockerfile: Dockerfile
    restart: always
    env_file:
      - ./.env
    volumes:
      - ./configs/{{ .Name }}-{{ .Version }}:/data/conf
#    ports:
#      - "{{ .Port }}:9000"
#      - "{{ .HttpPort }}:8000"
    depends_on:
      - postgres
      - jaeger
    networks:
      - backend
`
	config = `
server:
  grpc:
    addr: :9000
    timeout: 1s
  http:
    addr: :8000
    timeout: 1s
trace:
  endpoint: ${JAEGER_ENDPOINT}
#sentry:
#  dsn: ${SENTRY_DSN}
data:
  relational:
    dialect: postgres
    host: ${POSTGRES_HOST}
    port: ${POSTGRES_PORT}
    user: ${POSTGRES_USER}
    password: ${POSTGRES_PASSWORD}
    dbname: ${POSTGRES_DB}
    schema: {{ .Name }}
    ssl_mode: ${POSTGRES_SSL_MODE}
jwks:
  url: ${JWKS_URL}
  refresh_interval: 3600s
  refresh_rate_limit: 300s
  refresh_timeout: 10s
business:
`
)

func ArgEnvoy() error {
	services := settingsPackage.Services{}
	if err := services.Load(); err != nil {
		return err
	}
	fmt.Println("---")
	for _, service := range services.Services {
		fmt.Printf("      - %s-%s\n", service.Name, service.Version)
	}
	fmt.Println("---")
	for _, service := range services.Services {
		fmt.Println()
		fmt.Printf("dc-%s:\n", service.Name)
		fmt.Printf("\tdocker-compose build %s-%s\n", service.Name, service.Version)
		fmt.Printf("\tdocker-compose stop %s-%s\n", service.Name, service.Version)
		fmt.Printf("\tdocker-compose up -d %s-%s\n", service.Name, service.Version)
	}
	fmt.Println("---")
	fmt.Println("clone-all:")
	for _, service := range services.Services {
		fmt.Printf("\tgit clone https://gitlab.com/mcsolutions/find-psy/back/%s-%s.git ../../back\n", service.Name, service.Version)
	}
	fmt.Println("")
	fmt.Println("pull-all:")
	for _, service := range services.Services {
		fmt.Printf("\tgit -C ../../back/%s-%s pull\n", service.Name, service.Version)
	}
	fmt.Println("---")
	tempService, err := template.New("service").Parse(service)
	if err != nil {
		return err
	}
	for p, service := range services.Services {
		service.Port = int16(p + 9001)
		service.HttpPort = int16(p + 8001)
		if err := tempService.Execute(os.Stdout, service); err != nil {
			return err
		}
	}
	fmt.Println("---")
	tempConfig, err := template.New("config").Parse(config)
	if err != nil {
		return err
	}
	for _, service := range services.Services {
		if err := tempConfig.Execute(os.Stdout, service); err != nil {
			return err
		}
	}
	fmt.Println("---")
	fmt.Println("                      routes:")
	tempRoute, err := template.New("route").Parse(route)
	if err != nil {
		return err
	}
	for _, service := range services.Services {
		if err := tempRoute.Execute(os.Stdout, service); err != nil {
			return err
		}
	}
	fmt.Println("---")
	fmt.Println("  clusters:")
	tempCluster, err := template.New("cluster").Parse(cluster)
	if err != nil {
		return err
	}
	for _, service := range services.Services {
		if err := tempCluster.Execute(os.Stdout, service); err != nil {
			return err
		}
	}
	return nil
}
