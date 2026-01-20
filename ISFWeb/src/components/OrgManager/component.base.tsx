import * as React from 'react';
import { noop, uniqBy } from 'lodash';
import { ShareMgnt } from '@/core/thrift';
import WebComponent from '../webcomponent';

interface Props {
    // 已设置的部门
    deps: Array<any>;
    // 是否是管理员
    isManager: boolean;
    // 当前登录的用户id
    userid: string;
    // 选择事件
    onChange: (deps: Array<Department>) => Array<Department>; // deps是要设置的部门
}

interface State {
    // 当前选择的部门
    deps: Array<Department>;
    // 是否显示增加弹窗
    showOrganizationPicker: false;
}

interface Department {
    // 部门名
    departmentName: string;
    // 部门id
    departmentId: string;
    // 部门负责人
    responsiblePerson: any;
}

export default class OrgManagerBase extends WebComponent<any, any> {
    static defaultProps = {
        isManager: false,
        deps: [],
        onChange: noop,
        userid: '',
    }

    state = {
        deps: this.props.deps.concat(),
        showOrganizationPicker: false,
    }

    componentDidMount() {
        if (!(this.props.isManager && typeof (this.props.deps[0]) !== 'string')) {
            this.getDirectDeps().then((depInfos) => {
                this.setState({
                    deps: depInfos.map((depInfo) => {
                        return {
                            departmentId: depInfo.departmentId,
                            departmentName: depInfo.departmentName,
                            responsiblePerson: depInfo.responsiblePerson,
                        }
                    }),
                })
                this.props.onChange(this.state.deps);
            })
        }
    }

    /**
     * 获取直属部门信息
     * @return 返回部门信息的数组
     */
    getDirectDeps(): PromiseLike<Array<Department>> {
        return Promise.all(this.props.deps.map((value) => ShareMgnt('Usrm_GetOrgDepartmentById', [value])))
    }

    /**
     * 获取部门名
     * @param dep 部门
     * @return 返回部门名
     */
    depDataFormatter(dep: Department): string {
        return dep.departmentName;
    }

    /**
     * 设置部门数据
     * @param data 部门信息的数组
     */
    setDepsData(data: Array<Department>): void {
        this.setState({
            deps: data,
        })
        this.props.onChange(data);
    }

    /**
     * 显示增加弹窗
     */
    showOrganizationPicker(): void {
        this.setState({
            showOrganizationPicker: true,
        })
    }

    /**
     * 取消增加
     */
    cancelAddDeps(): void {
        this.setState({
            showOrganizationPicker: false,
        })
    }

    /**
     * 确定本次增加
     * @param deps 部门信息的数组
     */
    confirmAddDeps(deps): void {
        this.setState({
            showOrganizationPicker: false,
            deps: uniqBy(this.state.deps.concat(deps), (value) => value.departmentId),
        }, () => {
            this.props.onChange(this.state.deps);
        })
    }

    /**
     * 数据转换
     * @param
     */
    convererData(value) {
        return {
            departmentId: value.id,
            departmentName: value.name,
            responsiblePerson: value.original.responsiblePerson,
        }
    }

}