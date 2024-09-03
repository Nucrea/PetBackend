from locust import FastHttpUser, task

from api import BackendApi

class ShortlinkCreate(FastHttpUser):
    api: BackendApi

    @task
    def user_create_test(self):
        self.api.shortlink_create("https://ya.ru")

    def on_start(self):
        self.api = BackendApi(self)