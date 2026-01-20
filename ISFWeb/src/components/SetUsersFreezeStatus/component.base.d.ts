// eslint-disable-next-line spaced-comment
/// <reference path="@/core/thrift/sharemgnt/sharemgnt.d.ts"/>
declare namespace Components {
    namespace SetUsersFreezeStatusBase {

        interface Props {

            /**
             *  是否冻结用户 true：冻结，false：解冻
             */
            freezeStatus: boolean;

            /**
             *  用户id
             */
            userid: string;

            /**
             *  完成
             */
            onComplete: () => any;

            /**
             * 成功事件
             */
            onSuccess: () => any;
        }

        interface State {

            /**
             *  选中要处理的用户
             */
            selecteds: Array<any>;

            /**
             *  当前状态
             */
            status: number;

        }

        interface Data {

            /**
             *  id 用户或者部门
             */
            id: string;

            /**
             *  用户名或者部门名
             */
            name: string;

            /**
             *  数据类型
             */
            type: number;

            /**
             * 父节点名称
             */
            parentName: string;

            /**
             * 源数据
             */
            data: Core.ShareMgnt.ncTUsrmUserInfo | any
        }
    }
}