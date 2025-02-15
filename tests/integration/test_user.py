import pytest
from api import BackendApi, Requests, User

backendUrl = "http://localhost:8080"

class TestUser:
    def test_create_user(self):
        api = BackendApi(Requests(backendUrl))

        user = User("user@example.com", "aaaaaA1!", "SomeName")
        userWithBadEmail = User("example.com", "aaaaaA1!", "SomeName")
        userWithBadPassword = User("user@example.com", "badPassword", "SomeName")
        userWithBadName = User("user@example.com", "aaaaaA1!", "")

        with pytest.raises(Exception) as e:
            api.user_create(userWithBadEmail)
            raise e
        
        with pytest.raises(Exception) as e:
            api.user_create(userWithBadPassword)
            raise e
        
        with pytest.raises(Exception) as e:
            api.user_create(userWithBadName)
            raise e
        
        api.user_create(user)
        api.user_login(user)

    def test_login_user(self):
        api = BackendApi(Requests(backendUrl))
        user = api.user_create()
        api.user_login(user)