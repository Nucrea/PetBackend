from locust import FastHttpUser, task

from api import BackendApi, Auth

class DummyGet(FastHttpUser):
    api: BackendApi
    auth: Auth

    @task
    def dummy_test(self):
        self.api.dummy_get(self.auth)

    def on_start(self):
        self.api = BackendApi(self)
        
        user = self.api.user_create()
        self.auth = self.api.user_login(user)