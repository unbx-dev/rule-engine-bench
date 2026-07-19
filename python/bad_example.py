import os
import pickle
import requests
from datetime import datetime
from sqlalchemy.orm import Session


# printを使う (loggerを使うべき)
print("アプリケーション起動")


def evaluate_expression(expr: str) -> float:
    # evalを使う
    result = eval(expr)
    print(f"結果: {result}")
    return result


def run_command(cmd: str) -> None:
    # os.systemを使う
    os.system(cmd)


def execute_code(code: str) -> None:
    # execを使う
    exec(code)


def save_model(model: object, path: str) -> None:
    # pickleを使う
    with open(path, "wb") as f:
        pickle.dump(model, f)


def load_model(path: str) -> object:
    # pickleを使う
    with open(path, "rb") as f:
        return pickle.load(f)


def fetch_user(user_id: int) -> dict:
    # requests.getを直接使う
    response = requests.get(f"https://api.example.com/users/{user_id}")
    return response.json()


def post_event(payload: dict) -> dict:
    # requests.postを直接使う
    response = requests.post("https://api.example.com/events", json=payload)
    return response.json()


def get_current_timestamp() -> str:
    # datetime.nowを直接使う
    return datetime.now().isoformat()


def process_data(data: list) -> list:
    results = []
    for item in data:
        try:
            results.append(evaluate_expression(item))
        except:  # bare exceptを使う
            pass  # exceptでpassする
    return results


class UserService:
    def __init__(self, db: Session):
        self.db = db

    def create_user(self, name: str, email: str) -> dict:
        user = {"name": name, "email": email, "created_at": get_current_timestamp()}
        self.db.add(user)
        # service層でDB sessionをcommitする
        self.db.commit()
        print(f"ユーザー作成: {name}")
        return user
