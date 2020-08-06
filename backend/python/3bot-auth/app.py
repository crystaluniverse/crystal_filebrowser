from flask import Flask, request
import threebotlogin
from config import config
from filebrowser import FilebrowserAuthenticator

app = Flask(__name__)

if __name__ == '__main__':
    threebotlogin.configure(app, config['threebot-appid'], config['threebot-privatekey'])
    app.run(debug=True, host='0.0.0.0', port=config['port'])