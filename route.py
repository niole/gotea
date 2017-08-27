"""
The api endpoint for getting the top ten
tea matches for a user query
"""
from flask import Flask, request, jsonify
import json
import dbConnection
from scorer import getRankedMatches
from dbConnection import Tea
app = Flask(__name__)

@app.route('/match', methods=['POST'])
def matcher():
    """
    request.data - { userQuery: <querystring> }
    """

    if request.data:
        query = json.loads(request.data)
        print query
        userQuery = query[u'userQuery']
        teas = Tea.select().execute()
        teaMetadata = [{
            'name': t.name,
            'link': t.link
        } for t in teas]
        teaData = [t.data for t in teas]
        scores = getRankedMatches(teaData, userQuery)
        rankedTeas = zip(scores, teaMetadata)
        sorted(
            rankedTeas,
            key=lambda t: t[0],
            reverse=True
        )
        topTenMatches = [t[1] for t in rankedTeas[:10]]
        return jsonify(topTenMatches)
