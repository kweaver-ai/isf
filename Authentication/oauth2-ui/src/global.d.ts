declare module 'classnames';

declare module "*.gif" {
    const src: string;
    export default src;
}

declare module "*.jpg" {
    const src: string;
    export default src;
}

declare module "*.jpeg" {
    const src: string;
    export default src;
}

declare module "*.png" {
    const src: string;
    export default src;
}

declare module '*.svg' {
    import * as React from 'react';
  
    export const ReactComponent: React.FunctionComponent<
      React.SVGProps<SVGSVGElement> & {
        title?: string;
        className?: string;
        style?: Record<string, any>;
      }
    >;
  
    export default ReactComponent;
  }

declare module "*.css" {
    const classes: { readonly [key: string]: string };
    export default classes;
}

