FROM alpine:3.17
LABEL description="Ansible on Alpine 3.17 and Python 3.10"
ARG ANSIBLE_VERSION=7.3.0

COPY --from=golang:1.13-alpine /usr/local/go/ /usr/local/go/
COPY ./operator/ansible /app/ansible
COPY ./operator/app /app
ENV PATH="/usr/local/go/bin:${PATH}"

ENV USER_ID=2001
ENV GROUP_ID=2001
ENV USER_NAME=ansible
ENV GROUP_NAME=ansible

RUN /bin/sh -c set -xe && \
    apk update --no-cache && \
    apk upgrade --no-cache && \
    apk add --no-cache gcc make python3 python3-dev openssl-dev \
    py3-cffi py3-bcrypt py-cryptography py3-pynacl py3-pip bash curl && \
    pip3 install --no-cache-dir pip && \
    pip3 install --no-cache-dir ansible==${ANSIBLE_VERSION} && \
    addgroup -g $GROUP_ID $GROUP_NAME && \
    adduser --uid $USER_ID --disabled-password --home /home/ansible \
    --shell /bin/bash --ingroup $GROUP_NAME $USER_NAME  && \
    chown ansible:ansible -R /home/ansible && chown ansible:ansible -R /app && chown ansible:ansible -R /usr/local/go && \
    mkdir -p /home/ansible/.ssh && mkdir /app/upload  && chown ansible:ansible -R /app/upload && \
    echo "Host *\n\tStrictHostKeyChecking no\n" >> /home/ansible/.ssh/config

WORKDIR /app

USER ansible
RUN go install && cd /app/ansible/ && ansible-galaxy collection install -r requirements.yml
ENTRYPOINT [ "go", "run", "main.go" ]