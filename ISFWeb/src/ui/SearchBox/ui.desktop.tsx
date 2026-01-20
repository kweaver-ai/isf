import React from 'react';
import classnames from 'classnames';
import UIIcon from '../UIIcon/ui.desktop';
import Control from '../Control/ui.desktop';
import SearchInput from '../SearchInput/ui.desktop';
import SearchBoxBase from './ui.base';
import styles from './styles.desktop';

/**
 * 是否显示清空按钮
 */
const showClear = ({ disabled }, { value }): boolean => !disabled && value !== '';

export default class SearchBox extends SearchBoxBase {
    render() {
        const { role } = this.props

        return (
            <Control
                role={role}
                className={classnames(styles['searchbox'], this.props.className)}
                width={this.props.width}
                style={this.props.style}
                disabled={this.props.disabled}
                focus={this.state.focus}
            >
                {
                    this.props.icon ?
                        <div className={styles['icon']}>
                            <UIIcon size={14} color={'#c0c0c0'} code={this.props.icon} />
                        </div> :
                        null
                }
                <div className={classnames({ [styles['icon-indent']]: this.props.icon, [styles['clear-indent']]: showClear(this.props, this.state) })}>
                    <SearchInput
                        ref={(searchInput) => this.searchInput = searchInput}
                        value={this.state.value}
                        disabled={this.props.disabled}
                        autoFocus={this.props.autoFocus}
                        placeholder={this.props.placeholder}
                        validator={this.props.validator}
                        maxLength={this.props.maxLength}
                        loader={this.props.loader}
                        delay={this.props.delay}
                        onChange={this.handleChange.bind(this)}
                        onFetch={this.props.onFetch}
                        onLoad={this.props.onLoad}
                        onLoadFailed={this.props.onLoadFailed}
                        onFocus={this.handleFocus.bind(this)}
                        onBlur={this.handleBlur.bind(this)}
                        onClick={this.props.onClick && this.props.onClick.bind(this)}
                        onEnter={this.props.onEnter && this.props.onEnter.bind(this)}
                        onKeyDown={this.props.onKeyDown && this.props.onKeyDown.bind(this)}
                    />
                </div>
                {
                    showClear(this.props, this.state) ?
                        <div className={styles['clear']}>
                            <UIIcon
                                size={16}
                                code={'\uf013'}
                                className={styles['chip-x-icon']}
                                onClick={this.clearInput.bind(this)}
                            />
                        </div>
                        : null
                }
            </Control >
        )
    }
}