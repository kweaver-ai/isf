import React from 'react';

export default class TriggerPopMenuBase extends React.Component<UI.TriggerPopMenu.Props, UI.TriggerPopMenu.State> {
    static defaultProps = {
        timeout: 0,
        onRequestCloseWhenClick: (close) => { },
        numberOfChars: 20,
    }

    state = {
        clickStatus: false,
    }

    // 延迟定时器
    timeout: number | null = null;

    protected handleClickBtn() {
        this.setState({
            clickStatus: true,
        })
    }

    protected handleRequestCloseWhenBlur(close) {
        if (this.timeout) {
            clearTimeout(this.timeout)
        }

        this.timeout = setTimeout(() => {
            this.setState({
                clickStatus: false,
            })

        }, this.props.timeout)

        this.props.onRequestCloseWhenBlur(close)

    }

    protected handleRequestCloseWhenClick(close, open) {
        this.setState({
            clickStatus: open,
        })
        this.props.onRequestCloseWhenClick(close)
    }
}