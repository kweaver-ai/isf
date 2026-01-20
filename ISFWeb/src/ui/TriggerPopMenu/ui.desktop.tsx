import React from 'react';
import classnames from 'classnames';
import { decorateText } from '@/util/formatters';
import { Trigger } from '@/sweet-ui';
import Button from '../Button/ui.desktop';
import Title from '../Title/ui.desktop';
import UIIcon from '../UIIcon/ui.desktop';
import TriggerPopMenuBase from './ui.base';
import styles from './styles.desktop';

export default class TriggerPopMenu extends TriggerPopMenuBase {
    render() {
        let {
            popMenuClassName,
            label,
            children,
            numberOfChars,
            title,
            disabled = false,
            className,
            freezable = false,
        } = this.props;
        let { clickStatus } = this.state;
        return !disabled ? (
            <Trigger
                anchorOrigin={['left', 'bottom']}
                alignOrigin={['left', 'top']}
                freeze={freezable}
                renderer={({ setPopupVisibleOnClick }) => (
                    <div
                        key={'TriggerPopMenu'}
                        className={classnames(styles['button-container'], className)}
                        onMouseDown={setPopupVisibleOnClick}
                        onClick={this.props.onClick}
                    >
                        <Title content={title || label}>
                            <Button className={classnames(styles['button-btn'], { [styles['clicked']]: clickStatus })}>
                                {numberOfChars ? decorateText(label, { limit: numberOfChars }) : label}
                                <UIIcon className={classnames(styles['expand-icon'])} code={'\uF04C'} size="16px" />
                            </Button>
                        </Title>
                    </div>
                )}
                triggerEvent={'click'}
                onPopupVisibleChange={(event) =>
                    this.setState({
                        clickStatus: event.detail,
                    })}
                onBeforePopupClose={this.props.onRequestCloseWhenBlur}
            >
                {({ close, open }) => (
                    <ul
                        className={classnames(styles['menu'], popMenuClassName)}
                        onClick={() => this.handleRequestCloseWhenClick(close, open)}
                    >
                        {children}
                    </ul>
                )}
            </Trigger>
        ) : (
            <div className={classnames(styles['button-container'])}>
                <Title content={title || label}>
                    <div className={classnames(styles['button-btn'], styles['button-disabled'])}>
                        {numberOfChars ? decorateText(label, { limit: numberOfChars }) : label}
                        <UIIcon className={classnames(styles['expand-icon'])} code={'\uF04C'} size="16px" />
                    </div>
                </Title>
            </div>
        );
    }
}