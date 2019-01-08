# server.py
from flask import Flask, render_template
from flask_oauthlib.provider import OAuth2Provider


app = Flask(__name__, static_folder="../static/dist", template_folder="../static")
oauth = OAuth2Provider(app)


# render our react app here
@app.route("/")
def index():
    return render_template("index.html")


# other api endpoints here...
@app.route("/hello")
def hello():
    return "Hello World!"


if __name__ == "__main__":
    app.run()
