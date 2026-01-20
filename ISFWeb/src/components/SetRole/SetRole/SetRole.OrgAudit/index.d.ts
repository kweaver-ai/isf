declare namespace Console {
    namespace SetOrgAudit {
        interface Props extends React.Props<void> {
            /**
             * 登录用户id
             */
            userid: string;

            /**
             * 编辑角色
             */
            editRateInfo: Object;

            /**
             * 选择设置的角色
             */
            roleInfo: Object;

            /**
             * 选择配置角色的用户
             */
            userInfo: Object;

            /**
             * 直属部门
             */
            directDeptInfo: Object;

            /**
             * 确定
             */
            onConfirmSetRoleConfig: () => void;

            /**
             * 取消
             */
            onCancelSetRoleConfig: () => void;
        }

        interface State {
            /**
             * 管辖的部门
             */
            selectDeps: Array<any>;
        }
    }
}