import React from 'react';
import { noop } from 'lodash';

export default class RadioBoxBase extends React.Component<UI.RadioBox.Props, any> {
    static defaultProps = {
        onChange: noop,

        onCheck: noop,

        onUncheck: noop,
    }

    state: UI.RadioBox.State = {
        checked: this.props.checked,
    }

    static getDerivedStateFromProps({ checked }, prevState) {
        if(checked !== prevState.checked) {
            return {
                checked,
            }
        }
        return null
    }

    changeHandler(event) {
        if (!this.props.disabled) {
            const checked = event.target.checked;

            if (checked) {
                this.props.onCheck(this.props.value);
            } else {
                this.props.onUncheck(this.props.value);
            }

            if (this.state.checked !== checked) {
                this.props.onChange(checked, this.props.value);
                this.setState({
                    checked,
                })
            }
        }
    }
}