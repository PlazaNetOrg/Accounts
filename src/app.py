import os
from flask import Flask, request, redirect, url_for, render_template, jsonify, session
from flask_sqlalchemy import SQLAlchemy
from flask_jwt_extended import JWTManager, create_access_token, jwt_required, get_jwt_identity
from dotenv import load_dotenv
from passlib.hash import bcrypt
import os
from datetime import timedelta

from .models import db, User, FriendRequest

from .auth import bp as auth_bp
from .routes import bp as api_bp

jwt = JWTManager()

def create_app():
    load_dotenv()

    app = Flask(__name__)
    app.config["SECRET_KEY"] = os.getenv("SECRET_KEY", "changeme")
    app.config["SQLALCHEMY_DATABASE_URI"] = os.getenv("DATABASE_URL", "sqlite:///plazanet.db")
    app.config["SQLALCHEMY_TRACK_MODIFICATIONS"] = False
    app.config["JWT_SECRET_KEY"] = os.getenv("JWT_SECRET_KEY", "jwt-secret")
    app.config["PAL_PARTS_DIR"] = os.getenv("PAL_PARTS_DIR", "assets/pal_parts")
    app.config["JWT_TOKEN_LOCATION"] = ["cookies"]
    app.config["JWT_ACCESS_COOKIE_NAME"] = "access_token"

    db.init_app(app)
    jwt.init_app(app)

    app.register_blueprint(auth_bp)
    app.register_blueprint(api_bp)

    @app.route("/")
    def index():
        token = request.cookies.get("access_token")
        if token:
            return redirect(url_for("me_page"))
        return redirect(url_for("login_page"))

    @app.route("/login", methods=["GET"])
    def login_page():
        return render_template("login.html", error=None)

    @app.route("/register", methods=["GET"])
    def register_page():
        return render_template("register.html", error=None)
    
    def get_pal_parts():
        parts_dir = os.path.join(app.static_folder, "assets", "pal_parts")
        categories = [d for d in os.listdir(parts_dir) if os.path.isdir(os.path.join(parts_dir, d))]
        pal_parts = {}
        for cat in categories:
            cat_dir = os.path.join(parts_dir, cat)
            pal_parts[cat] = [f for f in os.listdir(cat_dir) if not f.startswith('.') and os.path.isfile(os.path.join(cat_dir, f))]
        return pal_parts

    @app.route("/pal_creator", methods=["GET"])
    @jwt_required()
    def pal_creator():
        pal_parts = get_pal_parts()
        return render_template("pal_creator.html", pal_parts=pal_parts)

    @app.route("/me")
    @jwt_required()
    def me_page():
        try:
            username = get_jwt_identity()
            user = User.query.filter_by(username=username).first()
            if not user:
                return redirect(url_for("login_page"))
        except Exception:
            return redirect(url_for("login_page"))
        return render_template("me.html", user=user)

    return app
