FROM go1.18

COPY ./etc /opt/mall/user/rpc/etc/
COPY userrpc /opt/mall/user/rpc/bin/

WORKDIR /opt/mall/user/rpc

USER normal

EXPOSE 9000

ENTRYPOINT ["bin/userrpc", "-f","./etc/user.yaml"]