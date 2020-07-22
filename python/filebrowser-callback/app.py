import json
import urllib.request 
from flask import Flask, request
from config import PORT

app = Flask(__name__)

@app.route('/callback',  methods=['GET', 'POST'])
def documentcallback():
    args = request.args
    if request.data:
        reqdata = json.loads(request.data)
        print(reqdata)
        if (reqdata["status"] == 2):
            saveFile(reqdata["url"], args["filename"])
        
    
    data = { 'error': 0 }

    return json.dumps(data)

def saveFile(url, filename):
    try:
        urllib.request.urlretrieve(url, f'../datafiles/{filename}')
    except Exception as e:
        print(e)


if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=PORT)