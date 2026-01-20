import { PickerRangeEnum } from "@/core/apis/console/authorization/type";

export enum RoleClassEnum {
    Business = 'business',
    System = 'system',
    User = 'user',
}

export enum OperationTypeEnum {
    AddRole = 'add-role',
    EditRoleInfo = 'edit-role-info',
    EditRolePerm = 'edit-role-perm', 
}

export interface RoleInfoType {
    id: string;
    name: string;
    type: string;
    parent_deps?: any[];
}

export interface VisitorType {
    id: string;
    name: string;
    type: PickerRangeEnum;
}

export enum InnerRoleIdEnum {
    /**
      * 数据管理员
      */
    DataAdmin = "00990824-4bf7-11f0-8fa7-865d5643e61f",
    /**
     * AI管理员 
     */
    AIAdmin = "3fb94948-5169-11f0-b662-3a7bdba2913f",
    /**
     * 应用管理员
     */
    AppAdmin = "1572fb82-526f-11f0-bde6-e674ec8dde71"

}

export enum PermOperationEnum {
    SetAllResource = "set-all-resource",
    EditAllResource = "edit-all-resource",
    SetSpecifiedResource = "set-specified-resource",
    EditSpecifiedResource = "edit-specified-resource"
}
