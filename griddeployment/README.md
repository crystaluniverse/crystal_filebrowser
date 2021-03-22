## Required env vars
THREEBOT_APP_ID -> appid to redirect to = url without protocol and trailing /
DOCUMENTSERVER_URL -> url from the documentserver


docker run -it -d --name filebrowser -e THREEBOT_APP_ID='filebrowser.jimbertesting.be' -e DOCUMENTSERVER_URL='https://office.jimber.org' jimbersoftware/filebrowser_grid:0.4