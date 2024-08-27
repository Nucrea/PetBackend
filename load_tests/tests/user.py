from locust import FastHttpUser, task

from api import BackendApi, User

class UserCreate(FastHttpUser):
    api: BackendApi

    @task
    def user_create_test(self):
        self.api.user_create()

    def on_start(self):
        self.api = BackendApi(self)

class UserLogin(FastHttpUser):
    api: BackendApi
    user: User

    @task
    def user_create_test(self):
        self.api.user_login(self.user)

    def on_start(self):
       self.api = BackendApi(self)
       self.user = self.api.user_create()