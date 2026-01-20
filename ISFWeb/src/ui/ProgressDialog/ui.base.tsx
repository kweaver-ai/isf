import React from 'react';
import { noop } from 'lodash';
import { chain } from '@/util/accessor';

interface Props {
    data: Array<any>;                     // 需要操作的对象数组
    title: string;                        // dialog的标题
    detailTemplate(item: any): string;    // 显示的文字详情
    loader(item: any): Promise<any>;      // 队列操作
    onSuccess(): void;                    // 队列操作成功
    onError(err: any, item: any): void;   // 队列操作失败
    onSingleSuccess(item: any): void;     // 单个成功
}

interface State {
    progress: number;   // 进度值
    item: any;          // 正在进行中的对象
}

export default class ProgressDialogBase extends React.Component<any, any> {
    static defaultProps = {
        data: null,
        detailTemplate: (item) => item.name,
        loader: noop,
        onSuccess: noop,
        onError: noop,
        onSingleSuccess: noop,
    }

    props: Props;

    state: State = {
        progress: 0,
        item: (this.props.data && this.props.data.length) ? this.props.data[0] : null,
    }

    interrupt: boolean = false;    // 队列进行中是否中断队列

    componentDidMount() {
        this.startQueue();
    }

    /**
     * 开始队列操作
     */
    startQueue(data: Array<any> = this.props.data): void {
        const queue = chain((item, index, array) => {
            this.setState({
                item,
            })

            if (!this.interrupt) {
                return this.props.loader(item).then(() => {
                    this.props.onSingleSuccess(item);
                    this.setState({
                        progress: ((index + 1) / array.length).toFixed(2),
                    })
                })
            }
            return Promise.resolve(1);
        })

        queue(data).then(() => {
            // 成功
            this.props.onSuccess();
        }, (err) => {
            // 失败
            this.props.onError(err, this.state.item)
        })
    }

    /**
     * 中断队列
     */
    handleCancel() {
        this.interrupt = true;
    }
}