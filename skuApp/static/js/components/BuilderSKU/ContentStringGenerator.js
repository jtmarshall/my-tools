// App.jsx
import React from "react";
import skuCodes from "../skuCodes";

class ContentStringGenerator extends React.Component {

    constructor(props) {
        super(props);
    }

    state = {
        network: 'A1',
        targeting: 'B1',
        format: 'C1',
        message: 'D1',
        age: 'E9',
        ethnicity: 'F9',
        familyRole: 'G13',
        gender: 'H4',
        income: 'I12',
        interests: 'J21',
        language: 'K5',
        education: 'L11',
        occupation: 'M17',
        relationship: 'N6',
        religion: 'O11',
    };

    // update local state content options
    handleLocalContent = name => event => {
        let contentString = '';

        this.setState({
            [name]: event.target.value,
        }, () => {
            // Concatenate all values from state to contentString
            Object.keys(this.state).forEach((key) => {
                contentString += this.state[key];
            });

            // Pass updated content string to parent component
            this.props.updateContent(contentString);
        });
    };

    render() {
        // Use content string from parent component state
        let utmContent = this.props.parentContent;

        return (
            <div className="componentContainer">

                <div className="mb-3">
                    <label htmlFor="content">Content</label>
                    <ul className="list-group">
                        <li className="list-group-item disabled"
                            style={{padding: '.375rem .75rem'}}>{utmContent}</li>
                    </ul>
                </div>
                <div className="row" style={{marginLeft: '4px', marginRight: '4px'}}>
                    <div className="col-md-4 mb-3">
                        <label htmlFor="network">Network</label>
                        <select
                            className="custom-select d-block w-100"
                            id="network"
                            required=""
                            value={this.state.network}
                            onChange={this.handleLocalContent('network')}
                        >
                            {Object.keys(skuCodes.Network).sort().map( key =>
                                <option value={skuCodes.Network[key]}>{key}</option>
                            )}
                        </select>
                    </div>
                    <div className="col-md-4 mb-3">
                        <label htmlFor="targeting">Targeting Method</label>
                        <select
                            className="custom-select d-block w-100"
                            id="targeting"
                            required=""
                            value={this.state.targeting}
                            onChange={this.handleLocalContent('targeting')}
                        >
                            {Object.keys(skuCodes.TargetingMethod).sort().map( key =>
                                <option value={skuCodes.TargetingMethod[key]}>{key}</option>
                            )}
                        </select>
                    </div>
                    <div className="col-md-4 mb-3">
                        <label htmlFor="format">Format</label>
                        <select
                            className="custom-select d-block w-100"
                            id="format"
                            required=""
                            value={this.state.format}
                            onChange={this.handleLocalContent('format')}
                        >
                            {Object.keys(skuCodes.Format).sort().map( key =>
                                <option value={skuCodes.Format[key]}>{key}</option>
                            )}
                        </select>
                    </div>
                </div>
                <div className="row" style={{marginLeft: '4px', marginRight: '4px'}}>
                    <div className="col-md-4 mb-3">
                        <label htmlFor="message">Message</label>
                        <select
                            className="custom-select d-block w-100"
                            id="message"
                            required=""
                            value={this.state.message}
                            onChange={this.handleLocalContent('message')}
                        >
                            {Object.keys(skuCodes.Message).sort().map( key =>
                                <option value={skuCodes.Message[key]}>{key}</option>
                            )}
                        </select>
                    </div>
                    <div className="col-md-4 mb-3">
                        <label htmlFor="age">Age</label>
                        <select
                            className="custom-select d-block w-100"
                            id="age"
                            required=""
                            value={this.state.age}
                            onChange={this.handleLocalContent('age')}
                        >
                            {Object.keys(skuCodes.AgeRange).sort().map( key =>
                                <option value={skuCodes.AgeRange[key]}>{key}</option>
                            )}
                        </select>
                    </div>
                    <div className="col-md-4 mb-3">
                        <label htmlFor="ethnicity">Ethnicity</label>
                        <select
                            className="custom-select d-block w-100"
                            id="ethnicity"
                            required=""
                            value={this.state.ethnicity}
                            onChange={this.handleLocalContent('ethnicity')}
                        >
                            {Object.keys(skuCodes.Ethnicity).sort().map( key =>
                                <option value={skuCodes.Ethnicity[key]}>{key}</option>
                            )}
                        </select>
                    </div>
                </div>
                <div className="row" style={{marginLeft: '4px', marginRight: '4px'}}>
                    <div className="col-md-4 mb-3">
                        <label htmlFor="familyRole">Family Role</label>
                        <select
                            className="custom-select d-block w-100"
                            id="familyRole"
                            required=""
                            value={this.state.familyRole}
                            onChange={this.handleLocalContent('familyRole')}
                        >
                            {Object.keys(skuCodes.FamilyRole).sort().map( key =>
                                <option value={skuCodes.FamilyRole[key]}>{key}</option>
                            )}
                        </select>
                    </div>
                    <div className="col-md-4 mb-3">
                        <label htmlFor="gender">Gender</label>
                        <select
                            className="custom-select d-block w-100"
                            id="gender"
                            required=""
                            value={this.state.gender}
                            onChange={this.handleLocalContent('gender')}
                        >
                            {Object.keys(skuCodes.Gender).sort().map( key =>
                                <option value={skuCodes.Gender[key]}>{key}</option>
                            )}
                        </select>
                    </div>
                    <div className="col-md-4 mb-3">
                        <label htmlFor="income">Income</label>
                        <select
                            className="custom-select d-block w-100"
                            id="income"
                            required=""
                            value={this.state.income}
                            onChange={this.handleLocalContent('income')}
                        >
                            {Object.keys(skuCodes.Income).sort().map( key =>
                                <option value={skuCodes.Income[key]}>{key}</option>
                            )}
                        </select>
                    </div>
                </div>
                <div className="row" style={{marginLeft: '4px', marginRight: '4px'}}>
                    <div className="col-md-4 mb-3">
                        <label htmlFor="interests">Interests/Behaviors</label>
                        <select
                            className="custom-select d-block w-100"
                            id="interests"
                            required=""
                            value={this.state.interests}
                            onChange={this.handleLocalContent('interests')}
                        >
                            {Object.keys(skuCodes.InterestsBehaviors).sort().map( key =>
                                <option value={skuCodes.InterestsBehaviors[key]}>{key}</option>
                            )}
                        </select>
                    </div>
                    <div className="col-md-4 mb-3">
                        <label htmlFor="language">Language</label>
                        <select
                            className="custom-select d-block w-100"
                            id="language"
                            required=""
                            value={this.state.language}
                            onChange={this.handleLocalContent('language')}
                        >
                            {Object.keys(skuCodes.Language).sort().map( key =>
                                <option value={skuCodes.Language[key]}>{key}</option>
                            )}
                        </select>
                    </div>
                    <div className="col-md-4 mb-3">
                        <label htmlFor="education">Level of Education</label>
                        <select
                            className="custom-select d-block w-100"
                            id="education"
                            required=""
                            value={this.state.education}
                            onChange={this.handleLocalContent('education')}
                        >
                            {Object.keys(skuCodes.Education).sort().map( key =>
                                <option value={skuCodes.Education[key]}>{key}</option>
                            )}
                        </select>
                    </div>
                </div>
                <div className="row" style={{marginLeft: '4px', marginRight: '4px'}}>
                    <div className="col-md-4 mb-3">
                        <label htmlFor="occupation">Occupation</label>
                        <select
                            className="custom-select d-block w-100"
                            id="occupation"
                            required=""
                            value={this.state.occupation}
                            onChange={this.handleLocalContent('occupation')}
                        >
                            {Object.keys(skuCodes.Occupation).sort().map( key =>
                                <option value={skuCodes.Occupation[key]}>{key}</option>
                            )}
                        </select>
                    </div>
                    <div className="col-md-4 mb-3">
                        <label htmlFor="relationship">Relationship</label>
                        <select
                            className="custom-select d-block w-100"
                            id="relationship"
                            required=""
                            value={this.state.relationship}
                            onChange={this.handleLocalContent('relationship')}
                        >
                            {Object.keys(skuCodes.Relationship).sort().map( key =>
                                <option value={skuCodes.Relationship[key]}>{key}</option>
                            )}
                        </select>
                    </div>
                    <div className="col-md-4 mb-3">
                        <label htmlFor="religion">Religion</label>
                        <select
                            className="custom-select d-block w-100"
                            id="religion"
                            required=""
                            value={this.state.religion}
                            onChange={this.handleLocalContent('religion')}
                        >
                            {Object.keys(skuCodes.Religion).sort().map( key =>
                                <option value={skuCodes.Religion[key]}>{key}</option>
                            )}
                        </select>
                    </div>
                </div>

            </div>
        );
    }
}

export default ContentStringGenerator;