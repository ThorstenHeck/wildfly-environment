FROM quay.io/wildfly/wildfly:latest

ARG WILDFLY_ADMIN_PW=_OVERRIDE_--build-arg
ARG SSH_PRIV_KEY=_OVERRIDE_--build-arg
ARG SSH_PUB_KEY=_OVERRIDE_--build-arg

RUN /opt/jboss/wildfly/bin/add-user.sh -u 'adminuser1' -p ${WILDFLY_ADMIN_PW} -g 'admin'

USER root

RUN yum -y install openssh-server sudo

RUN mkdir -p /root/.ssh && chmod 0700 /root/.ssh && \
    echo "$SSH_PUB_KEY" > /root/.ssh/authorized_keys && \
    sed -i 's/#PermitRootLogin yes/PermitRootLogin yes/' /etc/ssh/sshd_config && \
    # SSH login fix. Otherwise user is kicked off after login
    sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd

RUN mkdir /var/run/sshd && ssh-keygen -A
# New added for disable sudo password
RUN groupadd -g 110 sudo && usermod -aG sudo jboss
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

USER jboss

CMD ["/usr/bin/sudo", "/usr/sbin/sshd","-D"]