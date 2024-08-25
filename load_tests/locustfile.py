import random
import string
from locust import HttpUser, FastHttpUser, task

class DummyRoute(FastHttpUser):
    @task
    def dummy_test(self):
        self.client.get("/dummy", headers=self.headers)

    def on_start(self):
        randEmail = ''.join(random.choices(string.ascii_lowercase + string.digits, k=10)) + '@test.test'
        randName = ''.join(random.choices(string.ascii_letters, k=10))
        password = 'Abcdef1!!1'

        response = self.client.post(
            "/user/create",
            json={
                "email": randEmail,
                "password": password,
                "name": randName,
            },
        )
        if response.status_code != 200:
            raise AssertionError('can not create user')

        response = self.client.post(
            "/user/login", 
            json={
                "email": randEmail,
                "password": password,
            },
        )
        if response.status_code != 200:
            raise AssertionError('can not login user')
        
        token = response.json()['token']
        if token == '':
            raise AssertionError('empty user token')

        self.headers = {"X-Auth": token}