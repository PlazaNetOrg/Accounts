from dotenv import load_dotenv
import os

load_dotenv()

class Config:
    APP_NAME = os.getenv("APP_NAME", "PlazaNet")
    PLAZANET_ALLOWED_ORIGINS = os.getenv("PLAZANET_ALLOWED_ORIGINS", "")
    SECRET_KEY = os.getenv("SECRET_KEY", "dev-secret")
    JWT_SECRET_KEY = os.getenv("JWT_SECRET_KEY", "jwt-dev-secret")
    SQLALCHEMY_DATABASE_URI = os.getenv("DATABASE_URL", "sqlite:///plazanet.db")
    SQLALCHEMY_TRACK_MODIFICATIONS = False
    PAL_PARTS_DIR = "static/assets/pal_parts"
    # Allow tokens to be read from cookies (browser) and headers (API / game clients)
    _jwt_loc = os.getenv("JWT_TOKEN_LOCATION", "cookies,headers")
    JWT_TOKEN_LOCATION = [x.strip() for x in _jwt_loc.split(",") if x.strip()]
    JWT_ACCESS_COOKIE_NAME = os.getenv("JWT_ACCESS_COOKIE_NAME", "access_token")
    JWT_REFRESH_COOKIE_NAME = os.getenv("JWT_REFRESH_COOKIE_NAME", "refresh_token")
    JWT_COOKIE_CSRF_PROTECT = os.getenv("JWT_COOKIE_CSRF_PROTECT", "False").lower()
    JWT_ACCESS_TOKEN_EXPIRES = int(os.getenv("JWT_ACCESS_TOKEN_EXPIRES", "604800"))  # 7 days
    JWT_COOKIE_DOMAIN = os.getenv("JWT_COOKIE_DOMAIN", None)
    JWT_COOKIE_SECURE = os.getenv("JWT_COOKIE_SECURE", "False").lower()
    JWT_COOKIE_SAMESITE = os.getenv("JWT_COOKIE_SAMESITE", "Lax")

    # Settings
    PLAZANET_ENABLE = os.getenv("PLAZANET_ENABLE", "true")
    PLAZANET_DOMAIN = os.getenv("PLAZANET_DOMAIN", "app.plazanet.org")

