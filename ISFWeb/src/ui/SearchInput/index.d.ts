declare namespace UI {
    namespace SearchInput {

        interface Props extends React.Props<any> {
            /**
             * 关键字
             */
            value?: string | number;

            /**
             * 是否被禁用
             */
            disabled?: boolean;

            /**
             * 占位提示
             */
            placeholder?: string;

            /**
             * 是否自动聚焦
             */
            autoFocus?: boolean;

            /**
             * 触发搜索的延迟时间
             */
            delay?: number;

            /**
             * 搜索字符最大长度
             */
            maxLength?: number;

            /**
             * 验证输入
             */
            validator?: (value: any) => boolean;

            /**
             * 搜索函数
             * @param key 关键字
             */
            loader: (key: string) => Promise<any>;

            /**
             * 文本框变化时触发
             */
            onChange?: (value: string) => any;

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
             * 搜索加载失败
             */
            onLoadFailed?: (ex: any) => any;

            /**
             * 点击时触发
             */
            onClick?: (event: MouseEvent) => any;

            /**
             * 聚焦时时触发
             * @param ref 文本框对象
             */
            onFocus?: (event: FocusEvent) => any;

            /**
             * 失去焦点时触发
             * @param ref 文本框对象
             */
            onBlur?: (event: FocusEvent) => any;

            /**
             * 回车事件
             */
            onEnter?: (event: KeyboardEvent) => void;

            /**
             * 键盘输入时触发
             */
            onKeyDown?: (event: KeyboardEvent) => any;
        }

        interface State {
            /**
             * 文本框当前值
             */
            value?: string | number;
        }
    }
}