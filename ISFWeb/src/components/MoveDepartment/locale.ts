import i18n from '@/core/i18n';

export default i18n([
    [
        '移动部门',
        '移動部門',
        'Move department',
    ],
    [
        '确定',
        '確定',
        'OK',
    ],
    [
        '取消',
        '取消',
        'Cancel',
    ],
    [
        '您可以将部门 “',
        '您可以將部門 “',
        'You can move department "',
    ],
    [
        '” 移动至以下选中的部门下面：',
        '” 移動至以下選中的部門：',
        '" to departments below:',
    ],
    [
        '正在替换存储位置，请稍候...',
        '正在替換儲存位置，請稍候...',
        'Replacing now, please wait ...',
    ],
    [
        '移动部门 “${name}” 至部门 “${targetDepName}” 成功',
        '移動部門 “${name}” 至部門 “${targetDepName}” 成功',
        'Department "${name}" was successfully moved to department "${targetDepName}"',
    ],
    [
        '移动部门 “${name}” 至组织 “${targetDepName}” 成功',
        '移動部門 “${name}” 至組織 “${targetDepName}” 成功',
        'Department "${name}" was successfully moved to organization "${targetDepName}"',
    ],
    [
        '源组织路径 “${srcDepPath}”；新组织路径 “${newDepPath}”',
        '源組織路徑 “${srcDepPath}”；新組織路徑 “${newDepPath}”',
        'Source location “${srcDepPath}”; New location “${newDepPath}”',
    ],
    [
        '移动部门成功',
        '移動部門成功',
        'The department is moved successfully',
    ],
    [
        '编辑部门 “${depName}”的存储位置 成功',
        '編輯部門 “${depName}” 的存儲位置成功',
        'The storage location of Department “${depName}” is changed successfully',
    ],
    [
        '存储位置 “${storage}”',
        '儲存位置 “${storage}”',
        'Storage location "${storage}"',
    ],
    [
        '编辑用户 "${displayName}(${loginName})"的存储位置 成功',
        '編輯使用者 "${displayName}(${loginName})"的儲存位置 成功',
        'The storage location of User “${displayName}(${loginName})” is changed successfully',
    ],
    [
        '无法移动“${depName}”，此部门已不存在。',
        '無法移動“${depName}”，此部門已不存在。',
        'Failed to move "${depName}", as this department does not exist.',
    ],
    [
        '无法移动“${depName}”，您选中的目标部门“${targetDepName}”已不存在，请重新选择。',
        '無法移動“${depName}”，您選中的目標部門“${targetDepName}”已不存在，請重選。',
        'Failed to move "${depName}", as target department "${targetDepName}" does not exist, please make a new selection. ',
    ],
    [
        '无法移动“${depName}”，您选中的目标部门下存在与待移动部门同名的子部门，请重新选择或修改部门名称。',
        '無法移動“${depName}”，您選中的目標部門下存在與待移動部門同名的子部門，請重選或變更部門名稱。',
        'Failed to move "${depName}", as the same name department exists in target department. Please make a new selection or change department name.',
    ],
    [
        '移动部门“${depName}”失败。错误原因：${messge}',
        '移動部門“${depName}”失敗。錯誤原因：${messge}',
        'Failed to move "${depName}". Cause：${messge}',
    ],
    [
        '目标部门的存储位置已不可用，无法替换，仍保留部门当前的存储位置。',
        '目標部門的 儲存位置已不可用，無法替換，仍保留部門當前的儲存位置。',
        'Replace failed. The target location does not exist. This department will remain in the same location.',
    ],
    [
        '替换',
        '替換',
        'OK',
    ],
    [
        '不替换',
        '不替換',
        'Cancel',
    ],
])