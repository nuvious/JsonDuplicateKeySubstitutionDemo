from flask import Flask, render_template, request, jsonify
import requests
import uuid

app = Flask(__name__)
app.secret_key = uuid.uuid4().hex

# This server authenticates the user based on the username/password supplied to
# this service.
API_URL = 'http://auth-server:8080/auth'

@app.route('/', methods=['GET', 'POST'])
def login():
    if request.method == 'POST':
        # Get the json payload, it better have "username" and "password" in it.
        data = request.stream.read()
        
        # Send the payload to the auth server and get a response
        response = requests.post(API_URL, data=data)
        auth_response = response.json()
        
        # If the "username" key is in the response, it's authenticated, the
        # auth server returns {} on a failed authentication
        if "username" in auth_response:
            username = auth_response["username"]
            
            # TODO: Removed this, admin isn't allowed to log in on the new
            #       authentication server.
            if username == "admin":
                # If it's the admin give them the flag.
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
