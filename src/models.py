from flask_sqlalchemy import SQLAlchemy
from datetime import datetime
import json

db = SQLAlchemy()

friends = db.Table(
    'friends',
    db.Column('user_id', db.Integer, db.ForeignKey('user.id'), primary_key=True),
    db.Column('friend_id', db.Integer, db.ForeignKey('user.id'), primary_key=True)
)

class FriendRequest(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    from_user_id = db.Column(db.Integer, db.ForeignKey('user.id'), nullable=False)
    to_user_id = db.Column(db.Integer, db.ForeignKey('user.id'), nullable=False)
    status = db.Column(db.String(20), default='pending')  # pending, accepted, declined
    created_at = db.Column(db.DateTime, default=datetime.utcnow)

class User(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    username = db.Column(db.String(80), unique=True, nullable=False)
    password_hash = db.Column(db.String(256), nullable=False)
    pal_json = db.Column(db.Text, nullable=True)
    status = db.Column(db.String(200), default="")
    created_at = db.Column(db.DateTime, default=datetime.utcnow)

    friends = db.relationship(
        'User',
        secondary=friends,
        primaryjoin=id==friends.c.user_id,
        secondaryjoin=id==friends.c.friend_id,
        backref='friended_by'
    )

    def to_dict(self, include_pal=False):
        d = {
            'id': self.id,
            'username': self.username,
            'status': self.status,
            'created_at': self.created_at.isoformat()
        }
        if include_pal:
            try:
                d['pal'] = json.loads(self.pal_json) if self.pal_json else None
            except Exception:
                d['pal'] = None
        return d