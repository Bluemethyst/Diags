from flask import Flask, request, render_template, abort
from flask_cors import CORS
import json, os

app = Flask(__name__)
CORS(app)


@app.route('/<string:client_ip>/')
def root(client_ip):
    file_path = "data/" + client_ip + '.json'
    if not os.path.isfile(file_path):
        abort(404, description="Resource not found")
    with open(file_path, 'r') as f:
        data = json.load(f)
        client_data = json.dumps(data)
    return render_template("index.html", client_data=data)


@app.route('/api/upload_stats', methods=['POST', 'GET'])
def add_post():
    if request.method == 'GET':
        return {"status_code": 404, "message": "Page not found"}, 404
    elif request.method == 'POST':
        try:
            content = request.json
            print(json.dumps(content, indent=4))
            with open("data/" + request.remote_addr + '.json', 'w') as f:
                json.dump(content, f, indent=4)
            return {"status_code": 200, "message": "success, post added"}, 200
        except Exception as e:
            print(e)
            return {"status_code": 500, "message": "error, post not added"}, 500


@app.errorhandler(404)
def page_not_found(e):
    return {"status_code": 404, "message": "Page not found"}, 404


if __name__ == '__main__':
    app.run(host="0.0.0.0", debug=True)
