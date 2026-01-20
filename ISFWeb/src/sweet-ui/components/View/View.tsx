import React from 'react';
import classnames from 'classnames';
import styles from './styles';

interface ViewProps extends Testable {
    className?: string;

    /**
     * 是否是行内显示，为true则渲染为span，否则渲染为div
     */
    inline?: boolean;

    /**
     * CSS样式
     */
    style?: React.CSSProperties;

    /**
     * 挂载到DOM上时执行
     */
    onMounted?: (node: HTMLElement) => void;

    [key: string]: any;
}

/**
 * 行内视图元素
 * @param props 所有props会传递给span元素
 * @returns 返回span元素
 */
const InlineView: React.FunctionComponent<ViewProps> = ({ children, onMounted, ...props }) => <span ref={onMounted} {...props}>{children}</span>

/**
 * 块级视图元素
 * @param props 所有props会传递给div元素
 * @returns 返回div元素
 */
const BlockView: React.FunctionComponent<ViewProps> = ({ children, onMounted, ...props }) => <div ref={onMounted} {...props}>{children}</div>

const View: React.FunctionComponent<ViewProps> = ({ children, className, style, onMounted, testID, inline = false, ...restProps }) => {
    /**
     * 根据是否是行内，渲染不同的元素
     */
    const ViewElement = inline ? InlineView : BlockView

    return (
        <ViewElement
            className={
                classnames(styles['view'],
                    {
                        [styles['block']]: !inline,
                        [styles['inline']]: !!inline,
                    },
                    className,
                )
            }
            style={style}
            onMounted={onMounted}
            data-test-id={testID}
            {...restProps}
        >
            {children}
        </ViewElement>
    );
};

export default View;
