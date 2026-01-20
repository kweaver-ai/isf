import * as React from 'react'
import { noop } from 'lodash'
import * as WebUploader from '@/libs/webuploader';
import { ErrorCode } from '@/core/thrift/sharemgnt/errcode';
import { getProgress } from '@/core/thrift/sharemgnt/sharemgnt';
import { timer } from '@/util/timer';
import { getHeaders } from '@/core/token'
import { Message } from '@/sweet-ui';
import __ from './locale';
import AppConfigContext from '@/core/context/AppConfigContext';

interface State {
    /**
     * 上传的文件
     */
    packageFile: any;
    /**
     * 覆盖/同步同名用户
     */
    operationStatus: number;
    /**
     * 导入进度条
     */
    progress: number;
    /**
     * 导入成功数
     */
    successNum: number;
    /**
     * 导入失败数
     */
    failNum: number;
    /**
     * 导入总数
     */
    totalNum: number;
    /**
     * 选择文件按钮是否可用
     */
    isDisable: boolean;
    /**
     * 获取进度接口是否调用
     */
    isGetprogress: boolean;
}

interface Props {
    /**
     * 管理员id
     */
    userid: string;
    /**
     * 取消
     */
    onCancel: () => any;
    /**
     * 导出操作
     */
    onExportItem: () => any;
    /**
     * 导入操作
     */
    onImportItem: () => any;
}

export default class MainScreenBase extends React.PureComponent<Props, State> {
    static  contextType = AppConfigContext;

    static defautProps = {
        userid: '',

        onCancel: noop,

        onExportItem: noop,

        onImportItem: noop,
    }

    state = {
        packageFile: null,

        operationStatus: null,

        progress: 0,

        successNum: 0,

        failNum: 0,

        totalNum: 0,

        isDisable: true,

        isGetprogress: true,
    }
    /**
     * WebUploader 实例
     */
    uploader = null;

    /**
     * 导入进度是否完成
     */
    progressStatus: boolean;
    /**
     * 导入操作是否开始
     */
    importStatus: boolean;
    /**
     * 获取导入进度定时器
     */
    stopTimer = noop;

    /**
     * 导入文件时导入操作的div
     */
    select: HTMLDivElement | undefined = undefined

    componentDidMount() {
        this.initUploader();
    }

    componentWillUnmount() {
        // 撤销定时器
        this.stopTimer();
        // 在组件销毁后设置state，防止内存泄漏
        this.setState = (state, callback) => {
            return;
        };
    }

    /**
     * 初始化上传组件
     */
    private initUploader() {
        // 获取应用主节点
        const self = this;
        self.uploader = new WebUploader.create({
            swf: '/res/libs/webuploader/Uploader.swf',
            server: `${self.context?.prefix || ''}/isfweb/api/user/importuser/`,
            auto: false,
            pick: {
                id: self.select,
                innerHTML: __('选择文件'),
                multiple: false,
            },
            accept: {
                title: '*.xlsx',
                extensions: 'xls,xlsx',
                mimeTypes: 'application/vnd.ms-excel,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
            },
            timeout: 0,
            onBeforeFileQueued: () => {
                self.progressStatus = false;
                self.importStatus = false;
                self.uploader.reset();
            },
            onFileQueued: (file: any) => {
                self.file = file
                self.setState({
                    packageFile: file,
                    isGetprogress: true,
                })
            },
            onUploadBeforeSend: (object, data, headers) => {
                const { userid } = self.props;
                const { operationStatus, packageFile } = self.state;
                Object.assign(headers, getHeaders(self.context?.getToken?.()).headers);

                data.fileName = packageFile ? packageFile.name : null;
                data.userCover = operationStatus === 1 ? true : false;
                data.responsiblePersonId = userid;
            },
            onUploadStart: () => {
                setTimeout(() => {
                    if (self.state.isGetprogress) {
                        self.getProgress();
                    }
                }, 2000)
                self.setState({
                    isDisable: false,
                })
            },
            onUploadAccept: async (object, ret: any) => {
                // 上传出错了
                if (ret.error) {
                    switch (ret.error.errID) {
                        case ErrorCode.UserCreateError:
                            Message.alert({ message: __('文件内容格式错误，用户组织列表不能新建、删除、修改') });
                            self.importStatus = false;
                            self.setState({
                                isGetprogress: false,
                            })
                            break;

                        case ErrorCode.ImportTaskExist:
                            Message.alert({ message: __('另一个批量导入用户任务正在被执行，请稍后再试') });
                            self.importStatus = false;
                            self.setState({
                                isGetprogress: false,
                            })
                            break;

                        case ErrorCode.ExportTaskExist:
                            Message.alert({ message: __('另一个批量导出用户任务正在被执行，请稍后再试') });
                            self.importStatus = false;
                            self.setState({
                                isGetprogress: false,
                            })
                            break;

                        default:
                            Message.alert({ message: ret.error.errMsg });
                            self.importStatus = false;
                            self.setState({
                                isGetprogress: false,
                            })
                            break;
                    }
                }

                setTimeout(() => {
                    self.stopTimer();
                    self.setState({
                        isDisable: true,
                        progress: 0,
                    })
                }, 2000)
            },
            onUploadSuccess: async function (file, response) {
                const { failNum, successNum } = self.state;

                self.importStatus = true
                if (self.progressStatus) {
                    self.props.onImportItem(failNum, successNum)

                }
            },
        });
    }

    /**
     * 取消导入操作
     */
    protected async onCancel() {
        const { progress, isGetprogress } = this.state;
        const { onCancel } = this.props;

        if (progress !== 0 && isGetprogress) {
            const confirm = await Message.confirm(
                {
                    message: __('正在导入选中的用户信息，若关闭窗口，则会继续完成导入，您确定要执行次操作吗'),
                },
            )
            if (confirm) {
                onCancel();
            }
        } else {
            onCancel();
        }
    }
    /**
     * 点击“导入”按钮
     */

    protected onImportItem = () => {
        const { packageFile } = this.state;

        if (packageFile.size / 1024 / 1024 > 25) {
            Message.error({
                message: __('文件大小不能超过25MB'),
            })
            return false;
        }
        this.upload();
    }

    /**
     * 导入上传文件
     */
    private upload() {
        const { isGetprogress } = this.state;

        if (isGetprogress) {
            this.uploader.upload();
        } else {
            this.uploader.retry();
            this.setState({
                isGetprogress: true,
            })
        }
    }

    /**
     * 导出操作
     */
    protected changeApprovalStatus(operationStatus: number) {
        this.setState({
            operationStatus,
        })
    }

    /**
    * 获取进度条进度
    */
    private getProgress = async () => {
        this.stopTimer = timer(async () => {
            if (!this.state.isDisable) {

                const { successNum, failNum, totalNum } = await getProgress();

                this.setState({
                    progress: totalNum === 0 && !this.importStatus ? 0 : (successNum + failNum) / totalNum,
                    successNum,
                    failNum,
                    totalNum,
                })
                if (successNum + failNum === totalNum) {
                    this.stopTimer();
                    this.progressStatus = true;
                    this.setState({
                        isDisable: true,
                    })
                    if (this.importStatus) {
                        this.props.onImportItem(failNum, successNum)
                    }
                }
            }
        }, 300)
    }
}