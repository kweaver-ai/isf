import React from 'react';
import { noop } from 'lodash';

interface Props {
    crumbs: Array<any>;

    /**
     * 后退按钮是否禁用
     */
    backDisabled?: boolean;

    onClick?(crumb: any): void;

    formatter?(crumb: any): string;

    onChange?(crumbs: Array<any>): any;
}

export default class CrumbsBase extends React.PureComponent<Props, any> {

    static defaultProps = {
        crumbs: [],

        backDisabled: false,

        onClick: noop,

        onChange: noop,

        formatter: (crumb) => crumb,
    }

    state = {
        crumbs: this.props.crumbs,
    }

    componentWillReceiveProps({ crumbs }) {
        if (crumbs !== this.props.crumbs && crumbs !== this.state.crumbs) {
            this.setState({ crumbs: crumbs })
        }
    }

    protected back(crumb) {
        const nextCrumbs = this.state.crumbs.slice(0, this.state.crumbs.indexOf(crumb) + 1)

        this.setState({
            crumbs: nextCrumbs,
        }, () => this.fireOnChangeEvent(nextCrumbs))
    }

    protected clickCrumb(crumb) {
        this.fireOnClickEvent(crumb);
        this.back(crumb)
    }

    private fireOnClickEvent(crumb) {
        this.props.onClick(crumb);
    }

    private fireOnChangeEvent(crumbs) {
        this.props.onChange(crumbs);
    }
}