import React from 'react';
import GalleryBase from './ui.base';
import styles from './style.base.css';

export default class Gallery extends GalleryBase {
    render() {
        return (
            <div className={styles['container']}>
                {/* 使用 LinkIcon*/}
                <button className={styles['btn-prev']} disabled={this.state.currentGroupIndex === 0 ? 'disabled' : ''} onClick={this.prevHandler.bind(this)}></button>
                <ul className={styles['list-ul']}>
                    {this.updateGallery()}
                </ul>
                <button className={styles['btn-next']} disabled={this.state.currentGroupIndex === this.galleryGroups.length - 1 ? 'disabled' : ''} onClick={this.nextHandler.bind(this)}></button>
            </div>
        )
    }
}