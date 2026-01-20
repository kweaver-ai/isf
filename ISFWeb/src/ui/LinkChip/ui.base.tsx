import React from 'react';
import { noop } from 'lodash';

export default class LinkChipBase extends React.Component<any, any> {
    static defaultProps = {
        onClick: noop,

        onTouchStart: noop,

        onTouchEnd: noop,

        onDoubleClick: noop,
    }

    clickHandler(e) {
        if (!this.props.disabled) {
            this.props.onClick(e);
        }
    }

    touchStartHandler(e) {
        if (!this.props.disabled) {
            this.props.onTouchStart(e);
        }
    }

    touchEndHandler(e) {
        if (!this.props.disabled) {
            this.props.onTouchEnd(e);
        }
    }

    doubleClickHandler(e) {
        if (!this.props.disabled) {
            this.props.onDoubleClick(e);
        }
    }
}