import React from 'react';
import classnames from 'classnames';
import TabsBase from './ui.base';
import TabsNavigator from '../Tabs.Navigator/ui.desktop';
import TabsTab from '../Tabs.Tab/ui.desktop';
import TabsMain from '../Tabs.Main/ui.desktop';
import TabsContent from '../Tabs.Content/ui.desktop';
import styles from './styles.desktop';

export default class Tabs extends TabsBase {
    static Navigator = TabsNavigator;

    static Tab = TabsTab;

    static Main = TabsMain;

    static Content = TabsContent;

    render() {
        return (
            <div role={this.props.role} className={classnames(styles['container'], this.props.className)} style={{ height: this.props.height || '100%' }} >
                {
                    this.createChildren(...this.props.children)
                }
            </div>
        )
    }
}