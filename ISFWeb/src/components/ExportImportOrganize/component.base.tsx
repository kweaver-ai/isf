import { noop } from 'lodash';
import WebComponent from '../webcomponent';

interface State {
    /**
     * 主弹窗状态
     */
    mainScreen: boolean;
    /**
     * 导出组织弹框
     */
    exportOrganize: boolean;
    /**
     * 导入组织信息弹框
     */
    importOrganize: boolean;
    /**
     * 导入失败信息个数
     */
    failNum: number;
    /**
     * 导入成功信息个数
     */
    successNum: number;
    /**
     * 导出组织用户信息的任务id
     */
    taskid: string;
    /**
     * 是否有导入成功操作
     */
    import: boolean;
}

interface Props {
    /**
     * 管理员id
     */
    userid: string;
    /**
     * 移除结束的事件
     */
    onComplete: () => any;
    /**
     * 移除成功
     */
    onSuccess: () => any;
}

export default class ExportImportOrganizeBase extends WebComponent<Props, State> {

    static defaultProps = {
        userid: null,

        onComplete: noop,

        onSuccess: noop,
    }

    state = {

        mainScreen: true,

        exportOrganize: false,

        importOrganize: false,

        failNum: null,

        successNum: null,

        taskid: null,

        'import': false,
    }

    /**
     * 导出
     */
    protected handleExportItem() {
        this.setState({
            mainScreen: false,
            exportOrganize: true,
        })
    }

    /**
     * 导入
     */
    protected handleImportItem(failNum: number, successNum: number) {
        this.setState({
            mainScreen: false,
            'import': true,
            importOrganize: true,
            failNum,
            successNum,
        })
    }

    /**
     * 继续导入
     */
    protected handlecontinue() {
        this.setState({
            mainScreen: true,
            importOrganize: false,
        })
    }

    /**
     * 导入成功确定
     */
    protected handleImportSuccess() {
        this.setState({
            importOrganize: false,
        })
        this.props.onSuccess();
    }

    /**
     * 下载已导出的文件
     */
    protected handleDownloadFile() {
        this.handleCancelOperation()
    }

    /**
     * 取消弹框内的操作
     */
    protected handleCancelOperation() {
        this.setState({
            exportOrganize: false,
            importOrganize: false,
            mainScreen: false,
        })
        if (this.state.import) {
            this.props.onSuccess();
        } else {
            this.props.onComplete();
        }
    }
}