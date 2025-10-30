from flask import Blueprint, request, jsonify, current_app, make_response
from passlib.hash import pbkdf2_sha256
from .models import db, User
from flask_jwt_extended import (
    create_access_token,
    create_refresh_token,
    set_access_cookies,
    set_refresh_cookies,
    jwt_required,
    get_jwt_identity,
    get_jwt,
    decode_token,
    verify_jwt_in_request
)
from datetime import timedelta

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
    data = request.get_json() or {}
    username = data.get('username')
    password = data.get('password')
    audience = data.get('aud', 'plazanet')
    # Response Mode: 'cookie' (default) - sets httpOnly cookies for browser
    #                'json' - returns tokens in JSON payload for external clients (games, other servers)
    response_mode = data.get('response', data.get('response_type', 'cookie'))

    if not username or not password:
        return jsonify({'msg': 'username and password required'}), 400

    user = User.query.filter_by(username=username).first()
    if not user:
        return jsonify({'msg': 'username not found'}), 404
    if not pbkdf2_sha256.verify(password, user.password_hash):
        return jsonify({'msg': 'wrong password'}), 401

    additional_claims = {"aud": audience}
    access_token = create_access_token(
        identity=user.username,
        additional_claims=additional_claims,
        expires_delta=timedelta(seconds=current_app.config['JWT_ACCESS_TOKEN_EXPIRES'])
    )
    refresh_token = create_refresh_token(identity=user.username, additional_claims=additional_claims)

    # Browser flow: Set cookies and return user info
    if response_mode == 'cookie':
        resp = make_response(jsonify({'msg': 'login successful', 'user': user.to_dict()}))
        cookie_domain = current_app.config.get('JWT_COOKIE_DOMAIN')
        cookie_secure = current_app.config.get('JWT_COOKIE_SECURE')
        cookie_samesite = current_app.config.get('JWT_COOKIE_SAMESITE')
        set_access_cookies(resp, access_token)
        resp.set_cookie(
            current_app.config['JWT_ACCESS_COOKIE_NAME'],
            access_token,
            domain=cookie_domain,
            secure=cookie_secure,
            httponly=True,
            samesite=cookie_samesite
        )
        set_refresh_cookies(resp, refresh_token)
        resp.set_cookie(
            current_app.config['JWT_REFRESH_COOKIE_NAME'],
            refresh_token,
            domain=cookie_domain,
            secure=cookie_secure,
            httponly=True,
            samesite=cookie_samesite
        )
        return resp, 200

    # JSON flow (API clients / GamePlaza): Return tokens in JSON for storing client-side
    return jsonify({
        'msg': 'login successful',
        'user': user.to_dict(),
        'access_token': access_token,
        'refresh_token': refresh_token,
        'aud': audience
    }), 200

@bp.route('/refresh', methods=['POST'])
def refresh():
    data = request.get_json(silent=True) or {}
    refresh_token = data.get('refresh_token')
    try:
        if refresh_token:
            decoded = decode_token(refresh_token)
            identity = decoded.get('sub') or decoded.get('identity') or decoded.get('user')
            tok_type = decoded.get('type') or decoded.get('token_type')
            if tok_type and tok_type != 'refresh':
                return jsonify({'msg': 'provided token is not a refresh token'}), 401
            audience = decoded.get('aud', 'plazanet')
            additional_claims = {"aud": audience}
            new_token = create_access_token(identity=identity, additional_claims=additional_claims, expires_delta=timedelta(seconds=current_app.config['JWT_ACCESS_TOKEN_EXPIRES']))
            return jsonify({'access_token': new_token}), 200

        try:
            verify_jwt_in_request(refresh=True)
        except Exception as e:
            return jsonify({'msg': 'refresh token missing or invalid', 'error': str(e)}), 401

        identity = get_jwt_identity()
        claims = get_jwt()
        audience = claims.get('aud', 'plazanet')
        additional_claims = {"aud": audience}
        new_token = create_access_token(identity=identity, additional_claims=additional_claims, expires_delta=timedelta(seconds=current_app.config['JWT_ACCESS_TOKEN_EXPIRES']))
        return jsonify({'access_token': new_token}), 200
    except Exception as e:
        current_app.logger.exception('Failed to refresh token')
        return jsonify({'msg': 'failed to refresh', 'error': str(e)}), 500


@bp.route('/verify', methods=['GET'])
def verify():
    try:
        verify_jwt_in_request(optional=True)
        claims = get_jwt()
        identity = get_jwt_identity()
        if not identity:
            return jsonify({'valid': False}), 401
        return jsonify({'valid': True, 'identity': identity, 'claims': claims}), 200
    except Exception as e:
        return jsonify({'valid': False, 'error': str(e)}), 401