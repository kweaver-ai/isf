
export interface PermType {
  id: string;
  name: string;
  obligations?: ObligationType[];
}
export interface PermissionType {
  allow: PermType[];
  deny: PermType[];
}

export enum PickerRangeEnum {
    /* 仅可选择部门 */
    Dept = 'department',
    /* 仅可选择用户 */
    User = 'user',
    /* 仅可选择用户组 */
    Group = 'group',
    /* 仅可选择应用账号 */
    App = 'app',
    /* 角色 */
    Role= "role"
}

export interface RoleType {
  id: string;
  name: string;
  description?: string;
  resource_type_scopes: {
    unlimited: boolean;
    types: { id: string; name: string; }[]; 
  }
}

export interface VisitorType {
  id: string;
  name: string;
  type: PickerRangeEnum;
  parent_deps?: any[];
}

export interface SchemaType {
  type: string;
  properties: {
    [key: string]: any;
  };
  required?: string[];
}

export interface ObligationConfigType {
  id: string;
  name: string;
  schema: SchemaType;
  description?: string;
  default_value?: any;
  value?: any;
  ui_schema?: any;
}

export interface OperationObligationType {
  operation_id: string; 
  obligation_types: ObligationConfigType[];
}

export interface ObligationType {
  type_id: string;
  id?: string;
  name?: string;
  value: any;
  description?: string;
}

export interface PolicyConfigListType {
  id: string;
  name: string;
  allow: boolean;
  deny: boolean;
  obligations?: ObligationType[];
  description?: string;
  obligation_types?: ObligationConfigType[];
}

export interface ResourceType {
  id: string;
  name: string;
  type: string;
  children?: ResourceType[];
}

export interface PolicyConfigType {
  id: string;
  accessor: VisitorType;
  resource: ResourceType;
  operation: PermissionType;
  expires_at: string;
  condition?: any;
}

export interface addPolicyConfigBodyType {
  accessor: { id: string; type: string; };
  resource: ResourceType;
  operation: PermissionType;
  expires_at: string;
  condition?: any;
}

export interface updatePolicyConfigBodyType {
  operation: PermissionType;
  expires_at: string;
}

export interface permissionConfigType {
  id: string;
  name: string;
}

export enum AuthorizationErrorCodeEnum {
  /*
   * 角色名称重复
   */
  RoleNameConflict = 'Authorization.Conflict.RoleNameConflict',
  /**
   * 角色不存在
   */
  RoleNotFound = 'Authorization.NotFound.RoleNotFound'
}
