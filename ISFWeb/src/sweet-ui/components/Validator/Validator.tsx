import React from 'react';
import { isFunction } from 'lodash';
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import ValidateTip from '../ValidateTip';
import Placement from '../ValidateTip/ValidateTip';

function createChainedFunction() {
    const args = [].slice.call(arguments, 0);
    if (args.length === 1) {
        return args[0];
    }

    return function chainedFunction() {
        for (let i = 0; i < args.length; i++) {
            if (args[i] && args[i].apply) {
                args[i].apply(this, arguments);
            }
        }
    };
}

interface ValidatorProps {
    // todo 错误提示 支持气泡 / 红色文字提示
    /**
     * 气泡提示的位置
     */
    placement?: Placement;

    /**
     * 校验输入值的时机
     */
    validateTrigger?: 'onChange' | 'onBlur';

    /**
     * 验证规则
     */
    rules: Array<{
        /**
         * 提示信息
         */
        message: string;

        /**
         * 输入不允许为空
         */
        required?: boolean;

        /**
         * 自定义验证方法
         */
        validator?: (val: string) => boolean;
    }>;

    /**
     * 执行校验后的回调事件
     */
    afterValidate: (event: SweetUIEvent<boolean>) => void;
}

interface ValidatorState {
    validateMessage: string | null;
}

export default class Validator extends React.Component<ValidatorProps, ValidatorState> {
    static defaultProps = {
        validateTrigger: 'onChange',
    };

    state = {
        validateMessage: null,
    };

    public getUnVerifiedRule = (value: string) => {
        // todo 设置为public用于上层组件主动调用验证方法
        return this.props.rules.find(
            ({ required, validator, message }) =>
                (required && value === '') || (isFunction(validator) && !validator(value)),
        );
    };

    private handleValueChange = (event: SweetUIEvent<string>) => {
        const { detail } = event;
        if (this.props.validateTrigger === 'onChange') {
            const unVerifiedRule = this.getUnVerifiedRule(detail);

            // console.log(unVerifiedRule);

            this.setState({
                validateMessage: unVerifiedRule ? unVerifiedRule.message : null,
            });

            this.dispatchAfterValidateEvent(unVerifiedRule ? false : true);
        } else {
            if (this.state.validateMessage) {
                this.setState({
                    validateMessage: null,
                });
            }
        }
    };

    private handleBlur = (value: any) => {
        const unVerifiedRule = this.getUnVerifiedRule(String(value));

        this.setState({
            validateMessage: unVerifiedRule ? unVerifiedRule.message : null,
        });

        this.dispatchAfterValidateEvent(unVerifiedRule ? false : true);
    };

    private attachValidators = (children: any) => {
        if (children) {
            let child = children;
            if (React.isValidElement(child) && React.Children.only(child)) {
                let extraProps = {};
                if (this.props.validateTrigger === 'onChange') {
                    extraProps.onValueChange = createChainedFunction(child.props.onValueChange, this.handleValueChange);
                }

                if (this.props.validateTrigger === 'onBlur') {
                    extraProps.onBlur = createChainedFunction(child.props.onBlur, () =>
                        this.handleBlur(child.props.value),
                    );
                }
                child = React.cloneElement(child, {
                    ...extraProps,
                });
            }

            return child;
        }
    };

    dispatchAfterValidateEvent = createEventDispatcher(this.props.afterValidate);

    render() {
        const { validateMessage } = this.state;

        return (
            <ValidateTip
                placement={this.props.placement}
                content={validateMessage}
                visible={!!validateMessage}
                tipStatus={'error'}
            >
                {this.attachValidators(this.props.children)}
            </ValidateTip>
        );
    }
}
