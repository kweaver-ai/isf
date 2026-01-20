import i18n from '@/core/i18n';

export default i18n([
    [
        '添加失败，指定的站点无法连接。',
        '新增失敗，指定的站台無法連接。',
        'Add failed. The site cannot be connected.',
    ],
    [
        '该站点名已被占用，请重新输入。',
        '該站台名已被佔用，請重新輸入。',
        'This site name has been taken.',
    ],
    [
        '该站点已存在，不能重复添加。',
        '該站台已存在，不能重複新增。',
        'This site already exists.',
    ],
    [
        '不能添加本站点。',
        '無法新增此站台。',
        'This site cannot be added.',
    ],
    [
        '当前模板已不存在。',
        '當前範本已不存在。',
        'The template does not exist.',
    ],
    [
        '当前模板中的共享者已不存在。',
        '當前範本中的共用者已不存在。',
        'Sharers in the template do not exist. ',
    ],
    [
        '当前模板中的授权者已不存在。',
        '當前範本中的共用者已不存在。',
        'Sharers in the template do not exist. ',
    ],
    [
        'SMTP配置有误，请到“[运营和审计]-[系统配置]-[第三方服务器]” 页面修改配置信息',
        'SMTP設定有誤，請到“[運營與審計]-[系統設定]-[協力廠商伺服器]” 頁面修改設定資訊',
        'SMTP configuration is incorrect, please modify the configuration information in [Dashboards & Logs Audit] - [System Configuration] - [Third-Party Server] page',
    ],
    [
        '测试邮件发送失败，请检查您的网络配置',
        '測試郵件傳送失敗，請檢查您的網路設定',
        'Failed to send test email, please check your network settings.',
    ],
    [
        '该站点已被移除，请重新选择归属站点。',
        '該站台已被移除，請重新選擇歸屬站台。',
        'The site has been removed, please make a new selection.',
    ],
    [
        '该用户名已被占用，请重新输入。',
        '該使用者名已被佔用，請重新輸入。',
        'The site has been removed, please make a new selection.',
    ],
    [
        '用户名不合法，可能字符过长或包含 \\ / : * ? " < > | 特殊字符。',
        '使用者名不合法，可能字元過長或包含 \\ / : * ? " < > | 特殊字元。',
        'User name is illegal, maybe it is too long or contains special characters(\\ / : * ? " < > |).',
    ],
    [
        '该用户名已被管理员占用，请重新输入。',
        '該使用者名稱已被管理員占用，請重新輸入。',
        'This user name has been taken by system. Please re-enter.',
    ],
    [
        '该用户名不可用，请重新输入。',
        '該使用者名稱不可用，請重新輸入。',
        'This user name is not available. Please re-enter.',
    ],
    [
        '添加失败，指定的站点尚未激活授权码。',
        '新增失敗，指定的站台尚未啟動授權碼。',
        'Add failed, as the license has not been activated in this site.',
    ],
    [
        '添加失败，指定的站点未开启多站点模式。',
        '新增失敗，指定的站台未開啟多站台模式。',
        'Add failed, as multi-site mode has not been enabled in this site.',
    ],
    [
        '添加失败，指定的站点已属于其它的分布式系统。',
        '新增失敗，指定的站台已屬於其它分布式系統。',
        'Add failed, as this site belongs to another distributed system.',
    ],
    [
        '添加失败，您输入的站点密钥不正确。',
        '新增失敗，您輸入的站台金鑰不正確。',
        'Add failed, because site key is incorrect.',
    ],
    [
        '组织或部门不存在。',
        '組織或部門不存在。',
        'Department or organization does not exist.',
    ],
    [
        '当前站点不是总站点，可能是未开启多站点模式或已被其他站点添加。',
        '當前站台不是總站台，可能是未開啟多站台模式或已被其他站台添加。',
        'This site is not central site, maybe multi-site mode has not been enabled or it has been added by other sites.',
    ],
    [
        '授权码不存在。',
        '授權碼不存在。',
        'License does not exist.',
    ],
    [
        '未知的授权码类型。',
        '未知的授權碼類型。',
        'Unknown license type.',
    ],
    [
        '授权码 ${licenseCode} 与当前产品型号不匹配。',
        '授權碼 ${licenseCode} 與當前產品型號不匹配。',
        'License ${licenseCode} does not match this product model.',
    ],
    [
        '已存在相同的授权码 ${licenseCode} ，不能重复添加。',
        '已存在相同的授權碼 ${licenseCode} ，不能重複添加。',
        'The same license ${licenseCode} already exists and cannot be added repeatedly.',
    ],
    [
        '当前已存在基本件，无法添加新的基本件。',
        '基本件已存在，無法新增新的基本件。',
        'The essential module already exists and cannot add new essential module.',
    ],
    [
        '已存在测试授权码，请先删除测试授权码再添加正式授权码。',
        '測試授權碼已存在，請先刪除測試授權碼再新增正式授權碼。',
        'Trial license already exists. Please delete trial license before adding official license.',
    ],
    [
        '当前已存在测试授权码，无法添加新的测试授权码。',
        '測試授權碼已存在，無法新增新的測試授權碼。',
        'Trial license already exists and cannot add new trial license.',
    ],
    [
        '已存在正式授权码，无法添加测试授权码。',
        '正式授權碼已存在，無法新增測試授權碼。',
        'Official license already exists and cannot add trial license.',
    ],
    [
        '必须先添加相匹配的基本件，才能添加代理或选件。',
        '必須先添加相匹配的基本件，才能新增代理或選件。',
        'You should add essential module before adding agents or options.',
    ],
    [
        '添加的节点代理授权数，已超过当前产品型号限定的最大节点数。',
        '新增的節點代理數，已超過當前產品型號限定的最大節點數。',
        'Node agents added exceed the max nodes of this model.',
    ],
    [
        '添加的用户代理授权数，已超过当前产品型号限定的最大用户数。',
        '新增的使用者代理授權數，已超過當前產品型號限定的最大使用者數。',
        'User package added exceeds the max users of this model.',
    ],
    [
        '只允许授权一个云盘NAS网关选件。',
        '只允許授權一個雲盤NAS閘道選件。',
        'Only one NAS gateway option is allowed.',
    ],
    [
        '只允许授权一个选件。',
        '只允許授權一個選件。',
        'Only one option is allowed.',
    ],
    [
        '此用户授权包仅允许同时添加 6 个。',
        '此使用者授權碼只允許同時添加 6 個。',
        'This license code only allows 6 additions.',
    ],
    [
        '授权码已被激活。',
        '授權碼已被啟動。',
        'License has been activated.',
    ],
    [
        '激活码与授权码或机器码不匹配。',
        '啟動碼與授權碼或機器碼不匹配。',
        'Activation code does not match license or machine code.',
    ],
    [
        '必须先激活相匹配的基本件，才能激活代理或选件。',
        '必須先啟用相匹配的基本件，才能啟用代理或選件。',
        'You should activate essential module before activating agents or options.',
    ], [
        '短信服务器配置不合法。',
        'SMS伺服器設定不合法。',
        'Illegal SMS server configuration.',
    ],
    [
        '不支持的短信服务器类型。',
        '不支援的SMS伺服器類型。',
        'Unsupported SMS server type.',
    ],
    [
        '连接短信服务器失败。',
        '連接SMS服務器失敗。',
        'The SMS server connection failed.',
    ],
    [
        '不能重复导入组织',
        '不能重複匯入組織',
        'could not import again',
    ],
    [
        '另一个第三方用户组织正在导入',
        '另一個協力廠商使用者組織正在導入',
        'Another third-party user organization is importing',
    ],
    [
        '用户不存在。',
        '使用者不存在。',
        'User does not exist.',
    ],
    [
        '用户已拥有的角色与当前角色存在冲突。',
        '使用者已擁有的角色與當前角色存在衝突。',
        'Roles that user already owns are in conflict with this role.',
    ],
    [
        '角色不存在。',
        '角色不存在。',
        'Role does not exist.',
    ],
    [
        '邮箱地址不合法。',
        '電子郵件位址不合法。',
        'Illegal Email address.',
    ],
    [
        '已存在该角色名。',
        '已存在該角色名。',
        'The role name already exists.',
    ],
    [
        '角色名称不合法，可能字符过长或包含 \\ / : * ? " < > | 特殊字符；',
        '角色名稱不合法，可能字元過長或包含 \\ / : * ? " < > | 特殊字元。',
        'Role name is invalid, maybe it is too long or contains special characters \\ / : * ? " < > |.',
    ],
    [
        '不允许移除您管辖范围以外的部门或组织。',
        '不允許移除您管轄範圍以外的部門或組織。',
        'You are not allowed to remove departments or organizations out your range.',
    ],
    [
        '操作者非法。',
        '操作者非法。',
        'Illegal operator.',
    ],
    [
        '角色成员不存在。',
        '角色成員不存在。',
        'Member does not exist.',
    ],
    [
        '账号或密码不正确。',
        '帳戶或密碼不正確。',
        'Invalid username or password.',
    ],
    [
        '验证码已过期，请重新输入。',
        '驗證碼已過期，請重新輸入。',
        'Invalid verification code, please re-enter.',
    ],
    [
        '验证码不正确，请重新输入。',
        '驗證碼不正確，請重新輸入。',
        'Incorrect verification code, please re-enter.',
    ],
    [
        '您的登录密码已失效，是否立即修改密码？',
        '您的登入密碼已失效，是否立即變更密碼？',
        'Your password is expired, change now?',
    ],
    [
        '您的密码安全系数过低，是否立即修改密码？',
        '您的密碼安全係數過低，是否立即變更密碼？',
        'Weak password, change now?',
    ],
    [
        '无法使用初始密码登录，是否立即修改密码？',
        '無法使用初始密碼登入，是否立即變更密碼？',
        'You are not allowed to log in with initial password, change now?',
    ],
    [
        '您已输错1次密码，连续输错3次将导致账号被锁定。',
        '您已輸錯1次密碼，連續輸錯3次將導致帳戶被鎖定。',
        'Wrong password Once, account will be locked after three failed login attempts.',
    ],
    [
        '您已输错2次密码，连续输错3次将导致账号被锁定。',
        '您已輸錯2次密碼，連續輸錯3次將導致帳戶被鎖定。',
        'Wrong password Twice, account will be locked after three failed login attempts.',
    ],
    [
        '您输入错误次数超过限制，账号已被锁定，${time}分钟内无法登录，请稍后重试。',
        '您輸入錯誤次數超過限制，帳戶已被鎖定，${time}分鐘內無法登入，請稍後重試。',
        'Your account has been locked for ${time} minute(s) due to multiple failed login attempts. Please try again later.',
    ],
    [
        '您的密码已过期，请联系管理员。',
        '您的密碼已過期，請聯繫管理員。',
        'Your password has expired, please contact your administrator.',
    ],
    [
        '您输入错误次数超过限制，账号已被锁定。',
        '您輸入錯誤次數超過限制，帳戶已被鎖定。',
        'Your account is locked due to multiple failed login attempts.',
    ],
    [
        '当前系统未完成初始化配置，请联系系统管理员登录系统。',
        '當前系統未完成初始化設定，請聯繫系統管理員登入系統。',
        'System Initialization has not completed yet, please contact your administrator.',
    ],
    [
        '您是普通用户，无法登录控制台。',
        '您是普通使用者，无法登录控制台。',
        'You are ordinary user and cannot login to console.',
    ],
    [
        '验证码不能为空，请重新输入。',
        '驗證碼不能為空，請重新輸入。',
        'Verification code cannot be empty.',
    ],
    [
        '服务器资源不足，无法访问。',
        '伺服器資源不足，無法存取。',
        'Insufficient server resources to access.',
    ],
    [
        '您受到IP 网段限制，无法登录，请联系管理员。',
        '您受到IP 網路區段限制，無法登入，請聯繫管理員。',
        'Your account is not allowed to login to this IP segment, please contact your administrator.',
    ],
    [
        '非组织文档管理员不能创建组织文档',
        '非組織文件管理員不能建立組織文件',
        'You are not general administrator and cannot create new group document.',
    ],
    [
        '该文档库名已被占用，请重新输入。',
        '該文件庫名已被佔用，請重新輸入。',
        'This name has been taken, please re-enter.',
    ],
    [
        '组织文档对应的ID不存在',
        '組織文件對應的ID不存在',
        'ID for the group document does not exist',
    ],
    [
        '请求的文件或目录不存在',
        '請求的檔案或目錄不存在',
        'The requested file or directory does not exist',
    ],
    [
        '操作失败，请点击确定返回上一级页面后再进入。',
        '操作失敗，請點擊確定返回上一級頁面後再進入。',
        'Operation failed, please click OK to back to the previous page and try again.',
    ],
    [
        '个人文档库已存在。',
        '個人文件庫已存在。',
        'My documents already exists',
    ],
    [
        '入口文档记录不存在',
        '入口文件記錄不存在',
        'The entry file records do not exist',
    ],
    [
        '请求的文件或目录不存在。',
        '請求的文件或目錄不存在。',
        'The requested file or directory does not exist.',
    ],
    [
        '个人文档库不存在',
        '個人文件庫不存在',
        'My document does not exist',
    ],
    [
        '配额空间不足',
        '配額空間不足',
        'Insufficient quota.',
    ],
    [   // 以下 繁体 由 简体 代替
        '该节点已被其他的用户接入过。',
        '此節點已被其他的使用者接入過。',
        'The node has been accessed by other users.',
    ],
    [
        '手动获取CMP授权失败。',
        '手動獲取CMP授權失敗。',
        'Failed to get CMP license manually.',
    ],
    [
        '测试授权使用已达上限，无法再添加测试授权。',
        '測試授權使用已達上限，無法再添加測試授權。',
        'Usage limit of test authorization reached. The test authorization cannot be added again.',
    ],
    [
        '您不是超级管理员，无法修改邮箱或用户名。',
        '您不是超級管理員，無法變更電子郵件或使用者名稱。',
        'Only super admin can change the email address and user name.',
    ],
    [
        '包名长度不能超过150个字符。',
        '封裝名稱長度不能超過150個字元。',
        'The pack name should be less than 150 characters.',
    ],
    [
        '文件已不存在。',
        '檔案已不存在。',
        'File does not exist.',
    ],
    [
        '升级包名的类型验证失败。',
        '升級封裝名稱的類型驗證失敗。',
        'Failed to verify the type of upgrade pack.',
    ],
    [
        '升级包的升级模式验证失败。',
        '升級封裝的升級模式驗證失敗。',
        'Failed to verify the upgrade mode of upgrade pack.',
    ],
    [
        '上传升级包失败。',
        '上傳升級封裝失敗。',
        'Upgrade pack upload failed.',
    ],
    [
        '无法识别文件类型的适配参数。',
        '無法識別檔案類型的適配參數。',
        'Unable to recognize the parameters adapted to the file type.',
    ],
    [
        '删除升级包信息失败。',
        '刪除升級封裝的資訊失敗。',
        'Upgrade pack info deleted failed.',
    ],
    [
        '读取升级包URL失败。',
        '讀取升級封裝URL失敗。',
        'Upgrade pack URL read failed.',
    ],
    [
        '非法参数。',
        '非法參數。',
        'Invalid parameters.',
    ],
    [
        '内部错误',
        '內部錯誤',
        'Internal error',
    ],
    [
        '当前节点不是应用节点。',
        '此節點不是應用程式節點。',
        'The node is node application node.',
    ],
    [
        '没有可用的安装包。',
        '無可用的升級封裝。',
        'No available upgrade package.',
    ],
    [
        '服务已安装。',
        '服務已安裝。',
        'Service is installed.',
    ],
    [
        '服务未安装。',
        '服務未安裝。',
        'Service is not installed.',
    ],
    [
        '将要安装的版本低于已安装版本。',
        '即將安裝的版本低於已安裝的版本。',
        'The version you are about to install is lower than current version.',
    ],
    [
        '安装包损坏。',
        '安裝封裝損毀。',
        'Installation package is corrupt.',
    ],
    [
        '节点离线。',
        '節點已離線。',
        'Node is offline.',
    ],
    [
        '该日期已过期，请重新选择。',
        '該日期已過期，請重新選擇。',
        'The date has expired, please make a new selection.',
    ],
    [
        '该邮箱已被占用。',
        '該電子郵件已被使用。',
        'This Email already exists.',
    ],
    [
        '请输入正确的手机号。',
        '請輸入正確的手機號。',
        'Please enter a valid phone number.',
    ],
    [
        '该手机号已被占用。',
        '該手機號已被佔用。',
        'The mobile number already exists.',
    ],
    [
        '请输入正确的身份证号。',
        '請輸入正確的身份證號。',
        'Please enter a valid ID number.',
    ],
    [
        '备注不能包含 \\ / : * ? " < > | 特殊字符，长度不能超过128个字符。',
        '備註不能包含\\ / : * ? " < > |特殊字元，長度不能超過128個字元。',
        'Note cannot contain space or \\ / : * ? " < > |. Its length cannot exceed 128 characters.',
    ],
    [
        '该身份证号已被占用。',
        '該身份證號已被佔用。',
        'The ID number has been register.',
    ],
    [
        '用户密级不能高于系统密级',
        '使用者密級不能高於系統密級',
        'The security level of user cannot exceed that of the system',
    ],
    [
        '所指定的存储位置已不可用，请更换。',
        '所指定的儲存位置已不可用，請更換。',
        'The specified location is unavailable.',
    ],
])