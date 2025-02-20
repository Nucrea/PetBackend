import random
import string
import pytest
from kafka import KafkaConsumer
from api import BackendApi, Requests, User, rand_email

backendUrl = "http://localhost:8080"

def test_create_user():
    backend = BackendApi(Requests(backendUrl))

    user = User.random()
    userWithBadEmail = User("sdfsaadsfgdf", user.password, user.name)
    userWithBadPassword = User(user.email, "badPassword", user.name)

    with pytest.raises(Exception):
        backend.user_create(userWithBadEmail)
    with pytest.raises(Exception):
        backend.user_create(userWithBadPassword)
    
    resultUser = backend.user_create(user)

    #should not create user with same email
    with pytest.raises(Exception):
        backend.user_create(user)

    assert resultUser.email == user.email
    assert resultUser.id != ""

def test_login_user():
    backend = BackendApi(Requests(backendUrl))

    # consumer = KafkaConsumer(
    #     'backend_events',
    #     group_id='test-group',
    #     bootstrap_servers=['localhost:9092'],
    #     consumer_timeout_ms=1000)
    # consumer.seek_to_end()

    user = backend.user_create(User.random())

    with pytest.raises(Exception):
        backend.user_login(user.email, "badpassword")
    with pytest.raises(Exception):
        backend.user_login(rand_email(), user.password)
    
    #should not login without verified email 
    with pytest.raises(Exception):
        backend.user_login(user.email, user.password)
    
    # msgs = consumer.poll(timeout_ms=100)
    # print(msgs)

