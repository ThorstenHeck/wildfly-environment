# HPKV local dev environment

This project aims to create a lightweight containerized shareable environment for a Wildfly, Database (OracleDB or Postgresql) stack with a focus on deployment of AIS Software.

The overall goal is to create a service which we can feed with a defined zip archive and the output will be a running AIS App.

## Architecture

Internal Container:

- Wildfly 
- OracleDB/PostgreSQL 
- Operator

External Services:

- Gitlab (SCM)
- Nexus (Registry)

### Explaination

We will spin up our multi-container Docker container environment with Docker Compose, the Operator Container will trigger commands via ansible to deploy the Java Application to Wildfly and prepare the database.

To be able to choose between the databases Postgresql and Oracle, we have to set profiles, since its the only way to set conditional statements via docker-compose itself.  
https://github.com/compose-spec/compose-spec/blob/master/spec.md#profiles

### Installation

    git clone https://github.com/ThorstenHeck/wildfly-environment.git
    cd wildfly-environment
    docker-compose --profile postgres --build up -d

### operator

The operator contains of a basic golang app, which provides a RESTFUL API to execute ansible.

#### Install golang

https://go.dev/doc/install

##### alpine (dockerfile)

COPY --from=golang:1.13-alpine /usr/local/go/ /usr/local/go/
 
ENV PATH="/usr/local/go/bin:${PATH}"




docker build -t operator .

docker run -it -p 10000:10000 --rm operator /bin/bash

docker-compose --profile postgres up --build -d

docker-compose logs -f operator

ansible-galaxy collection install -r requirements.yml

ansible-galaxy collection install middleware_automation.wildfly

ansible-galaxy collection install middleware_automation.jcliff
https://ansiblemiddleware.com/ansible_collections_jcliff/main/plugins/jcliff_module.html#ansible-collections-middleware-automation-jcliff-jcliff-module




docker-compose --profile oracle up -d

docker rm -f wildfly-environment-operator-1
docker rm -f wildfly-environment-db-postgres-1
docker rm -f wildfly-environment-wildfly-1
docker network rm wildfly-environment_hpkv

docker-compose --profile postgres up -d
docker-compose --profile postgres down
docker-compose --profile postgres ps

docker rm -f wildfly-environment-operator-1
docker rm -f wildfly-environment-db-oracle-1
docker rm -f wildfly-environment-wildfly-1
docker network rm wildfly-environment_hpkv

docker exec -it wildfly-environment-wildfly-1 /opt/jboss/wildfly/bin/add-user.sh -u 'adminuser1' -p 'adminuser1' -g 'admin'

docker exec -it wildfly-environment-db-postgres-1 /bin/bash


docker-compose build wildfly --build-arg ADMINPW=CLI




docker-compose build wildfly operator db-postgres db-oracle --build-arg WILDFLY_ADMIN_PW=Password!


ssh-keygen -t ed25519 -f ./ssh/ -q -N "" 