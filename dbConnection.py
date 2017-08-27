"""
connects database and establishes Models
"""

import secrets
import os
from flask import Flask
from peewee import *
from playhouse.flask_utils import FlaskDB
from playhouse.db_url import connect

dbUrl = 'mysql://{0}:{1}@{2}:{3}/{4}'.format(
    os.environ['DB_USER'],
    os.environ['DB_PW'],
    os.environ['DB_HOST'],
    os.environ['DB_PORT'],
    os.environ['DB_NAME']
)

db = connect(dbUrl)
app = Flask(__name__)
app.config.from_object(__name__)

#db_wrapper = FlaskDB(app)

class BaseModel(Model):
    class Meta:
        database = db

class Tea(BaseModel):
    """
    tea table
    """
    name = CharField()
    link = CharField()
    data = TextField()

db.create_tables([Tea], safe=True)
