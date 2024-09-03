import random
import string
from locust import HttpUser, FastHttpUser

class Auth():
    token: string

    def __init__(self, token):
        self.token = token

class User():
    email: string
    name: string
    password: string

    def __init__(self, email, password, name, token = ""):
        self.email = email
        self.password = password
        self.name = name
        self.token = token


class BackendApi():
    http: FastHttpUser

    def __init__(self, http: FastHttpUser):
        self.http = http

    def user_create(self) -> User:
        email = ''.join(random.choices(string.ascii_lowercase + string.digits, k=10)) + '@test.test'
        name = ''.join(random.choices(string.ascii_letters, k=10))
        password = 'Abcdef1!!1'
    
        response = self.http.client.post(
            "/user/create",
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
        response = self.http.client.post(
            "/user/login", 
            json={
                "email": user.email,
                "password": user.password,
            },
        )
        if response.status_code != 200:
            raise AssertionError('can not login user')
        
        token = response.json()['token']
        if token == '':
            raise AssertionError('empty user token')
        
        return Auth(token)        

    def dummy_get(self, auth: Auth):
        headers = {"X-Auth": auth.token}
        response = self.http.client.get("/dummy", headers=headers)
        if response.status_code != 200:
            raise AssertionError('something wrong')
        
    def health_get(self):
        response = self.http.client.get("/health")
        if response.status_code != 200:
            raise AssertionError('something wrong')
        
    def shortlink_create(self, url: string) -> string:
        response = self.http.client.post("/s/new?url=" + url)
        if response.status_code != 200:
            raise AssertionError('can not login user')
        
        link = response.json()['link']
        if link == '':
            raise AssertionError('empty user token')
        
        return link