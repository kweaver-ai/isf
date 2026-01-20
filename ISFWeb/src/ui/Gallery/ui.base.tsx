import React from 'react';

export default class GalleryBase extends React.Component<UI.Gallery.Props, any> implements UI.Gallery.Base {
    static defaultProps = {
        /**
         * 分组大小
         */
        groupSize: 8,

        /**
         * 当前显示的组序号
         */
        groupIndex: 0,
    }

    state = {
        currentGroupIndex: this.props.groupIndex,
    }

    galleryGroups = []

    /**
     * 分组缩略图
     */
    updateGallery() {
        if (!Array.isArray(this.props.children)) {
            // 子节点数目为1时 children类型为object 而不是array

            // props不允许修改
            this.props.children = [this.props.children];
        }

        // children 不允许直接修改
        // 通过样式来实现翻页
        this.galleryGroups = _.chunk(this.props.children, this.props.groupSize);

        return this.galleryGroups[this.state.currentGroupIndex];
    }

    /**
     * 向前按钮
     */
    prevHandler() {
        if (this.state.currentGroupIndex > 0) {
            const { currentGroupIndex } = this.state;

            this.setState({
                currentGroupIndex: currentGroupIndex - 1,
            });
        }
    }

    /**
     * 向后按钮
     */
    nextHandler() {
        if (this.state.currentGroupIndex < this.galleryGroups.length - 1) {
            const { currentGroupIndex } = this.state;

            this.setState({
                currentGroupIndex: currentGroupIndex + 1,
            });
        }
    }
}