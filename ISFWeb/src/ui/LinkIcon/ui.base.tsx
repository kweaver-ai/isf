import React from 'react';
import { noop } from 'lodash';

export default class LinkIconBase extends React.Component<any, any> {
    static defaultProps = {
        onClick: noop,
    }

    clickHandler(e) {
        if (!this.props.disabled) {
            this.props.onClick(e);
        }
    }
}