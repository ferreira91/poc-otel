# POC otel
This repository contains a demonstration application that collects and stores telemetry data. The application was created as a proof of concept (POC) to show the ability to collect and process real-time telemetry data.

### Requirements
* [cURL](https://curl.se/)
* [Docker](https://www.docker.com/)
* [Docker Compose](https://docs.docker.com/compose/)
* [Make](https://www.gnu.org/software/make/)

Este é um projeto de prova de conceito (POC) para demonstração de uma solução de open telemetry. A aplicação é capaz de coletar e enviar métricas, rastros e logs para um sistema de monitoramento.

### Setup
Clone the repository:
```$ git clone https://github.com/ferreira91/poc-otel.git```

### Run
Run: 
```$ make up```

### Usage
Run:
```$ curl http:\\localhost:1323\test

### Show results
Jaeger:
```$ http://localhost:16686```
Prometheus:
```$ http://localhost:9090```
Grafana:
```$ http://localhost:3000```

### Final considerations
This is a simple example of how to collect and store telemetry data using Go. It is important to remember that the POC was developed only for demonstration purposes and additional improvements may be necessary to meet production requirements. Additionally, it is important to test and evaluate the scalability of the application before implementing it in a production environment.
