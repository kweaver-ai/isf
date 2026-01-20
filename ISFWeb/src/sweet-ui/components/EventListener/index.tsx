import React from 'react';
import { passiveOption } from './supports';

const defaultEventOptions = {
    capture: false,
    passive: false,
};

function mergeDefaultEventOptions(options) {
    return { ...defaultEventOptions, ...options };
}

function getEventListenerArgs(eventName, callback, options) {
    let args = [eventName, callback];
    args = [...args, passiveOption ? options : options.capture]

    return args;
}

function on(target, eventName, callback, options) {
    // eslint-disable-next-line prefer-spread
    target.addEventListener.apply(target, getEventListenerArgs(eventName, callback, options));
}

function off(target, eventName, callback, options) {
    // eslint-disable-next-line prefer-spread
    target.removeEventListener.apply(target, getEventListenerArgs(eventName, callback, options));
}

function forEachListener(props, iteratee) {
    const {
        children, // eslint-disable-line no-unused-vars
        target, // eslint-disable-line no-unused-vars
        ...eventProps
    } = props;

    Object.keys(eventProps).forEach((name) => {
        if (name.substring(0, 2) !== 'on') {
            return;
        }

        const prop = eventProps[name];
        const type = typeof prop;
        const isObject = type === 'object';
        const isFunction = type === 'function';

        if (!isObject && !isFunction) {
            return;
        }

        const capture = name.substr(-7).toLowerCase() === 'capture';
        let eventName = name.substring(2).toLowerCase();
        eventName = capture ? eventName.substring(0, eventName.length - 7) : eventName;

        if (isObject) {
            iteratee(eventName, prop.handler, prop.options);
        } else {
            iteratee(eventName, prop, mergeDefaultEventOptions({ capture }));
        }
    });
}

export function withOptions(handler, options) {
    return {
        handler,
        options: mergeDefaultEventOptions(options),
    };
}

interface EventListenerProps {
    target: object | string;
}

export default class EventListener extends React.Component<EventListenerProps, any> {
    componentDidMount() {
        this.applyListeners(on);
    }

    componentDidUpdate(prevProps) {
        this.applyListeners(off, prevProps);
        this.applyListeners(on);
    }

    componentWillUnmount() {
        this.applyListeners(off);
    }

    applyListeners(onOrOff, props = this.props) {
        const { target } = props;

        if (target) {
            let element = target;

            if (typeof target === 'string') {
                element = window[target];
            }

            forEachListener(props, onOrOff.bind(null, element));
        }
    }

    render() {
        return this.props.children || null;
    }
}
