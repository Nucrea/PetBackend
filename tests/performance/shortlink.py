from locust import FastHttpUser, task

from api import BackendApi

class ShortlinkCreate(FastHttpUser):
    def on_start(self):
        self.api = BackendApi(self.client)

    @task
    def shortlink_create_test(self):
        self.api.shortlink_create("https://example.com")