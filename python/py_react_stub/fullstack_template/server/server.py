# server.py
from flask import Flask, render_template
from flask_oauthlib.provider import OAuth2Provider
import pandas as pd
import numpy as np

app = Flask(__name__, static_folder="../static/dist", template_folder="../static")
oauth = OAuth2Provider(app)


# render our react app here
@app.route("/")
def index():
    return render_template("index.html")


df = pd.DataFrame(np.random.randint(0, 100, size=(7, 5)), columns=['Direct', 'Email', 'Organic', 'PaidAd', 'Referring'])


# spit out panda data
@app.route("/pandas")
def pandas():
    return df.to_json()


# other api endpoints here...
@app.route("/hello")
def hello():
    return "Hello World!"


if __name__ == "__main__":
    app.run()
