import React from 'react';
import classnames from 'classnames';
import { isEqual } from 'lodash';
import { ClassName } from '../helper';
import PopMenuItem from '../PopMenu.Item/ui.desktop';
import TriggerPopMenu from '../TriggerPopMenu/ui.desktop';
import { decorateText } from '@/util/formatters';
import SelectMenuBase from './ui.base';
import styles from './styles.desktop';

export default class SelectMenu extends SelectMenuBase {

    render() {

        let { label, className, candidateItems, numberOfChars = 20, freezable = false, disabled } = this.props;

        let { selectValue, hover } = this.state;

        return (
            <div className={classnames(styles['select-menu'], className)}>

                <TriggerPopMenu
                    className={styles['button-style']}
                    popMenuClassName={styles['condition-select-menu']}
                    title={this.props.customAttrs ? label + selectValue.name : selectValue.name}
                    label={this.props.customAttrs ? label + decorateText(selectValue.name, { limit: numberOfChars }) : decorateText(selectValue.name, { limit: numberOfChars })}
                    onRequestCloseWhenClick={(close) => close()}
                    timeout={150}
                    freezable={freezable}
                    disabled={disabled}
                    onRequestCloseWhenBlur={() => this.handleCloseMenuWhenBlur()}
                >
                    {
                        candidateItems.map((item, index) =>
                            <PopMenuItem
                                label={item.name}
                                key={index}
                                className={classnames(styles['condition-select-items'], { [styles['selected']]: isEqual(item, selectValue) && hover })}
                                labelClassName={classnames({ [ClassName.Color]: isEqual(item, selectValue) && hover })}
                                onClick={(e) => { this.handleClickCandidateItem(e, item) }}
                                onMouseEnter={() => { this.handleMouseEnter() }}
                            >
                            </PopMenuItem>,
                        )
                    }
                </TriggerPopMenu>
            </div >
        )
    }
}