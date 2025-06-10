import uuid
import requests
import pendulum


def create_user(name: str = "") -> dict:
    response = requests.post(
        "http://localhost:8080/signup",
        json={"name": name or random_name()},
    )
    response.raise_for_status()

    return response.json()


def random_name() -> str:
    return str(uuid.uuid4())


def make_auth_header(token):
    return {"Authorization": f"Bearer {token}"}
