import React from 'react';
import ContentStringGenerator from './ContentStringGenerator';
import CSVHandler from './CSVHandler';
import {facilityList, domainShortnames} from "../../facilityList";
import axios from "axios";


class URLForm extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            facility: '',
            path: '',
            domain: '',
            shortname: '',
            medium: '',
            source: '',
            campaign: '',
            term: '',
            content: 'A1B1C1D1E9F9G13H4I12J21K5L11M17N6O11',
            processedURL: '',
        }
    }

    // Update state facility and domain on user selection
    handleFacilitySelect = event => {
        let sname = domainShortnames[event.target.value] ? domainShortnames[event.target.value] : '';
        this.setState({
            facility: event.target.key,
            domain: event.target.value,
            shortname: sname,
        }, () => {
            console.log(this.state);
        });
    };

    handleSelect = name => event => {
        this.setState({
            [name]: event.target.value,
        });
    };

    // Update content sku value
    updateContentString = (val) => {
        this.setState({
            content: val
        }, () => {
            console.log(this.state);
        });
    };

    // Swap user input spaces with '+'
    handleCampaignTerm = (str) => {
        // regex for whitespace
        str = str.replace(/\s+/g, '+');
        return str;
    };

    processURL = () => {
        this.validateURL();
        // setup base url path
        let url = "https://" + this.state.domain + this.state.path + "?";

        // check source
        if (this.state.source.length < 1) return alert("Invalid utm_source");
        url += "utm_source=" + this.state.source;

        // check medium
        if (this.state.medium.length < 1) return alert("Invalid utm_medium");
        url += "&utm_medium=" + this.state.medium;

        // check shortname
        if (this.state.shortname.length < 1) return alert("Invalid shortname");
        url += "&sf_shortname=" + this.state.shortname;

        // check campaign
        if (this.state.campaign.length < 1) return alert("Invalid utm_campaign");
        url += "&utm_campaign=" + this.state.campaign;

        // check term
        if (this.state.term.length > 0) url += "&utm_term=" + this.state.term;

        // check content
        if (this.state.content.length < 1) return alert("Invalid utm_content");
        url += "&utm_content=" + this.state.content;

        this.setState({
            processedURL: url,
        });
    };

    validateURL = () => {
        let url = "https://" + this.state.domain + this.state.path;
        axios.get(url)
            .then(function (response) {
                console.log(response.status);
            })
            .catch(function (error) {
                // handle error
                console.log(error);
                return alert("Invalid URL/Path\n" + error);
            });
    };

    render() {
        let facilityDomain = this.state.domain;
        let sfShortname = this.state.shortname;

        return (
            <div className="container">
                <small className="text-muted help-summary">
                    Build a URL from scratch <em>OR</em> upload from Kenshoo and let SKUle automate sku's for you.
                </small>
                <br/>

                <div className="row">
                    <div className="col-lg-4 order-lg-2">
                        <CSVHandler/>
                    </div>

                    <div className="col-lg-8 order-lg-1">
                        <h4 className="mb-3">Build URL</h4>
                        <form>
                            <div className="row">
                                <div className="col-md-6 mb-3">
                                    <label htmlFor="facility">Facility</label>
                                    <select
                                        className="custom-select d-block w-100"
                                        id="facility"
                                        required=""
                                        value={this.state.domain}
                                        onChange={this.handleFacilitySelect}
                                    >
                                        <option value="">Choose...</option>
                                        {
                                            facilityList.map((facility) =>
                                                <option key={facility.facility_name} value={facility.domain}>
                                                    {facility.facility_name}
                                                </option>
                                            )
                                        }
                                    </select>
                                    <div className="invalid-feedback">
                                        Please select a valid facility.
                                    </div>
                                </div>
                                <div className="col-md-6 mb-3">
                                    <label htmlFor="path">Path</label>
                                    <input type="text" className="form-control" id="path" placeholder="/example/path"
                                           onChange={this.handleSelect('path')} required=""/>
                                    <div className="invalid-feedback">
                                        Valid path is required.
                                    </div>
                                </div>
                            </div>

                            <div className="mb-3">
                                <label htmlFor="domain">Domain</label>
                                <ul className="list-group">
                                    <li className="list-group-item disabled"
                                        style={{padding: '.375rem .75rem'}}>{facilityDomain}</li>
                                </ul>
                                {/*<input type="email" className="form-control" id="domain" placeholder="autofilled from facility"/>*/}
                            </div>
                            <div className="mb-3">
                                <label htmlFor="shortname">SF Shortname</label>
                                <input type="text" className="form-control" id="shortname" value={this.state.shortname}
                                       onChange={this.handleSelect('shortname')} required=""/>
                                <div className="invalid-feedback">
                                    Valid sfShortname is required.
                                </div>
                            </div>

                            <div className="row">
                                <div className="col-md-6 mb-3">
                                    <label htmlFor="medium">Medium</label>
                                    <select
                                        className="custom-select d-block w-100"
                                        id="medium"
                                        required=""
                                        value={this.state.medium}
                                        onChange={this.handleSelect('medium')}
                                    >
                                        <option value="">Choose...</option>
                                        <option value={'ad-video'}>Ad-Video</option>
                                        <option value={'cpc'}>CPC</option>
                                        <option value={'direct'}>Direct</option>
                                        <option value={'email'}>Email</option>
                                        <option value={'offline'}>Offline</option>
                                        <option value={'organic'}>Organic</option>
                                        <option value={'referral'}>Referral</option>
                                    </select>
                                    <div className="invalid-feedback">
                                        Please select a valid utm_medium.
                                    </div>
                                </div>
                                <div className="col-md-6 mb-3">
                                    <label htmlFor="source">Source</label>
                                    <select
                                        className="custom-select d-block w-100"
                                        id="source"
                                        required=""
                                        value={this.state.source}
                                        onChange={this.handleSelect('source')}
                                    >
                                        <option value="">Choose...</option>
                                        <option value={'addiction-center'}>Addiction Center</option>
                                        <option value={'bing'}>Bing</option>
                                        <option value={'business+development'}>Business Development</option>
                                        <option value={'campaign-monitor'}>Campaign Monitor</option>
                                        <option value={'consumer'}>Consumer</option>
                                        <option value={'eating-disorder-hope'}>Eating Disorder Hope</option>
                                        <option value={'google'}>Google</option>
                                        <option value={'guidedoc'}>Guidedoc</option>
                                        <option value={'healthy-place'}>Healthy Place</option>
                                        <option value={'luxury-rehab'}>Luxury Rehab</option>
                                        <option value={'psych-today'}>Psych Today</option>
                                        <option value={'quora'}>Quora</option>
                                        <option value={'self-growth'}>Self Growth</option>
                                        <option value={'stackadapt'}>Stackadapt</option>
                                        <option value={'theravive'}>Theravive</option>
                                        <option value={'webmd'}>WebMD</option>
                                        <option value={'yellow-pages'}>Yellow Pages</option>
                                        <option value={'yelp'}>Yelp</option>
                                    </select>
                                    <div className="invalid-feedback">
                                        Please provide a valid utm_source.
                                    </div>
                                </div>
                            </div>
                            <div className="row">
                                <div className="col-md-6 mb-3">
                                    <label htmlFor="campaign">Campaign</label>
                                    <input type="text" className="form-control" id="campaign"
                                           onChange={this.handleSelect('campaign')} required=""/>
                                    <div className="invalid-feedback">
                                        Valid campaign is required.
                                    </div>
                                </div>
                                <div className="col-md-6 mb-3">
                                    <label htmlFor="term">Term</label>
                                    <input type="text" className="form-control" id="term"
                                           onChange={this.handleSelect('term')} required=""/>
                                    <div className="invalid-feedback">
                                        Valid term is required.
                                    </div>
                                </div>
                            </div>

                            <hr className="mb-4"/>
                            <ContentStringGenerator parentContent={this.state.content}
                                                    updateContent={this.updateContentString}/>
                            <hr className="mb-4"/>

                            <span className="btn btn-primary btn-lg btn-block" onClick={this.processURL}>
                            Generate
                        </span>
                        </form>
                        {this.state.processedURL.length > 0 &&
                        <div>
                            <h5 className="card-title" style={{color: '#F5B961'}}>Processed URL</h5>
                            <div className="card">
                                <div className="card-body">
                                    <samp className="processedURL">{this.state.processedURL}</samp>
                                </div>
                            </div>
                        </div>
                        }
                    </div>
                </div>
            </div>
        );
    }
}

export default URLForm;