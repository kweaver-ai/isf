import React from 'react'
import PopOverBase from './ui.base'

export default class PopOver extends PopOverBase {
    render() {
        const { trigger, triggerEvent, role } = this.props
        if (trigger) {
            switch (triggerEvent) {
                case 'click':
                    return React.cloneElement(trigger, { onClick: (e) => this.handleTriggerClick(e, trigger.props), role })
                case 'mouseover':
                    return React.cloneElement(
                        trigger,
                        {
                            onMouseEnter: (e) => this.handleTriggerMouseEnter(e, trigger.props),
                            onMouseLeave: (e) => this.handleTriggerMouseLeave(e, trigger.props),
                            role,
                        },
                    )
                default:
                    return React.cloneElement(trigger, { role })
            }
        }
        return null
    }
}