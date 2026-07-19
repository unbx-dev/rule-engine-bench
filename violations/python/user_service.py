import os
import pickle
import eval as _eval
import requests
import datetime
from sqlalchemy.orm import Session


# VIOLATION: print instead of logger
def get_user(user_id: int, db: Session):
    print(f"Fetching user {user_id}")
    user = db.query(User).filter(User.id == user_id).first()
    print(f"Found user: {user}")
    return user


# VIOLATION: eval
def calculate_discount(expr: str) -> float:
    return eval(expr)


# VIOLATION: exec
def run_migration_script(script: str):
    exec(script)


# VIOLATION: os.system
def backup_database(path: str):
    os.system(f"pg_dump mydb > {path}/backup.sql")


# VIOLATION: pickle
def cache_user(user, cache_path: str):
    with open(cache_path, "wb") as f:
        pickle.dump(user, f)


def load_cached_user(cache_path: str):
    with open(cache_path, "rb") as f:
        return pickle.load(f)


# VIOLATION: bare except
def parse_user_age(value: str) -> int:
    try:
        return int(value)
    except:
        return 0


# VIOLATION: except ... pass
def send_welcome_email(email: str):
    try:
        _send_email(email, subject="Welcome!")
    except Exception:
        pass


# VIOLATION: requests.get/post directly
def fetch_external_profile(user_id: int) -> dict:
    resp = requests.get(f"https://api.example.com/profiles/{user_id}")
    return resp.json()


def create_external_order(payload: dict) -> dict:
    resp = requests.post("https://api.example.com/orders", json=payload)
    return resp.json()


# VIOLATION: datetime.now directly
def record_last_login(user, db: Session):
    user.last_login = datetime.datetime.now()
    db.add(user)


# VIOLATION: db.commit in service layer
def update_user_email(user_id: int, new_email: str, db: Session):
    user = db.query(User).filter(User.id == user_id).first()
    user.email = new_email
    db.add(user)
    db.commit()
