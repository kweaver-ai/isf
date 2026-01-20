declare namespace Console {
    namespace UserExpireTimeTips {
        interface Props extends React.Props<void> {
            /**
             * 过了有效期限的用户
             */
            expireTimeUsers: ReadonlyArray<any>;

            /**
             * 当前登录的用户
             */
            userid: string;

            /**
             * 【提示】弹窗点击【确定】事件
             */
            completeExpireTimeTips: () => void

            /**
             * 【提示】弹窗点击【取消】事件
             */
            cancelExpireTimeTips: () => void
        }
    }
}
