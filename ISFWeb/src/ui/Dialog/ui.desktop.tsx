import React from 'react';
import Mask from '../Mask/ui.desktop';
import Header from '../Dialog.Header/ui.desktop';
import DialogMain from '../Dialog.Main/ui.desktop';
import Footer from '../Dialog.Footer/ui.desktop';
import DialogButton from '../Dialog.Button/ui.desktop';
import DialogBase from './ui.base';
import styles from './styles.desktop';

export default class Dialog extends DialogBase {

    static Header = Header;

    static Main = DialogMain;

    static Footer = Footer;

    static Button = DialogButton;

    render() {
        return (
            <div role={this.props.role}>
                {
                    // 通过props传入hide = true实现组件隐藏， 兼容对话框中的flash问题
                    this.props.hide ? null : <Mask />
                }
                <div
                    ref="container"
                    className={ styles['container'] }
                    style={ this.props.hide ? { top: '100%', left: '100%' } : { width: this.props.width, top: this.state.top, left: this.state.left } }
                    onMouseDown={ this.props.draggable && this.startDrag.bind(this) }
                    onMouseUp={ this.props.draggable && this.endDrag.bind(this) }
                >
                    {
                        this.props.children
                    }
                </div>
            </div>
        )
    }
}