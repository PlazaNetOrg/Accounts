from flask import Blueprint, jsonify, request, current_app, send_file, session
from flask_jwt_extended import jwt_required, get_jwt_identity
from .models import User, db, FriendRequest
from .utils.pal import compose_pal
import io
import json

bp = Blueprint('api', __name__, url_prefix='/api')

@bp.route('/me', methods=['GET'])
@jwt_required()
def me():
    username = get_jwt_identity()
    user = User.query.filter_by(username=username).first_or_404()
    return jsonify(user.to_dict(include_pal=True))

@bp.route('/user/<username>', methods=['GET'])
def get_user(username):
    user = User.query.filter_by(username=username).first_or_404()
    return jsonify(user.to_dict(include_pal=True))

@bp.route('/user/<username>/friends', methods=['GET'])
def get_friends(username):
    user = User.query.filter_by(username=username).first_or_404()
    friends = [{'username': f.username, 'status': f.status} for f in user.friends]
    return jsonify({'friends': friends})

@bp.route('/me/friends/request', methods=['POST'])
@jwt_required()
def send_friend_request():
    me_username = get_jwt_identity()
    data = request.get_json() or {}
    to_username = data.get('to')
    me = User.query.filter_by(username=me_username).first_or_404()
    to_user = User.query.filter_by(username=to_username).first()
    if not to_user:
        return jsonify({'msg': 'user not found'}), 404
    fr = FriendRequest(from_user_id=me.id, to_user_id=to_user.id)
    db.session.add(fr)
    db.session.commit()
    return jsonify({'msg': 'friend request sent'})

@bp.route('/me/friends/requests', methods=['GET'])
@jwt_required()
def list_friend_requests():
    me_username = get_jwt_identity()
    me = User.query.filter_by(username=me_username).first_or_404()
    incoming = FriendRequest.query.filter_by(to_user_id=me.id, status='pending').all()
    outgoing = FriendRequest.query.filter_by(from_user_id=me.id, status='pending').all()
    return jsonify({
        'incoming': [{'id': r.id, 'from': User.query.get(r.from_user_id).username} for r in incoming],
        'outgoing': [{'id': r.id, 'to': User.query.get(r.to_user_id).username} for r in outgoing]
    })

@bp.route('/me/friends/requests/<int:req_id>/accept', methods=['POST'])
@jwt_required()
def accept_request(req_id):
    me_username = get_jwt_identity()
    me = User.query.filter_by(username=me_username).first_or_404()
    fr = FriendRequest.query.get_or_404(req_id)
    if fr.to_user_id != me.id:
        return jsonify({'msg': 'not authorized'}), 403
    fr.status = 'accepted'
    user_from = User.query.get(fr.from_user_id)
    me.friends.append(user_from)
    db.session.commit()
    return jsonify({'msg': 'friend request accepted'})

@bp.route('/me/friends/requests/<int:req_id>/decline', methods=['POST'])
@jwt_required()
def decline_request(req_id):
    me_username = get_jwt_identity()
    me = User.query.filter_by(username=me_username).first_or_404()
    fr = FriendRequest.query.get_or_404(req_id)
    if fr.to_user_id != me.id:
        return jsonify({'msg': 'not authorized'}), 403
    fr.status = 'declined'
    db.session.commit()
    return jsonify({'msg': 'friend request declined'})

@bp.route('/me/status', methods=['POST'])
@jwt_required()
def set_status():
    me_username = get_jwt_identity()
    data = request.get_json() or {}
    status = data.get('status', '')[:200]
    me = User.query.filter_by(username=me_username).first_or_404()
    me.status = status
    db.session.commit()
    return jsonify({'msg': 'status updated'})

@bp.route("/save_pal", methods=["POST"])
@jwt_required()
def save_pal():
    data = request.get_json(silent=True) or {}
    pal = data.get("pal")
    username = get_jwt_identity()
    user = User.query.filter_by(username=username).first()
    if not user:
        return jsonify({"success": False, "error": "User not found"}), 404
    import json as _json
    try:
        user.pal_json = _json.dumps(pal)
        db.session.commit()
    except Exception as e:
        current_app.logger.exception("Failed saving pal for %s", username)
        return jsonify({"success": False, "error": "internal error"}), 500
    return jsonify({"success": True})

@bp.route('/pal/<username>.svg')
def pal_image(username):
    user = User.query.filter_by(username=username).first_or_404()
    parts_dir = current_app.config.get('PAL_PARTS_DIR')
    svg = compose_pal(user.pal_json, parts_dir)
    return current_app.response_class(svg, mimetype='image/svg+xml')

@bp.route('/pal_preview', methods=['POST'])
def pal_preview():
    data = request.get_json() or {}
    pal = data.get('pal')
    parts_dir = current_app.config.get('PAL_PARTS_DIR')
    from .utils.pal import compose_pal
    svg = compose_pal(pal, parts_dir)
    return current_app.response_class(svg, mimetype='image/svg+xml')