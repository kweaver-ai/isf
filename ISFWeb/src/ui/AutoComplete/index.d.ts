declare namespace UI {
    namespace AutoComplete {
        interface Props extends React.Props<any> {
            /**
             * 角色
             */
            role?: string;

            /**
             * 控件宽度
             */
            width?: number;

            /**
             * 组件被禁用
             */
            disabled?: boolean;

            /**
             * class，只针对搜索框生效
             */
            className?: string;

            /**
             * 内联样式，只针对搜索框生效
             */
            style?: Object;

            /**
             * 输入限制函数
             */
            validator?: (value) => boolean;

            /**
             * 搜索图标
             */
            icon: React.ReactElement<any>;

            /**
             * 是否自动聚焦
             */
            autoFocus: boolean;

            /**
             * 搜索关键字
             */
            value?: string;

            /**
             * 占位提示
             */
            placeholder: string;

            /**
             * 搜索结果为空提示
             */
            missingMessage: string;

            /**
             * 触发搜索的延迟时间
             */
            delay: number;

            /**
             * 搜索函数
             * @param key 搜索关键字
             */
            loader: (key: string) => any;

            /**
             * 加载失败
             */
            loaderFailed: (ex: any) => any;

            /**
             * 数据加载完成触发
             */
            onLoad: (results: any) => any;

            /**
             * 数据加载完成触发
             */
            onFetch?: (key: string, process: Promise<any>) => any;

            /**
             * 文本框变化时触发
             */
            onChange?: (key: string) => any;

            /**
             * 聚焦时时触发
             * @param event 文本框对象
             */
            onFocus?(event: FocusEvent): any;

            /**
             * 失去焦点时触发
             * @param event 文本框对象
             */
            onBlur?(event: FocusEvent): any;

            /**
             * 键盘输入时触发
             */
            onKeyDown?: (event: KeyboardEvent) => any;

            /**
             * 回车触发
             */
            onEnter?: (event: KeyboardEvent, selectIndex: number) => any;
        }

        interface State {
            value: string; // 搜索关键字

            active: boolean; // 下拉菜单是否激活

            status: number; // 组件状态,

            selectIndex: number;   // 当前选中项

            keyDown: number;   // 按下键盘的状态
        }

    }
}