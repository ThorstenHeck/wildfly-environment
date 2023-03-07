FROM postgres:latest

RUN chown postgres.postgres -R /var/lib/postgresql

ARG SSH_PUB_KEY=_OVERRIDE_--build-arg

USER root

RUN apt-get update && apt-get install -y openssh-server openssh-client sudo

RUN mkdir -p /root/.ssh && chmod 0700 /root/.ssh && \
    echo "$SSH_PUB_KEY" > /root/.ssh/authorized_keys && \
    sed -i 's/#PermitRootLogin yes/PermitRootLogin yes/' /etc/ssh/sshd_config && \
    # SSH login fix. Otherwise user is kicked off after login
    sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd

RUN mkdir /var/run/sshd && ssh-keygen -A
# New added for disable sudo password
RUN usermod -aG sudo postgres
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
COPY ./docker-entrypoint.sh /etc/docker-entrypoint.sh
RUN chown postgres:postgres /etc/docker-entrypoint.sh
RUN chmod 770 /etc/docker-entrypoint.sh
USER postgres

ENTRYPOINT ["/etc/docker-entrypoint.sh"]

CMD ["/usr/local/bin/docker-entrypoint.sh", "postgres"]
