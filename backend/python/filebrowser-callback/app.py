import json
import urllib.request 
from flask import Flask, request
from config import PORT
import requests

app = Flask(__name__)

@app.route('/callback',  methods=['GET', 'POST'])
def documentcallback():
    args = request.args
    if request.data:
        reqdata = json.loads(request.data)
        print(reqdata)
        if (reqdata["status"] == 2):
            saveFile(reqdata["url"], args["auth"], args["filepath"])
        
    data = { 'error': 0 }
    return json.dumps(data)

def saveFile(newFile, auth, filepath):
    try:
        headers = {
            'X-Auth': auth,
        }

        params = (
            ('override', 'true'),
        )

        file= urllib.request.urlopen(newFile).read()

        response = requests.post(f'http://127.0.0.1:81/api/resources{filepath}', headers=headers, params=params, data=file)
    except Exception as e:
        print(e)


if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=PORT)