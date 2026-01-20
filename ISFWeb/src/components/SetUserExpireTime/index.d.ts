
declare namespace Console {
    namespace SetUserExpireTime {
        interface Props extends React.Props<void> {
            /**
             * 当前选择的用户
             */
            users?: ReadonlyArray<any>,

            /**
             * 选中的部门
             */
            dep?: any,

            /**
             * 当前登录的用户
             */
            userid: string,

            /**
             * 是否是在启用时，需要设置有效期限
             */
            shouldEnableUsers?: boolean,

            /**
             * 设置有效期【取消】事件
             */
            onCancel: () => void,

            /**
             * 设置有效期成功事件
             */
            onSuccess: () => void
        }
        interface State {
            /**
             * 选择设置有效期的对象
             * 包括：个人、组织、组织及其子部门
             */
            selected: Range,

            /**
             * 管理员设置的有效期限
             * 从 1970-01-01 开始算，到截止日期之间的时间，单位：毫秒
             */
            expireTime: number,

            /**
             * 报错信息
             */
            errors: ReadonlyArray<any>,

            /**
             * 启用用户时转圈圈状态
             */
            status: number,

            /**
             * 日期组件选中日期时触发
             */
            onChange?: (date: Date) => any;

            /**
             * 无效的有效期
             */
            invalidExpireTime: boolean
        }

        interface Range {
            /**
             * 部门及其子部门
             */
            DEPARTMENT_DEEP,

            /**
             * 当前部门
             */
            DEPARTMENT,

            /**
             * 当前选中的用户
             */
            USERS
        }
    }
}
