FROM go1.18

USER root

COPY ./etc /opt/userrpc/etc/
COPY userrpc /opt/userrpc/bin/
RUN chown -R normal /opt/userrpc

WORKDIR /opt/userrpc

USER normal

EXPOSE 9000

ENTRYPOINT ["bin/userrpc", "-f","./etc/user.yaml"]