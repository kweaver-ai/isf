import * as React from 'react';
import { noop, uniqBy } from 'lodash';
import { Message } from '@/sweet-ui';
import { timer } from '@/util/timer';
import { ErrorCode } from '@/core/thrift/sharemgnt/errcode';
import { concatWithToken } from '@/core/token';
import { usrmExportBatchUsers } from '@/core/thrift/sharemgnt/sharemgnt';
import { usrmGetExportBatchUsersTaskStatus, getProgress } from '@/core/thrift/sharemgnt/sharemgnt';
import __ from './locale';
import AppConfigContext from '@/core/context/AppConfigContext';

interface State {
    /**
     * 已选择的组织或部门
     */
    selectedMember: ReadonlyArray<any>;

    /**
     * 导出状态
     */
    exportStatus: boolean;

    /**
     * 进度条
     */
    progress: number;

    /**
     * taskid任务id
     */
    taskid: string;

    /**
     * 不存在的部门
     */
    unExistMember: ReadonlyArray<any>;
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
     * 下载导出的组织信息
     */
    onDownloadFile: () => any;
}

export default class ExportOrganizeBase extends React.PureComponent<Props, State> {
    static contextType = AppConfigContext;

    static defautProps = {
        userid: '',

        onCancel: noop,

        onDownloadFile: noop,
    }

    state = {
        selectedMember: [],

        exportStatus: false,

        progress: 0,

        unExistMember: [],
    }

    /**
     * 获取导出进度定时器
     */
    stopTimer = noop;

    componentWillUnmount() {
        // 撤销定时器
        this.stopTimer();
    }
    /**
     * 选择成员
     */
    protected addMember = (value: any) => {
        this.setState({
            selectedMember: uniqBy([...this.state.selectedMember, value], 'id'),
        })
    }

    /**
     * 清空已选用户
     */
    protected clearSelectDep = () => {
        this.setState({
            selectedMember: [],
        })
    }

    /**
     * 导出excel列表
     */
    protected onSaveMember = async () => {
        const departmentIds = this.state.selectedMember.map(({ id }) => id)
        const { userid } = this.props;

        try {
            const taskid = await usrmExportBatchUsers([departmentIds, userid])
            this.setState({
                taskid,
                exportStatus: true,
            })
            this.getProgress(taskid);
        } catch ({ error }) {
            if (error && error.errDetail) {
                const errDetail = JSON.parse(error.errDetail).unexist_depart_ids;
                let { selectedMember, unExistMember } = this.state;

                selectedMember.map((item) => {
                    for (let i in errDetail) {
                        if (errDetail[i] === item.id) {
                            unExistMember = [...unExistMember, item.displayName]
                        }
                    }
                })
                this.setState({
                    unExistMember,
                })
            } else if (error && error.errMsg) {
                if (error.errID === ErrorCode.ExportTaskExist) {
                    await Message.error({
                        message: __('另一个批量导出用户任务正在被执行，请稍后再试'),
                    })
                } else if (error.errID === ErrorCode.ImportTaskExist) {
                    await Message.error({
                        message: __('另一个批量导入用户任务正在被执行，请稍后再试'),
                    })
                } else {
                    await Message.error({
                        message: error.errMsg,
                    })
                }
            }
        }
    }

    /**
     * 删除已选成员
     */
    protected deleteSelectDep = (value: any) => {
        this.setState({
            selectedMember: this.state.selectedMember.filter((user) => user.id !== value.id),
        })
    }

    /**
     * 获取进度条进度
     */
    private getProgress = async (taskid) => {

        this.stopTimer = timer(async () => {

            const { successNum, failNum, totalNum } = await getProgress();

            if (successNum + failNum === totalNum) {
                const status = await usrmGetExportBatchUsersTaskStatus([taskid]);
                if (status) {
                    this.setState({
                        progress: totalNum === 0 ? 1 : (successNum + failNum) / totalNum,
                    })
                    this.stopTimer();
                }
            } else {
                this.setState({
                    progress: (successNum + failNum) / totalNum,
                })
            }
        }, 300)
    }

    /**
     * 下载带有用户信息的exel表
     */
    protected async downloadFile() {
        const { onDownloadFile } = this.props;
        const { taskid } = this.state;

        try {
            const status = await usrmGetExportBatchUsersTaskStatus([taskid])

            if (status) {
                window.location.assign(concatWithToken(this.context?.prefix || '', this.context?.getToken?.(), `/isfweb/api/user/downloaduser/?taskId=${taskid}`))


                onDownloadFile()
            }
        } catch ({ error }) {
            if (error && error.errMsg) {
                await Message.error({
                    message: error.errMsg,
                })
            }
        }
    }

    /**
     * 关闭导出进度弹框
     */
    protected async cancelExport() {
        const { progress } = this.state;
        const { onCancel } = this.props;

        if (progress === 1) {
            onCancel();
        } else {
            const confirm = await Message.confirm(
                {
                    message: __('正在导出选中的用户信息，若关闭窗口，则会继续完成导出，您确定要执行此操作吗？'),
                },
            )
            if (confirm) {
                onCancel();
            }
        }
    }

    /**
     * 确定提示错误弹窗
    */
    protected closeErrorMessage = () => {
        this.setState({
            unExistMember: [],
        })
    }
}