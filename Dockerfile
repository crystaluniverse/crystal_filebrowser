FROM debian:buster

RUN apt-get update && apt-get upgrade -y
RUN apt-get install -y python3 python3-pip
RUN pip3 --version

VOLUME /srv

COPY .docker.json /.filebrowser.json
COPY filebrowser /filebrowser

COPY ./python/filebrowser-callback /python
COPY startup.sh /startup.sh

RUN pip3 install -r /python/requirements.txt

RUN chmod +x startup.sh

RUN mkdir /datafiles

EXPOSE 80
EXPOSE 5000

CMD [ "./startup.sh" ]