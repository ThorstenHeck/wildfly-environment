FROM container-registry.oracle.com/database/express:21.3.0-xe

ARG SSH_PUB_KEY=_OVERRIDE_--build-arg

USER root

RUN yum -y install openssh-server sudo python3

RUN mkdir -p /root/.ssh && chmod 0700 /root/.ssh && \
    echo "$SSH_PUB_KEY" > /root/.ssh/authorized_keys && \
    sed -i 's/#PermitRootLogin yes/PermitRootLogin yes/' /etc/ssh/sshd_config && \
    # SSH login fix. Otherwise user is kicked off after login
    sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd

RUN mkdir /var/run/sshd && ssh-keygen -A
RUN groupadd -g 110 sudo && usermod -aG sudo oracle
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
COPY ./docker-entrypoint.sh /etc/docker-entrypoint.sh
RUN chown oracle:root /etc/docker-entrypoint.sh
RUN chmod 770 /etc/docker-entrypoint.sh
USER oracle

ENTRYPOINT ["/etc/docker-entrypoint.sh"]

CMD [ "/bin/sh", "-c", "exec $ORACLE_BASE/$RUN_FILE" ]