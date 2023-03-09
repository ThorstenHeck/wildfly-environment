FROM alpine:3.17
LABEL description="Ansible on Alpine 3.17 and Python 3.10"
ARG ANSIBLE_VERSION=7.3.0

ARG SSH_PRIV_KEY=
ARG SSH_PUB_KEY=

COPY --from=golang:1.13-alpine /usr/local/go/ /usr/local/go/

ENV PATH="/usr/local/go/bin:${PATH}"

ENV USER_ID=2001
ENV GROUP_ID=2001
ENV USER_NAME=ansible
ENV GROUP_NAME=ansible

ENV GROUP_ID_SUDO=110
ENV GROUP_NAME_SUDO=sudo

RUN /bin/sh -c set -e && \
    apk update --no-cache && \
    apk upgrade --no-cache && \
    apk add --no-cache gcc make python3 python3-dev openssl-dev \
    py3-cffi py3-bcrypt py-cryptography py3-pynacl py3-pip bash curl && \
    pip3 install --no-cache-dir pip && \
    pip3 install --no-cache-dir ansible==${ANSIBLE_VERSION} && \
    addgroup -g $GROUP_ID $GROUP_NAME && \
    adduser --uid $USER_ID --disabled-password --home /home/ansible \
    --shell /bin/bash --ingroup $GROUP_NAME $USER_NAME  && \
    chown ansible:ansible -R /home/ansible

WORKDIR /app

RUN apk add --no-cache openssh sudo && \
    mkdir -p /home/ansible/.ssh && chmod 0700 /home/ansible/.ssh && \
    echo "$SSH_PRIV_KEY" > /home/ansible/.ssh/id_ed25519 && \
    echo "$SSH_PUB_KEY" > /home/ansible/.ssh/id_ed25519.pub && \
    chmod 600 /home/ansible/.ssh/id_ed25519 && \
    chmod 600 /home/ansible/.ssh/id_ed25519.pub && \
    echo "Host *\n\tStrictHostKeyChecking no\n" >> /home/ansible/.ssh/config && \
    chown ansible:ansible -R /home/ansible && \
    mkdir -p /root/.ssh && chmod 0700 /root/.ssh && \
    echo "$SSH_PUB_KEY" > /root/.ssh/authorized_keys && \
    sed -i 's/#PermitRootLogin yes/PermitRootLogin yes/' /etc/ssh/sshd_config


RUN mkdir /var/run/sshd && ssh-keygen -A
# New added for disable sudo password
RUN addgroup -g $GROUP_ID_SUDO $GROUP_NAME_SUDO 
RUN addgroup -S $USER_NAME $GROUP_NAME_SUDO
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
COPY ./docker-entrypoint.sh /app/docker-entrypoint.sh
COPY ./operator/ansible /app/ansible
COPY ./operator/app /app
RUN chown ansible:ansible /app/docker-entrypoint.sh && chmod 770 /app/docker-entrypoint.sh 
RUN mkdir /app/logs && \ 
    chown ansible:ansible -R /app && chown ansible:ansible -R /usr/local/go
# RUN mkdir /app/upload  && mkdir /app/logs && \ 
#     chown ansible:ansible -R /app && chown ansible:ansible -R /usr/local/go

USER ansible

RUN go install && cd /app/ansible/ && ansible-galaxy collection install -r requirements.yml

ENTRYPOINT ["./docker-entrypoint.sh"]
CMD [ "go", "run", "main.go"]