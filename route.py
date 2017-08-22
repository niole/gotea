from flask import Flask, request
import dbConnection
import scorer
app = Flask(__name__)

@app.route('/match', methods=['POST'])
def matcher():
    if request.data:
        print request.data
        return "you did it"
