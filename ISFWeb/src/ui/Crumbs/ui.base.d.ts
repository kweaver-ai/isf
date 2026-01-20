declare namespace UI {
    /**
     * 面包屑组件
     */
    namespace Crumbs {
        /**
         * props
         * @prop crumbs 面包屑列表
         * @prop onClick 点击面包屑回调
         */
        interface Props extends React.Props {
            crumbs: Array<Object>;

            onClick(crumb: Object): void;

            formatter(crumb: Object): string;
        }
    }
}