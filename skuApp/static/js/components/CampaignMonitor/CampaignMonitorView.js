import React from 'react';
import CMHandler from './CMHandler';

export default class CampaignMonitorView extends React.Component {
    render() {
        return (
            <div className="container">
                <small className="text-muted help-summary">
                    Expected Column Headers:<br/>
                    (Network, Targeting Method, Format, Message, Age, Ethnicity, Family Role, Gender, Income,<br/>
                    Interests/Behaviors, Language, Level of Education, Occupation, Relationship, Religion)
                </small>
                <br/>

                <div className="row">
                    <div style={{margin: 'auto'}}>
                        <CMHandler/>
                    </div>
                </div>
            </div>
        );
    }
}