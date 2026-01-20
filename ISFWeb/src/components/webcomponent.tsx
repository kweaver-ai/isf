import * as React from 'react';

export default class WebComponent<P, T> extends React.PureComponent<P, T> {
    /**
     * 在组件销毁后设置state，防止内存泄漏
     */
    componentWillUnmount() {
        this.setState = () => {
            return;
        }
    }
}