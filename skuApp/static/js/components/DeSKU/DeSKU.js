import React from 'react';
import codes from '../skuCodes';
import UploadDeSKU from './UploadDeSKU';


class DeSKU extends React.Component {

    constructor(props) {
        super(props);
    }

    state = {
        showBreakdown: false,
        inputSKU: '',
        network: '',
        targetingMethod: '',
        format: '',
        message: '',
        ageRange: '',
        ethnicity: '',
        familyRole: '',
        gender: '',
        income: '',
        interestsBehaviors: '',
        language: '',
        education: '',
        occupation: '',
        relationship: '',
        religion: ''
    };

    handleSelect = name => event => {
        this.setState({
            [name]: event.target.value,
        });
        console.log(this.state);
    };

    deSKU = () => {
        let skuString = this.state.inputSKU;
        // let skuString = 'A1B1C1D1E9F9G13H4I12J21K5L11M17N6O11';
        let explodeSKU = skuString.split(/(?=[A-Z])/);
        let skuList = [
            "network", "targetingMethod", "format", "message", "ageRange", "ethnicity", "familyRole", "gender",
            "income", "interestsBehaviors", "language", "education", "occupation", "relationship", "religion"
        ];

        // Make sure full sku exists
        if (explodeSKU.length !== skuList.length) {
            return alert("Sku is not complete.");
        }

        for (let i = 0; i < explodeSKU.length; i++) {
            this.setState({
                [skuList[i]]: codes.InverseCodes[explodeSKU[i]]
            });
            console.log(codes.InverseCodes[explodeSKU[i]]);
        }
        console.log(this.state);
        this.setState({
            showBreakdown: true
        });
    };

    render() {
        const showBreakdown = this.state.showBreakdown;

        return (
            <div className="container">
                <small className="text-muted help-summary">
                    Reverse-engineer a single sku code <em>OR</em> upload for bulk de-skuing.
                </small>
                <br/>

                <div className="row">
                    <div className="col-lg-4 order-lg-2">
                        <UploadDeSKU/>
                    </div>

                    <div className="col-lg-8 order-lg-1" style={{minWidth: '600px'}}>
                        <h4 className="mb-3">DeSKU</h4>
                        <form>
                            <div className="row">
                                <div className="mb-3" style={{width: '100%'}}>
                                    <input type="text" className="form-control" id="deSKU" placeholder="SKU Input"
                                           onChange={this.handleSelect('inputSKU')} required=""/>
                                    <div className="invalid-feedback">
                                        Valid sku is required.
                                    </div>
                                </div>
                            </div>
                            <span className="btn btn-primary btn-lg btn-block" onClick={this.deSKU}>DeSKU</span>
                        </form>

                        {showBreakdown &&
                        <ul className="list-group" style={{width: '100%', marginTop: '20px'}}>
                            <li className="list-group-item">
                                <b>Network:</b>
                                <span className="listCode"> {this.state.network}</span>
                            </li>
                            <li className="list-group-item">
                                <b>Targeting Method:</b>
                                <span className="listCode"> {this.state.targetingMethod}</span>
                            </li>
                            <li className="list-group-item">
                                <b>Format:</b>
                                <span className="listCode"> {this.state.format}</span>
                            </li>
                            <li className="list-group-item">
                                <b>Message:</b>
                                <span className="listCode"> {this.state.message}</span>
                            </li>
                            <li className="list-group-item">
                                <b>Age Range:</b>
                                <span className="listCode"> {this.state.ageRange}</span>
                            </li>
                            <li className="list-group-item">
                                <b>Ethnicity:</b>
                                <span className="listCode"> {this.state.ethnicity}</span>
                            </li>
                            <li className="list-group-item">
                                <b>Family Role:</b>
                                <span className="listCode">{this.state.familyRole}</span>
                            </li>
                            <li className="list-group-item">
                                <b>Gender:</b>
                                <span className="listCode"> {this.state.gender}</span>
                            </li>
                            <li className="list-group-item">
                                <b>Income:</b>
                                <span className="listCode"> {this.state.income}</span>
                            </li>
                            <li className="list-group-item">
                                <b>Interests/Behaviors:</b>
                                <span className="listCode"> {this.state.interestsBehaviors}</span>
                            </li>
                            <li className="list-group-item">
                                <b>Language:</b>
                                <span className="listCode"> {this.state.language}</span>
                            </li>
                            <li className="list-group-item">
                                <b>Level of Education:</b>
                                <span className="listCode"> {this.state.education}</span>
                            </li>
                            <li className="list-group-item">
                                <b>Occupation:</b>
                                <span className="listCode"> {this.state.occupation}</span>
                            </li>
                            <li className="list-group-item">
                                <b>Relationship:</b>
                                <span className="listCode"> {this.state.relationship}</span>
                            </li>
                            <li className="list-group-item">
                                <b>Religion:</b>
                                <span className="listCode"> {this.state.religion}</span>
                            </li>
                        </ul>
                        }
                    </div>
                </div>
            </div>
        );
    }
}

export default DeSKU;