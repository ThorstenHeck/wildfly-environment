FROM quay.io/wildfly/wildfly:latest

ARG WILDFLY_ADMIN_PW=adminuser1
ARG SSH_PRIV_KEY=
ARG SSH_PUB_KEY=

RUN whoami && ls -lisa /root
RUN yum –y install openssh-server openssh-clients
RUN /opt/jboss/wildfly/bin/add-user.sh -u 'adminuser1' -p ${WILDFLY_ADMIN_PW} -g 'admin'

RUN echo ${SSH_PRIV_KEY} ${SSH_PUB_KEY} 

RUN mkdir -p /root/.ssh && chmod 0700 /root/.ssh && \
    echo "$SSH_PRIV_KEY" > /root/.ssh/id_ed25519 && \
    echo "$SSH_PUB_KEY" > /root/.ssh/id_ed25519.pub && \
    chmod 600 /root/.ssh/id_ed25519 && \
    chmod 600 /root/.ssh/id_ed25519.pub
