#!/bin/sh
export THREEBOT_APP_ID=172.17.0.2
# export DOCUMENTSERVER_URL=https://office.jimber.org/
export DOCUMENTSERVER_URL=http://172.17.0.3
sed -i 's@REPLACE_APPID@'"$THREEBOT_APP_ID"'@' /3bot_auth/config.py
sed -i 's@REPLACE_DOCUMENTSERVERURL@'"$DOCUMENTSERVER_URL"'@' /filebrowser_frontend/dist/index.html

service nginx start
exec /usr/bin/supervisord