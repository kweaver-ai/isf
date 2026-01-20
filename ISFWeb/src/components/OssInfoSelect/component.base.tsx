import * as React from 'react'
import { map, noop, isEqual } from 'lodash';
import WebComponent from '../webcomponent';
import { Message } from '@/sweet-ui';
import { getErrorMessage } from '@/core/exception';
import { getObjectStorageInfoByApp } from '@/core/apis/console/ossgateway'
import { LocationType } from './helper';
import __ from './locale'

interface OssInfoSelectProps extends React.Props<void> {

    /**
     * 宽度
     */
    width?: number;

    /**
     * 类型
     */
    type: string;

    /**
     * 是否禁用
     */
    disabled?: boolean;

    /**
     * 选择的存储位置
     */
    ossInfo: Core.ShareMgnt.ncTUsrmOSSInfo;

    /**
     * 校验状态
     */
    validateState: number;

    /**
     * 气泡提示语
     */
    validateMessages: {
        [key: string]: string;
    };

    /**
     * 选择存储位置
     */
    onRequestSelectOssInfo: (ossInfo) => void;

    /**
     * 失焦事件
     */
    onBlur: () => void;
}

interface OssInfoSelectState {

    /**
     * 选择的存储位置
     */
    ossInfo: Core.ShareMgnt.ncTUsrmOSSInfo;

    /**
     * 存储位置列表
     */
    ossInfos: ReadonlyArray<Core.ShareMgnt.ncTUsrmOSSInfo>;
}

export default class OssInfoSelectBase extends WebComponent<OssInfoSelectProps, OssInfoSelectState> {

    static defaultProps = {
        ossInfo: null,
        disabled: false,
        validateState: 0,
        validateMessages: {},
        onRequestSelectOssInfo: noop,
        onBlur: noop,
    }

    state: OssInfoSelectState = {
        ossInfo: (this.props.ossInfo && this.props.ossInfo.ossId) ? this.props.ossInfo : { enabled: true, ossId: '', ossName: '' },
        ossInfos: [],
    }

    async componentDidMount() {
        await this.getOssInfo();
    }

    componentDidUpdate(prevProps, prevState) {
        const { ossInfo } = this.props;
        if (!isEqual(prevProps.ossInfo, ossInfo) || !isEqual(prevState.ossInfos, this.state.ossInfos)) {
            const { ossInfos } = this.state;

            if (ossInfo && ossInfo.ossId && !ossInfo.enabled) {
                let lastOssInfos = map(ossInfos, (oss) => {
                    if (ossInfo.ossId === oss.ossId) {
                        return {
                            ...ossInfo,
                            displayName: this.displayOssInfo(ossInfo),
                        };
                    } else {
                        return oss;
                    }
                })
                if (!ossInfos.some((oss) => oss.ossId === ossInfo.ossId)) {
                    lastOssInfos = [...lastOssInfos.slice(0, 1), { ...ossInfo, displayName: this.displayOssInfo(ossInfo) }, ...lastOssInfos.slice(1)];
                }
                this.setState({
                    ossInfos: lastOssInfos,
                })
            }
            // 编辑用户/部门/组织传入的ossInfo为{ossId: null, ossName: null,...}需做转换
            this.setState({
                ossInfo: ossInfo && ossInfo.hasOwnProperty('ossId') && !ossInfo.ossId ?
                    { enabled: true, ossId: '', ossName: '' }
                    : ossInfo,
            })

        }
    }

    /**
     * 初始化存储位置
     */
    private async getOssInfo(): Promise<void> {
        const { ossInfo } = this.props;

        try {
            let ossDatas = await getObjectStorageInfoByApp({ app: 'as', enabled: true })

            let ossInfos = ossDatas.map((item) => {
                return {
                    enabled: item.enabled,
                    ossId: item.id,
                    ossName: item.name,
                }
            })

            // ossId为-1时，select的label为空(批量编辑文档库时)
            if (ossInfo && ossInfo.ossId && !ossInfos.some((oss) => oss.ossId === ossInfo.ossId)) {
                // 下拉列表添加已选择的，但被禁用的对象存储
                ossInfos = [ossInfo, ...ossInfos]
            }
            ossInfos = [{ enabled: true, ossId: '', ossName: '' }, ...ossInfos];
            this.setState({
                ossInfos: map(ossInfos, (oss) => {
                    return {
                        ...oss,
                        displayName: this.displayOssInfo(oss),
                    }
                }),
            })
        } catch (error) {
            let lastOssInfos = [{ enabled: true, ossId: '', ossName: '' }]

            if (ossInfo && ossInfo.ossId) {
                lastOssInfos = [...lastOssInfos, ossInfo]
            }

            this.setState({
                ossInfos: map(lastOssInfos, (oss) => ({
                    ...oss,
                    displayName: this.displayOssInfo(oss),
                })),
            })

            if (error && error.errID) {
                Message.info({ message: getErrorMessage(error.errID) })
            }
        }

    }

    /**
     * 修改存储位置
     */
    protected updateSelectedOss = (ossId: string): void => {
        this.setState({
            ossInfo: this.state.ossInfos.find((oss) => oss.ossId === ossId),
        }, () => {
            this.props.onRequestSelectOssInfo(this.state.ossInfo);
        })
    }

    /**
     * 格式化对象存储展示的信息
     * @param ossInfo 对象存储
     */
    protected displayOssInfo(ossInfo) {
        const { ossId } = ossInfo;
        return !ossId ?
            (this.props.type === LocationType.DocLib ? __('未指定（跟随文件上传者的指定存储位置）') : __('未指定（使用默认存储）'))
            :
            ossInfo.ossName
    }
}