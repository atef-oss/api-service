import uvicorn
from fastapi import FastAPI
from pydantic import BaseModel
from typing import Optional

class User(BaseModel):
    id: int
    name: str
    email: str

app = FastAPI()

@app.post("/users/")
def create_user(user: User):
    return {"message": "User created successfully"}

@app.get("/users/{user_id}")
def read_user(user_id: int):
    return {"user_id": user_id}

@app.put("/users/{user_id}")
def update_user(user_id: int, user: User):
    return {"message": "User updated successfully"}

@app.delete("/users/{user_id}")
def delete_user(user_id: int):
    return {"message": "User deleted successfully"}

if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=8000, reload=True)