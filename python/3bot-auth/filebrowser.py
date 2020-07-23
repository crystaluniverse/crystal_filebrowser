import requests
import urllib.parse
import json 
from config import config
from admincredentials import credentials

class FilebrowserAuthenticator:
    def __init__(self):
        self.admin_authkey = self.getAuthenticationToken(credentials['admin_login'],credentials['admin_password'])

    def registerUser(self, userName, password):
        headers = {
            'X-Auth': self.admin_authkey,
            'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8'
        }
        scope = f"./{userName}"
        user = {
            "what": "user",
            "which": [],
            "data": {
                "scope": scope,
                "locale": "en",
                "viewMode": "mosaic",
                "sorting": {
                    "by": "",
                    "asc": False
                },
                "perm": {
                    "admin": False,
                    "execute": True,
                    "create": True,
                    "rename": True,
                    "modify": True,
                    "delete": True,
                    "share": True,
                    "download": True
                },
                "commands": [],
                "username": userName,
                "password": password,
                "rules": [],
                "lockPassword": False,
                "id": 0
            }  
        }
        try:
            response = requests.post(f"{config['filebrowserUrl']}/api/users", headers=headers, json=user)
            print(response)
            if response.status_code == 401:
                self.admin_authkey = self.getAuthenticationToken(credentials['admin_login'],credentials['admin_password'])
        except Exception:
            print('failed to register user')

    def getAuthenticationToken(self, userName, password):
        response = requests.post(f"{config['filebrowserUrl']}/api/login", json={'username': userName, 'password':password})
        if response.status_code != 200:
            return 'failed login request'
        authkey = response.text
        return authkey