import { reduce, assign } from 'lodash';

export function align(keys = 'bottom left') {
    return reduce(keys.split(/\s+/), (result, key) => {
        switch (key) {
            case 'top':
                return assign(result, { bottom: 33 });

            case 'right':
                return assign(result, { right: 0 });

            case 'bottom':
                return assign(result, { top: 33 });

            case 'left':
                return assign(result, { left: 0 });
        }
    }, {})
}