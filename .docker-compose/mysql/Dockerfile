FROM mysql:8.0.36-oracle

COPY ./init.sh /init/init.sh

RUN microdnf update -y
RUN microdnf install -y openssh-clients openssh-server
RUN chmod +x /init/init.sh
