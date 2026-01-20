import * as React from 'react';
import FlexBox from '@/ui/FlexBox/ui.desktop';
import GroupGrid from './GroupGrid/component.view';
import MemberGrid from './MemberGrid/component.view';
import UserGroupBase from './component.base';
import styles from './styles.view';

export default class UserGroup extends UserGroupBase {
    render() {
        const { selectedGroup, groupStatus } = this.state

        return (
            <div  className={styles["user-group"]}>
                <FlexBox role={'ui-flexbox'}>
                    <FlexBox.Item role={'ui-flexbox.item'} className={styles['flex-box']}>
                        <GroupGrid
                            ref={(groupGrid) => this.groupGrid = groupGrid}
                            onRequestSelectGroup={this.selectGroup}
                            onRequestChangeStatus={this.changeGroupStatus}
                        />
                    </FlexBox.Item>
                    <FlexBox.Item role={'ui-flexbox.item'} width={10}></FlexBox.Item>
                    <FlexBox.Item role={'ui-flexbox.item'} className={styles['flex-box']}>
                        <MemberGrid
                            groupStatus={groupStatus}
                            selectedGroup={selectedGroup}
                            onRequestUpdateGroup={this.updateGroup}
                        />
                    </FlexBox.Item>
                </FlexBox>
            </div>
        )
    }
}