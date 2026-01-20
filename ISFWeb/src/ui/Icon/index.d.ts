declare namespace UI {
    namespace Icon {
        interface Props {
            /**
             * role
             */
            role?: string;

            /**
             * 图标url，支持base64
             */
            url: string;
            /**
             * 图标尺寸
             */
            size: number | string;
        }
    }
}