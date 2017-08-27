from flask import Flask, request, jsonify
import dbConnection
import scorer
from dbConnection import Tea
app = Flask(__name__)

@app.route('/match', methods=['POST'])
def matcher():
    if request.data:
        print request.data
        teas = Tea.select().execute()
        formattedTeas = [{
            'data': t.data,
            'name': t.name,
            'link': t.link
        } for t in teas]
        return jsonify(formattedTeas)
