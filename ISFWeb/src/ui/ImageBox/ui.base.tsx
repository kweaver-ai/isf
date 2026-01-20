import React from 'react';
import __ from './locale';

export default class ImageBoxBase extends React.Component<UI.ImageBox.Props, any> implements UI.ImageBox.Base {
    state = {
        /**
         * 图片资源
         */
        src: this.props.src,

        /**
         * 图片加载状态
         */
        loadState: 0,

        /**
         * 图片style样式
         */
        imgStyle: {
            opacity: '0',
        },
    }

    static getDerivedStateFromProps({ src }, prevState) {
        if(src !== prevState.src) {
            return {
                src,
                imgStyle: {
                    opacity: '0',
                },
            }
        }
        return null
    }

    protected handleError(e: React.SyntheticEvent<any>) {
        this.props.onError(e)
    }
}