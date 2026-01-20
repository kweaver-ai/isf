import React from 'react';
import UserIcon from '../../icons/user.svg';
import DepIcon from '../../icons/dep.svg';
import GroupIcon from '../../icons/group.svg';
import AppCountIcon from '../../icons/app-count.svg';
import RoleIcon from '../../icons/role.svg';
import { PickerRangeEnum } from "@/core/apis/console/authorization/type";

export const getIcon = (type: PickerRangeEnum, size = 13) => {
    switch (type) {
        case PickerRangeEnum.User:
            return <UserIcon style={{ width: size, height: size }}/>
        case PickerRangeEnum.Dept:
            return <DepIcon style={{ width: size, height: size }}/>
        case PickerRangeEnum.Group:
            return <GroupIcon style={{ width: size, height: size }}/>
        case PickerRangeEnum.App:
            return <AppCountIcon style={{ width: size, height: size }}/>
        case PickerRangeEnum.Role:
            return <RoleIcon style={{ width: size, height: size }}/>
    }
};