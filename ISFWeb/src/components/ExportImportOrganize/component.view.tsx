import * as React from 'react';
import ExportImportOrganizeBase from './component.base';
import ExportOrganize from './ExportOrganize/component.view'
import MainScreen from './MainScreen/component.view';
import ImportOrganize from './ImportOrganize/component.view';
import styles from './styles.view';

export default class ExportImportOrganize extends ExportImportOrganizeBase {

    render() {
        const { mainScreen, exportOrganize, importOrganize, failNum, successNum } = this.state;
        const { userid } = this.props;

        return (
            <div className={styles['export-import-screen']}>
                {
                    mainScreen ?
                        <MainScreen
                            userid={userid}
                            onCancel={this.handleCancelOperation.bind(this)}
                            onExportItem={this.handleExportItem.bind(this)}
                            onImportItem={this.handleImportItem.bind(this)}
                        />
                        : null
                }
                {
                    exportOrganize ?
                        <ExportOrganize
                            userid={userid}
                            onCancel={this.handleCancelOperation.bind(this)}
                            onDownloadFile={this.handleDownloadFile.bind(this)}
                        />
                        : null
                }
                {
                    importOrganize ?
                        <ImportOrganize
                            failNum={failNum}
                            successNum={successNum}
                            onCancel={this.handleCancelOperation.bind(this)}
                            onContinue={this.handlecontinue.bind(this)}
                            onImportSuccess={this.handleImportSuccess.bind(this)}
                        />
                        : null
                }
            </div>
        )
    }
}