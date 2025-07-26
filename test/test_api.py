from call_api import *
from hamcrest import assert_that, contains_inanyorder, equal_to, not_none
import uuid


def test_create_user():
    name = random_name()
    password = random_name()
    token = create_user(name, password)["token"]
    assert_that(token, not_none())

    user = get_user(token)
    assert_that(token, not_none())
