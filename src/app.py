
from flask import Flask, render_template, redirect, url_for, request
from .config import Config
from .models import db, User
from .auth import bp as auth_bp
from .routes import bp as api_bp
from flask_jwt_extended import JWTManager, jwt_required, get_jwt_identity
import os

def get_pal_parts(app):
    parts_dir = os.path.join(app.static_folder, "assets", "pal_parts")
    categories = [d for d in os.listdir(parts_dir) if os.path.isdir(os.path.join(parts_dir, d))]
    pal_parts = {}
    for cat in categories:
        cat_dir = os.path.join(parts_dir, cat)
        pal_parts[cat] = [f for f in os.listdir(cat_dir) if not f.startswith('.') and os.path.isfile(os.path.join(cat_dir, f))]
    return pal_parts

def create_app():
    app = Flask(__name__, static_folder="static")
    app.config.from_object(Config)

    db.init_app(app)
    jwt = JWTManager(app)

    app.register_blueprint(auth_bp)
    app.register_blueprint(api_bp)

    @app.route("/")
    def index():
        return redirect(url_for("me_page")) if request.cookies.get("access_token") else redirect(url_for("login_page"))

    @app.route("/login", methods=["GET"])
    def login_page():
        return render_template("login.html")

    @app.route("/register", methods=["GET"])
    def register_page():
        return render_template("register.html")

    @app.route("/pal_creator", methods=["GET"])
    @jwt_required()
    def pal_creator():
        pal_parts = get_pal_parts(app)
        # Load current pal if exists
        username = get_jwt_identity()
        user = User.query.filter_by(username=username).first()
        pal_json = None
        if user and user.pal_json:
            pal_json = user.pal_json
        return render_template("pal_creator.html", pal_parts=pal_parts, pal_json=pal_json)

    @app.route("/me")
    @jwt_required()
    def me_page():
        username = get_jwt_identity()
        user = User.query.filter_by(username=username).first()
        if not user:
            return redirect(url_for("login_page"))
        return render_template(
            "me.html",
            user=user,
            PLAZANET_ENABLE=app.config.get("PLAZANET_ENABLE"),
            APP_NAME=app.config.get("APP_NAME"),
            PLAZANET_DOMAIN=app.config.get("PLAZANET_DOMAIN")
        )

    @app.errorhandler(401)
    def handle_unauthorized(e):
        from flask import request
        if request.accept_mimetypes.accept_json and not request.accept_mimetypes.accept_html:
            return jsonify({'msg': 'unauthorized'}), 401
        return redirect(url_for("login_page"))

    @app.errorhandler(422)
    def handle_unprocessable_entity(e):
        from flask import request
        if request.accept_mimetypes.accept_json and not request.accept_mimetypes.accept_html:
            return jsonify({'msg': 'token error'}), 422
        return redirect(url_for("login_page"))

    return app
