import React from 'react';
import classnames from 'classnames';
import Dialog from '../Dialog2/ui.desktop';
import Panel from '../Panel/ui.desktop';
import UIIcon from '../UIIcon/ui.desktop';
import WizardStep from '../Wizard.Step/ui.desktop';
import WizardBase from './ui.base';
import __ from './locale';
import styles from './styles.desktop';

export default class Wizard extends WizardBase {
    static Step = WizardStep;

    render() {
        const { children, role } = this.props
        const activeChildren = React.Children.toArray(children)[this.state.activeIndex]

        return (
            <Dialog
                role={role}
                title={this.props.title}
                onClose={this.props.onCancel}
            >
                <Panel>
                    <Panel.Main>
                        <div className={styles['crumbs']}>
                            {
                                React.Children.map(children, (child, index) => (
                                    <span
                                        className={classnames(
                                            [styles['crumb']],
                                            { [styles['active']]: this.state.activeIndex === index },
                                            { [styles['configured']]: this.state.activeIndex > index },
                                        )}
                                        title={child.props.title}
                                        key={index}
                                    >
                                        {
                                            child.props.title
                                        }
                                        {
                                            React.Children.count(children) > index + 1 ?
                                                <UIIcon
                                                    className={styles['icon']}
                                                    size="16"
                                                    code={'\uf002'}
                                                    color={'#999'}
                                                />
                                                : null
                                        }
                                    </span>
                                ))
                            }
                        </div>
                        <div>
                            {
                                React.Children.map(children, (child, index) => (
                                    React.cloneElement(child, {
                                        active: this.state.activeIndex === index,
                                    })
                                ))
                            }
                        </div>
                    </Panel.Main>
                    <Panel.Footer>
                        {
                            this.state.activeIndex !== 0 ?
                                <Panel.Button onClick={this.navigate.bind(this, WizardBase.Direction.BACKWARD)}>{__('上一步')}</Panel.Button> : null
                        }
                        {
                            this.state.activeIndex === React.Children.count(children) - 1 ?
                                <Panel.Button
                                    type="submit"
                                    onClick={this.onFinish.bind(this)}
                                    disabled={activeChildren && activeChildren.props.disabled}
                                >
                                    {__('完成')}
                                </Panel.Button> :
                                <Panel.Button
                                    onClick={this.navigate.bind(this, WizardBase.Direction.FORWARD)}
                                    disabled={activeChildren && activeChildren.props.disabled}
                                >
                                    {__('下一步')}
                                </Panel.Button>
                        }
                        <Panel.Button onClick={this.onCancel.bind(this)}>{__('取消')}</Panel.Button>
                    </Panel.Footer>
                </Panel>
            </Dialog>
        )
    }
}