import React from 'react';
import GalleryImageBase from './ui.base';

export default class GalleryImage extends GalleryImageBase {
    render() {
        return (
            <img src={this.props.src} />
        )
    }
}