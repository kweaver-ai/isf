import * as React from 'react';
import classNames from 'classnames';
import SearchDep from '../SearchDep/component.desktop';
import OrganizationTree from '../OrganizationTree/component.view';
import FlexBox from '@/ui/FlexBox/ui.desktop';
import session from '@/util/session';
import { NodeType } from '@/core/organization';
import OrganizationTreeBoxBase from './component.base';
import styles from './styles.view';
import __ from './locale';

export default class OrganizationTreeBox extends OrganizationTreeBoxBase {
    render() {
        const {
            isShowSearch,
            selectType,
            searchWidth,
            onSelectionChange,
            isShowUndistributed,
        } = this.props;
        return (
            <div className={styles['container']}>
                <FlexBox>
                    <FlexBox.Item align="top" className={styles['flexbox-item']}>
                        <div className={styles['selecor']}>
                            {
                                isShowSearch ?
                                    <SearchDep
                                        onSelectDep={(value) => { this.selectDep(value) }}
                                        width={searchWidth}
                                        selectType={selectType}
                                        placeholder={__('查找用户或部门')}
                                        userid={session.get('isf.userid')}
                                    />
                                    : null
                            }
                            <div className={classNames(
                                styles['tree-wrapper'],
                                {
                                    [styles['show-search']]: isShowSearch,
                                },
                            )}>
                                <OrganizationTree
                                    userid={session.get('isf.userid')}
                                    isShowUndistributed={isShowUndistributed}
                                    selectType={[NodeType.DEPARTMENT, NodeType.ORGANIZATION, NodeType.USER]}
                                    onSelectionChange={(value) => onSelectionChange(value)}
                                />
                            </div>
                        </div>
                    </FlexBox.Item>
                </FlexBox>
            </div>
        )
    }
}