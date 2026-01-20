import React from 'react';
import classnames from 'classnames';
import View from '../View';
import styles from './styles';

interface ModalProps {
    /**
     * 覆盖方式
     * @argument fullScreen 全屏覆盖
     * @argument currentContext 覆盖最近的非静态定位元素
     * @default 'fullScreen'
     */
    presentationStyle?: 'fullScreen' | 'currentContext';

    /**
     * 是否透视，目前仅支持右上角透视
     */
    transparentRegion?: {
        width: number;
        height: number;
    };
}

const Modal: React.SFC<ModalProps> = function Modal({
    children,
    presentationStyle = 'fullScreen',
    transparentRegion,
}) {
    const hasTransparentRegion = !!transparentRegion

    return (
        <View className={classnames(styles['root'], [styles[presentationStyle]])}>
            <View className={styles['mask']} />
            {
                hasTransparentRegion && (
                    <View
                        className={styles['transparent-region']}
                        style={{ right: transparentRegion.width, bottom: `calc(100% - ${transparentRegion.height}px)` }}
                    />
                )
            }
            <View
                className={classnames(styles['content'], hasTransparentRegion ? styles['content-with-transparent-region'] : null)}
                style={hasTransparentRegion ? { top: transparentRegion.height } : {}}
            >
                {children}
            </View>
        </View>
    );
};

export default Modal;
