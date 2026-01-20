import React from 'react'
import classnames from 'classnames'
import { isString } from 'lodash'
import { decorateText } from '@/util/formatters'
import { isBrowser, Browser } from '@/util/browser'
import Title from '../Title/ui.desktop';
import { getTextStyle } from '../helper';
import styles from './styles.desktop'

// 判断是否为Safari浏览器，是时，添加空的伪元素，解决Safari浏览器下显示双tooltip
const isSafari = isBrowser({ app: Browser.Safari });

const Text: React.FunctionComponent<UI.Text.Props> = function Text({
    selectable = true,
    ellipsizeMode = 'tail',
    numberOfChars = 70,
    className,
    titleClassName,
    children,
    fontSize,
    role,
    element,
    titleInline = false,
}) {
    return (
        <Title
            role={role}
            element={element}
            content={isString(children) ? children : ''}
            className={titleClassName}
            inline={titleInline}
        >
            {
                <div className={classnames(
                    styles['text'],
                    {
                        [styles['ellisize-tail']]: !!(ellipsizeMode === 'tail'),
                    },
                    {
                        [styles['selectable']]: selectable,
                    },
                    {
                        [styles['safari']]: isSafari,
                    },
                    className,
                )}
                style={getTextStyle({ fontSize })}
                >
                    <span className={styles['text-layout']}>
                        {
                            ellipsizeMode === 'tail' ? children : decorateText(children, { limit: numberOfChars })
                        }
                    </span>
                </div>
            }
        </Title>
    )
}

export default Text