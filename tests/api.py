import random
import string
import requests

class Requests():
    def __init__(self, baseUrl):
        self.baseUrl = baseUrl
    
    def post(self, path, json = {}):
        return requests.post(self.baseUrl + path, json=json)

class Auth():
    token: string
    def __init__(self, token):
        self.token = token

class User():
    id: string
    email: string
    name: string
    password: string

    def __init__(self, email, password, name, id=""):
        self.id = id
        self.email = email
        self.password = password
        self.name = name

    @classmethod
    def random(cls):
        email = ''.join(random.choices(string.ascii_lowercase + string.digits, k=10)) + '@example.com'
        name = ''.join(random.choices(string.ascii_letters, k=10))
        password = 'Abcdef1!!1'
        return cls(email, password, name)

def rand_email():
    return ''.join(random.choices(string.ascii_lowercase + string.digits, k=10)) + '@example.com'

class BackendApi():
    def __init__(self, httpClient):
        self.httpClient = httpClient

    def parse_response(self, response):
        if response.status_code != 200:
            raise AssertionError('Request error')
            
        json = response.json()
        if json['status'] == 'success':
            if 'result' in json:
                return json['result']
            return None
        
        error = json['error']
        raise AssertionError(error['id'], error['message'])

    def user_create(self, user: User) -> User:
        res = self.parse_response(
            self.httpClient.post(
                "/api/v1/user/create", json={
                    "email": user.email,
                    "password": user.password,
                    "name": user.name,
                }
            )
        )
        return User(res['email'], user.password, res['name'], id=res['id'])
        
    def user_login(self, email, password) -> Auth:        
        res = self.parse_response(
            self.httpClient.post(
                "/api/v1/user/login", json={
                    "email": email,
                    "password": password,
                },
            )
        )
        return Auth(res['status'])        

    def dummy_get(self, auth: Auth):
        headers = {"X-Auth": auth.token}
        response = self.httpClient.get("/v1/dummy", headers=headers)
        if response.status_code != 200:
            raise AssertionError('something wrong')
        
    def health_get(self):
        response = self.httpClient.get("/health")
        if response.status_code != 200:
            raise AssertionError('something wrong')
