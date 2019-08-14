import axios from "axios";
import React from "react";

export default class CMHandler extends React.Component {

    pathName = location.pathname === "/" ? "" : location.pathname;
    csvEndpoint = location.origin + this.pathName + "/api/buildCSV";

    constructor(props) {
        super(props);
    }

    state = {
        selectedFile: null,
        loaded: 0,
        showLoading: false,
        triggeredLargeUpload: false,
        largeUploadFinish: false,
        largeReturnURL: '',
        buildType: 'CM',
    };

    // Exists so we can bind state update to inner scope of secondary requests
    largeReturnState = (url) => {
        this.setState({
            triggeredLargeUpload: false,
            largeUploadFinish: true,
            largeReturnURL: url
        });
    };

    // update selectedFile on user input
    handleSelectedFile = event => {
        // Reset state when new file selected
        this.setState({
            selectedFile: event.target.files[0],
            loaded: 0,
            triggeredLargeUpload: false,
            largeUploadFinish: false,
            largeReturnURL: ''
        })
    };

    // send api request with user file
    uploadRequest = (data) => {
        // Show loading dots
        this.setState({
            showLoading: true,
        });
        return axios.post(this.csvEndpoint, data, {
            headers: {'Content-Type': 'multipart/form-data'},
            params: {
                useS3: false,
                buildType: this.state.buildType,
            }
        })
            .then((response) => {
                console.log("Headers: ", response.headers);
                if (response.headers['content-disposition']) {
                    // get filename from resp
                    const filename = response.headers['content-disposition'].split('filename=')[1];
                    // create blob to start browser download
                    const url = window.URL.createObjectURL(new Blob([response.data]));
                    const link = document.createElement('a');
                    link.href = url;
                    link.setAttribute('download', filename);
                    document.body.appendChild(link);
                    link.click();
                } else {
                    // If we are expecting an S3 link back
                    window.open(response.data, '_blank');
                    this.largeReturnState(response.data);
                }
                return response.data;
            })
            .catch((err) => {
                console.log("Upload ERR: ", err);
                if (err.response) {
                    console.log(err.response);
                    console.log(err.response.data);
                    alert("Upload ERR: " + err.response.data);
                } else {
                    alert("Upload ERR: " + err.message);
                }
            })
            .finally(() => {
                this.setState({
                    showLoading: false,
                });
            })
    };

    // send api request with user file
    largeUploadRequest = () => {
        // Allow show user that the upload is large
        this.setState({
            triggeredLargeUpload: true,
        });
        // get s3 signed url to upload file into s3 first
        const s3Endpoint = location.origin + this.pathName + "/api/s3import";
        const csvEndpoint = this.csvEndpoint;
        let file = this.state.selectedFile;
        let signedUrl = '';
        let buildType = this.state.buildType;
        // Bind largeReturnState
        let largeReturnState = this.largeReturnState;

        axios.get(s3Endpoint, {
            params: {
                filename: file.name,
                filetype: file.type
            }
        })
            .then(function (result) {
                signedUrl = result.data;
                console.log(signedUrl);

                let options = {
                    headers: {
                        'Content-Type': file.type,
                        'x-amz-acl': 'public-read'
                    }
                };
                // PUT the file in S3
                axios.put(signedUrl, file, options).then(function (result) {
                    console.log(result);
                    // After file uploaded into S3, kickoff processing
                    return axios.post(csvEndpoint, null, {
                        headers: {'Content-Type': 'multipart/form-data'},
                        params: {
                            useS3: true,
                            filename: file.name,
                            signedUrl: signedUrl,
                            buildType: buildType,
                        }
                    })
                        .then((response) => {
                            // If no errors update the state
                            largeReturnState(response.data);
                            // Trigger download by opening return url in new tab
                            window.open(response.data, '_blank');
                            return response.data;
                        })
                        .catch((err) => {
                            console.log("Upload ERR: ", err);
                            if (err.response) {
                                console.log(err.response);
                                console.log(err.response.data);
                                alert("Upload ERR: " + err.response.data);
                            } else {
                                alert("Upload ERR: " + err.message);
                            }
                        })
                })
                    .catch(function (err) {
                        return console.log(err);
                    });
                // END PUT
            })
            .catch(function (err) {
                return console.log(err);
            });
    };

    // on button click
    handleUpload = () => {
        // Check file
        let fileExt = this.state.selectedFile.name.split('.').pop();
        if (fileExt !== 'csv' && fileExt !== 'xlsx') {
            alert("Upload file must be .csv or .xlsx");
            return
        }
        this.setState({
            largeReturnURL: '',
        });
        // large file size ~6MB or larger
        if (this.state.selectedFile.size > 5900000) {
            console.log("large upload");
            this.largeUploadRequest();
        } else {
            // small file size
            const data = new FormData();
            data.append('file', this.state.selectedFile, this.state.selectedFile.name);

            // start request
            this.uploadRequest(data);
        }
    };

    render() {
        return (
            <div style={{width: '460px'}}>
                <h4 className="d-flex justify-content-between align-items-center mb-3">
                    <span className="text-muted"
                          style={{display: "block", margin: "auto"}}>CM: Build SKU</span>
                </h4>

                <div className="input-group">
                    <input type="file" className="form-control" placeholder="CSV Input"
                           onChange={this.handleSelectedFile}/>
                    <div className="input-group-append">
                        <span className="btn btn-secondary" onClick={this.handleUpload}>Upload</span>
                    </div>
                </div>
                <p style={{textAlign: 'center', color: '#6c757d', padding: '8px', fontSize: '0.9em'}}>
                    Expecting .csv or .xlsx
                </p>

                {this.state.showLoading &&
                <p className="loading">
                    <span><code>Grading</code>.</span><span>.</span><span>.</span>
                </p>
                }

                {this.state.triggeredLargeUpload &&
                <p className="loading">
                    <code>It looks like your upload is quite <b>large</b>, consulting the cloud</code>
                    <span>.</span><span>.</span><span>.</span>
                </p>
                }

                {this.state.largeReturnURL.length > 1 &&
                <p className="loading">
                    <code>
                        It's done! Download should commence automatically.<br/>
                        But if it doesn't you can also just click
                        <a target="_blank" href={this.state.largeReturnURL}> here.</a>
                    </code>
                </p>
                }

            </div>
        );
    }
}