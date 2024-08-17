import random
import string
from locust import HttpUser, task

class DummyRoute(HttpUser):
    @task
    def dummy_test(self):
        self.client.get("/dummy")

    # @task
    # def user_create_test(self):
    #     randEmail = ''.join(random.choices(string.ascii_lowercase + string.digits, k=10)) + '@test.test'
    #     randPassword = ''.join(random.choices(string.ascii_letters + string.digits, k=10))
    #     randName = ''.join(random.choices(string.ascii_letters, k=10))
    #     self.client.post(
    #         "/user/create",
    #         json={
    #             "email":randEmail,
    #             "password":randPassword,
    #             "name": randName,
    #         },
    #     )

    def on_start(self):
        randEmail = ''.join(random.choices(string.ascii_lowercase + string.digits, k=10)) + '@test.test'
        randPassword = ''.join(random.choices(string.ascii_letters + string.digits, k=10))
        randName = ''.join(random.choices(string.ascii_letters, k=10))

        response = self.client.post(
            "/user/create",
            json={
                "email":randEmail,
                "password":randPassword,
                "name": randName,
            },
        )
        if response.status_code != 200:
            raise AssertionError('can not create user')

        response = self.client.post(
            "/user/login", 
            json={
                "email":randEmail, 
                "password":randPassword,
            },
        )
        if response.status_code != 200:
            raise AssertionError('can not login user')
        
        token = response.json()['token']
        if token == '':
            raise AssertionError('empty user token')

        self.client.headers = {"X-Auth": token}