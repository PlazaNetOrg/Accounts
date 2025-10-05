from dotenv import load_dotenv
import os

load_dotenv()

class Config:
    APP_NAME = os.getenv("APP_NAME", "PlazaNet")
    SECRET_KEY = os.getenv("SECRET_KEY", "dev-secret")
    JWT_SECRET_KEY = os.getenv("JWT_SECRET_KEY", "jwt-dev-secret")
    SQLALCHEMY_DATABASE_URI = os.getenv("DATABASE_URL", "sqlite:///plazanet.db")
    SQLALCHEMY_TRACK_MODIFICATIONS = False
    PAL_PARTS_DIR = "static/assets/pal_parts"
    JWT_TOKEN_LOCATION = os.getenv("JWT_TOKEN_LOCATION", "cookies").split(",")
    JWT_ACCESS_COOKIE_NAME = os.getenv("JWT_ACCESS_COOKIE_NAME", "access_token")
    JWT_REFRESH_COOKIE_NAME = os.getenv("JWT_REFRESH_COOKIE_NAME", "refresh_token")
    JWT_COOKIE_CSRF_PROTECT = os.getenv("JWT_COOKIE_CSRF_PROTECT", "False") == "True"
    JWT_ACCESS_TOKEN_EXPIRES = int(os.getenv("JWT_ACCESS_TOKEN_EXPIRES", "604800"))  # 7 days

    # Settings
    ## At least one needs to be enabled:
    PLAZANET_ENABLE = os.getenv("PLAZANET_ENABLE", "true")
    PLAZANET_DOMAIN = os.getenv("PLAZANET_DOMAIN", "app.plazanet.org")
    ENABLE_GAMEPLAZA = os.getenv("GAMEPLAZA_ENABLE", "true")

