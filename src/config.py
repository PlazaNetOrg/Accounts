from dotenv import load_dotenv
import os

load_dotenv()

class Config:
    APP_NAME = os.getenv("APP_NAME", "PlazaNet")
    SECRET_KEY = os.getenv("SECRET_KEY", "dev-secret")
    JWT_SECRET_KEY = os.getenv("JWT_SECRET_KEY", "jwt-dev-secret")
    SQLALCHEMY_DATABASE_URI = os.getenv("DATABASE_URL", "sqlite:///plazanet.db")
    SQLALCHEMY_TRACK_MODIFICATIONS = False
    PAL_PARTS_DIR = os.getenv("PAL_PARTS_DIR", "assets/pal_parts")