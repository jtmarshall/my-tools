import React from 'react';
import AuditCSV from './AuditCSV';

class Audit extends React.Component {

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

    render() {
        return (
            <div className="container">
                <small className="text-muted help-summary">
                    Audit check for Kenshoo URLs.
                </small>
                <br/>

                <div className="row">
                    <div className="col-lg-4 order-lg-2">
                        <AuditCSV/>
                    </div>
                </div>
            </div>
        );
    }
}

export default Audit;