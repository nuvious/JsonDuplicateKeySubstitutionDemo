from flask import Flask, render_template, request, jsonify
import requests

app = Flask(__name__)
app.secret_key = 'some_secret_key'
API_URL = 'http://auth-server:8080/auth'

@app.route('/', methods=['GET', 'POST'])
def login():
    if request.method == 'POST':
        data = request.stream.read()
        response = requests.post(API_URL, data=data)

        auth_response = response.json()
        if "username" in auth_response:
            username = auth_response["username"]
            if username == "admin":
                with open("flag.txt", "r") as f:
                    username = f.read().strip()
            return jsonify({"status": "success", "username": username})
        else:
            return jsonify({"status": "failure"})

    return render_template('login.html')


@app.route('/welcome/<username>', methods=['GET'])
def welcome(username):
    return render_template('welcome.html', username=username)


if __name__ == '__main__':
    app.run(debug=True)
