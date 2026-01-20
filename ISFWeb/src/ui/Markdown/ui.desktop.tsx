import React from 'react';
import MarkdownBase from './ui.base';
import styles from './styles.desktop';

export default class Markdown extends MarkdownBase {
    render() {
        return (
            <div className={styles['markdown']} dangerouslySetInnerHTML={{ __html: this.props.children }}></div>
        )
    }
}