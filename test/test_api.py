from call_api import *
from hamcrest import assert_that, contains_inanyorder, equal_to
import uuid


def test_create_user():
    name = random_name()
    token = create_user(name)["token"]
