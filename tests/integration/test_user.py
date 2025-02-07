from api import BackendApi, Requests
import requests

backendUrl = "http://localhost:8080"

class TestUser:
    def test_create_user(self):
        api = BackendApi(Requests(backendUrl))
        api.user_create()

    def test_login_user(self):
        api = BackendApi(Requests(backendUrl))
        user = api.user_create()
        api.user_login(user)