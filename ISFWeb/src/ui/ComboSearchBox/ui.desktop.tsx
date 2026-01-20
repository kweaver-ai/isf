import React from 'react'
import { trim } from 'lodash'
import classnames from 'classnames'
import { UIIcon, Chip, TextInput } from '../ui.desktop'
import { PopMenu } from '@/sweet-ui'
import ComboSearchBoxBase from './ui.base'
import styles from './styles.desktop'

export default class ComboSearchBox extends ComboSearchBoxBase {
    render() {
        const { keys, renderOption, renderComboItem, placeholder, className } = this.props
        const { value, searchValue, searchAnchor, isSearchFocus, isSearchMenu, filterInputValue } = this.state
        return (
            <div role={this.props.role} className={classnames(styles['search'], className)}>
                <div
                    className={classnames(styles['search-container'], { [styles['search-foucs']]: isSearchFocus })}
                    style={this.props.style}
                    onFocus={this.handleSearchBoxFocus.bind(this)}
                    onBlur={this.handleSearchBoxBlur.bind(this)}
                >

                    <UIIcon
                        className={(styles['search-icon'])}
                        code={'\uf01E'}
                        size={14}
                    />

                    <div className={classnames(styles['search-content'])}>
                        <div className={styles['search-filters']}>
                            {
                                [...searchValue].reverse().map((item, index) =>
                                    <div key={index} className={styles['search-filter-box']}>
                                        {
                                            index === 0 ?
                                                <input
                                                    className={styles['search-filter-input']}
                                                    type="text"
                                                    ref={(firstFilterInput) => this.firstFilterInput = firstFilterInput}
                                                    value={filterInputValue}
                                                    onChange={(e) => { this.handleStopInput(e); }}
                                                />
                                                :
                                                null
                                        }

                                        <Chip
                                            removeHandler={(e) => { this.handleItemDelete(e, item, index) }}
                                            className={styles['search-filter-item']}
                                            actionClassName={styles['search-delete']}
                                        >
                                            {
                                                keys && keys.length !== 0 ? renderComboItem(item.key, item.value) : item.value
                                            }
                                        </Chip>

                                        <input
                                            className={styles['search-filter-input']}
                                            type="text"
                                            value={filterInputValue}
                                            ref={(ref) => this.searchFilterInput[index] = ref}
                                            onKeyDown={(e) => { this.handleSearchFilterDelete(e, index); }}
                                            onChange={(e) => { this.handleStopInput(e); }}
                                        />
                                    </div>,
                                )
                            }
                        </div>

                        <TextInput
                            className={styles['search-input']}
                            type="text"
                            value={value}
                            placeholder={searchValue.length === 0 ? placeholder : ''}
                            onChange={this.handleSearchInputChange.bind(this)}
                            onKeyDown={this.handleInputKeyDown.bind(this)}
                            ref={(searchInput) => this.searchInput = searchInput}
                        />

                        {
                            (value || searchValue.length !== 0) ?
                                <UIIcon
                                    className={(styles['empty-icon'])}
                                    code={'\uf013'}
                                    size={16}
                                    onClick={this.handleTotalDelete.bind(this)}
                                />
                                :
                                <span className={(styles['blank-icon'])}></span>
                        }
                    </div>
                </div >

                {
                    keys && keys.length !== 0
                        ?
                        <PopMenu
                            anchor={searchAnchor}
                            anchorOrigin={['right', 'bottom']}
                            alignOrigin={['right', 'top']}
                            className={classnames(styles['search-menu'])}
                            freeze={false}
                            open={trim(value) && isSearchMenu}
                            element={this.props.element}
                        >
                            {
                                keys.map((key) => {
                                    return <PopMenu.Item
                                        key={key}
                                        className={classnames(styles['search-item'])}
                                        onClick={() => { this.handleSearchItemClick(key, trim(value)) }}
                                    >
                                        {
                                            renderOption(key, trim(value))
                                        }
                                    </PopMenu.Item>
                                })
                            }
                        </PopMenu> : null
                }

            </div >
        )
    }
}