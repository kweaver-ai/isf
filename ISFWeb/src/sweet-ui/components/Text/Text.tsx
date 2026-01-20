import React from 'react';
import { pick } from 'lodash'
import View from '../View'

/**
 * 悬浮提示
 * 使用`boolean`类型时，`true`表示使用文本控件内的`textContent`作为浮动提示，`false`表示不显示。
 * 使用`string`类型时，浮动提示根据传入的字符串来显示。
 * 使用`JSX`类型时，浮动提示可以接受任意JSX片段，此时可以显示更丰富的提示内容。
 */
type Tooltip = boolean | string | JSX.Element

/**
 * 支持的CSS样式
 */
type ValidStyleProps =
    | 'fontFamily'
    | 'fontSize'
    | 'color'
    | 'fontWeight'
    | 'lineHeight'
    | 'textAlign'
    | 'whiteSpace'
    | 'letterSpacing'

/**
 * 支持的字体样式
 */
type TextStyle = Pick<React.CSSProperties, ValidStyleProps>

/**
 * 缩略显示模式
 */
enum EllipsizeMode {
    /**
     * 从中间省略
     */
    Middle = 'middle',

    /**
     * 省略末尾
     */
    Tail = 'tail',
}

interface TextProps extends Testable {
    /**
     * 浮动提示
     */
    tooltip?: Tooltip;

    /**
     * 缩略显示模式
     */
    ellipsizeMode?: EllipsizeMode;

    /**
     * CSS样式
     */
    textStyle?: TextStyle;

    /**
     * 是否是行内文本
     */
    inline?: boolean;
}

/**
 * 有效的文本样式
 */
const validTextStyle: ReadonlyArray<ValidStyleProps> = [
    'fontFamily',
    'fontSize',
    'color',
    'fontWeight',
    'lineHeight',
    'textAlign',
    'whiteSpace',
    'letterSpacing',
]

/**
 * 过滤出支持的样式的style对象
 * @param style 传入的style对象
 * @returns 返回只保留支持的样式的style对象
 */
const filterValidTextStyle = (style: Partial<React.CSSProperties> = {}): TextStyle => {
    return pick(style, validTextStyle)
}

const Text: React.FunctionComponent<TextProps> = ({ inline, testID, tooltip, ellipsizeMode, textStyle, children, role }) => {
    return (
        <View
            role={role}
            inline={inline}
            style={filterValidTextStyle(textStyle)}
            testID={testID}
        >
            {
                children
            }
        </View>
    )
}

export default Text