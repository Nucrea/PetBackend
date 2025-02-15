from locust import FastHttpUser, task

from api import BackendApi

class HealthGet(FastHttpUser):
    def on_start(self):
        self.api = BackendApi(self.client)

    @task
    def user_create_test(self):
        self.api.health_get()