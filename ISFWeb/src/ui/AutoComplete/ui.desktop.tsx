import React from 'react';
import Text from '../Text/ui.desktop';
import LazyLoader from '../LazyLoader/ui.desktop'
import SearchBox from '../SearchBox/ui.desktop';
import Locator from '../Locator/ui.desktop';
import AutoCompleteBase from './ui.base';
import styles from './styles.desktop';

export default class AutoComplete extends AutoCompleteBase {
    render() {
        const { keyDown, selectIndex } = this.state;
        const { lazyLoader, role } = this.props

        return (
            <div
                role={role}
                ref={(container) => this.container = container}
                className={styles['container']}
                style={{ width: this.props.width }}
            >
                <SearchBox
                    ref={(searchBox) => this.searchBox = searchBox}
                    value={this.state.value}
                    style={this.props.style}
                    width={this.props.width}
                    className={this.props.className}
                    disabled={this.props.disabled}
                    icon={this.props.icon}
                    delay={this.props.delay}
                    autoFocus={this.props.autoFocus}
                    placeholder={this.props.placeholder}
                    validator={this.props.validator}
                    loader={this.props.loader.bind(this)}
                    onChange={this.handleChange.bind(this)}
                    onFetch={this.handleFetch.bind(this)}
                    onLoad={this.handleLoad.bind(this)}
                    onLoadFailed={this.props.loaderFailed.bind(this)}
                    onFocus={this.handleFocus.bind(this)}
                    onBlur={this.handleBlur.bind(this)}
                    onEnter={this.handleEnter.bind(this)}
                    onKeyDown={this.handleKeyDown.bind(this)}
                />
                {
                    this.state.active && this.state.status !== AutoCompleteBase.Status.FETCHING ?
                        React.Children.count(this.props.children) ?
                            (
                                <Locator>
                                    <div
                                        onMouseDown={this.preventHideResults.bind(this)}
                                        className={styles['results']}
                                        style={{ width: this.container.offsetWidth }}
                                    >
                                        {
                                            lazyLoader ?
                                                <LazyLoader
                                                    ref={(lazyloadListDom) => this.lazyloadListDom = lazyloadListDom}
                                                    limit={lazyLoader.limit}
                                                    trigger={lazyLoader.trigger}
                                                    onChange={lazyLoader.onChange}
                                                >
                                                    {
                                                        React.Children.map(this.props.children, (child) => React.cloneElement(child, {
                                                            selectIndex,
                                                            keyDown,
                                                            onSelectionChange: this.handleSelectionChange.bind(this),
                                                        }))
                                                    }
                                                </LazyLoader>
                                                :
                                                React.Children.map(this.props.children, (child) => React.cloneElement(child, {
                                                    selectIndex: selectIndex,
                                                    keyDown: keyDown,
                                                    onSelectionChange: this.handleSelectionChange.bind(this),
                                                }))
                                        }
                                    </div>
                                </Locator>
                            ) :
                            this.state.value !== '' && this.props.missingMessage ?
                                (
                                    <Locator>
                                        <div
                                            style={{ width: this.container.offsetWidth }}
                                            className={styles['missing-message']}
                                            onClick={() => this.toggleActive(false)}
                                        >
                                            <div className={styles['padding']}>
                                                <Text>
                                                    {
                                                        this.props.missingMessage
                                                    }
                                                </Text>
                                            </div>
                                        </div>
                                    </Locator>
                                ) : null
                        : null
                }
            </div>
        )
    }
}
