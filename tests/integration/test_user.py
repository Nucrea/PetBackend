from api import BackendApi
import requests

class TestUser:
    def test_create_user(self):
        api = BackendApi(requests)
        api.user_create()

    def test_login_user(self):
        api = BackendApi(requests)
        user = api.user_create()
        api.user_login(user)