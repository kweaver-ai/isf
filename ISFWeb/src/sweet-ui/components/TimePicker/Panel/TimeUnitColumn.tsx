import React from 'react';
import classnames from 'classnames';
import { scrollTo } from '../util';
import styles from './styles';

export interface Unit {
    label: React.ReactText;
    value: number;
    disabled?: boolean;
}

export interface TimeUnitColumnProps {
    open?: boolean;
    units?: Unit[];
    value?: number;
    hideDisabledOptions?: boolean;
    onSelect?: (value: number) => void;
}

function TimeUnitColumn(props: TimeUnitColumnProps) {
    const { units, onSelect, value, hideDisabledOptions, open } = props;

    const ulRef = React.useRef<HTMLUListElement>(null);
    const liRefs = React.useRef<Map<number, HTMLElement | null>>(new Map());

    React.useLayoutEffect(
        () => {
            const li = liRefs.current.get(value!);
            if (li && open !== false) {
                scrollTo(ulRef.current!, li.offsetTop, 120);
            }
        },
        // eslint-disable-next-line react-hooks/exhaustive-deps
        [ value ],
    );

    React.useLayoutEffect(
        () => {
            if (open) {
                const li = liRefs.current.get(value!);
                if (li) {
                    scrollTo(ulRef.current!, li.offsetTop, 0);
                }
            }
        },
        // eslint-disable-next-line react-hooks/exhaustive-deps
        [ open ],
    );

    return (
        <ul
            className={styles['column']}
            ref={ulRef}
        >
            {units!.map((unit) => {
                if (hideDisabledOptions && unit.disabled) {
                    return null;
                }

                return (
                    <li
                        key={unit.value}
                        ref={(element) => {
                            liRefs.current.set(unit.value, element);
                        }}
                        className={classnames(styles['cell'], {
                            [styles['cell-disabled']]: unit.disabled,
                            [styles['cell-selected']]: value === unit.value,
                        })}
                        onMouseDown={() => {
                            if (unit.disabled) {
                                return;
                            }
                            onSelect!(unit.value);
                        }}
                    >
                        <div className={classnames(styles['cell-inner'], { [styles['disabled']]: unit.disabled })}>
                            {unit.label}
                        </div>
                    </li>
                );
            })}
        </ul>
    );
}

export default TimeUnitColumn;
