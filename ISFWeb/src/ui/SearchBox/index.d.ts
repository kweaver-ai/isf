declare namespace UI {
    namespace SearchBox {
        interface Props extends React.Props<any> {
            /**
             * 是否被禁用
             */
            disabled?: boolean;

            /**
             * class
             */
            className?: string;

            /**
             * role
             */
            role?: string;

            /**
             * 内联样式
             */
            style?: {
                [key: string]: string | number,
            };

            /**
             * 宽度
             */
            width?: number | string;

            /**
             * 关键字
             */
            value: string;

            /**
             * 搜索图标
             */
            icon?: string;

            /**
             * 是否自动聚焦
             */
            autoFocus?: boolean;

            /**
             * 搜索字符最大长度
             */
            maxLength?: number;

            /**
             * 输入限制函数
             */
            validator?: (value) => boolean;

            /**
             * 占位提示
             */
            placeholder: string;

            /**
             * 触发搜索的延迟时间
             */
            delay: number;

            /**
             * 文本框变化时触发
             */
            onChange?: (value: string) => any;

            /**
             * 搜索函数
             * @param key 关键字
             */
            loader?: (key: string) => Promise<any>;

            /**
             * 搜索开始时触发
             * @param result 搜索结果
             */
            onFetch?: (key: string, process: Promise<any>) => any;

            /**
             * 搜索完成时触发
             * @param result 搜索结果
             */
            onLoad?: (result: any) => any;

            /**
             * 加载失败时触发
             */
            onLoadFailed?: (ex: any) => any;

            /**
             * 聚焦时时触发
             * @param event 文本框对象
             */
            onFocus?: (event: FocusEvent) => any;

            /**
             * 失去焦点时触发
             * @param event 文本框对象
             */
            onBlur?: (event: FocusEvent) => any;

            /**
             * 键盘输入时触发
             */
            onKeyDown?: (event: KeyboardEvent) => any;

            /**
             * 点击时触发
             */
            onClick?: (event: MouseEvent) => any;

            /**
             * 回车触发
             */
            onEnter?: (event: KeyboardEvent) => any;
        }

        interface State {
            /**
             * 输入值
             */
            value: string;
        }
    }
}