# Ddosify Self hosted

## Introduction

## Features

✅ **No Code Web UI**: No need to write code to run load tests. Simply specify the target URL and the load test parameters.

✅ **Distributed Mode**: Scale your load testing capabilities by utilizing multiple engines simultaneously.

✅ **Arm64 and Amd64 Support**: Broad architecture compatibility ensures the tool works seamlessly across different systems on both Linux and MacOS.

✅ **Dockerized**: Containerized solution simplifies deployment and reduces dependency management overhead.

✅ **Easy to Deploy**: Automated setup processes using Docker Compose.


## 🛠 Prerequisites

- [Git](https://git-scm.com/)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/) (`docker-compose` or `docker compose`)

## ⚡️ Quick Start (Recommended)

You can quickly deploy Ddosify Self Hosted by running the following command. This script clones the Ddosify repository to your `$HOME/.ddosify` directory, and deploys the services using Docker Compose. Please check the [install.sh](./install.sh) file to see what it does.

```bash
curl -sSL https://raw.githubusercontent.com/ddosify/ddosify/selfhosted_release/selfhosted/install.sh | bash
```

## 📖 Manual Installation

### 1. Clone the repository

```bash
git clone https://github.com/ddosify/ddosify.git
cd ddosify/selfhosted
```

### 2. Update the environment variables (optional)

The default values for the environment variables are set in the [.env](./.env) file. You can modify these values to suit your needs. The following environment variables are available:

- `DOCKER_INFLUXDB_INIT_USERNAME`: InfluxDB username. Default: `admin`
- `DOCKER_INFLUXDB_INIT_PASSWORD`: InfluxDB password. Default: `ChangeMe`
- `DOCKER_INFLUXDB_INIT_ADMIN_TOKEN`: InfluxDB admin token. Default: `5yR2qD5zCqqvjwCKKXojnPviQaB87w9JcGweVChXkhWRL`
- `POSTGRES_PASSWORD`: Postgres password. Default: `ChangeMe`

### 3. Deploy the services

```bash
docker-compose up -d
```
### 4. Access the dashboard

The dashboard is available at [http://localhost:8014](http://localhost:8014)

## 🔧 Add New Engine

The Ddosify Self Hosted includes a default engine out of the box. To integrate additional engines, simply run a Docker container for each new engine. These engine containers will automatically register with the service and become available for use. Before adding new engines, ensure that you have enabled the distributed mode by clicking the `Unlock the Distributed Mode` button in the dashboard.

In case you have modified the default values like InfluxDB password in the `.env` file, utilize the `--env` flag in the docker run command to establish the necessary environment variables.

Make sure the new engine server can access the service server. Use the `DDOSIFY_SERVICE_ADDRESS` environment variable to specify the service server address where the [install.sh](install.sh) script was executed.

The engine server must connect to the following ports on the `DDOSIFY_SERVICE_ADDRESS`:

- `9901`: Hammer Manager service. The service server utilizes this port to register the engine.
- `6672`: RabbitMQ server. The engine server connects to this port to send and receive messages to and from the service server.
- `9086`: InfluxDB server. The engine server accesses this port to transmit metrics to the backend.
- `9900`: Object storage server. The engine server uses this port to exchange files with the service server.

### **Example 1**: Adding the engine to the same server

```bash
NAME=ddosify_hammer_1
docker run --name $NAME -dit \
    --network selfhosted_ddosify \
    --restart always \
    ddosify/selfhosted_hammer
```

### **Example 2**: Adding the engine to a different server

Set `DDOSIFY_SERVICE_ADDRESS` to the IP address of the service server. Set `IP_ADDRESS` to the IP address of the engine server.

```bash
# Make sure to set the following environment variables
DDOSIFY_SERVICE_ADDRESS=16.62.190.13
IP_ADDRESS=local
NAME=ddosify_hammer_1

docker run --name $NAME -dit \
    --env DDOSIFY_SERVICE_ADDRESS=$DDOSIFY_SERVICE_ADDRESS \
    --env IP_ADDRESS=$IP_ADDRESS \
    --restart always \
    ddosify/selfhosted_hammer
```

You should see `mq_waiting_new_job` log in the engine container logs. This means that the engine is waiting for a job from the service server. After the engine is added, you can see it in the Engines page in the dashboard.


## 🧹 Clean Up

If you added new engines, stop the engines first. Change the name of the engine container to the name of your engine container.

```bash
docker rm -f ddosify_hammer_1
```

If you installed the project using the [install.sh](./install.sh) script, you must first change the directory to the `$HOME/.ddosify` directory before running the commands below.

```bash
cd $HOME/.ddosify/selfhosted
docker compose down
```

## 🧩 Services Overview

| Service              | Description                                                                                       |
|----------------------|---------------------------------------------------------------------------------------------------|
| Hammer               | The engine responsible for executing load tests. You can add multiple hammers to scale your load testing capabilities.                                                  |
| Hammer Manager       | Manages the engines (Hammers) involved in load testing.                                           |
| Backend              | Handles load test management and stores results.                                                  |
| InfluxDB             | Database that stores metrics collected during testing.                                            |
| Postgres             | Database that preserves load test results.                                                        |
| RabbitMQ             | Message broker enabling communication between Hammer Manager and Hammers.                         |
| Minio Object Storage | Object storage for multipart files and test data (CSV) used in load tests.                        |
| Nginx                | Reverse proxy for backend and frontend services.                                                  |

## 📝 License

Ddosify Self Hosted is licensed under the AGPLv3: https://www.gnu.org/licenses/agpl-3.0.html