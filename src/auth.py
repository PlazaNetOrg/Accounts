from flask import Blueprint, request, jsonify, current_app
from passlib.hash import pbkdf2_sha256
from .models import db, User
from flask_jwt_extended import create_access_token, create_refresh_token, set_access_cookies, set_refresh_cookies
from datetime import timedelta
from flask_jwt_extended import jwt_required, get_jwt_identity

bp = Blueprint('auth', __name__, url_prefix='/auth')

@bp.route('/register', methods=['POST'])
def register():
    data = request.get_json() or {}
    username = data.get('username')
    password = data.get('password')
    birthday = data.get('birthday')
    pal = data.get('pal')

    if not username or not password:
        return jsonify({'msg': 'username and password required'}), 400

    if User.query.filter_by(username=username).first():
        return jsonify({'msg': 'username taken'}), 409

    pw_hash = pbkdf2_sha256.hash(password[:72])
    user = User(username=username, password_hash=pw_hash, pal_json=None)
    if birthday:
        try:
            from datetime import datetime
            user.birthday = datetime.fromisoformat(birthday).date()
        except Exception:
            pass
    if pal:
        import json
        user.pal_json = json.dumps(pal)

    db.session.add(user)
    db.session.commit()

    return jsonify({'msg': 'registered', 'user': user.to_dict()}), 201

@bp.route('/login', methods=['POST'])
def login():
    from flask import make_response
    data = request.get_json() or {}
    username = data.get('username')
    password = data.get('password')
    audience = data.get('aud', 'plazanet')

    if not username or not password:
        return jsonify({'msg': 'username and password required'}), 400

    user = User.query.filter_by(username=username).first()
    if not user:
        return jsonify({'msg': 'username not found'}), 404
    if not pbkdf2_sha256.verify(password, user.password_hash):
        return jsonify({'msg': 'wrong password'}), 401

    additional_claims = {"aud": audience}
    access_token = create_access_token(identity=user.username, additional_claims=additional_claims, expires_delta=timedelta(seconds=current_app.config['JWT_ACCESS_TOKEN_EXPIRES']))
    refresh_token = create_refresh_token(identity=user.username, additional_claims=additional_claims)
    resp = make_response(jsonify({'msg': 'login successful', 'user': user.to_dict()}))
    set_access_cookies(resp, access_token)
    set_refresh_cookies(resp, refresh_token)
    return resp, 200

@bp.route('/refresh', methods=['POST'])
def refresh():
    from flask_jwt_extended import get_jwt
    identity = get_jwt_identity()
    claims = get_jwt()
    audience = claims.get("aud", "plazanet")
    additional_claims = {"aud": audience}
    new_token = create_access_token(identity=identity, additional_claims=additional_claims, expires_delta=timedelta(seconds=current_app.config['JWT_ACCESS_TOKEN_EXPIRES']))
    return jsonify({'access_token': new_token})