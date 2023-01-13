FROM go1.18

WORKDIR /opt/mall/user/api

COPY ./etc etc/
COPY userapi bin/

EXPOSE 8000

USER normal

ENTRYPOINT ["bin/userapi", "-f","etc/user.yaml"]