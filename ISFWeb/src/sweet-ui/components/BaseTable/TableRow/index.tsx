import React from 'react';
import classnames from 'classnames';
import { isFunction } from 'lodash';

export interface TableRowProps extends React.ClassAttributes<void> {
    className?: string;
    hoverClassName?: string;
    dropOverClassName?: string;
    isOver?: boolean;
    index?: number;
    onClick?: (event: React.MouseEvent<HTMLTableRowElement>) => void;
    onDoubleClick?: (event: React.MouseEvent<HTMLTableRowElement>) => void;
    onEnter?: (event: React.MouseEvent<HTMLTableRowElement>) => void;
    onLeave?: (event: React.MouseEvent<HTMLTableRowElement>) => void;
    onDragEnd?: (dragIndex: number, hoverIndex: number) => void;
}

export interface TableRowState {
    hover: boolean;
}

export default class TableRow extends React.PureComponent<TableRowProps, TableRowState> {
    state = {
        hover: false,
    };

    /**
     * 点击行时触发
     */
    handleRowClicked = (event: React.MouseEvent<HTMLTableRowElement>) => {
        isFunction(this.props.onClick) && this.props.onClick(event);
    };

    /**
     * 双击行时触发
     */
    handleRowDoubleClicked = (event: React.MouseEvent<HTMLTableRowElement>) => {
        isFunction(this.props.onDoubleClick) && this.props.onDoubleClick(event);
    };

    /**
     * 鼠标移入行时触发
     */
    handleRowEnter = (event: React.MouseEvent<HTMLTableRowElement>) => {
        this.setState({ hover: true });
        isFunction(this.props.onEnter) && this.props.onEnter(event);
    };

    /**
     * 鼠标移出行时触发
     */
    handleRowLeave = (event: React.MouseEvent<HTMLTableRowElement>) => {
        this.setState({ hover: false });
        isFunction(this.props.onLeave) && this.props.onLeave(event);
    };

    render() {
        const {
            className,
            hoverClassName = '',
            children,
            dropOverClassName = '',
            isOver,
        } = this.props;
        const { hover } = this.state;

        return  (
            <tr
                className={classnames(className, { [hoverClassName]: hover }, { [dropOverClassName]: isOver })}
                onClick={this.handleRowClicked}
                onDoubleClick={this.handleRowDoubleClicked}
                onMouseEnter={this.handleRowEnter}
                onMouseLeave={this.handleRowLeave}
            >
                {children}
            </tr>
        );
    }
}
