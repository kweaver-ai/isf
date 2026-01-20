declare namespace UI {
    namespace Stack {
        interface Props extends React.Props<any> {

            /**
             * backgroundColor
             */
            background?: string;

            /**
             * 所占比例
             */
            rate?: number;

            className?: string;

        }
    }
}