FROM ubuntu:latest

ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get upgrade -y
RUN apt-get install -y python3 python3-pip nodejs npm curl nginx supervisor

RUN npm install -g http-server
RUN rm -f /etc/nginx/sites-enabled/default


WORKDIR /filebrowser_frontend
COPY ./frontend .
RUN npm install
RUN npm run build

WORKDIR /3bot_auth
ADD ./backend/python/3bot-auth .
RUN pip3 install -r ./requirements.txt
EXPOSE 5001

WORKDIR /callback
ADD ./backend/python/filebrowser-callback .
RUN pip3 install -r ./requirements.txt
EXPOSE 5000

WORKDIR /filebrowser
ADD ./griddeployment/grid.json .filebrowser.json
ADD ./backend/filebrowser /filebrowser


ADD ./griddeployment/nginx.conf /etc/nginx/conf.d/filebrowser.conf

RUN mkdir -p /var/log/supervisor
ADD ./griddeployment/supervisor.conf /etc/supervisor/conf.d/supervisord.conf

ADD ./griddeployment/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]