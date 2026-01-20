declare namespace Console {

    namespace IDNumberBox {

        interface Props extends React.Props<void> {
            /**
             * 身份证号码
             */
            IDNumber: string,

            /**
             * 框宽度
             */
            width: number,

            /**
             * 初始身份证号码
             */
            defaultidcardNumber: string,

            /**
             * 现在的身份证号码
             */
            idcardNumber: string,

            /**
             * 确定按钮是否被点击
             */
            isClickBtn: boolean,

            /**
             * 文本框值改变时触发函数
             */
            onChange: (value: any) => any
        }

        interface State {
            /**
             * 身份证号码框是否可以编辑
             */
            status: number,

            /**
             * 身份证号码验证状态
             */
            IDCardStatus: number,

            /**
             * 身份证号码
             */
            IDNumber: string,

            /**
             * 身份证显示
             */
            showIDNumber: string,
        }
    }
}
