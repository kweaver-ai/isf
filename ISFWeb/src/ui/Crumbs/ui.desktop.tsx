import React from 'react';
import { last, dropRight } from 'lodash';
import InlineButton from '../InlineButton/ui.desktop';
import LinkChip from '../LinkChip/ui.desktop';
import UIIcon from '../UIIcon/ui.desktop';
import CrumbsBase from './ui.base';
import styles from './styles.desktop';

export default class Crumbs extends CrumbsBase {
    render() {
        return (
            <div className={styles['container']}>
                <div className={styles['back']}>
                    <InlineButton
                        code={'\uf0a0'}
                        color={'#000'}
                        size={24}
                        className={!this.props.backDisabled && styles['back-btn']}
                        disabled={this.props.backDisabled}
                        onClick={() => this.back(last(dropRight(this.state.crumbs)))}
                    />
                </div>
                <ol className={styles['crumbs']}>
                    {
                        this.state.crumbs.map((crumb, i) => (
                            <li
                                className={styles['crumb']}
                                key={i}
                            >
                                <div className={styles['crumbWrap']}>
                                    {
                                        this.state.crumbs.length > 1 ?
                                            (
                                                <UIIcon
                                                    code={'\uf04e'}
                                                    size={20}
                                                    color={'#0000008c'}
                                                    className={styles['joiner']}
                                                />

                                            )
                                            : null
                                    }
                                    <LinkChip className={styles['link']} onClick={this.clickCrumb.bind(this, crumb)}>
                                        {
                                            this.props.formatter(crumb)
                                        }
                                    </LinkChip>
                                </div>
                            </li>
                        ))
                    }
                </ol>
            </div>
        )
    }
}