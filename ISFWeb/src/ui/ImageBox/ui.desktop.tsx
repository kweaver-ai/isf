import React from 'react';
import __ from './locale';
import { LoadState } from './helper';
import ImageBoxBase from './ui.base';

import loading from './assets/loading.gif';
import warning from './assets/warning.png';

export default class ImageBox extends ImageBoxBase {
    /**
     * 根据加载状态渲染
     */
    render() {
        return this.state.loadState === LoadState.COMPLETE ?
            (
                <img src={this.state.src} />
            ) :
            this.state.loadState === LoadState.LOADING ?
                (
                    <div>
                        <img src={loading} />
                        <p>{__('正在加载...')}</p>
                    </div>
                ) :
                this.state.loadState === LoadState.ERROR ?
                    (
                        <div>
                            <img src={warning} />
                            <p>{__('加载失败')}</p>
                        </div>
                    ) :
                    null
    }
}