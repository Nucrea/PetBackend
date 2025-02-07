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

    def __init__(self, email, password, name, token = ""):
        self.email = email
        self.password = password
        self.name = name
        self.token = token


class BackendApi():
    def __init__(self, httpClient):
        self.httpClient = httpClient

    def user_create(self) -> User:
        email = ''.join(random.choices(string.ascii_lowercase + string.digits, k=10)) + '@test.test'
        name = ''.join(random.choices(string.ascii_letters, k=10))
        password = 'Abcdef1!!1'
    
        response = self.httpClient.post(
            "/v1/user/create",
            json={
                "email": email,
                "password": password,
                "name": name,
            },
        )
        if response.status_code != 200:
            raise AssertionError('can not create user')
        
        return User(email, password, name)
        
    def user_login(self, user: User) -> Auth:        
        response = self.httpClient.post(
            "/v1/user/login", 
            json={
                "email": user.email+"a",
                "password": user.password,
            },
        )

        status = response.json()['status']
        if status == 'error':
            raise AssertionError(response.json()['error']['message'])

        if response.status_code != 200:
            raise AssertionError('can not login user')

        token = response.json()['result']['token']
        if token == '':
            raise AssertionError('empty user token')
        
        return Auth(token)        

    def dummy_get(self, auth: Auth):
        headers = {"X-Auth": auth.token}
        response = self.httpClient.get("/v1/dummy", headers=headers)
        if response.status_code != 200:
            raise AssertionError('something wrong')
        
    def health_get(self):
        response = self.httpClient.get("/health")
        if response.status_code != 200:
            raise AssertionError('something wrong')
