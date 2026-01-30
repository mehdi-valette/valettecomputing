FROM docker.io/library/golang:1.25.6-trixie

RUN apt-get update \
  && apt-get install -y openssh-server sudo inotify-tools \
  && echo "PermitEmptyPasswords yes" >> /etc/ssh/sshd_config

RUN ssh-keygen -A

RUN adduser vscode -uid 1000 \
  && passwd -d vscode \
  && usermod vscode -aG sudo

WORKDIR /home/vscode/app

USER vscode

CMD sudo /usr/sbin/sshd && sleep infinity
