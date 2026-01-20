declare namespace Core {
    namespace EACP {

        /******************************数据类型定义********************************/

        /********************************** 函数声明*****************************/

        /**
         * 设置消息设置
         */
        type SetMessageNotifyStatus = Core.APIs.ThriftAPI<
            boolean,
            void
        >
        /**
         * 获取消息设置
         */
        type GetMessageNotifyStatus = Core.APIs.ThriftAPI<
            void,
            boolean
        >
        /**
         * 清除超出实名共享的权限配置
         */
        type ClearPermOutOfScope = Core.APIs.ThriftAPI<
            void,
            void
        >
        /**
         * 清除超出范围的历史匿名共享
         */
        type clearPermOutOfScope = Core.APIs.ThriftAPI<
            void,
            void
        >
    }
}