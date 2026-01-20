declare namespace UI {
    namespace ComboSearchBox {
        interface Props extends React.Props<any> {
            /**
             * 展开的下拉列表条目
             */
            keys?: ReadonlyArray<string>;

            /**
             * 搜索框为空时，默认显示文字
             */
            placeholder?: string;

            /**
             * 界面样式
             */
            className?: string;

            /**
             * 输入的搜索内容
             */
            searchKey?: string;

            /**
             *根据输入的搜索内容和用户的选择，拼接的搜索条件
             */
            searchValue?: Array<object>;

            /**
             * 自定义弹出的下拉列表中搜索关键字和搜索值之间的组合关系
             */
            renderOption: (key, value) => any;

            /**
             * 自定义生成的搜索条目搜索关键字和搜索值之间的组合关系
             */
            renderComboItem: (key, value) => any;

            /**
             * 搜索结果发生改变时通知父元素更新显示列表
             */
            onComboChange: (searchValue, value) => any;

            /**
             * 弹出层嵌入的位置
             */
            element: HTMLElement;
        }

        interface State {
            /**
             * 输入的搜索内容
             */
            value: string;

            /**
             *根据输入的搜索内容和用户的选择，拼接的搜索条件
             */
            searchValue: Array<object>;

            /**
             * 外围包绕的聚焦元素，用于PopMenu定位
             */
            searchAnchor: Element | null;

            /**
             * 搜索框聚焦和失焦状态下的样式
             */
            isSearchFocus: boolean;

            /**
             * 是否打开搜索框下拉列表
             */
            isSearchMenu: boolean;
        }
    }
}
