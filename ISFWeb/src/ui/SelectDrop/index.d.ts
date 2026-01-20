declare namespace UI {
    namespace SelectDrop {
        /**
        * 下拉选项
        */
        type Option = any;

        interface Props extends React.ClassAttributes<void>, UI.PopOver.Props {

            /**
             * 下拉选项数组
             */
            options: ReadonlyArray<Option>

            /**
             * 字体图标对应的code
             */
            icon?: string;

            /**
             * 字体图标大小
             */
            size?: number;

            /**
             * 字体图标颜色
             */
            color?: string;

            /**
             * IE8/9下字体图标对应的图片
             */
            iconFallback?: string;

            /**
             * 下拉选项对齐方式：left,center,right
             */
            align?: 'left' | 'center' | 'right';

            /**
             * 默认的选中项，数据必须来源于options中的一项
             */
            defaultOption?: Option;

            /**
             * 格式化label值
             */
            labelFormatter?: (option: Option) => string;

            /**
             * 选中下拉选项时触发
             */
            onChange?: (option: Option) => void;

        }

        interface State {
            /**
             * 当前选中项
             */
            option: Option;
        }
    }

}