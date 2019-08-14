# server.py
import os
import json
import boto3
import io
from functools import wraps
from flask import Flask, Response, render_template, request, jsonify
from werkzeug.utils import secure_filename
# fixes relative import issues
import sys; sys.path.append(os.path.dirname(os.path.realpath(__file__)))
# local packages
from site_packages.skule_simplelogin import SimpleLogin, login_required
from site_packages.processCSV import processCSV, s3processCSV, urlAuditCSV
from site_packages.deSKU import processDeSKU
from site_packages.campaign_monitor import processCM


# basic auth
creds = {
    'usr': "Username",
    'pwd': "password"
}
app = Flask(__name__, static_folder="../static", template_folder="../static")
UPLOAD_FOLDER = '/tmp/'
app.config['UPLOAD_FOLDER'] = UPLOAD_FOLDER
app.config['SIMPLELOGIN_USERNAME'] = creds['usr']
app.config['SIMPLELOGIN_PASSWORD'] = creds['pwd']
SimpleLogin(app)


# check if username/password combo is valid
def check_auth(username, password):
    return username == creds['usr'] and password == creds['pwd']


# sends 401 response that enables basic auth
def authenticate():
    return Response('Unable to verify access.\nMust have proper credentials', 401, {'WWW-Authenticate': 'Basic'})


# decorator for basic auth on routes
def requires_auth(f):
    @wraps(f)
    def decorated(*args, **kwargs):
        auth = request.authorization
        if not auth or not check_auth(auth.username, auth.password):
            return authenticate()
        return f(*args, **kwargs)

    return decorated


ALLOWED_EXTENSIONS = set(['csv', 'xls', 'xlsx'])


# helper func for csv upload to see if file is allowed
def allowed_file(filename):
    return '.' in filename and \
           filename.rsplit('.', 1)[1].lower() in ALLOWED_EXTENSIONS


"""
We no longer need this conversion function for processing,
but keeping it just in case.
"""
# @app.route("/api/jsonifyCSV", methods=['POST'])
# def csvToJSON():
#     # check if request has file we're looking for
#     if 'file' not in request.files:
#         print('No file part in request')
#         return 'No file part in request'
#
#     file = request.files['file']
#
#     # if user does not select file, browser also submit a empty part without filename
#     if file.filename == '':
#         return 'No selected file'
#
#     filename = secure_filename(file.filename)
#     filepath = os.path.join(app.config['UPLOAD_FOLDER'], filename)
#     file.save(filepath)
#
#     columns = []
#     with open(filepath, 'rU') as f:
#         reader = csv.reader(f)
#         for row in reader:
#             if columns:
#                 for i, value in enumerate(row):
#                     if value != '':
#                         columns[i].append(value)
#             else:
#                 # first row
#                 columns = [[value] for value in row]
#     # you now have a column-major 2D array of your file.
#     as_dict = {c[0]: c[1:] for c in columns}
#     print(as_dict)
#
#     json_data = json.dumps(as_dict)
#
#     print(json_data)
#     return jsonify(json_data)


"""
Digest & build SKU's for user uploaded csv file
"""
@app.route("/api/buildCSV", methods=['POST'])
def build_csv_handler():
    try:
        # get csv build type
        buildType = request.args['buildType']
        # string args for S3 or not
        useS3 = request.args['useS3']
        # check for S3
        if useS3 == 'true':
            # Get filenames, bucketnames, etc...
            filename = request.args['filename']
            bucketname = 'skule'
            s3 = boto3.client('s3')
            filepath = os.path.join(app.config['UPLOAD_FOLDER'], filename)
            reFilename = "SKUle-" + filename
            reFilepath = os.path.join(app.config['UPLOAD_FOLDER'], reFilename)

            # filepath = '/tmp/' + filename
            s3.download_file(bucketname, filename, filepath)
            # process the csv
            if buildType == 'CM':
                return processCM(filepath, reFilename, reFilepath)
            return processCSV(filepath, reFilename, reFilepath)

        else:
            # check if request has file we're looking for
            if 'file' not in request.files:
                print('No file part in request')
                return 'No file part in request'

            file = request.files['file']

            # if user does not select file, browser also submit a empty part without filename
            if file.filename == '':
                return 'No selected file'

            # save uploaded file to the uploads folder so we can access it later
            if file and allowed_file(file.filename):
                filename = secure_filename(file.filename)
                filepath = os.path.join(app.config['UPLOAD_FOLDER'], filename)
                reFilename = "SKUle-" + filename
                reFilepath = os.path.join(app.config['UPLOAD_FOLDER'], reFilename)
                file.save(filepath)

                if buildType == 'CM':
                    return processCM(filepath, reFilename, reFilepath)
                return processCSV(filepath, reFilename, reFilepath)
    except Exception as e:
        return Response("There was an issue reading file from request. \n%s" % e, status=400)


# user upload for deSKU process
@app.route("/api/deSKU", methods=['POST'])
def desku_handler():
    try:
        # check if request has file we're looking for
        if 'file' not in request.files:
            print('No file part in request')
            return 'No file part in request'

        file = request.files['file']

        # if user does not select file, browser also submit a empty part without filename
        if file.filename == '':
            return 'No selected file'

        # save uploaded file to the uploads folder so we can access it later
        if file and allowed_file(file.filename):
            filename = secure_filename(file.filename)
            filepath = os.path.join(app.config['UPLOAD_FOLDER'], filename)
            reFilename = "SKUle-" + filename
            reFilepath = os.path.join(app.config['UPLOAD_FOLDER'], reFilename)
            file.save(filepath)

            return processDeSKU(filepath, reFilename, reFilepath)
    except Exception as e:
        return Response("There was an issue reading file from request. \n%s" % e, status=400)


# If file is larger than 6MB
@app.route("/api/s3import", methods=['GET', 'POST'])
def s3import_handler():
    try:
        # Get the service client.
        s3 = boto3.client('s3')

        # Generate the presigned URL for put requests
        presigned_url = s3.generate_presigned_url(
            ClientMethod='put_object',
            Params={
                'Bucket': 'skule',
                'Key': request.args['filename'],
                'Expires': 60,
                'ContentType': request.args['filetype'],
                'ACL': 'public-read'
            }
        )

        # Return the presigned URL
        return presigned_url
    except Exception as e:
        return Response("There was an issue creating/returning the presigned URL. \n%s" % e, status=400)


# API route for auditing url's via csv upload
@app.route("/api/auditCSV", methods=['POST'])
def audit_csv_handler():
    try:
        # string args for S3 or not
        useS3 = request.args['useS3']
        # check for S3 (must send to S3 for status audit)
        if useS3 == "true":
            # Get filenames, bucketnames, etc...
            filename = request.args['filename']
            bucketname = 'skule'
            s3 = boto3.client('s3')
            filepath = os.path.join(app.config['UPLOAD_FOLDER'], filename)
            reFilename = "SKUle-" + filename
            reFilepath = os.path.join(app.config['UPLOAD_FOLDER'], reFilename)

            # filepath = '/tmp/' + filename
            s3.download_file(bucketname, filename, filepath)
            # process the csv
            return urlAuditCSV(filepath, reFilename, reFilepath)

        else:
            # check if request has file we're looking for
            if 'file' not in request.files:
                print('No file part in request')
                return 'No file part in request'

            file = request.files['file']

            # if user does not select file, browser also submit a empty part without filename
            if file.filename == '':
                return 'No selected file'

            # save uploaded file to the uploads folder so we can access it later
            if file and allowed_file(file.filename):
                filename = secure_filename(file.filename)
                filepath = os.path.join(app.config['UPLOAD_FOLDER'], filename)
                reFilename = "SKUle-" + filename
                reFilepath = os.path.join(app.config['UPLOAD_FOLDER'], reFilename)
                file.save(filepath)

                return urlAuditCSV(filepath, reFilename, reFilepath, app.config['UPLOAD_FOLDER'])
    except Exception as e:
        return Response("There was an issue reading file from request. \n%s" % e, status=400)


# render our react app here
@app.route("/")
@login_required
def home():
    return render_template("index.html", rootenv="dev")


if __name__ == "__main__":
    app.run()
