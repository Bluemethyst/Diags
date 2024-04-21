from flask import Flask, request
from flask_cors import CORS
import sqlite3
import json

app = Flask(__name__)
CORS(app)


@app.route('/')
def hello_world():
    return {"status_code": 200, "message": "you shouldn't be here :)"}


@app.route('/upload_stats', methods=['POST', 'GET'])
def add_post():
    if request.method == 'GET':
        return {"status_code": 200, "message": "you shouldn't be here :)"}
    elif request.method == 'POST':
        try:
            content = request.json
            print(json.dumps(content, indent=4))
            print(request.remote_addr)
            return {"status_code": 200, "message": "success, post added"}, 200
        except Exception as e:
            print(e)
            return {"status_code": 500, "message": "error, post not added"}, 500


@app.errorhandler(404)
def page_not_found(e):
    return {"status_code": 404, "message": "Page not found"}, 404


if __name__ == '__main__':
    app.run(debug=True)
