FROM go1.18

USER root

COPY ./etc /etc/
COPY productrpc /usr/bin

WORKDIR /usr/bin/

EXPOSE 9001

ENTRYPOINT ["/usr/bin/productrpc", "-f","/etc/product.yaml"]