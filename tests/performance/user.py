from locust import FastHttpUser, task

from api import BackendApi, User

class UserCreate(FastHttpUser):
    def on_start(self):
        self.api = BackendApi(self.client)

    @task
    def user_create_test(self):
        self.api.user_create()

class UserLogin(FastHttpUser):
    def on_start(self):
       self.api = BackendApi(self)
       self.user = self.api.user_create()

    @task
    def user_create_test(self):
        self.api.user_login(self.user)