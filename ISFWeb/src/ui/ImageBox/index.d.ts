declare namespace UI {
    namespace ImageBox {
        type Props = {
            /**
             * 图片资源url
             */
            src: string;

            /**
             * 横向全屏显示
             */
            fullScreen: boolean;

            /**
             * 预览图片出错
             */
            onError: (e: React.SyntheticEvent<any>) => void;
        }

        type State = {
            /**
             * 图片资源
             */
            src: string;

            /**
             * 图片加载状态
             */
            loadState: number;

            /**
             * 图片style样式
             */
            imgStyle: {
                maxWidth: string;

                maxHeight: string;

                width?: string;

                height?: string;
            };
        }

        interface Base {
            props: Props;

            state: State;
        }
    }
}