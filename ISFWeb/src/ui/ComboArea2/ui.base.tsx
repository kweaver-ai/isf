import React from 'react';
import { includes, noop } from 'lodash';
import { mapKeyCode, isBrowser, Browser } from '@/util/browser';

export default class ComboArea2Base extends React.PureComponent<UI.ComboArea2.Props, any> {
    static defaultProps = {
        value: [],

        readOnly: false,

        uneditable: false,

        disabled: false,

        minHeight: 50,

        maxHeight: 100,

        onChange: noop,

        placeholder: '',

        spliter: [],

        validator: (val) => true,

    }

    props: UI.ComboArea2.Props;

    state = {
        placeholder: this.props.placeholder,
    }

    isIE = false;

    componentDidMount() {
        this.isIE = isBrowser({ app: Browser.MSIE });
        const { children } = this.props

        this.setState({
            placeholder: this.props.placeholder,
        })
        this.setState({
            focus: true,
        })
        this.refs.input && this.refs.input.focus();
    }

    static getDerivedStateFromProps({ placeholder }, prevState) {
        if(placeholder !== prevState.placeholder) {
            return {
                placeholder,
            }
        }
        return null
    }

    protected focusInput() {
        this.refs.input && this.refs.input.focus();
        this.setState({
            focus: true,
            placeholder: this.isIE ? '' : this.props.placeholder,
        })
    }

    protected blurInput() {
        const input = this.refs.input.state.value;
        this.setState({
            focus: false,
            placeholder: this.props.placeholder,
        })
        if (this.props.validator(input)) {
            this.props.addChip(input);
        }
        this.clearInput();
    }

    protected keyDownHandler(e) {
        const input = this.refs.input.state.value;
        // 增加Chip
        if (includes(this.props.spliter, mapKeyCode(e.keyCode)) || e.keyCode === 13) {

            // const chips = input.split(new RegExp(this.props.spliter.join('|')))
            //     .filter(chip => this.props.validator(chip));
            // this.props.addChip(chips);
            // this.clearInput();

            if (this.props.validator(input)) {
                this.props.addChip(input);
                this.clearInput();
            } else {
                this.clearInput();
            }

            e.preventDefault ? e.preventDefault() : (e.returnValue = false);

        }
        // 删除输入或Chip
        else if (e.keyCode === 8) {
            if (!input) {
                this.props.removeChip(input);
            }
        }

        // 如果已经达到容量上限，禁止任何输入
        if(React.Children.toArray(this.props.children).length >= 30) {
            this.clearInput();

        }
        this.props.onChange(this.refs.input.value());
    }

    clearInput() {
        this.refs.input.clear();
    }
}