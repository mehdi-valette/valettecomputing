FROM golang:tip-trixie

RUN apt-get update \
  && apt-get install -y openssh-server sudo inotify-tools \
  && echo "PermitEmptyPasswords yes" >> /etc/ssh/sshd_config

RUN adduser vscode \
  && passwd -d vscode \
  && usermod vscode -aG sudo

WORKDIR /usr/src/app

RUN ssh-keygen -A

CMD /usr/sbin/sshd && sleep infinity