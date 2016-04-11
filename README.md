# Creating Kodayif Agent Image

```
cd docker/agent
docker build -t="kodayif-agent" .
```

# Creating Kodayif Controller Image

```
cd docker/controller
docker build -t="kodayif-controller" .
```

# Running RabbitMQ Container

```
docker run --name kodayif-rabbit -d rabbitmq:3.6
```

# Running Redis Container

```
docker run --name kodayif-redis -d redis:3.0
```

# Running Kodayif Controller Container

```
docker run --name kodayif-controller -p 127.0.0.1:8080:8080 --link kodayif-rabbit:rabbitmq --link kodayif-redis:redis -d kodayif-controller
```

# Running Kodayif Agent Container

```
docker run --name kodayif-agent-1 --link kodayif-rabbit:rabbitmq --link kodayif-controller:kodayif-controller -d kodayif-agent
docker run --name kodayif-agent-2 --link kodayif-rabbit:rabbitmq --link kodayif-controller:kodayif-controller -d kodayif-agent
docker run --name kodayif-agent-3 --link kodayif-rabbit:rabbitmq --link kodayif-controller:kodayif-controller -d kodayif-agent
```
