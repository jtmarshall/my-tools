// App.jsx
import React from 'react';
import URLForm from './components/BuilderSKU/URLForm';
import DeSKU from './components/DeSKU/DeSKU';
import Audit from './components/Audit/Audit';
import CampaignMonitorView from './components/CampaignMonitor/CampaignMonitorView';


class App extends React.Component {
    puns = [
        'Recess time is back! Let SKUle do the heavy lifting.',
        'Easy A.',
        'No homework here!',
        'Time for some extra credit!'
    ];

    state = {
        activeTab: 1,
        catchPhrase: '',
    };

    // Handles tab nav
    handleTabChange = (val) => {
        this.setState({
            activeTab: val,
        });
    };

    render() {
        const activeTab = this.state.activeTab;
        // Select random pun if catchPhrase empty
        if (this.state.catchPhrase === '') {
            this.state.catchPhrase = this.puns[Math.floor(Math.random() * this.puns.length)];
        }
        // Load random catchphrase
        let catchPhraseElement = document.getElementById('catchPhrase');
        catchPhraseElement.innerText = this.state.catchPhrase;

        return (
            <div className="row">
                <div className="tab-content">
                    <ul className="nav nav-tabs justify-content-center">
                        <li className="nav-item">
                            <a className={this.state.activeTab === 1 ? 'nav-link active' : 'nav-link'}
                               onClick={() => this.handleTabChange(1)}>
                                SKU
                            </a>
                        </li>
                        <li className="nav-item">
                            <a className={this.state.activeTab === 2 ? 'nav-link active' : 'nav-link'}
                               onClick={() => this.handleTabChange(2)}>
                                DeSKU
                            </a>
                        </li>
                        <li className="nav-item">
                            <a className={this.state.activeTab === 3 ? 'nav-link active' : 'nav-link'}
                               onClick={() => this.handleTabChange(3)}>
                                Audit
                            </a>
                        </li>
                        <li className="nav-item">
                            <a className={this.state.activeTab === 4 ? 'nav-link active' : 'nav-link'}
                               onClick={() => this.handleTabChange(4)}>
                                Campaign Monitor
                            </a>
                        </li>
                    </ul>

                    {activeTab === 1 &&
                        <URLForm/>
                    }
                    {activeTab === 2 &&
                        <DeSKU/>
                    }
                    {activeTab === 3 &&
                        <Audit/>
                    }
                    {activeTab === 4 &&
                        <CampaignMonitorView/>
                    }

                </div>
            </div>
        );
    }
}

export default App;