import React from 'react';
import { noop, isFunction } from 'lodash';
import { Direction } from './helper';

export default class WizardBase extends React.Component<UI.Wizard.Props, any> {
    static Direction = Direction;

    static defaultProps = {
        title: '',

        onFinish: noop,

        onCancel: noop,
    }

    state = {
        activeIndex: 0, // 激活页下标
    };

    /**
     * 导航到上一步／下一步
     * @param direction 向上一步／下一步前进
     */
    protected async navigate(direction: Direction) {
        if (this.props.children) {
            const currentStep = this.props.children[this.state.activeIndex];
            const beforeLeaveEvent = currentStep.props.onBeforeLeave;
            const leaveEvent = currentStep.props.onLeave;
            const nextStep = this.props.children[this.state.activeIndex + direction];

            let shouldLeave = true;

            if (direction === Direction.FORWARD && isFunction(beforeLeaveEvent)) {
                shouldLeave = await beforeLeaveEvent();
            }

            if (shouldLeave !== false) {
                isFunction(leaveEvent) && leaveEvent();

                this.setState({
                    activeIndex: this.state.activeIndex + direction,
                }, () => {

                    if (nextStep) {
                        const enterEvent = nextStep.props.onEnter;

                        isFunction(enterEvent) && enterEvent();
                    }
                })
            }
        }
    }

    /**
     * 点击完成时触发
     */
    protected onFinish() {
        isFunction(this.props.onFinish) && this.props.onFinish()
    }

    /**
     * 点击取消时触发
     */
    protected onCancel() {
        isFunction(this.props.onCancel) && this.props.onCancel()
    }
}