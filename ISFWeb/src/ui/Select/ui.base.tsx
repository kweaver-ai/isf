import React from 'react';
import { noop, isArray, reduce } from 'lodash';

export default class SelectBase extends React.PureComponent<UI.Select.Props, any> implements UI.Select.Component {
    static defaultProps = {
        onChange: noop,

        menu: {},
    }

    state: UI.Select.State = {
        value: this.props.value,
    }

    componentDidMount() {
        this.setSelected(this.props.children)
    }

    componentDidUpdate(prevProps, prevState) {
        if(this.props.children !== prevProps.children || this.state.value !== prevState.value) {
            this.setSelected(this.props.children);
        }
    }

    private setSelected(options) {
        const optionsList = isArray(options) ? options : [options];
        const selected = reduce(optionsList, (prev, option) => {
            const match = (option.props.value === this.props.value || option.props.selected) ? option : null;
            // 如果在Options中找到匹配项
            if (match) {
                // 如果之前已经有匹配，但是是由select属性匹配到的
                if (prev) {
                    // 此次匹配是value匹配，使用此次匹配
                    if (prev.props.selected) {
                        return match;
                    }
                    // 否则认为上一次的匹配是value匹配，抛弃此次匹配
                    else {
                        return prev;
                    }
                }
                // 如果之前没有匹配，使用此次匹配
                else {
                    return match;
                }
            } else {
                return prev;
            }
        }, null);

        if (selected) {
            this.setState({
                value: selected.props.value,
                text: selected.props.children,
            });
        }
    }

    private fireChangeEvent(value) {
        this.props.onChange(value);
    }

    /**
     * 选中选项
     */
    protected onOptionSelected({ value, text }) {
        this.setState({
            value,
            text,
            active: false,
        }, () => {
            this.fireChangeEvent(this.state.value);
        });
    }
}
