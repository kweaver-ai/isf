import i18n from '@/core/i18n';

export default i18n([
    [
        '在主域控制器无法正常运行时，备用域控制器可进行用户登录验证、域导入、域同步和反向同步。（最多可添加5个备用域控制器）',
        '在主網域控制站無法正常運行時，備用網域控制站可以進行使用者登入驗證、網域導入、網域同步和反向同步。（至多可以添加5個備用網域控制站）',
        'Backup controller can verify, import, sync users from domain, and perform reverse sync when the main controller is out of service.(You can add at most 5 backup controllers)',
    ],
    [
        '添加备用域控制器',
        '新增備用網域控制站',
        'Add Backup Domain Controller',
    ],
    [
        '备用域控制器',
        '備用網域控制站',
        'Backup Domain Controller',
    ],
    [
        '一',
        '一',
        '1',
    ],
    [
        '二',
        '二',
        '2',
    ],
    [
        '三',
        '三',
        '3',
    ],
    [
        '四',
        '四',
        '4',
    ],
    [
        '五',
        '五',
        '5',
    ],
    [
        '删除',
        '刪除',
        'Delete',
    ],
    [
        '（',
        '（',
        '(',
    ],
    [
        '）',
        '）',
        ')',
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
        '域控制器地址：',
        '網域控制站位址：',
        'Controller Address:',
    ],
    [
        '域控制器端口：',
        '網域控制站埠：',
        'Controller Port:',
    ],
    [
        '域管理员账号：',
        '網域管理員帳戶：',
        'Admin Account:',
    ],
    [
        '域管理员密码：',
        '網域管理員密碼：',
        'Admin Password:',
    ],
    [
        '此输入项不允许为空。',
        '必填欄位。',
        'This field is required.',
    ],
    [
        '您输入的端口号有误。端口必须是 1~65535 之间的数字，且不以 0 开头，请重新输入。',
        '您輸入的連接埠號有誤。埠號必須是數字在1~65535之間，且不能以0開頭，請重新輸入。',
        'Port number error. It should be a number within 1~65535, and cannot be start with 0, please re-enter.',
    ],
    [
        '您的输入有误。域管理员账号不允许使用下列特殊字符 \\ / : * ? " < > | ，且长度范围为 1~128 个字符，请重新输入。',
        '您的輸入有誤。網域管理員帳戶不允許使用下列特殊字元 \\ / : * ? " < > | ，且長度範圍在 1~128 個字元，請重新輸入',
        'Input error. Admin account should be 1~128 characters long, and cannot contain special characters(\\ / : * ? " < > |), please re-enter.',
    ],
    [
        '使用SSL',
        '使用SSL',
        'Use SSL',
    ],
    [
        '已存在相同的域控制器地址，请重新输入',
        '已存在相同的網域控制器位址，請重新輸入',
        'The domain controller address already exists, please re-enter.',
    ],
    [
        '连接LDAP服务器失败，请检查域控制器地址是否正确，或者域控制器是否开启。',
        '連接LDAP伺服器失敗，請檢查網域控制器位址是否正確，或者網域控制器是否已開啟。',
        'Failed to connect the LDAP server. Please ensure that your controller address is correct and the controller is enabled.',
    ],
    [
        '账号或密码不正确。',
        '帳戶或密碼不正確。',
        'Incorrect account or password.',
    ],
    [
        '备用域地址不能和主域地址相同。',
        '備用網域位址不能與主網域位址一樣。',
        'The backup domain address cannot be the same as that of the primary one.',
    ],
    [
        '域不存在',
        '網域不存在',
        'domain not exists',
    ],
    [
        '域控制器 “${domainName}” 已不存在。',
        '域網控制器 “${domainName}” 已不存在。',
        'Domain Controller “${domainName}” does not exist.',
    ],
    [
        '保存',
        '儲存',
        'Save',
    ],
    [
        '删除 备用域控 “${name}” 成功',
        '刪除 備用網域控制器 “${name}” 成功',
        'Delete Backup Domain Controller "${name}" successfully',
    ],
    [
        '添加 备用域控 “${name}” 成功',
        '新增 備用網域控制器 “${name}” 成功',
        'Add Backup Domain Controller "${name}" successfully',
    ],
])