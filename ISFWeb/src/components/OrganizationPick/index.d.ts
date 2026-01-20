declare namespace Components {
    namespace OrganizationPick {
        interface Props extends React.Props<void> {
            /**
             * 点击确定事件
             */
            onConfirm: (data: Array<Object>) => any;

            /**
             * 点击取消事件
             */
            onCancel: () => any;

            /**
             * 当前管理员id
             */
            userid: string;

            /**
             * 是否加载用户
             */
            selectType?: Array<number>;

            /**
             * 初始值
             */
            data?: Array<any>;

            /**
             * 选择的部门发生变化(添加/删除)
             */
            onSelectionChange: (data: Array<any>) => void;

            /**
             * 数据转换内部数据结构
             */
            converterIn?: (x) => Node;

            /**
             * 数据转换外部数据结构
             */
            convererOut: (Node) => any;

            /**
             * tree高度
             */
            height: string | number;

            /**
             * tree宽度
             */
            width: string | number;
        }

        interface State {
            /**
             * 新加的部门
             */
            data: Array<Node>
        }

        /**
         * 部门信息
         */
        interface Node {
            /**
             * 名称
             */
            name: string;

            /**
             * 部门id
             */
            id: string;

            /**
             * 类型
             */
            type: number;

            /**
             * 原始数据
             */
            original?: any;
        }
    }
}