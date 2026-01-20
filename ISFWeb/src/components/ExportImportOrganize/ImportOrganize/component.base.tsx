import * as React from 'react'
import { noop } from 'lodash';
import { concatWithToken } from '@/core/token';
import { usrmGetErrorInfos } from '@/core/thrift/sharemgnt/sharemgnt';
import __ from './locale'
import AppConfigContext from '@/core/context/AppConfigContext';

interface State {
    /**
     * 列表数据
     */
    defaultList: ReadonlyArray<any>;
    /**
     * 当前页面page
     */
    page: number;
    /**
     * 列表数据总数
     */
    count: number;
}

interface Props {
    /**
     * 导入失败个数
     */
    failNum: number;
    /**
     * 导入成功个数
     */
    successNum: number;
    /**
     * 取消
     */
    onCancel: () => any;
    /**
     * 继续导入
     */
    onContinue: () => any;
    /**
     * 导入成功
     */
    onImportSuccess: () => any;
}

export default class ImportOrganizeBase extends React.PureComponent<Props, State> {
    static contextType = AppConfigContext;

    static defautProps = {
        failNum: 0,

        successNum: 0,

        onCancel: noop,

        onContinue: noop,

        onImportSuccess: noop,
    }

    state = {

        defaultList: [],

        page: 0,

        count: this.props.failNum ? this.props.failNum : 0,

    }

    PageSize = 20; // 默认每页显示文档数

    componentDidMount() {
        const { failNum } = this.props;

        if (failNum !== 0) {
            this.getWatermarkDocByPage();
        }
    }

    /**
     * 翻页
     * @param page
     */
    protected handlePageChange(page: number) {
        this.setState({
            page: page - 1,
        }, this.getWatermarkDocByPage.bind(this))
    }

    /**
     * 翻页数据处理
     */
    private async getWatermarkDocByPage() {
        const defaultList = await usrmGetErrorInfos([this.PageSize * this.state.page, this.PageSize])

        this.setState({
            defaultList: defaultList,
        })
    }

    /**
     * 下载导入失败的记录
     * @param page
     */
    protected async downloadErrorList() {
        window.location.assign(concatWithToken(this.context?.prefix || '', this.context?.getToken?.(), `/isfweb/api/user/downloaduser/?taskId=0`))
    }
}